package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/pagination"
)

// AppointmentQueryHandler handles read operations for appointments (CQRS read side).
type AppointmentQueryHandler struct {
	repo repository.AppointmentRepository
}

// NewAppointmentQueryHandler creates a new AppointmentQueryHandler.
func NewAppointmentQueryHandler(repo repository.AppointmentRepository) *AppointmentQueryHandler {
	return &AppointmentQueryHandler{repo: repo}
}

// Get returns a single appointment by ID within the tenant scope.
func (h *AppointmentQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.AppointmentResponse, error) {
	appointment, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToAppointmentResponse(appointment)
	return &resp, nil
}

// List returns a paginated, filtered list of appointments.
func (h *AppointmentQueryHandler) List(ctx context.Context, orgID, appID uuid.UUID, q dto.ListAppointmentsQuery) (pagination.Result[dto.AppointmentResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	filter := repository.AppointmentFilter{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      q.ProjectID,
		OwnerID:        q.OwnerID,
		SortOrder:      constants.NormalizeSortOrder(q.SortOrder),
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}
	if q.Status != nil {
		s := model.AppointmentStatus(*q.Status)
		filter.Status = &s
	}
	if q.Upcoming {
		now := time.Now().UTC()
		filter.From = &now
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.AppointmentResponse]{}, err
	}
	return pagination.NewResult(dto.ToAppointmentResponseList(items), params, total), nil
}
