package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/iam/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/repository"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/service"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

// LoginUseCase handles user authentication.
type LoginUseCase struct {
	userRepo       repository.UserRepository
	orgUserRepo    repository.OrgUserRepository
	roleRepo       repository.RoleRepository
	rbac           service.RBACService
	auth           service.AuthService
	ensureOrgRoles *EnsureOrgRolesUseCase
	tokenTTLHours  int
}

// NewLoginUseCase creates a new LoginUseCase.
func NewLoginUseCase(
	userRepo repository.UserRepository,
	orgUserRepo repository.OrgUserRepository,
	roleRepo repository.RoleRepository,
	rbac service.RBACService,
	auth service.AuthService,
	ensureOrgRoles *EnsureOrgRolesUseCase,
	tokenTTLHours int,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		orgUserRepo:    orgUserRepo,
		roleRepo:       roleRepo,
		rbac:           rbac,
		auth:           auth,
		ensureOrgRoles: ensureOrgRoles,
		tokenTTLHours:  tokenTTLHours,
	}
}

// Execute authenticates a user and returns a JWT token with optional org-scoped RBAC claims.
func (uc *LoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrUnauthorized, "invalid credentials", nil)
	}

	if !user.IsActive() {
		return nil, domainErr.New(domainErr.ErrForbidden, "account is not active", nil)
	}

	if err := uc.auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, err
	}

	claims := service.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
	}

	var organizationID uuid.UUID
	if req.OrganizationID != nil && *req.OrganizationID != uuid.Nil {
		organizationID = *req.OrganizationID

		orgUser, err := uc.orgUserRepo.GetByOrgAndUser(ctx, organizationID, user.ID)
		if err != nil {
			return nil, ErrOrganizationAccessDenied
		}
		if orgUser.Status != model.OrgUserStatusActive {
			return nil, domainErr.New(domainErr.ErrForbidden, "organization membership is not active", nil)
		}

		if uc.ensureOrgRoles != nil {
			if err := uc.ensureOrgRoles.Execute(ctx, organizationID); err != nil {
				return nil, err
			}
		}

		roles, err := uc.roleRepo.GetUserRoleNames(ctx, user.ID, organizationID)
		if err != nil {
			return nil, err
		}

		permissions, err := uc.rbac.GetUserPermissions(ctx, user.ID, organizationID)
		if err != nil {
			return nil, err
		}

		claims.OrganizationID = organizationID
		claims.Roles = roles
		claims.Permissions = permissions
	}

	token, err := uc.auth.GenerateToken(ctx, claims)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to generate token", err)
	}

	return &dto.LoginResponse{
		Token:          token,
		TokenType:      "Bearer",
		ExpiresInHours: uc.tokenTTLHours,
		User: dto.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Status:    string(user.Status),
			CreatedAt: user.CreatedAt,
		},
		OrganizationID: organizationID,
		Roles:          claims.Roles,
		Permissions:    claims.Permissions,
	}, nil
}
