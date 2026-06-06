package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DocumentStatus represents the lifecycle status of a document.
type DocumentStatus string

const (
	DocumentStatusDraft       DocumentStatus = "draft"
	DocumentStatusSubmitted   DocumentStatus = "submitted"
	DocumentStatusUnderReview DocumentStatus = "under_review"
	DocumentStatusApproved    DocumentStatus = "approved"
	DocumentStatusRejected    DocumentStatus = "rejected"
	DocumentStatusMissing     DocumentStatus = "missing"
	DocumentStatusExpired     DocumentStatus = "expired"
	DocumentStatusInvalid     DocumentStatus = "invalid"
)

// IsValidDocumentStatus reports whether the given status is known.
func IsValidDocumentStatus(s DocumentStatus) bool {
	switch s {
	case DocumentStatusDraft, DocumentStatusSubmitted, DocumentStatusUnderReview,
		DocumentStatusApproved, DocumentStatusRejected, DocumentStatusMissing,
		DocumentStatusExpired, DocumentStatusInvalid:
		return true
	default:
		return false
	}
}

// Document represents an uploaded file tied to a project (and optionally
// building/unit/owner) within the transformation workflow.
type Document struct {
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
	Status         DocumentStatus  `json:"status"`
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

// DocumentReviewStatus represents the outcome of a document review.
type DocumentReviewStatus string

const (
	ReviewStatusApproved      DocumentReviewStatus = "approved"
	ReviewStatusRejected      DocumentReviewStatus = "rejected"
	ReviewStatusMissingItems  DocumentReviewStatus = "missing_items"
	ReviewStatusNeedsRevision DocumentReviewStatus = "needs_revision"
)

// IsValidReviewStatus reports whether the given review status is known.
func IsValidReviewStatus(s DocumentReviewStatus) bool {
	switch s {
	case ReviewStatusApproved, ReviewStatusRejected, ReviewStatusMissingItems, ReviewStatusNeedsRevision:
		return true
	default:
		return false
	}
}

// DocumentReview records a reviewer's decision on a document.
type DocumentReview struct {
	ID             uuid.UUID            `json:"id"`
	OrganizationID uuid.UUID            `json:"organization_id"`
	AppID          uuid.UUID            `json:"app_id"`
	DocumentID     uuid.UUID            `json:"document_id"`
	ReviewerID     uuid.UUID            `json:"reviewer_id"`
	Status         DocumentReviewStatus `json:"status"`
	MissingItems   json.RawMessage      `json:"missing_items,omitempty"`
	ReviewNotes    string               `json:"review_notes"`
	ReviewedAt     time.Time            `json:"reviewed_at"`
	CreatedAt      time.Time            `json:"created_at"`
}

// DocumentType is master data describing a kind of required document.
type DocumentType struct {
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
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
