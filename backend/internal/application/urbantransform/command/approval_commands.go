package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// ApprovalCommandHandler handles write operations for approvals (CQRS write side).
type ApprovalCommandHandler struct {
	repo     repository.ApprovalRepository
	eventBus events.EventBus
}

// NewApprovalCommandHandler creates a new ApprovalCommandHandler.
func NewApprovalCommandHandler(repo repository.ApprovalRepository, eventBus events.EventBus) *ApprovalCommandHandler {
	return &ApprovalCommandHandler{repo: repo, eventBus: eventBus}
}

// Create opens a new pending approval within the tenant scope.
func (h *ApprovalCommandHandler) Create(ctx context.Context, orgID, appID uuid.UUID, req dto.CreateApprovalRequest) (*dto.ApprovalResponse, error) {
	approvalType := model.ApprovalType(req.ApprovalType)
	if !model.IsValidApprovalType(approvalType) {
		return nil, domainErr.New(domainErr.ErrValidation, "invalid approval type", nil)
	}

	approval := &model.Approval{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      req.ProjectID,
		ApprovalType:   approvalType,
		ApproverRole:   req.ApproverRole,
		OwnerID:        req.OwnerID,
		Status:         model.ApprovalStatusPending,
		ExpiresAt:      req.ExpiresAt,
	}

	if err := h.repo.Create(ctx, approval); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeApproval, event.ActionCreated, approval.ID, orgID, appID))

	resp := dto.ToApprovalResponse(approval)
	return &resp, nil
}

// Decide approves or rejects a pending approval, recording the decider.
func (h *ApprovalCommandHandler) Decide(ctx context.Context, orgID, appID, id, deciderID uuid.UUID, req dto.DecideApprovalRequest) (*dto.ApprovalResponse, error) {
	approval, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}

	if approval.Status != model.ApprovalStatusPending {
		return nil, domainErr.New(domainErr.ErrConflict, "approval already decided", nil)
	}

	status := model.ApprovalStatus(req.Status)
	if status != model.ApprovalStatusApproved && status != model.ApprovalStatusRejected {
		return nil, domainErr.New(domainErr.ErrValidation, "decision must be approved or rejected", nil)
	}

	now := time.Now().UTC()
	approval.Status = status
	approval.DecisionNotes = req.DecisionNotes
	approval.DecidedAt = &now
	if deciderID != uuid.Nil {
		approval.ApproverID = &deciderID
	}

	if err := h.repo.Update(ctx, approval); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeApproval, event.ActionDecided, approval.ID, orgID, appID))

	resp := dto.ToApprovalResponse(approval)
	return &resp, nil
}

// Delete removes an approval within the tenant scope.
func (h *ApprovalCommandHandler) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}
	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}
	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeApproval, event.ActionDeleted, id, orgID, appID))
	return nil
}
