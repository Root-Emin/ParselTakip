package model

import (
	"time"

	"github.com/google/uuid"
)

// Message is a direct message between two users, optionally tied to a project.
type Message struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	AppID          uuid.UUID  `json:"app_id"`
	ProjectID      *uuid.UUID `json:"project_id,omitempty"`
	SenderID       uuid.UUID  `json:"sender_id"`
	RecipientID    uuid.UUID  `json:"recipient_id"`
	ParentID       *uuid.UUID `json:"parent_id,omitempty"`
	Subject        string     `json:"subject"`
	Body           string     `json:"body"`
	IsRead         bool       `json:"is_read"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
