package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

// EnsureOrgRolesUseCase syncs urban transformation system roles into organization-scoped RBAC roles.
type EnsureOrgRolesUseCase struct {
	systemRoleRepo repository.SystemRoleRepository
	roleRepo       repository.RoleRepository
}

// NewEnsureOrgRolesUseCase creates a new EnsureOrgRolesUseCase.
func NewEnsureOrgRolesUseCase(systemRoleRepo repository.SystemRoleRepository, roleRepo repository.RoleRepository) *EnsureOrgRolesUseCase {
	return &EnsureOrgRolesUseCase{systemRoleRepo: systemRoleRepo, roleRepo: roleRepo}
}

// Execute ensures all system role templates exist for the given organization.
func (uc *EnsureOrgRolesUseCase) Execute(ctx context.Context, organizationID uuid.UUID) error {
	definitions, err := uc.systemRoleRepo.ListAll(ctx)
	if err != nil {
		return err
	}

	for _, def := range definitions {
		_, err := uc.roleRepo.GetByNameAndScope(ctx, model.ScopeTypeOrganization, organizationID, def.Code)
		if err == nil {
			continue
		}
		if !errors.Is(err, domainErr.ErrNotFound) {
			return err
		}

		role := &model.Role{
			ScopeType:   model.ScopeTypeOrganization,
			ScopeID:     organizationID,
			Name:        def.Code,
			Description: def.Name,
		}
		if err := uc.roleRepo.Create(ctx, role); err != nil {
			return err
		}

		for _, permission := range def.Permissions {
			if err := uc.roleRepo.AddPermission(ctx, role.ID, permission); err != nil {
				return err
			}
		}
	}

	return nil
}

// ErrOrganizationAccessDenied indicates the user is not a member of the organization.
var ErrOrganizationAccessDenied = domainErr.New(domainErr.ErrForbidden, "user is not a member of this organization", nil)
