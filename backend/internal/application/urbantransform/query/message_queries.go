package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/pagination"
)

// MessageQueryHandler handles read operations for messages (CQRS read side).
type MessageQueryHandler struct {
	repo repository.MessageRepository
}

// NewMessageQueryHandler creates a new MessageQueryHandler.
func NewMessageQueryHandler(repo repository.MessageRepository) *MessageQueryHandler {
	return &MessageQueryHandler{repo: repo}
}

// Get returns a single message by ID within the tenant scope.
func (h *MessageQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.MessageResponse, error) {
	m, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToMessageResponse(m)
	return &resp, nil
}

// List returns the current user's inbox or sent messages, paginated.
func (h *MessageQueryHandler) List(ctx context.Context, orgID, appID, userID uuid.UUID, q dto.ListMessagesQuery) (pagination.Result[dto.MessageResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	box := q.Box
	if box != "sent" {
		box = "inbox"
	}

	filter := repository.MessageFilter{
		OrganizationID: orgID,
		AppID:          appID,
		UserID:         userID,
		Box:            box,
		ProjectID:      q.ProjectID,
		IsRead:         q.IsRead,
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.MessageResponse]{}, err
	}
	return pagination.NewResult(dto.ToMessageResponseList(items), params, total), nil
}

// UnreadCount returns the number of unread messages for the current user.
func (h *MessageQueryHandler) UnreadCount(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error) {
	return h.repo.CountUnread(ctx, orgID, appID, userID)
}
