package model

import (
	"time"

	"github.com/google/uuid"
)

// AppointmentStatus represents the state of an appointment.
type AppointmentStatus string

const (
	AppointmentStatusScheduled   AppointmentStatus = "scheduled"
	AppointmentStatusCompleted   AppointmentStatus = "completed"
	AppointmentStatusCancelled   AppointmentStatus = "cancelled"
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

// IsValidAppointmentStatus reports whether the given status is known.
func IsValidAppointmentStatus(s AppointmentStatus) bool {
	switch s {
	case AppointmentStatusScheduled, AppointmentStatusCompleted, AppointmentStatusCancelled, AppointmentStatusRescheduled:
		return true
	default:
		return false
	}
}

// Appointment is a scheduled meeting/informational date tied to a project/owner.
type Appointment struct {
	ID              uuid.UUID         `json:"id"`
	OrganizationID  uuid.UUID         `json:"organization_id"`
	AppID           uuid.UUID         `json:"app_id"`
	ProjectID       *uuid.UUID        `json:"project_id,omitempty"`
	OwnerID         *uuid.UUID        `json:"owner_id,omitempty"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Location        string            `json:"location"`
	ScheduledAt     time.Time         `json:"scheduled_at"`
	DurationMinutes *int              `json:"duration_minutes,omitempty"`
	Status          AppointmentStatus `json:"status"`
	CreatedBy       *uuid.UUID        `json:"created_by,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
