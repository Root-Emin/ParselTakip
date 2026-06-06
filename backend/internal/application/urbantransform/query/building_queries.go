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

// BuildingQueryHandler handles read operations for buildings (CQRS read side).
type BuildingQueryHandler struct {
	repo repository.BuildingRepository
}

// NewBuildingQueryHandler creates a new BuildingQueryHandler.
func NewBuildingQueryHandler(repo repository.BuildingRepository) *BuildingQueryHandler {
	return &BuildingQueryHandler{repo: repo}
}

// Get returns a single building by ID within the tenant scope.
func (h *BuildingQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.BuildingResponse, error) {
	building, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToBuildingResponse(building)
	return &resp, nil
}

// List returns a paginated, filtered and searchable list of buildings.
func (h *BuildingQueryHandler) List(ctx context.Context, orgID, appID uuid.UUID, q dto.ListBuildingsQuery) (pagination.Result[dto.BuildingResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	sortBy := q.SortBy
	if !constants.IsAllowedBuildingSortKey(sortBy) {
		sortBy = constants.BuildingSortKeyCreatedAt
	}

	filter := repository.BuildingFilter{
		OrganizationID: orgID,
		AppID:          appID,
		City:           q.City,
		District:       q.District,
		Search:         q.Search,
		SortBy:         sortBy,
		SortOrder:      constants.NormalizeSortOrder(q.SortOrder),
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}
	if q.Status != nil {
		st := model.BuildingStatus(*q.Status)
		filter.Status = &st
	}
	if q.RiskStatus != nil {
		rs := model.RiskStatus(*q.RiskStatus)
		filter.RiskStatus = &rs
	}
	if q.BuildingType != nil {
		bt := model.BuildingType(*q.BuildingType)
		filter.BuildingType = &bt
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.BuildingResponse]{}, err
	}
	return pagination.NewResult(dto.ToBuildingResponseList(items), params, total), nil
}
