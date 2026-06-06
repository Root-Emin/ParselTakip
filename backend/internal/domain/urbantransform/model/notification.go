package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// NotificationChannel represents the delivery channel of a notification.
type NotificationChannel string

const (
	NotificationChannelInApp NotificationChannel = "in_app"
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelSMS   NotificationChannel = "sms"
)

// IsValidNotificationChannel reports whether the given channel is known.
func IsValidNotificationChannel(c NotificationChannel) bool {
	switch c {
	case NotificationChannelInApp, NotificationChannelEmail, NotificationChannelSMS:
		return true
	default:
		return false
	}
}

// Valid notification type values (mirrors the DB CHECK constraint).
var validNotificationTypes = map[string]struct{}{
	"document_required": {}, "document_approved": {}, "document_rejected": {}, "document_missing": {},
	"step_completed": {}, "step_blocked": {}, "approval_required": {}, "objection_deadline": {},
	"meeting_scheduled": {}, "contract_signed": {}, "permit_issued": {}, "demolition_scheduled": {},
	"rent_assistance_update": {}, "construction_update": {}, "title_deed_ready": {}, "general": {},
}

// IsValidNotificationType reports whether the given type is known.
func IsValidNotificationType(t string) bool {
	_, ok := validNotificationTypes[t]
	return ok
}

// Notification is an in-app/email/sms message addressed to a user.
type Notification struct {
	ID               uuid.UUID           `json:"id"`
	OrganizationID   uuid.UUID           `json:"organization_id"`
	AppID            uuid.UUID           `json:"app_id"`
	ProjectID        *uuid.UUID          `json:"project_id,omitempty"`
	UserID           uuid.UUID           `json:"user_id"`
	NotificationType string              `json:"notification_type"`
	Title            string              `json:"title"`
	Message          string              `json:"message"`
	Channel          NotificationChannel `json:"channel"`
	IsRead           bool                `json:"is_read"`
	ReadAt           *time.Time          `json:"read_at,omitempty"`
	Metadata         json.RawMessage     `json:"metadata,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
}
