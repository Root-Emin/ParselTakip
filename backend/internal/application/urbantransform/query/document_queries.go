package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/pagination"
)

// DocumentQueryHandler handles read operations for documents (CQRS read side).
type DocumentQueryHandler struct {
	repo     repository.DocumentRepository
	typeRepo repository.DocumentTypeRepository
}

// NewDocumentQueryHandler creates a new DocumentQueryHandler.
func NewDocumentQueryHandler(repo repository.DocumentRepository, typeRepo repository.DocumentTypeRepository) *DocumentQueryHandler {
	return &DocumentQueryHandler{repo: repo, typeRepo: typeRepo}
}

// Get returns a single document by ID within the tenant scope.
func (h *DocumentQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.DocumentResponse, error) {
	doc, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToDocumentResponse(doc)
	return &resp, nil
}

// List returns a paginated, filtered and searchable list of documents.
func (h *DocumentQueryHandler) List(ctx context.Context, orgID, appID uuid.UUID, q dto.ListDocumentsQuery) (pagination.Result[dto.DocumentResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	sortBy := q.SortBy
	if !constants.IsAllowedDocSortKey(sortBy) {
		sortBy = constants.DocSortKeyCreatedAt
	}

	filter := repository.DocumentFilter{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      q.ProjectID,
		BuildingID:     q.BuildingID,
		OwnerID:        q.OwnerID,
		DocumentTypeID: q.DocumentTypeID,
		Search:         q.Search,
		SortBy:         sortBy,
		SortOrder:      constants.NormalizeSortOrder(q.SortOrder),
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}
	if q.Status != nil {
		st := model.DocumentStatus(*q.Status)
		filter.Status = &st
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.DocumentResponse]{}, err
	}
	return pagination.NewResult(dto.ToDocumentResponseList(items), params, total), nil
}

// ListReviews returns the review history for a document.
func (h *DocumentQueryHandler) ListReviews(ctx context.Context, orgID, appID, documentID uuid.UUID) ([]dto.DocumentReviewResponse, error) {
	if _, err := h.repo.GetByID(ctx, orgID, appID, documentID); err != nil {
		return nil, err
	}
	reviews, err := h.repo.ListReviews(ctx, orgID, appID, documentID)
	if err != nil {
		return nil, err
	}
	return dto.ToDocumentReviewResponseList(reviews), nil
}

// ListTypes returns document type master data, optionally filtered by category.
func (h *DocumentQueryHandler) ListTypes(ctx context.Context, category string) ([]dto.DocumentTypeResponse, error) {
	types, err := h.typeRepo.ListAll(ctx, category)
	if err != nil {
		return nil, err
	}
	return dto.ToDocumentTypeResponseList(types), nil
}
