package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// ProjectFilter holds optional filter criteria for querying projects.
// Nil fields are ignored. Multi-tenant scoping (OrganizationID/AppID) is always required.
type ProjectFilter struct {
	OrganizationID       uuid.UUID
	AppID                uuid.UUID
	Status               *model.ProjectStatus
	AssignedContractorID *uuid.UUID
	InitiatedBy          *uuid.UUID
	CurrentWorkflowStepID *uuid.UUID
	Search               string // free-text search on code/name/description
	SortBy               string // whitelisted column
	SortOrder            string // asc | desc
	Offset               int
	Limit                int
}

// ProjectRepository defines persistence operations for projects.
type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Project, error)
	GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*model.Project, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter ProjectFilter) ([]*model.Project, int, error)
}
