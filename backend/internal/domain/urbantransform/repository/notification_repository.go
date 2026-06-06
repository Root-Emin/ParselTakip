package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// NotificationFilter holds optional filter criteria for querying notifications.
// Notifications are always scoped to a single recipient (UserID).
type NotificationFilter struct {
	OrganizationID   uuid.UUID
	AppID            uuid.UUID
	UserID           uuid.UUID
	ProjectID        *uuid.UUID
	NotificationType *string
	IsRead           *bool
	SortOrder        string
	Offset           int
	Limit            int
}

// NotificationRepository defines persistence operations for notifications.
type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Notification, error)
	List(ctx context.Context, filter NotificationFilter) ([]*model.Notification, int, error)
	MarkRead(ctx context.Context, orgID, appID, userID, id uuid.UUID) error
	MarkAllRead(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error)
	CountUnread(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error)
}
