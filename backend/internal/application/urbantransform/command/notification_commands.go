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

// NotificationCommandHandler handles write operations for notifications (CQRS write side).
type NotificationCommandHandler struct {
	repo     repository.NotificationRepository
	eventBus events.EventBus
}

// NewNotificationCommandHandler creates a new NotificationCommandHandler.
func NewNotificationCommandHandler(repo repository.NotificationRepository, eventBus events.EventBus) *NotificationCommandHandler {
	return &NotificationCommandHandler{repo: repo, eventBus: eventBus}
}

// Create sends a notification to a recipient within the tenant scope.
func (h *NotificationCommandHandler) Create(ctx context.Context, orgID, appID uuid.UUID, req dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	if !model.IsValidNotificationType(req.NotificationType) {
		return nil, domainErr.New(domainErr.ErrValidation, "invalid notification type", nil)
	}
	channel := model.NotificationChannel(req.Channel)
	if channel == "" {
		channel = model.NotificationChannelInApp
	}
	if !model.IsValidNotificationChannel(channel) {
		return nil, domainErr.New(domainErr.ErrValidation, "invalid notification channel", nil)
	}

	notification := &model.Notification{
		OrganizationID:   orgID,
		AppID:            appID,
		ProjectID:        req.ProjectID,
		UserID:           req.UserID,
		NotificationType: req.NotificationType,
		Title:            req.Title,
		Message:          req.Message,
		Channel:          channel,
		Metadata:         req.Metadata,
	}

	if err := h.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeNotification, event.ActionCreated, notification.ID, orgID, appID))

	resp := dto.ToNotificationResponse(notification)
	return &resp, nil
}

// MarkRead marks a single notification (owned by the user) as read.
func (h *NotificationCommandHandler) MarkRead(ctx context.Context, orgID, appID, userID, id uuid.UUID) error {
	return h.repo.MarkRead(ctx, orgID, appID, userID, id)
}

// MarkAllRead marks all of the user's notifications as read and returns the count.
func (h *NotificationCommandHandler) MarkAllRead(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error) {
	return h.repo.MarkAllRead(ctx, orgID, appID, userID)
}
