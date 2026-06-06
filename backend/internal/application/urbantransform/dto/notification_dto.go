package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateNotificationRequest is the input for sending a notification (command).
type CreateNotificationRequest struct {
	UserID           uuid.UUID       `json:"user_id" validate:"required"`
	ProjectID        *uuid.UUID      `json:"project_id,omitempty"`
	NotificationType string          `json:"notification_type" validate:"required"`
	Title            string          `json:"title" validate:"required"`
	Message          string          `json:"message" validate:"required"`
	Channel          string          `json:"channel" validate:"omitempty,oneof=in_app email sms"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
}

// ListNotificationsQuery is the input for listing the current user's notifications.
type ListNotificationsQuery struct {
	ProjectID        *uuid.UUID
	NotificationType *string
	IsRead           *bool
	SortOrder        string
	Page             int
	PerPage          int
}

// NotificationResponse is the public representation of a notification (response model).
type NotificationResponse struct {
	ID               uuid.UUID       `json:"id"`
	OrganizationID   uuid.UUID       `json:"organization_id"`
	AppID            uuid.UUID       `json:"app_id"`
	ProjectID        *uuid.UUID      `json:"project_id,omitempty"`
	UserID           uuid.UUID       `json:"user_id"`
	NotificationType string          `json:"notification_type"`
	Title            string          `json:"title"`
	Message          string          `json:"message"`
	Channel          string          `json:"channel"`
	IsRead           bool            `json:"is_read"`
	ReadAt           *time.Time      `json:"read_at,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}

// ToNotificationResponse maps a domain notification to its response model.
func ToNotificationResponse(n *model.Notification) NotificationResponse {
	return NotificationResponse{
		ID:               n.ID,
		OrganizationID:   n.OrganizationID,
		AppID:            n.AppID,
		ProjectID:        n.ProjectID,
		UserID:           n.UserID,
		NotificationType: n.NotificationType,
		Title:            n.Title,
		Message:          n.Message,
		Channel:          string(n.Channel),
		IsRead:           n.IsRead,
		ReadAt:           n.ReadAt,
		Metadata:         n.Metadata,
		CreatedAt:        n.CreatedAt,
	}
}

// ToNotificationResponseList maps a slice of notifications to response models.
func ToNotificationResponseList(items []*model.Notification) []NotificationResponse {
	out := make([]NotificationResponse, 0, len(items))
	for _, n := range items {
		out = append(out, ToNotificationResponse(n))
	}
	return out
}
