package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// AppointmentCommandHandler handles write operations for appointments (CQRS write side).
type AppointmentCommandHandler struct {
	repo     repository.AppointmentRepository
	eventBus events.EventBus
}

// NewAppointmentCommandHandler creates a new AppointmentCommandHandler.
func NewAppointmentCommandHandler(repo repository.AppointmentRepository, eventBus events.EventBus) *AppointmentCommandHandler {
	return &AppointmentCommandHandler{repo: repo, eventBus: eventBus}
}

func (h *AppointmentCommandHandler) publish(ctx context.Context, action string, id, orgID, appID uuid.UUID) {
	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform, map[string]interface{}{
		"entity_type":     "appointment",
		"action":          action,
		"appointment_id":  id,
		"organization_id": orgID,
		"app_id":          appID,
	})
}

// Create schedules a new appointment within the tenant scope.
func (h *AppointmentCommandHandler) Create(ctx context.Context, orgID, appID, createdBy uuid.UUID, req dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error) {
	appointment := &model.Appointment{
		OrganizationID:  orgID,
		AppID:           appID,
		ProjectID:       req.ProjectID,
		OwnerID:         req.OwnerID,
		Title:           req.Title,
		Description:     req.Description,
		Location:        req.Location,
		ScheduledAt:     req.ScheduledAt,
		DurationMinutes: req.DurationMinutes,
		Status:          model.AppointmentStatusScheduled,
	}
	if createdBy != uuid.Nil {
		appointment.CreatedBy = &createdBy
	}
	if err := h.repo.Create(ctx, appointment); err != nil {
		return nil, err
	}
	h.publish(ctx, "created", appointment.ID, orgID, appID)
	resp := dto.ToAppointmentResponse(appointment)
	return &resp, nil
}

// Update applies a partial update to an appointment.
func (h *AppointmentCommandHandler) Update(ctx context.Context, orgID, appID, id uuid.UUID, req dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error) {
	appointment, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		appointment.Title = *req.Title
	}
	if req.Description != nil {
		appointment.Description = *req.Description
	}
	if req.Location != nil {
		appointment.Location = *req.Location
	}
	if req.ScheduledAt != nil {
		appointment.ScheduledAt = *req.ScheduledAt
	}
	if req.DurationMinutes != nil {
		appointment.DurationMinutes = req.DurationMinutes
	}
	if req.Status != nil {
		status := model.AppointmentStatus(*req.Status)
		if !model.IsValidAppointmentStatus(status) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid appointment status", nil)
		}
		appointment.Status = status
	}
	if err := h.repo.Update(ctx, appointment); err != nil {
		return nil, err
	}
	h.publish(ctx, "updated", appointment.ID, orgID, appID)
	resp := dto.ToAppointmentResponse(appointment)
	return &resp, nil
}

// Delete removes an appointment within the tenant scope.
func (h *AppointmentCommandHandler) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}
	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}
	h.publish(ctx, "deleted", id, orgID, appID)
	return nil
}
