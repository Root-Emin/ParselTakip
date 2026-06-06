package iam

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/iam/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"
	"github.com/masterfabric-go/masterfabric/internal/shared/response"
	"github.com/masterfabric-go/masterfabric/internal/shared/validator"
)

// --- User administration ---

// UpdateUser handles PATCH /users/{id}.
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	var req dto.UpdateUserRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	user, err := h.manageUsersUC.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// DeactivateUser handles POST /users/{id}/deactivate.
func (h *Handler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	user, err := h.manageUsersUC.SetStatus(r.Context(), id, model.UserStatusInactive)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// ActivateUser handles POST /users/{id}/activate.
func (h *Handler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	user, err := h.manageUsersUC.SetStatus(r.Context(), id, model.UserStatusActive)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// --- Role administration ---

// CreateRole handles POST /roles.
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization context required"})
		return
	}
	var req dto.CreateRoleRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	role, err := h.manageRolesUC.Create(r.Context(), orgID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, role)
}

// ListRoles handles GET /roles (organization-scoped).
func (h *Handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	orgID, ok := middleware.OrgIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization context required"})
		return
	}
	roles, err := h.manageRolesUC.List(r.Context(), orgID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, roles)
}

// GetRole handles GET /roles/{roleId}.
func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roleId"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role id"})
		return
	}
	role, err := h.manageRolesUC.Get(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, role)
}

// UpdateRole handles PATCH /roles/{roleId}.
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roleId"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role id"})
		return
	}
	var req dto.UpdateRoleRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	role, err := h.manageRolesUC.Update(r.Context(), id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, role)
}

// DeleteRole handles DELETE /roles/{roleId}.
func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roleId"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role id"})
		return
	}
	if err := h.manageRolesUC.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

// SetRolePermissions handles PUT /roles/{roleId}/permissions.
func (h *Handler) SetRolePermissions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roleId"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role id"})
		return
	}
	var req dto.SetPermissionsRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	role, err := h.manageRolesUC.SetPermissions(r.Context(), id, req.Permissions)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, role)
}

// RevokeRole handles POST /roles/revoke.
func (h *Handler) RevokeRole(w http.ResponseWriter, r *http.Request) {
	var req dto.RevokeRoleRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.manageRolesUC.RevokeRole(r.Context(), req); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
