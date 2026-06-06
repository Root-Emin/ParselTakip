package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/tenant/dto"
	tenantEvent "github.com/masterfabric-go/masterfabric/internal/domain/tenant/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/tenant/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/tenant/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"
)

// ManageOrgUseCase handles organization update and deletion (system admin).
type ManageOrgUseCase struct {
	orgRepo  repository.OrgRepository
	eventBus events.EventBus
}

// NewManageOrgUseCase creates a new ManageOrgUseCase.
func NewManageOrgUseCase(orgRepo repository.OrgRepository, eventBus events.EventBus) *ManageOrgUseCase {
	return &ManageOrgUseCase{orgRepo: orgRepo, eventBus: eventBus}
}

// Update applies a partial update to an organization.
func (uc *ManageOrgUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateOrgRequest) (*dto.OrgInfo, error) {
	org, err := uc.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Status != "" {
		status := model.OrgStatus(req.Status)
		switch status {
		case model.OrgStatusActive, model.OrgStatusSuspended, model.OrgStatusArchived:
			org.Status = status
		default:
			return nil, domainErr.New(domainErr.ErrValidation, "invalid organization status", nil)
		}
	}

	if err := uc.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	updatedBy, _ := middleware.UserIDFromContext(ctx)
	_ = uc.eventBus.Publish(ctx, events.TopicTenant, tenantEvent.OrganizationUpdated{
		OrganizationID: org.ID,
		UpdatedBy:      updatedBy,
		Timestamp:      time.Now().UTC(),
	})

	return &dto.OrgInfo{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		Status:    string(org.Status),
		CreatedAt: org.CreatedAt,
	}, nil
}

// Delete removes an organization.
func (uc *ManageOrgUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.orgRepo.GetByID(ctx, id); err != nil {
		return err
	}
	if err := uc.orgRepo.Delete(ctx, id); err != nil {
		return err
	}

	deletedBy, _ := middleware.UserIDFromContext(ctx)
	_ = uc.eventBus.Publish(ctx, events.TopicTenant, tenantEvent.OrganizationDeleted{
		OrganizationID: id,
		DeletedBy:      deletedBy,
		Timestamp:      time.Now().UTC(),
	})
	return nil
}
