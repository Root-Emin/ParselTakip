package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/iam/dto"
	iamEvent "github.com/masterfabric-go/masterfabric/internal/domain/iam/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/repository"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/service"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// ManageRolesUseCase handles role CRUD and permission editing (system admin).
type ManageRolesUseCase struct {
	roleRepo repository.RoleRepository
	rbac     service.RBACService
	eventBus events.EventBus
}

// NewManageRolesUseCase creates a new ManageRolesUseCase.
func NewManageRolesUseCase(roleRepo repository.RoleRepository, rbac service.RBACService, eventBus events.EventBus) *ManageRolesUseCase {
	return &ManageRolesUseCase{roleRepo: roleRepo, rbac: rbac, eventBus: eventBus}
}

func (uc *ManageRolesUseCase) toResponse(ctx context.Context, role *model.Role) (*dto.RoleResponse, error) {
	perms, err := uc.roleRepo.GetPermissions(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	if perms == nil {
		perms = []string{}
	}
	return &dto.RoleResponse{
		ID:          role.ID,
		ScopeType:   string(role.ScopeType),
		ScopeID:     role.ScopeID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: perms,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

// Create creates a new organization-scoped role with optional permissions.
func (uc *ManageRolesUseCase) Create(ctx context.Context, orgID uuid.UUID, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role := &model.Role{
		ScopeType:   model.ScopeTypeOrganization,
		ScopeID:     orgID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := uc.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}
	for _, p := range req.Permissions {
		if p == "" {
			continue
		}
		if err := uc.roleRepo.AddPermission(ctx, role.ID, p); err != nil {
			return nil, err
		}
	}
	return uc.toResponse(ctx, role)
}

// List returns all organization-scoped roles for the given org.
func (uc *ManageRolesUseCase) List(ctx context.Context, orgID uuid.UUID) ([]dto.RoleResponse, error) {
	roles, err := uc.roleRepo.ListByScope(ctx, model.ScopeTypeOrganization, orgID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		resp, err := uc.toResponse(ctx, role)
		if err != nil {
			return nil, err
		}
		out = append(out, *resp)
	}
	return out, nil
}

// Get returns a single role with its permissions.
func (uc *ManageRolesUseCase) Get(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, error) {
	role, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toResponse(ctx, role)
}

// Update updates a role's metadata.
func (uc *ManageRolesUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if err := uc.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}
	return uc.toResponse(ctx, role)
}

// Delete removes a role.
func (uc *ManageRolesUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.roleRepo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.roleRepo.Delete(ctx, id)
}

// SetPermissions replaces the full permission set of a role.
func (uc *ManageRolesUseCase) SetPermissions(ctx context.Context, id uuid.UUID, permissions []string) (*dto.RoleResponse, error) {
	role, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	current, err := uc.roleRepo.GetPermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	desired := make(map[string]struct{}, len(permissions))
	for _, p := range permissions {
		if p != "" {
			desired[p] = struct{}{}
		}
	}
	existing := make(map[string]struct{}, len(current))
	for _, p := range current {
		existing[p] = struct{}{}
	}

	for p := range desired {
		if _, ok := existing[p]; !ok {
			if err := uc.roleRepo.AddPermission(ctx, id, p); err != nil {
				return nil, err
			}
		}
	}
	for p := range existing {
		if _, ok := desired[p]; !ok {
			if err := uc.roleRepo.RemovePermission(ctx, id, p); err != nil {
				return nil, err
			}
		}
	}

	return uc.toResponse(ctx, role)
}

// RevokeRole removes a role from a user and invalidates their permission cache.
func (uc *ManageRolesUseCase) RevokeRole(ctx context.Context, req dto.RevokeRoleRequest) error {
	if err := uc.roleRepo.RemoveRoleFromUser(ctx, req.UserID, req.RoleID); err != nil {
		return err
	}
	_ = uc.rbac.InvalidateCache(ctx, req.UserID, req.OrganizationID)
	_ = uc.eventBus.Publish(ctx, events.TopicIAM, iamEvent.RoleRevoked{
		UserID:         req.UserID,
		RoleID:         req.RoleID,
		OrganizationID: req.OrganizationID,
		Timestamp:      time.Now().UTC(),
	})
	return nil
}
