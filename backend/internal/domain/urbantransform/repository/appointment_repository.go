package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// AppointmentFilter holds optional filter criteria for querying appointments.
type AppointmentFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	ProjectID      *uuid.UUID
	OwnerID        *uuid.UUID
	Status         *model.AppointmentStatus
	From           *time.Time // scheduled_at >= From (e.g. upcoming)
	SortOrder      string
	Offset         int
	Limit          int
}

// AppointmentRepository defines persistence operations for appointments.
type AppointmentRepository interface {
	Create(ctx context.Context, appointment *model.Appointment) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Appointment, error)
	Update(ctx context.Context, appointment *model.Appointment) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter AppointmentFilter) ([]*model.Appointment, int, error)
}
