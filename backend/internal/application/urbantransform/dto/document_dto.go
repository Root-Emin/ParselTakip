package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateDocumentRequest is the input for uploading/registering a document (command).
type CreateDocumentRequest struct {
	ProjectID      uuid.UUID       `json:"project_id" validate:"required"`
	DocumentTypeID uuid.UUID       `json:"document_type_id" validate:"required"`
	BuildingID     *uuid.UUID      `json:"building_id,omitempty"`
	UnitID         *uuid.UUID      `json:"unit_id,omitempty"`
	OwnerID        *uuid.UUID      `json:"owner_id,omitempty"`
	FileName       string          `json:"file_name" validate:"required"`
	FilePath       string          `json:"file_path" validate:"required"`
	FileSize       *int64          `json:"file_size,omitempty"`
	MimeType       string          `json:"mime_type"`
	IsNotarized    bool            `json:"is_notarized"`
	NotaryDate     *time.Time      `json:"notary_date,omitempty"`
	ExpiryDate     *time.Time      `json:"expiry_date,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
}

// UpdateDocumentRequest is the input for updating a document (partial command).
type UpdateDocumentRequest struct {
	FileName    *string         `json:"file_name,omitempty"`
	FilePath    *string         `json:"file_path,omitempty"`
	FileSize    *int64          `json:"file_size,omitempty"`
	MimeType    *string         `json:"mime_type,omitempty"`
	Status      *string         `json:"status,omitempty" validate:"omitempty,oneof=draft submitted under_review approved rejected missing expired invalid"`
	IsNotarized *bool           `json:"is_notarized,omitempty"`
	NotaryDate  *time.Time      `json:"notary_date,omitempty"`
	ExpiryDate  *time.Time      `json:"expiry_date,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
}

// ReviewDocumentRequest is the input for reviewing a document (command).
type ReviewDocumentRequest struct {
	Status       string          `json:"status" validate:"required,oneof=approved rejected missing_items needs_revision"`
	MissingItems json.RawMessage `json:"missing_items,omitempty"`
	ReviewNotes  string          `json:"review_notes"`
}

// ListDocumentsQuery is the input for listing/filtering/searching documents (query).
type ListDocumentsQuery struct {
	ProjectID      *uuid.UUID
	BuildingID     *uuid.UUID
	OwnerID        *uuid.UUID
	DocumentTypeID *uuid.UUID
	Status         *string
	Search         string
	SortBy         string
	SortOrder      string
	Page           int
	PerPage        int
}

// DocumentResponse is the public representation of a document (response model).
type DocumentResponse struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	AppID          uuid.UUID       `json:"app_id"`
	ProjectID      uuid.UUID       `json:"project_id"`
	DocumentTypeID uuid.UUID       `json:"document_type_id"`
	BuildingID     *uuid.UUID      `json:"building_id,omitempty"`
	UnitID         *uuid.UUID      `json:"unit_id,omitempty"`
	OwnerID        *uuid.UUID      `json:"owner_id,omitempty"`
	FileName       string          `json:"file_name"`
	FilePath       string          `json:"file_path"`
	FileSize       *int64          `json:"file_size,omitempty"`
	MimeType       string          `json:"mime_type"`
	Status         string          `json:"status"`
	IsNotarized    bool            `json:"is_notarized"`
	NotaryDate     *time.Time      `json:"notary_date,omitempty"`
	ExpiryDate     *time.Time      `json:"expiry_date,omitempty"`
	UploadedBy     *uuid.UUID      `json:"uploaded_by,omitempty"`
	UploadedByRole string          `json:"uploaded_by_role"`
	Version        int             `json:"version"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// DocumentReviewResponse is the public representation of a document review.
type DocumentReviewResponse struct {
	ID           uuid.UUID       `json:"id"`
	DocumentID   uuid.UUID       `json:"document_id"`
	ReviewerID   uuid.UUID       `json:"reviewer_id"`
	Status       string          `json:"status"`
	MissingItems json.RawMessage `json:"missing_items,omitempty"`
	ReviewNotes  string          `json:"review_notes"`
	ReviewedAt   time.Time       `json:"reviewed_at"`
}

// DocumentTypeResponse is the public representation of a document type.
type DocumentTypeResponse struct {
	ID                     uuid.UUID `json:"id"`
	Code                   string    `json:"code"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Category               string    `json:"category"`
	IsMandatory            bool      `json:"is_mandatory"`
	RequiresNotary         bool      `json:"requires_notary"`
	RequiresMunicipalStamp bool      `json:"requires_municipal_stamp"`
	IsValidWithoutNotary   bool      `json:"is_valid_without_notary"`
	InvalidReason          string    `json:"invalid_reason"`
}

