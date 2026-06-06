package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// DocumentFilter holds optional filter criteria for querying documents.
type DocumentFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	ProjectID      *uuid.UUID
	BuildingID     *uuid.UUID
	OwnerID        *uuid.UUID
	DocumentTypeID *uuid.UUID
	Status         *model.DocumentStatus
	Search         string // free-text on file_name
	SortBy         string
	SortOrder      string
	Offset         int
	Limit          int
}

// DocumentRepository defines persistence operations for documents and reviews.
type DocumentRepository interface {
	Create(ctx context.Context, doc *model.Document) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Document, error)
	Update(ctx context.Context, doc *model.Document) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter DocumentFilter) ([]*model.Document, int, error)

	// Reviews
	CreateReview(ctx context.Context, review *model.DocumentReview) error
	ListReviews(ctx context.Context, orgID, appID, documentID uuid.UUID) ([]*model.DocumentReview, error)
}

// DocumentTypeRepository defines read access to document type master data.
type DocumentTypeRepository interface {
	ListAll(ctx context.Context, category string) ([]*model.DocumentType, error)
}
