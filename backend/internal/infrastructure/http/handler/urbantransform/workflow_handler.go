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

// WorkflowHandler exposes endpoints for workflow steps, project states and history.
type WorkflowHandler struct {
	cmd *command.WorkflowCommandHandler
	qry *query.WorkflowQueryHandler
}

// NewWorkflowHandler creates a new WorkflowHandler.
func NewWorkflowHandler(cmd *command.WorkflowCommandHandler, qry *query.WorkflowQueryHandler) *WorkflowHandler {
	return &WorkflowHandler{cmd: cmd, qry: qry}
}

// ListSteps handles GET /workflow-steps (global step definitions).
func (h *WorkflowHandler) ListSteps(w http.ResponseWriter, r *http.Request) {
	result, err := h.qry.ListSteps(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgWorkflowStepsListed, result)
}

// ListProjectWorkflow handles GET /projects/{projectId}/workflow.
func (h *WorkflowHandler) ListProjectWorkflow(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid project id"})
		return
	}
	result, err := h.qry.ListProjectStates(r.Context(), orgID, appID, projectID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgWorkflowStatesListed, result)
}

// ListProjectHistory handles GET /projects/{projectId}/workflow/history.
func (h *WorkflowHandler) ListProjectHistory(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid project id"})
		return
	}
	result, err := h.qry.ListProjectHistory(r.Context(), orgID, appID, projectID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgWorkflowHistoryListed, result)
}

// Advance handles POST /projects/{projectId}/workflow/advance.
func (h *WorkflowHandler) Advance(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid project id"})
		return
	}
	var req dto.AdvanceWorkflowRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	changedBy, _ := middleware.UserIDFromContext(r.Context())
	result, err := h.cmd.Advance(r.Context(), orgID, appID, projectID, changedBy, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgWorkflowAdvanced, result)
}
