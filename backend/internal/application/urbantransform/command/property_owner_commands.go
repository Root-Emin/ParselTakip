package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// PropertyOwnerCommandHandler handles write operations for property owners (CQRS write side).
type PropertyOwnerCommandHandler struct {
	repo     repository.PropertyOwnerRepository
	eventBus events.EventBus
}

// NewPropertyOwnerCommandHandler creates a new PropertyOwnerCommandHandler.
func NewPropertyOwnerCommandHandler(repo repository.PropertyOwnerRepository, eventBus events.EventBus) *PropertyOwnerCommandHandler {
	return &PropertyOwnerCommandHandler{repo: repo, eventBus: eventBus}
}

// Create creates a new property owner within the tenant scope.
func (h *PropertyOwnerCommandHandler) Create(ctx context.Context, orgID, appID uuid.UUID, req dto.CreatePropertyOwnerRequest) (*dto.PropertyOwnerResponse, error) {
	ratio := 1.0
	if req.OwnershipRatio != nil {
		ratio = *req.OwnershipRatio
	}

	owner := &model.PropertyOwner{
		OrganizationID:   orgID,
		AppID:            appID,
		UserID:           req.UserID,
		UnitID:           req.UnitID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		IdentityNumber:   req.IdentityNumber,
		Phone:            req.Phone,
		Email:            req.Email,
		Address:          req.Address,
		IBAN:             req.IBAN,
		OwnershipRatio:   ratio,
		IsPrimaryContact: req.IsPrimaryContact,
		Status:           model.OwnerStatusActive,
	}

	if err := h.repo.Create(ctx, owner); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypePropertyOwner, event.ActionCreated, owner.ID, orgID, appID))

	resp := dto.ToPropertyOwnerResponse(owner)
	return &resp, nil
}

// Update applies a partial update to a property owner.
func (h *PropertyOwnerCommandHandler) Update(ctx context.Context, orgID, appID, id uuid.UUID, req dto.UpdatePropertyOwnerRequest) (*dto.PropertyOwnerResponse, error) {
	owner, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		owner.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		owner.LastName = *req.LastName
	}
	if req.IdentityNumber != nil {
		owner.IdentityNumber = *req.IdentityNumber
	}
	if req.Phone != nil {
		owner.Phone = *req.Phone
	}
	if req.Email != nil {
		owner.Email = *req.Email
	}
	if req.Address != nil {
		owner.Address = *req.Address
	}
	if req.IBAN != nil {
		owner.IBAN = *req.IBAN
	}
	if req.OwnershipRatio != nil {
		owner.OwnershipRatio = *req.OwnershipRatio
	}
	if req.IsPrimaryContact != nil {
		owner.IsPrimaryContact = *req.IsPrimaryContact
	}
	if req.Status != nil {
		st := model.OwnerStatus(*req.Status)
		if !model.IsValidOwnerStatus(st) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid owner status", nil)
		}
		owner.Status = st
	}

	if err := h.repo.Update(ctx, owner); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypePropertyOwner, event.ActionUpdated, owner.ID, orgID, appID))

	resp := dto.ToPropertyOwnerResponse(owner)
	return &resp, nil
}

// Delete removes a property owner within the tenant scope.
func (h *PropertyOwnerCommandHandler) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}
	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}
	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypePropertyOwner, event.ActionDeleted, id, orgID, appID))
	return nil
}
