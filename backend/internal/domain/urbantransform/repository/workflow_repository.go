package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// WorkflowRepository defines persistence for workflow definitions, states and history.
type WorkflowRepository interface {
	// Step definitions (global master data)
	ListStepDefinitions(ctx context.Context) ([]*model.WorkflowStepDefinition, error)
	GetStepDefinition(ctx context.Context, id uuid.UUID) (*model.WorkflowStepDefinition, error)

	// Project workflow states
	ListStates(ctx context.Context, orgID, appID, projectID uuid.UUID) ([]*model.ProjectWorkflowState, error)
	GetState(ctx context.Context, orgID, appID, projectID, stepID uuid.UUID) (*model.ProjectWorkflowState, error)
	UpsertState(ctx context.Context, state *model.ProjectWorkflowState) error

	// Project workflow history
	AddHistory(ctx context.Context, history *model.ProjectWorkflowHistory) error
	ListHistory(ctx context.Context, orgID, appID, projectID uuid.UUID) ([]*model.ProjectWorkflowHistory, error)
}
