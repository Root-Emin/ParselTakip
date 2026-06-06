package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// WorkflowCommandHandler handles workflow state transitions (CQRS write side).
type WorkflowCommandHandler struct {
	repo     repository.WorkflowRepository
	eventBus events.EventBus
}

// NewWorkflowCommandHandler creates a new WorkflowCommandHandler.
func NewWorkflowCommandHandler(repo repository.WorkflowRepository, eventBus events.EventBus) *WorkflowCommandHandler {
	return &WorkflowCommandHandler{repo: repo, eventBus: eventBus}
}

// Advance creates or updates a project's step state and records history.
func (h *WorkflowCommandHandler) Advance(ctx context.Context, orgID, appID, projectID, changedBy uuid.UUID, req dto.AdvanceWorkflowRequest) (*dto.ProjectWorkflowStateResponse, error) {
	status := model.WorkflowStateStatus(req.Status)
	if !model.IsValidWorkflowStatus(status) {
		return nil, domainErr.New(domainErr.ErrValidation, "invalid workflow status", nil)
	}

	// Ensure the target step definition exists.
	if _, err := h.repo.GetStepDefinition(ctx, req.WorkflowStepID); err != nil {
		return nil, err
	}

	// Determine the previously active step (for history "from").
	var fromStepID *uuid.UUID
	if existing, err := h.repo.ListStates(ctx, orgID, appID, projectID); err == nil {
		for _, s := range existing {
			if s.Status == model.WorkflowStatusInProgress && s.WorkflowStepID != req.WorkflowStepID {
				step := s.WorkflowStepID
				fromStepID = &step
				break
			}
		}
	}

	now := time.Now().UTC()
	state := &model.ProjectWorkflowState{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      projectID,
		WorkflowStepID: req.WorkflowStepID,
		Status:         status,
		DueAt:          req.DueAt,
		BlockedReason:  req.BlockedReason,
	}
	if changedBy != uuid.Nil {
		state.UpdatedBy = &changedBy
	}
	switch status {
	case model.WorkflowStatusInProgress:
		state.StartedAt = &now
	case model.WorkflowStatusCompleted:
		state.CompletedAt = &now
	}

	if err := h.repo.UpsertState(ctx, state); err != nil {
		return nil, err
	}

	history := &model.ProjectWorkflowHistory{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      projectID,
		FromStepID:     fromStepID,
		ToStepID:       req.WorkflowStepID,
		Action:         req.Status,
		Notes:          req.Notes,
	}
	if changedBy != uuid.Nil {
		history.ChangedBy = &changedBy
	}
	if err := h.repo.AddHistory(ctx, history); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform, map[string]interface{}{
		"entity_type":      "workflow",
		"action":           "advanced",
		"project_id":       projectID,
		"workflow_step_id": req.WorkflowStepID,
		"status":           req.Status,
		"organization_id":  orgID,
		"app_id":           appID,
		"timestamp":        now,
	})

	resp := dto.ToWorkflowStateResponse(state)
	return &resp, nil
}
