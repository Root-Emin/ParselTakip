package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/pagination"
)

// NotificationQueryHandler handles read operations for notifications (CQRS read side).
type NotificationQueryHandler struct {
	repo repository.NotificationRepository
}

// NewNotificationQueryHandler creates a new NotificationQueryHandler.
func NewNotificationQueryHandler(repo repository.NotificationRepository) *NotificationQueryHandler {
	return &NotificationQueryHandler{repo: repo}
}

// List returns the current user's notifications, filtered and paginated.
func (h *NotificationQueryHandler) List(ctx context.Context, orgID, appID, userID uuid.UUID, q dto.ListNotificationsQuery) (pagination.Result[dto.NotificationResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	filter := repository.NotificationFilter{
		OrganizationID:   orgID,
		AppID:            appID,
		UserID:           userID,
		ProjectID:        q.ProjectID,
		NotificationType: q.NotificationType,
		IsRead:           q.IsRead,
		SortOrder:        constants.NormalizeSortOrder(q.SortOrder),
		Offset:           params.Offset(),
		Limit:            params.Limit(),
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.NotificationResponse]{}, err
	}
	return pagination.NewResult(dto.ToNotificationResponseList(items), params, total), nil
}

// UnreadCount returns the number of unread notifications for the current user.
func (h *NotificationQueryHandler) UnreadCount(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error) {
	return h.repo.CountUnread(ctx, orgID, appID, userID)
}