// ToDocumentResponse maps a domain document to its response model.
func ToDocumentResponse(d *model.Document) DocumentResponse {
	return DocumentResponse{
		ID:             d.ID,
		OrganizationID: d.OrganizationID,
		AppID:          d.AppID,
		ProjectID:      d.ProjectID,
		DocumentTypeID: d.DocumentTypeID,
		BuildingID:     d.BuildingID,
		UnitID:         d.UnitID,
		OwnerID:        d.OwnerID,
		FileName:       d.FileName,
		FilePath:       d.FilePath,
		FileSize:       d.FileSize,
		MimeType:       d.MimeType,
		Status:         string(d.Status),
		IsNotarized:    d.IsNotarized,
		NotaryDate:     d.NotaryDate,
		ExpiryDate:     d.ExpiryDate,
		UploadedBy:     d.UploadedBy,
		UploadedByRole: d.UploadedByRole,
		Version:        d.Version,
		Metadata:       d.Metadata,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

// ToDocumentResponseList maps a slice of documents to response models.
func ToDocumentResponseList(items []*model.Document) []DocumentResponse {
	out := make([]DocumentResponse, 0, len(items))
	for _, d := range items {
		out = append(out, ToDocumentResponse(d))
	}
	return out
}

// ToDocumentReviewResponse maps a domain review to its response model.
func ToDocumentReviewResponse(r *model.DocumentReview) DocumentReviewResponse {
	return DocumentReviewResponse{
		ID:           r.ID,
		DocumentID:   r.DocumentID,
		ReviewerID:   r.ReviewerID,
		Status:       string(r.Status),
		MissingItems: r.MissingItems,
		ReviewNotes:  r.ReviewNotes,
		ReviewedAt:   r.ReviewedAt,
	}
}

// ToDocumentReviewResponseList maps a slice of reviews to response models.
func ToDocumentReviewResponseList(items []*model.DocumentReview) []DocumentReviewResponse {
	out := make([]DocumentReviewResponse, 0, len(items))
	for _, r := range items {
		out = append(out, ToDocumentReviewResponse(r))
	}
	return out
}

// ToDocumentTypeResponse maps a domain document type to its response model.
func ToDocumentTypeResponse(t *model.DocumentType) DocumentTypeResponse {
	return DocumentTypeResponse{
		ID:                     t.ID,
		Code:                   t.Code,
		Name:                   t.Name,
		Description:            t.Description,
		Category:               t.Category,
		IsMandatory:            t.IsMandatory,
		RequiresNotary:         t.RequiresNotary,
		RequiresMunicipalStamp: t.RequiresMunicipalStamp,
		IsValidWithoutNotary:   t.IsValidWithoutNotary,
		InvalidReason:          t.InvalidReason,
	}
}

// ToDocumentTypeResponseList maps a slice of document types to response models.
func ToDocumentTypeResponseList(items []*model.DocumentType) []DocumentTypeResponse {
	out := make([]DocumentTypeResponse, 0, len(items))
	for _, t := range items {
		out = append(out, ToDocumentTypeResponse(t))
	}
	return out
}
