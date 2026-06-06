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

// DocumentHandler exposes CQRS-based HTTP endpoints for documents, reviews and types.
type DocumentHandler struct {
	cmd *command.DocumentCommandHandler
	qry *query.DocumentQueryHandler
}

// NewDocumentHandler creates a new DocumentHandler.
func NewDocumentHandler(cmd *command.DocumentCommandHandler, qry *query.DocumentQueryHandler) *DocumentHandler {
	return &DocumentHandler{cmd: cmd, qry: qry}
}

// Create handles POST /documents (upload/register).
func (h *DocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreateDocumentRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	roles, _ := middleware.RolesFromContext(r.Context())
	uploadedByRole := ""
	if len(roles) > 0 {
		uploadedByRole = roles[0]
	}
	result, err := h.cmd.Create(r.Context(), orgID, appID, userID, uploadedByRole, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgDocumentCreated, result)
}

// Update handles PATCH /documents/{documentId}.
func (h *DocumentHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamDocumentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidDocumentID})
		return
	}
	var req dto.UpdateDocumentRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Update(r.Context(), orgID, appID, id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgDocumentUpdated, result)
}

// Delete handles DELETE /documents/{documentId}.
func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamDocumentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidDocumentID})
		return
	}
	if err := h.cmd.Delete(r.Context(), orgID, appID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgDocumentDeleted, nil)
}

// Get handles GET /documents/{documentId}.
func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamDocumentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidDocumentID})
		return
	}
	result, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgDocumentFetched, result)
}

// List handles GET /documents (list + filter + search).
func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	q := dto.ListDocumentsQuery{
		Search:    r.URL.Query().Get(constants.QueryKeyDocSearch),
		SortBy:    r.URL.Query().Get(constants.QueryKeySortBy),
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if id := parseUUIDQuery(r, constants.QueryKeyDocProject); id != nil {
		q.ProjectID = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyDocBuilding); id != nil {
		q.BuildingID = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyDocOwner); id != nil {
		q.OwnerID = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyDocType); id != nil {
		q.DocumentTypeID = id
	}
	if v := r.URL.Query().Get(constants.QueryKeyDocStatus); v != "" {
		q.Status = &v
	}
	result, err := h.qry.List(r.Context(), orgID, appID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgDocumentListed, result)
}

// Review handles POST /documents/{documentId}/reviews (approve/reject/mark-missing).
func (h *DocumentHandler) Review(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamDocumentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidDocumentID})
		return
	}
	var req dto.ReviewDocumentRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	reviewerID, _ := middleware.UserIDFromContext(r.Context())
	result, err := h.cmd.Review(r.Context(), orgID, appID, id, reviewerID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgDocumentReviewed, result)
}

// ListReviews handles GET /documents/{documentId}/reviews.
func (h *DocumentHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamDocumentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidDocumentID})
		return
	}
	result, err := h.qry.ListReviews(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgReviewsListed, result)
}

// ListTypes handles GET /document-types (master data, optional ?category=).
func (h *DocumentHandler) ListTypes(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get(constants.QueryKeyDocCategory)
	result, err := h.qry.ListTypes(r.Context(), category)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgDocumentTypesListed, result)
}
