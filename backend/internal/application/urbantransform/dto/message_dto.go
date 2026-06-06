package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateMessageRequest is the input for sending a message (command).
type CreateMessageRequest struct {
	RecipientID uuid.UUID  `json:"recipient_id" validate:"required"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body" validate:"required"`
}

// ListMessagesQuery is the input for listing the current user's messages.
type ListMessagesQuery struct {
	Box       string
	ProjectID *uuid.UUID
	IsRead    *bool
	Page      int
	PerPage   int
}

// MessageResponse is the public representation of a message (response model).
type MessageResponse struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
	SenderID    uuid.UUID  `json:"sender_id"`
	RecipientID uuid.UUID  `json:"recipient_id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body"`
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToMessageResponse maps a domain message to its response model.
func ToMessageResponse(m *model.Message) MessageResponse {
	return MessageResponse{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		SenderID:    m.SenderID,
		RecipientID: m.RecipientID,
		ParentID:    m.ParentID,
		Subject:     m.Subject,
		Body:        m.Body,
		IsRead:      m.IsRead,
		ReadAt:      m.ReadAt,
		CreatedAt:   m.CreatedAt,
	}
}

// ToMessageResponseList maps a slice of messages to response models.
func ToMessageResponseList(items []*model.Message) []MessageResponse {
	out := make([]MessageResponse, 0, len(items))
	for _, m := range items {
		out = append(out, ToMessageResponse(m))
	}
	return out
}
