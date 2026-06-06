package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateAppointmentRequest is the input for creating an appointment (command).
type CreateAppointmentRequest struct {
	ProjectID       *uuid.UUID `json:"project_id,omitempty"`
	OwnerID         *uuid.UUID `json:"owner_id,omitempty"`
	Title           string     `json:"title" validate:"required"`
	Description     string     `json:"description"`
	Location        string     `json:"location"`
	ScheduledAt     time.Time  `json:"scheduled_at" validate:"required"`
	DurationMinutes *int       `json:"duration_minutes,omitempty"`
}

// UpdateAppointmentRequest is the input for updating an appointment (partial command).
type UpdateAppointmentRequest struct {
	Title           *string    `json:"title,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Location        *string    `json:"location,omitempty"`
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty"`
	DurationMinutes *int       `json:"duration_minutes,omitempty"`
	Status          *string    `json:"status,omitempty" validate:"omitempty,oneof=scheduled completed cancelled rescheduled"`
}

// ListAppointmentsQuery is the input for listing/filtering appointments (query).
type ListAppointmentsQuery struct {
	ProjectID *uuid.UUID
	OwnerID   *uuid.UUID
	Status    *string
	Upcoming  bool
	SortOrder string
	Page      int
	PerPage   int
}

// AppointmentResponse is the public representation of an appointment (response model).
type AppointmentResponse struct {
	ID              uuid.UUID  `json:"id"`
	ProjectID       *uuid.UUID `json:"project_id,omitempty"`
	OwnerID         *uuid.UUID `json:"owner_id,omitempty"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Location        string     `json:"location"`
	ScheduledAt     time.Time  `json:"scheduled_at"`
	DurationMinutes *int       `json:"duration_minutes,omitempty"`
	Status          string     `json:"status"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ToAppointmentResponse maps a domain appointment to its response model.
func ToAppointmentResponse(a *model.Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:              a.ID,
		ProjectID:       a.ProjectID,
		OwnerID:         a.OwnerID,
		Title:           a.Title,
		Description:     a.Description,
		Location:        a.Location,
		ScheduledAt:     a.ScheduledAt,
		DurationMinutes: a.DurationMinutes,
		Status:          string(a.Status),
		CreatedBy:       a.CreatedBy,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

// ToAppointmentResponseList maps a slice of appointments to response models.
func ToAppointmentResponseList(items []*model.Appointment) []AppointmentResponse {
	out := make([]AppointmentResponse, 0, len(items))
	for _, a := range items {
		out = append(out, ToAppointmentResponse(a))
	}
	return out
}
