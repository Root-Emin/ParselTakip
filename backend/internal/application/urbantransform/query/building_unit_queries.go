package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/pagination"
)

// BuildingUnitQueryHandler handles read operations for building units (CQRS read side).
type BuildingUnitQueryHandler struct {
	repo repository.BuildingUnitRepository
}

// NewBuildingUnitQueryHandler creates a new BuildingUnitQueryHandler.
func NewBuildingUnitQueryHandler(repo repository.BuildingUnitRepository) *BuildingUnitQueryHandler {
	return &BuildingUnitQueryHandler{repo: repo}
}

// Get returns a single building unit by ID within the tenant scope.
func (h *BuildingUnitQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.BuildingUnitResponse, error) {
	unit, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToBuildingUnitResponse(unit)
	return &resp, nil
}

// List returns a paginated, filtered and searchable list of building units.
func (h *BuildingUnitQueryHandler) List(ctx context.Context, orgID, appID uuid.UUID, q dto.ListBuildingUnitsQuery) (pagination.Result[dto.BuildingUnitResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	sortBy := q.SortBy
	if !constants.IsAllowedUnitSortKey(sortBy) {
		sortBy = constants.UnitSortKeyCreatedAt
	}

	filter := repository.BuildingUnitFilter{
		OrganizationID: orgID,
		AppID:          appID,
		BuildingID:     q.BuildingID,
		Search:         q.Search,
		SortBy:         sortBy,
		SortOrder:      constants.NormalizeSortOrder(q.SortOrder),
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}
	if q.Status != nil {
		st := model.UnitStatus(*q.Status)
		filter.Status = &st
	}
	if q.OwnershipType != nil {
		ot := model.OwnershipType(*q.OwnershipType)
		filter.OwnershipType = &ot
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.BuildingUnitResponse]{}, err
	}
	return pagination.NewResult(dto.ToBuildingUnitResponseList(items), params, total), nil
}
