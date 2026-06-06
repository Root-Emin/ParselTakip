package urbantransform

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/command"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/query"
	"github.com/masterfabric-go/masterfabric/internal/shared/response"
	"github.com/masterfabric-go/masterfabric/internal/shared/validator"
)

// PropertyOwnerHandler exposes CQRS-based HTTP endpoints for property owners.
type PropertyOwnerHandler struct {
	cmd *command.PropertyOwnerCommandHandler
	qry *query.PropertyOwnerQueryHandler
}

// NewPropertyOwnerHandler creates a new PropertyOwnerHandler.
func NewPropertyOwnerHandler(cmd *command.PropertyOwnerCommandHandler, qry *query.PropertyOwnerQueryHandler) *PropertyOwnerHandler {
	return &PropertyOwnerHandler{cmd: cmd, qry: qry}
}

// Create handles POST /property-owners.
func (h *PropertyOwnerHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreatePropertyOwnerRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Create(r.Context(), orgID, appID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgOwnerCreated, result)
}

// Update handles PATCH /property-owners/{ownerId}.
func (h *PropertyOwnerHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamOwnerID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidOwnerID})
		return
	}
	var req dto.UpdatePropertyOwnerRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Update(r.Context(), orgID, appID, id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgOwnerUpdated, result)
}

// Delete handles DELETE /property-owners/{ownerId}.
func (h *PropertyOwnerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamOwnerID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidOwnerID})
		return
	}
	if err := h.cmd.Delete(r.Context(), orgID, appID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgOwnerDeleted, nil)
}

// Get handles GET /property-owners/{ownerId}.
func (h *PropertyOwnerHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamOwnerID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidOwnerID})
		return
	}
	result, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgOwnerFetched, result)
}

// List handles GET /property-owners (list + filter + search).
func (h *PropertyOwnerHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	q := dto.ListPropertyOwnersQuery{
		Search:    r.URL.Query().Get(constants.QueryKeyOwnerSearch),
		SortBy:    r.URL.Query().Get(constants.QueryKeySortBy),
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if id := parseUUIDQuery(r, constants.QueryKeyOwnerUnit); id != nil {
		q.UnitID = id
	}
	if v := r.URL.Query().Get(constants.QueryKeyOwnerStatus); v != "" {
		q.Status = &v
	}
	result, err := h.qry.List(r.Context(), orgID, appID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgOwnerListed, result)
}
