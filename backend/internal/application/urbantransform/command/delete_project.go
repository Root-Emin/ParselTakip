package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// DeleteProjectHandler handles the delete-project command (CQRS write side).
type DeleteProjectHandler struct {
	repo     repository.ProjectRepository
	eventBus events.EventBus
}

// NewDeleteProjectHandler creates a new DeleteProjectHandler.
func NewDeleteProjectHandler(repo repository.ProjectRepository, eventBus events.EventBus) *DeleteProjectHandler {
	return &DeleteProjectHandler{repo: repo, eventBus: eventBus}
}

// Execute removes a project within the given tenant scope.
func (h *DeleteProjectHandler) Execute(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}

	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform, event.ProjectDeleted{
		ProjectID:      id,
		OrganizationID: orgID,
		AppID:          appID,
		Timestamp:      time.Now().UTC(),
	})

	return nil
}
