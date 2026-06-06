package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// DocumentCommandHandler handles write operations for documents (CQRS write side).
type DocumentCommandHandler struct {
	repo     repository.DocumentRepository
	eventBus events.EventBus
}

// NewDocumentCommandHandler creates a new DocumentCommandHandler.
func NewDocumentCommandHandler(repo repository.DocumentRepository, eventBus events.EventBus) *DocumentCommandHandler {
	return &DocumentCommandHandler{repo: repo, eventBus: eventBus}
}

// Create registers a new document upload within the tenant scope.
func (h *DocumentCommandHandler) Create(ctx context.Context, orgID, appID, uploadedBy uuid.UUID, uploadedByRole string, req dto.CreateDocumentRequest) (*dto.DocumentResponse, error) {
	doc := &model.Document{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      req.ProjectID,
		DocumentTypeID: req.DocumentTypeID,
		BuildingID:     req.BuildingID,
		UnitID:         req.UnitID,
		OwnerID:        req.OwnerID,
		FileName:       req.FileName,
		FilePath:       req.FilePath,
		FileSize:       req.FileSize,
		MimeType:       req.MimeType,
		Status:         model.DocumentStatusSubmitted,
		IsNotarized:    req.IsNotarized,
		NotaryDate:     req.NotaryDate,
		ExpiryDate:     req.ExpiryDate,
		UploadedByRole: uploadedByRole,
		Version:        1,
		Metadata:       req.Metadata,
	}
	if uploadedBy != uuid.Nil {
		doc.UploadedBy = &uploadedBy
	}

	if err := h.repo.Create(ctx, doc); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeDocument, event.ActionCreated, doc.ID, orgID, appID))

	resp := dto.ToDocumentResponse(doc)
	return &resp, nil
}

// Update applies a partial update to a document.
func (h *DocumentCommandHandler) Update(ctx context.Context, orgID, appID, id uuid.UUID, req dto.UpdateDocumentRequest) (*dto.DocumentResponse, error) {
	doc, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}

	if req.FileName != nil {
		doc.FileName = *req.FileName
	}
	if req.FilePath != nil {
		doc.FilePath = *req.FilePath
	}
	if req.FileSize != nil {
		doc.FileSize = req.FileSize
	}
	if req.MimeType != nil {
		doc.MimeType = *req.MimeType
	}
	if req.Status != nil {
		st := model.DocumentStatus(*req.Status)
		if !model.IsValidDocumentStatus(st) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid document status", nil)
		}
		doc.Status = st
	}
	if req.IsNotarized != nil {
		doc.IsNotarized = *req.IsNotarized
	}
	if req.NotaryDate != nil {
		doc.NotaryDate = req.NotaryDate
	}
	if req.ExpiryDate != nil {
		doc.ExpiryDate = req.ExpiryDate
	}
	if req.Metadata != nil {
		doc.Metadata = req.Metadata
	}

	if err := h.repo.Update(ctx, doc); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeDocument, event.ActionUpdated, doc.ID, orgID, appID))

	resp := dto.ToDocumentResponse(doc)
	return &resp, nil
}

// Delete removes a document within the tenant scope.
func (h *DocumentCommandHandler) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}
	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}
	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeDocument, event.ActionDeleted, id, orgID, appID))
	return nil
}

// Review records a reviewer decision and transitions the document status accordingly.
func (h *DocumentCommandHandler) Review(ctx context.Context, orgID, appID, documentID, reviewerID uuid.UUID, req dto.ReviewDocumentRequest) (*dto.DocumentReviewResponse, error) {
	doc, err := h.repo.GetByID(ctx, orgID, appID, documentID)
	if err != nil {
		return nil, err
	}

	reviewStatus := model.DocumentReviewStatus(req.Status)
	if !model.IsValidReviewStatus(reviewStatus) {
		return nil, domainErr.New(domainErr.ErrValidation, "invalid review status", nil)
	}

	review := &model.DocumentReview{
		OrganizationID: orgID,
		AppID:          appID,
		DocumentID:     documentID,
		ReviewerID:     reviewerID,
		Status:         reviewStatus,
		MissingItems:   req.MissingItems,
		ReviewNotes:    req.ReviewNotes,
	}
	if err := h.repo.CreateReview(ctx, review); err != nil {
		return nil, err
	}

	// Transition the document status based on the review outcome.
	switch reviewStatus {
	case model.ReviewStatusApproved:
		doc.Status = model.DocumentStatusApproved
	case model.ReviewStatusRejected:
		doc.Status = model.DocumentStatusRejected
	case model.ReviewStatusMissingItems:
		doc.Status = model.DocumentStatusMissing
	case model.ReviewStatusNeedsRevision:
		doc.Status = model.DocumentStatusUnderReview
	}
	if err := h.repo.Update(ctx, doc); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeDocument, event.ActionReviewed, documentID, orgID, appID))

	resp := dto.ToDocumentReviewResponse(review)
	return &resp, nil
}
