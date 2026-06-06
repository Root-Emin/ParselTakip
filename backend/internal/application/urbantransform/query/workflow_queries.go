package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
)

// WorkflowQueryHandler handles workflow read operations (CQRS read side).
type WorkflowQueryHandler struct {
	repo repository.WorkflowRepository
}

// NewWorkflowQueryHandler creates a new WorkflowQueryHandler.
func NewWorkflowQueryHandler(repo repository.WorkflowRepository) *WorkflowQueryHandler {
	return &WorkflowQueryHandler{repo: repo}
}

// ListSteps returns all global workflow step definitions, ordered.
func (h *WorkflowQueryHandler) ListSteps(ctx context.Context) ([]dto.WorkflowStepResponse, error) {
	steps, err := h.repo.ListStepDefinitions(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ToWorkflowStepResponseList(steps), nil
}

// ListProjectStates returns the workflow states of a single project.
func (h *WorkflowQueryHandler) ListProjectStates(ctx context.Context, orgID, appID, projectID uuid.UUID) ([]dto.ProjectWorkflowStateResponse, error) {
	states, err := h.repo.ListStates(ctx, orgID, appID, projectID)
	if err != nil {
		return nil, err
	}
	return dto.ToWorkflowStateResponseList(states), nil
}

// ListProjectHistory returns the workflow transition history of a project.
func (h *WorkflowQueryHandler) ListProjectHistory(ctx context.Context, orgID, appID, projectID uuid.UUID) ([]dto.WorkflowHistoryResponse, error) {
	history, err := h.repo.ListHistory(ctx, orgID, appID, projectID)
	if err != nil {
		return nil, err
	}
	return dto.ToWorkflowHistoryResponseList(history), nil
}
