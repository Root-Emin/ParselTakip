package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// MessageCommandHandler handles write operations for messages (CQRS write side).
type MessageCommandHandler struct {
	repo     repository.MessageRepository
	eventBus events.EventBus
}

// NewMessageCommandHandler creates a new MessageCommandHandler.
func NewMessageCommandHandler(repo repository.MessageRepository, eventBus events.EventBus) *MessageCommandHandler {
	return &MessageCommandHandler{repo: repo, eventBus: eventBus}
}

// Create sends a message from the sender to the recipient.
func (h *MessageCommandHandler) Create(ctx context.Context, orgID, appID, senderID uuid.UUID, req dto.CreateMessageRequest) (*dto.MessageResponse, error) {
	message := &model.Message{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      req.ProjectID,
		SenderID:       senderID,
		RecipientID:    req.RecipientID,
		ParentID:       req.ParentID,
		Subject:        req.Subject,
		Body:           req.Body,
	}
	if err := h.repo.Create(ctx, message); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform, map[string]interface{}{
		"entity_type":     "message",
		"action":          "created",
		"message_id":      message.ID,
		"recipient_id":    req.RecipientID,
		"organization_id": orgID,
		"app_id":          appID,
	})

	resp := dto.ToMessageResponse(message)
	return &resp, nil
}

// MarkRead marks a message (owned as recipient) as read.
func (h *MessageCommandHandler) MarkRead(ctx context.Context, orgID, appID, userID, id uuid.UUID) error {
	return h.repo.MarkRead(ctx, orgID, appID, userID, id)
}
