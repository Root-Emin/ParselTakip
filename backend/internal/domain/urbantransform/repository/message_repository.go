package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// MessageFilter holds optional filter criteria for querying messages.
type MessageFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	UserID         uuid.UUID // the current user
	Box            string    // "inbox" (recipient=user) or "sent" (sender=user)
	ProjectID      *uuid.UUID
	IsRead         *bool
	Offset         int
	Limit          int
}

// MessageRepository defines persistence operations for messages.
type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Message, error)
	List(ctx context.Context, filter MessageFilter) ([]*model.Message, int, error)
	MarkRead(ctx context.Context, orgID, appID, userID, id uuid.UUID) error
	CountUnread(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error)
}
