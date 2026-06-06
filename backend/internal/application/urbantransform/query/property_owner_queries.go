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

// PropertyOwnerQueryHandler handles read operations for property owners (CQRS read side).
type PropertyOwnerQueryHandler struct {
	repo repository.PropertyOwnerRepository
}

// NewPropertyOwnerQueryHandler creates a new PropertyOwnerQueryHandler.
func NewPropertyOwnerQueryHandler(repo repository.PropertyOwnerRepository) *PropertyOwnerQueryHandler {
	return &PropertyOwnerQueryHandler{repo: repo}
}

// Get returns a single property owner by ID within the tenant scope.
func (h *PropertyOwnerQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.PropertyOwnerResponse, error) {
	owner, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToPropertyOwnerResponse(owner)
	return &resp, nil
}

// List returns a paginated, filtered and searchable list of property owners.
func (h *PropertyOwnerQueryHandler) List(ctx context.Context, orgID, appID uuid.UUID, q dto.ListPropertyOwnersQuery) (pagination.Result[dto.PropertyOwnerResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	sortBy := q.SortBy
	if !constants.IsAllowedOwnerSortKey(sortBy) {
		sortBy = constants.OwnerSortKeyCreatedAt
	}

	filter := repository.PropertyOwnerFilter{
		OrganizationID: orgID,
		AppID:          appID,
		UnitID:         q.UnitID,
		Search:         q.Search,
		SortBy:         sortBy,
		SortOrder:      constants.NormalizeSortOrder(q.SortOrder),
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}
	if q.Status != nil {
		st := model.OwnerStatus(*q.Status)
		filter.Status = &st
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.PropertyOwnerResponse]{}, err
	}
	return pagination.NewResult(dto.ToPropertyOwnerResponseList(items), params, total), nil
}
