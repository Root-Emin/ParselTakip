package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	iamUC "github.com/masterfabric-go/masterfabric/internal/application/iam/usecase"
	iamModel "github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	tenantModel "github.com/masterfabric-go/masterfabric/internal/domain/tenant/model"
	infraAuth "github.com/masterfabric-go/masterfabric/internal/infrastructure/auth"
	pgIam "github.com/masterfabric-go/masterfabric/internal/infrastructure/postgres/iam"
	pgTenant "github.com/masterfabric-go/masterfabric/internal/infrastructure/postgres/tenant"
	"github.com/masterfabric-go/masterfabric/internal/shared/config"
	"github.com/masterfabric-go/masterfabric/internal/shared/crypto"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// seedSuperAdmin ensures a fixed super-admin user (with every system role) exists
// on every startup. It is fully idempotent: existing rows are reused, never reset,
// so an operator who later changes the password is not overwritten on restart.
func seedSuperAdmin(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
	db *pgxpool.Pool,
	enc *crypto.Encryptor,
	eventBus events.EventBus,
) error {
	orgRepo := pgTenant.NewOrgRepo(db)
	appRepo := pgTenant.NewAppRepo(db)
	userRepo := pgIam.NewUserRepo(db, enc)
	roleRepo := pgIam.NewRoleRepo(db)
	orgUserRepo := pgIam.NewOrgUserRepo(db)
	systemRoleRepo := pgIam.NewSystemRoleRepo(db)
	auth := infraAuth.NewJWTService(cfg.JWT, cfg.Security.BcryptCost)
	ensureRoles := iamUC.NewEnsureOrgRolesUseCase(systemRoleRepo, roleRepo)

	// 1. Default organization.
	org, err := orgRepo.GetBySlug(ctx, cfg.SeedAdmin.OrgSlug)
	if err != nil || org == nil {
		org = &tenantModel.Organization{
			Name:   cfg.SeedAdmin.OrgName,
			Slug:   cfg.SeedAdmin.OrgSlug,
			Status: tenantModel.OrgStatusActive,
		}
		if err := orgRepo.Create(ctx, org); err != nil {
			return fmt.Errorf("create seed org: %w", err)
		}
		log.Info("seed: organization created", "slug", org.Slug, "id", org.ID)
	}

	// 2. Default app (domain entities are scoped by app_id).
	app, err := appRepo.GetBySlug(ctx, org.ID, cfg.SeedAdmin.AppSlug)
	if err != nil || app == nil {
		app = &tenantModel.App{
			OrganizationID: org.ID,
			Name:           cfg.SeedAdmin.AppName,
			Slug:           cfg.SeedAdmin.AppSlug,
			Status:         tenantModel.AppStatusActive,
			SLATier:        "enterprise",
		}
		if err := appRepo.Create(ctx, app); err != nil {
			return fmt.Errorf("create seed app: %w", err)
		}
		log.Info("seed: app created", "slug", app.Slug, "id", app.ID)
	}

	// 3. Sync the system roles into the organization.
	if err := ensureRoles.Execute(ctx, org.ID); err != nil {
		return fmt.Errorf("ensure org roles: %w", err)
	}

	// 4. Super-admin user (created once).
	user, err := userRepo.GetByEmail(ctx, cfg.SeedAdmin.Email)
	if err != nil || user == nil {
		hash, herr := auth.HashPassword(cfg.SeedAdmin.Password)
		if herr != nil {
			return fmt.Errorf("hash seed password: %w", herr)
		}
		now := time.Now().UTC()
		user = &iamModel.User{
			Email:              cfg.SeedAdmin.Email,
			PasswordHash:       hash,
			FirstName:          "Super",
			LastName:           "Admin",
			Status:             iamModel.UserStatusActive,
			UserType:           iamModel.UserTypeSystemAdmin,
			KVKKConsentAt:      &now,
			KVKKConsentVersion: "1.0",
		}
		if err := userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("create seed user: %w", err)
		}
		log.Info("seed: super-admin user created", "email", user.Email, "id", user.ID)
	}

	// 5. Organization membership.
	if _, err := orgUserRepo.GetByOrgAndUser(ctx, org.ID, user.ID); err != nil {
		if err := orgUserRepo.Add(ctx, &iamModel.OrganizationUser{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Status:         iamModel.OrgUserStatusActive,
		}); err != nil {
			return fmt.Errorf("add seed org membership: %w", err)
		}
		log.Info("seed: org membership granted", "org", org.Slug, "user", user.Email)
	}

	// 6. Assign every organization role (idempotent via ON CONFLICT DO NOTHING).
	roles, err := roleRepo.ListByScope(ctx, iamModel.ScopeTypeOrganization, org.ID)
	if err != nil {
		return fmt.Errorf("list org roles: %w", err)
	}
	appID := app.ID
	for _, role := range roles {
		if err := roleRepo.AssignRoleToUser(ctx, &iamModel.UserRole{
			UserID:         user.ID,
			RoleID:         role.ID,
			OrganizationID: org.ID,
			AppID:          &appID,
		}); err != nil {
			log.Warn("seed: assign role failed", "role", role.Name, "error", err)
		}
	}

	log.Info("seed: super-admin ready",
		"email", cfg.SeedAdmin.Email,
		"org", org.Slug,
		"app", app.Slug,
		"roles_assigned", len(roles),
	)
	return nil
}
