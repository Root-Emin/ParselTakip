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

// ContractorHandler exposes CQRS-based HTTP endpoints for contractors.
type ContractorHandler struct {
	cmd *command.ContractorCommandHandler
	qry *query.ContractorQueryHandler
}

// NewContractorHandler creates a new ContractorHandler.
func NewContractorHandler(cmd *command.ContractorCommandHandler, qry *query.ContractorQueryHandler) *ContractorHandler {
	return &ContractorHandler{cmd: cmd, qry: qry}
}

// Create handles POST /contractors.
func (h *ContractorHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreateContractorRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Create(r.Context(), orgID, appID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgContractorCreated, result)
}

// Update handles PATCH /contractors/{contractorId}.
func (h *ContractorHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamContractorID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidContractorID})
		return
	}
	var req dto.UpdateContractorRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Update(r.Context(), orgID, appID, id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgContractorUpdated, result)
}

// Delete handles DELETE /contractors/{contractorId}.
func (h *ContractorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamContractorID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidContractorID})
		return
	}
	if err := h.cmd.Delete(r.Context(), orgID, appID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgContractorDeleted, nil)
}

// Get handles GET /contractors/{contractorId}.
func (h *ContractorHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamContractorID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidContractorID})
		return
	}
	result, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgContractorFetched, result)
}

// List handles GET /contractors (list + filter + search).
func (h *ContractorHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	q := dto.ListContractorsQuery{
		Search:    r.URL.Query().Get(constants.QueryKeyContractorSearch),
		SortBy:    r.URL.Query().Get(constants.QueryKeySortBy),
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if v := r.URL.Query().Get(constants.QueryKeyContractorStatus); v != "" {
		q.Status = &v
	}
	result, err := h.qry.List(r.Context(), orgID, appID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgContractorListed, result)
}
