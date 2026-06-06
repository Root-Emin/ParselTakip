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

// BuildingUnitCommandHandler handles write operations for building units (CQRS write side).
type BuildingUnitCommandHandler struct {
	repo     repository.BuildingUnitRepository
	eventBus events.EventBus
}

// NewBuildingUnitCommandHandler creates a new BuildingUnitCommandHandler.
func NewBuildingUnitCommandHandler(repo repository.BuildingUnitRepository, eventBus events.EventBus) *BuildingUnitCommandHandler {
	return &BuildingUnitCommandHandler{repo: repo, eventBus: eventBus}
}

// Create creates a new building unit within the tenant scope.
func (h *BuildingUnitCommandHandler) Create(ctx context.Context, orgID, appID uuid.UUID, req dto.CreateBuildingUnitRequest) (*dto.BuildingUnitResponse, error) {
	ownershipType := model.OwnershipType(req.OwnershipType)
	if ownershipType == "" {
		ownershipType = model.OwnershipTypeKatMulkiyeti
	}

	unit := &model.BuildingUnit{
		OrganizationID: orgID,
		AppID:          appID,
		BuildingID:     req.BuildingID,
		UnitNo:         req.UnitNo,
		FloorNo:        req.FloorNo,
		AreaSqm:        req.AreaSqm,
		RoomCount:      req.RoomCount,
		OwnershipType:  ownershipType,
		TitleDeedNo:    req.TitleDeedNo,
		Status:         model.UnitStatusActive,
	}

	if err := h.repo.Create(ctx, unit); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeBuildingUnit, event.ActionCreated, unit.ID, orgID, appID))

	resp := dto.ToBuildingUnitResponse(unit)
	return &resp, nil
}

// Update applies a partial update to a building unit.
func (h *BuildingUnitCommandHandler) Update(ctx context.Context, orgID, appID, id uuid.UUID, req dto.UpdateBuildingUnitRequest) (*dto.BuildingUnitResponse, error) {
	unit, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}

	if req.UnitNo != nil {
		unit.UnitNo = *req.UnitNo
	}
	if req.FloorNo != nil {
		unit.FloorNo = req.FloorNo
	}
	if req.AreaSqm != nil {
		unit.AreaSqm = req.AreaSqm
	}
	if req.RoomCount != nil {
		unit.RoomCount = *req.RoomCount
	}
	if req.OwnershipType != nil {
		ot := model.OwnershipType(*req.OwnershipType)
		if !model.IsValidOwnershipType(ot) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid ownership type", nil)
		}
		unit.OwnershipType = ot
	}
	if req.TitleDeedNo != nil {
		unit.TitleDeedNo = *req.TitleDeedNo
	}
	if req.Status != nil {
		st := model.UnitStatus(*req.Status)
		if !model.IsValidUnitStatus(st) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid unit status", nil)
		}
		unit.Status = st
	}

	if err := h.repo.Update(ctx, unit); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeBuildingUnit, event.ActionUpdated, unit.ID, orgID, appID))

	resp := dto.ToBuildingUnitResponse(unit)
	return &resp, nil
}

// Delete removes a building unit within the tenant scope.
func (h *BuildingUnitCommandHandler) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}
	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}
	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeBuildingUnit, event.ActionDeleted, id, orgID, appID))
	return nil
}
