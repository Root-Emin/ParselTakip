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

// UpdateProjectHandler handles the update-project command (CQRS write side).
type UpdateProjectHandler struct {
	repo     repository.ProjectRepository
	eventBus events.EventBus
}

// NewUpdateProjectHandler creates a new UpdateProjectHandler.
func NewUpdateProjectHandler(repo repository.ProjectRepository, eventBus events.EventBus) *UpdateProjectHandler {
	return &UpdateProjectHandler{repo: repo, eventBus: eventBus}
}

// Execute applies a partial update to an existing project.
func (h *UpdateProjectHandler) Execute(ctx context.Context, orgID, appID, id uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	project, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Status != nil {
		status := model.ProjectStatus(*req.Status)
		if !model.IsValidProjectStatus(status) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid project status", nil)
		}
		project.Status = status
		switch status {
		case model.ProjectStatusInProgress:
			if project.StartedAt == nil {
				now := time.Now().UTC()
				project.StartedAt = &now
			}
		case model.ProjectStatusCompleted:
			now := time.Now().UTC()
			project.CompletedAt = &now
		}
	}
	if req.AssignedContractorID != nil {
		project.AssignedContractorID = req.AssignedContractorID
	}
	if req.TargetCompletionAt != nil {
		project.TargetCompletionAt = req.TargetCompletionAt
	}

	if err := h.repo.Update(ctx, project); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform, event.ProjectUpdated{
		ProjectID:      project.ID,
		OrganizationID: orgID,
		AppID:          appID,
		Status:         string(project.Status),
		Timestamp:      time.Now().UTC(),
	})

	resp := dto.ToProjectResponse(project)
	return &resp, nil
}
