package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateRoleRequest is the input for creating an organization-scoped role.
type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required,min=2"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// UpdateRoleRequest is the input for updating a role's metadata.
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SetPermissionsRequest replaces the full permission set of a role.
type SetPermissionsRequest struct {
	Permissions []string `json:"permissions" validate:"required"`
}

// RevokeRoleRequest is the input for removing a role from a user.
type RevokeRoleRequest struct {
	UserID         uuid.UUID `json:"user_id" validate:"required"`
	RoleID         uuid.UUID `json:"role_id" validate:"required"`
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
}

// RoleResponse is the public representation of a role with its permissions.
type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	ScopeType   string    `json:"scope_type"`
	ScopeID     uuid.UUID `json:"scope_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
