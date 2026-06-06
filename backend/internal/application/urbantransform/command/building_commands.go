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

// BuildingCommandHandler handles write operations for buildings (CQRS write side).
type BuildingCommandHandler struct {
	repo     repository.BuildingRepository
	eventBus events.EventBus
}

// NewBuildingCommandHandler creates a new BuildingCommandHandler.
func NewBuildingCommandHandler(repo repository.BuildingRepository, eventBus events.EventBus) *BuildingCommandHandler {
	return &BuildingCommandHandler{repo: repo, eventBus: eventBus}
}

// Create creates a new building within the tenant scope.
func (h *BuildingCommandHandler) Create(ctx context.Context, orgID, appID, createdBy uuid.UUID, req dto.CreateBuildingRequest) (*dto.BuildingResponse, error) {
	unitCount := req.UnitCount
	if unitCount < 1 {
		unitCount = 1
	}
	buildingType := model.BuildingType(req.BuildingType)
	if buildingType == "" {
		buildingType = model.BuildingTypeResidential
	}
	riskStatus := model.RiskStatus(req.RiskStatus)
	if riskStatus == "" {
		riskStatus = model.RiskStatusUnknown
	}

	building := &model.Building{
		OrganizationID:   orgID,
		AppID:            appID,
		Name:             req.Name,
		Address:          req.Address,
		City:             req.City,
		District:         req.District,
		Neighborhood:     req.Neighborhood,
		BlockNo:          req.BlockNo,
		ParcelNo:         req.ParcelNo,
		IslandNo:         req.IslandNo,
		FloorCount:       req.FloorCount,
		UnitCount:        unitCount,
		ConstructionYear: req.ConstructionYear,
		BuildingType:     buildingType,
		RiskStatus:       riskStatus,
		Status:           model.BuildingStatusActive,
	}
	if createdBy != uuid.Nil {
		building.CreatedBy = &createdBy
	}

	if err := h.repo.Create(ctx, building); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeBuilding, event.ActionCreated, building.ID, orgID, appID))

	resp := dto.ToBuildingResponse(building)
	return &resp, nil
}

// Update applies a partial update to a building.
func (h *BuildingCommandHandler) Update(ctx context.Context, orgID, appID, id uuid.UUID, req dto.UpdateBuildingRequest) (*dto.BuildingResponse, error) {
	building, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		building.Name = *req.Name
	}
	if req.Address != nil {
		building.Address = *req.Address
	}
	if req.City != nil {
		building.City = *req.City
	}
	if req.District != nil {
		building.District = *req.District
	}
	if req.Neighborhood != nil {
		building.Neighborhood = *req.Neighborhood
	}
	if req.BlockNo != nil {
		building.BlockNo = *req.BlockNo
	}
	if req.ParcelNo != nil {
		building.ParcelNo = *req.ParcelNo
	}
	if req.IslandNo != nil {
		building.IslandNo = *req.IslandNo
	}
	if req.FloorCount != nil {
		building.FloorCount = req.FloorCount
	}
	if req.UnitCount != nil {
		building.UnitCount = *req.UnitCount
	}
	if req.ConstructionYear != nil {
		building.ConstructionYear = req.ConstructionYear
	}
	if req.BuildingType != nil {
		bt := model.BuildingType(*req.BuildingType)
		if !model.IsValidBuildingType(bt) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid building type", nil)
		}
		building.BuildingType = bt
	}
	if req.RiskStatus != nil {
		rs := model.RiskStatus(*req.RiskStatus)
		if !model.IsValidRiskStatus(rs) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid risk status", nil)
		}
		building.RiskStatus = rs
	}
	if req.Status != nil {
		st := model.BuildingStatus(*req.Status)
		if !model.IsValidBuildingStatus(st) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid building status", nil)
		}
		building.Status = st
	}

	if err := h.repo.Update(ctx, building); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeBuilding, event.ActionUpdated, building.ID, orgID, appID))

	resp := dto.ToBuildingResponse(building)
	return &resp, nil
}

// Delete removes a building within the tenant scope.
func (h *BuildingCommandHandler) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}
	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}
	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeBuilding, event.ActionDeleted, id, orgID, appID))
	return nil
}
