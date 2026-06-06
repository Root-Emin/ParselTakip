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

// CreateProjectHandler handles the create-project command (CQRS write side).
type CreateProjectHandler struct {
	repo     repository.ProjectRepository
	eventBus events.EventBus
}

// NewCreateProjectHandler creates a new CreateProjectHandler.
func NewCreateProjectHandler(repo repository.ProjectRepository, eventBus events.EventBus) *CreateProjectHandler {
	return &CreateProjectHandler{repo: repo, eventBus: eventBus}
}

// Execute creates a new project within the given tenant scope.
func (h *CreateProjectHandler) Execute(ctx context.Context, orgID, appID, initiatedBy uuid.UUID, req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	if existing, _ := h.repo.GetByCode(ctx, orgID, req.Code); existing != nil {
		return nil, domainErr.New(domainErr.ErrAlreadyExists, "project code already exists", nil)
	}

	project := &model.Project{
		OrganizationID:       orgID,
		AppID:                appID,
		Code:                 req.Code,
		Name:                 req.Name,
		Description:          req.Description,
		Status:               model.ProjectStatusDraft,
		AssignedContractorID: req.AssignedContractorID,
		TargetCompletionAt:   req.TargetCompletionAt,
	}
	if initiatedBy != uuid.Nil {
		project.InitiatedBy = &initiatedBy
	}

	if err := h.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform, event.ProjectCreated{
		ProjectID:      project.ID,
		OrganizationID: orgID,
		AppID:          appID,
		Code:           project.Code,
		Name:           project.Name,
		Timestamp:      time.Now().UTC(),
	})

	resp := dto.ToProjectResponse(project)
	return &resp, nil
}
