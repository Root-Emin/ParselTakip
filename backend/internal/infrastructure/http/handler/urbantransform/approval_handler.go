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

// ApprovalHandler exposes CQRS-based HTTP endpoints for approvals.
type ApprovalHandler struct {
	cmd *command.ApprovalCommandHandler
	qry *query.ApprovalQueryHandler
}

// NewApprovalHandler creates a new ApprovalHandler.
func NewApprovalHandler(cmd *command.ApprovalCommandHandler, qry *query.ApprovalQueryHandler) *ApprovalHandler {
	return &ApprovalHandler{cmd: cmd, qry: qry}
}

// Create handles POST /approvals.
func (h *ApprovalHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreateApprovalRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Create(r.Context(), orgID, appID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgApprovalCreated, result)
}

// Decide handles PATCH /approvals/{approvalId}/decision (approve/reject).
func (h *ApprovalHandler) Decide(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamApprovalID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidApprovalID})
		return
	}
	var req dto.DecideApprovalRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	deciderID, _ := middleware.UserIDFromContext(r.Context())
	result, err := h.cmd.Decide(r.Context(), orgID, appID, id, deciderID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgApprovalDecided, result)
}

// Delete handles DELETE /approvals/{approvalId}.
func (h *ApprovalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamApprovalID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidApprovalID})
		return
	}
	if err := h.cmd.Delete(r.Context(), orgID, appID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgApprovalDeleted, nil)
}

// Get handles GET /approvals/{approvalId}.
func (h *ApprovalHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamApprovalID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidApprovalID})
		return
	}
	result, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgApprovalFetched, result)
}

// List handles GET /approvals (list + filter). Use ?status=pending for the pending-approvals view.
func (h *ApprovalHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	q := dto.ListApprovalsQuery{
		SortBy:    r.URL.Query().Get(constants.QueryKeySortBy),
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if id := parseUUIDQuery(r, constants.QueryKeyApprovalProject); id != nil {
		q.ProjectID = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyApprovalApprover); id != nil {
		q.ApproverID = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyApprovalOwner); id != nil {
		q.OwnerID = id
	}
	if v := r.URL.Query().Get(constants.QueryKeyApprovalType); v != "" {
		q.ApprovalType = &v
	}
	if v := r.URL.Query().Get(constants.QueryKeyApprovalStatus); v != "" {
		q.Status = &v
	}
	result, err := h.qry.List(r.Context(), orgID, appID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgApprovalListed, result)
}
