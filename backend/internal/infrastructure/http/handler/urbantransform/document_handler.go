package urbantransform

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/command"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/query"
	domainStorage "github.com/masterfabric-go/masterfabric/internal/domain/storage"
	"github.com/masterfabric-go/masterfabric/internal/shared/jobs"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"
	"github.com/masterfabric-go/masterfabric/internal/shared/response"
	"github.com/masterfabric-go/masterfabric/internal/shared/validator"
)

// maxUploadBytes caps in-memory multipart parsing for document uploads (32 MiB).
const maxUploadBytes = 32 << 20

// DocumentHandler exposes CQRS-based HTTP endpoints for documents, reviews and types.
type DocumentHandler struct {
	cmd           *command.DocumentCommandHandler
	qry           *query.DocumentQueryHandler
	storage       domainStorage.ObjectStorage
	jobs          jobs.Enqueuer
	presignExpiry time.Duration
}

// NewDocumentHandler creates a new DocumentHandler. storage and enqueuer may be
// nil when object storage / background jobs are disabled.
func NewDocumentHandler(
	cmd *command.DocumentCommandHandler,
	qry *query.DocumentQueryHandler,
	storage domainStorage.ObjectStorage,
	enqueuer jobs.Enqueuer,
	presignExpiry time.Duration,
) *DocumentHandler {
	if presignExpiry <= 0 {
		presignExpiry = 10 * time.Minute
	}
	return &DocumentHandler{cmd: cmd, qry: qry, storage: storage, jobs: enqueuer, presignExpiry: presignExpiry}
}

// sanitizeFilename strips path separators to keep object keys safe.
func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	if name == "" || name == "." || name == "/" {
		return "file"
	}
	return name
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

// Upload handles POST /documents/upload (multipart/form-data). The server streams
// the file into object storage (MinIO), computes a sha256 checksum, persists the
// document metadata and enqueues post-processing.
func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	if h.storage == nil {
		response.JSON(w, http.StatusServiceUnavailable, map[string]string{"error": "object storage is not configured"})
		return
	}
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "file is required"})
		return
	}
	defer file.Close()

	projectID, err := uuid.Parse(r.FormValue("project_id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "valid project_id is required"})
		return
	}
	docTypeID, err := uuid.Parse(r.FormValue("document_type_id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "valid document_type_id is required"})
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, maxUploadBytes))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read file"})
		return
	}
	sum := sha256.Sum256(data)
	checksum := hex.EncodeToString(sum[:])

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	key := fmt.Sprintf("org/%s/project/%s/%s-%s", orgID, projectID, uuid.NewString(), sanitizeFilename(header.Filename))

	if _, err := h.storage.Upload(r.Context(), domainStorage.UploadInput{
		Key:         key,
		Reader:      bytes.NewReader(data),
		Size:        int64(len(data)),
		ContentType: contentType,
	}); err != nil {
		response.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to store file"})
		return
	}

	size := int64(len(data))
	req := dto.CreateDocumentRequest{
		ProjectID:      projectID,
		DocumentTypeID: docTypeID,
		BuildingID:     parseUUIDForm(r, "building_id"),
		UnitID:         parseUUIDForm(r, "unit_id"),
		OwnerID:        parseUUIDForm(r, "owner_id"),
		FileName:       header.Filename,
		FilePath:       key,
		FileSize:       &size,
		MimeType:       contentType,
		StorageBucket:  h.storage.Bucket(),
		StorageKey:     key,
		Checksum:       checksum,
		IsNotarized:    r.FormValue("is_notarized") == "true",
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

	if h.jobs != nil {
		_ = h.jobs.Enqueue(r.Context(), jobs.TopicDocumentProcess, jobs.DocumentProcessJob{
			DocumentID:     result.ID.String(),
			OrganizationID: orgID.String(),
			AppID:          appID.String(),
			StorageBucket:  h.storage.Bucket(),
			StorageKey:     key,
			Checksum:       checksum,
		})
	}

	response.CreatedEnvelope(w, constants.MsgDocumentUploaded, result)
}

// PresignUpload handles POST /documents/presign-upload. It returns a short-lived
// URL the client uses to PUT the file directly to object storage, then registers
// metadata via POST /documents with the returned storage_key.
func (h *DocumentHandler) PresignUpload(w http.ResponseWriter, r *http.Request) {
	orgID, _, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	if h.storage == nil {
		response.JSON(w, http.StatusServiceUnavailable, map[string]string{"error": "object storage is not configured"})
		return
	}
	var req dto.PresignUploadRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	key := fmt.Sprintf("org/%s/project/%s/%s-%s", orgID, req.ProjectID, uuid.NewString(), sanitizeFilename(req.FileName))
	url, err := h.storage.PresignPut(r.Context(), key, h.presignExpiry)
	if err != nil {
		response.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to presign upload"})
		return
	}
	response.Success(w, constants.MsgDocumentPresigned, dto.PresignUploadResponse{
		UploadURL:     url,
		StorageBucket: h.storage.Bucket(),
		StorageKey:    key,
		ExpiresIn:     int(h.presignExpiry.Seconds()),
	})
}

// Download handles GET /documents/{documentId}/download and returns a short-lived
// presigned URL to fetch the stored file.
func (h *DocumentHandler) Download(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	if h.storage == nil {
		response.JSON(w, http.StatusServiceUnavailable, map[string]string{"error": "object storage is not configured"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamDocumentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidDocumentID})
		return
	}
	doc, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	key := doc.StorageKey
	if key == "" {
		key = doc.FilePath
	}
	if key == "" {
		response.JSON(w, http.StatusNotFound, map[string]string{"error": "no stored file for this document"})
		return
	}
	url, err := h.storage.PresignGet(r.Context(), key, h.presignExpiry, doc.FileName)
	if err != nil {
		response.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to presign download"})
		return
	}
	response.Success(w, constants.MsgDocumentDownloadURL, dto.DownloadResponse{
		URL:       url,
		ExpiresIn: int(h.presignExpiry.Seconds()),
	})
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
