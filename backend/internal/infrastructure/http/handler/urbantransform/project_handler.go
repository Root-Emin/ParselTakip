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

// ProjectHandler exposes CQRS-based HTTP endpoints for urban transformation projects.
// Write operations delegate to command handlers; reads delegate to query handlers.
type ProjectHandler struct {
	createCmd *command.CreateProjectHandler
	updateCmd *command.UpdateProjectHandler
	deleteCmd *command.DeleteProjectHandler
	getQuery  *query.GetProjectHandler
	listQuery *query.ListProjectsHandler
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(
	createCmd *command.CreateProjectHandler,
	updateCmd *command.UpdateProjectHandler,
	deleteCmd *command.DeleteProjectHandler,
	getQuery *query.GetProjectHandler,
	listQuery *query.ListProjectsHandler,
) *ProjectHandler {
	return &ProjectHandler{
		createCmd: createCmd,
		updateCmd: updateCmd,
		deleteCmd: deleteCmd,
		getQuery:  getQuery,
		listQuery: listQuery,
	}
}

// resolveTenant extracts organization and app scoping from the request.
// Organization comes from JWT/tenant context; app from the X-App-ID header.
func resolveTenant(r *http.Request) (orgID, appID uuid.UUID, ok bool) {
	orgID, found := middleware.OrgIDFromContext(r.Context())
	if !found || orgID == uuid.Nil {
		if tid, tok := middleware.TenantIDFromContext(r.Context()); tok {
			orgID = tid
		}
	}
	appID, _ = uuid.Parse(r.Header.Get("X-App-ID"))
	if orgID == uuid.Nil || appID == uuid.Nil {
		return uuid.Nil, uuid.Nil, false
	}
	return orgID, appID, true
}

// Create handles POST /projects (command: add).
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}

	var req dto.CreateProjectRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())

	result, err := h.createCmd.Execute(r.Context(), orgID, appID, userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgProjectCreated, result)
}

// Update handles PATCH /projects/{projectId} (command: update).
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamProjectID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidProjectID})
		return
	}

	var req dto.UpdateProjectRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.updateCmd.Execute(r.Context(), orgID, appID, id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgProjectUpdated, result)
}

// Delete handles DELETE /projects/{projectId} (command: delete).
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamProjectID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidProjectID})
		return
	}

	if err := h.deleteCmd.Execute(r.Context(), orgID, appID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgProjectDeleted, nil)
}

// Get handles GET /projects/{projectId} (query: single).
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamProjectID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidProjectID})
		return
	}

	result, err := h.getQuery.Execute(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgProjectFetched, result)
}

// List handles GET /projects (query: list + filter + search + sort + paginate).
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}

	q := dto.ListProjectsQuery{
		Search:    r.URL.Query().Get(constants.QueryKeySearch),
		SortBy:    r.URL.Query().Get(constants.QueryKeySortBy),
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if v := r.URL.Query().Get(constants.QueryKeyStatus); v != "" {
		q.Status = &v
	}
	if id := parseUUIDQuery(r, constants.QueryKeyContractor); id != nil {
		q.AssignedContractorID = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyInitiator); id != nil {
		q.InitiatedBy = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyStep); id != nil {
		q.CurrentWorkflowStepID = id
	}

	result, err := h.listQuery.Execute(r.Context(), orgID, appID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgProjectListed, result)
}
