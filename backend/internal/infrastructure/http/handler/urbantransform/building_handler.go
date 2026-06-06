package urbantransform

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/command"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/query"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"
	"github.com/masterfabric-go/masterfabric/internal/shared/response"
	"github.com/masterfabric-go/masterfabric/internal/shared/validator"
)

// BuildingHandler exposes CQRS-based HTTP endpoints for buildings.
type BuildingHandler struct {
	cmd *command.BuildingCommandHandler
	qry *query.BuildingQueryHandler
}

// NewBuildingHandler creates a new BuildingHandler.
func NewBuildingHandler(cmd *command.BuildingCommandHandler, qry *query.BuildingQueryHandler) *BuildingHandler {
	return &BuildingHandler{cmd: cmd, qry: qry}
}

// Create handles POST /buildings.
func (h *BuildingHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreateBuildingRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	result, err := h.cmd.Create(r.Context(), orgID, appID, userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgBuildingCreated, result)
}

// Update handles PATCH /buildings/{buildingId}.
func (h *BuildingHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamBuildingID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidBuildingID})
		return
	}
	var req dto.UpdateBuildingRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Update(r.Context(), orgID, appID, id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgBuildingUpdated, result)
}

// Delete handles DELETE /buildings/{buildingId}.
func (h *BuildingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamBuildingID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidBuildingID})
		return
	}
	if err := h.cmd.Delete(r.Context(), orgID, appID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgBuildingDeleted, nil)
}

// Get handles GET /buildings/{buildingId}.
func (h *BuildingHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamBuildingID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidBuildingID})
		return
	}
	result, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgBuildingFetched, result)
}

// List handles GET /buildings (list + filter + search).
func (h *BuildingHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	q := dto.ListBuildingsQuery{
		City:      r.URL.Query().Get(constants.QueryKeyBuildingCity),
		District:  r.URL.Query().Get(constants.QueryKeyBuildingDistrict),
		Search:    r.URL.Query().Get(constants.QueryKeyBuildingSearch),
		SortBy:    r.URL.Query().Get(constants.QueryKeySortBy),
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if v := r.URL.Query().Get(constants.QueryKeyBuildingStatus); v != "" {
		q.Status = &v
	}
	if v := r.URL.Query().Get(constants.QueryKeyBuildingRiskStatus); v != "" {
		q.RiskStatus = &v
	}
	if v := r.URL.Query().Get(constants.QueryKeyBuildingType); v != "" {
		q.BuildingType = &v
	}
	result, err := h.qry.List(r.Context(), orgID, appID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgBuildingListed, result)
}
