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

// ContractorQueryHandler handles read operations for contractors (CQRS read side).
type ContractorQueryHandler struct {
	repo repository.ContractorRepository
}

// NewContractorQueryHandler creates a new ContractorQueryHandler.
func NewContractorQueryHandler(repo repository.ContractorRepository) *ContractorQueryHandler {
	return &ContractorQueryHandler{repo: repo}
}

// Get returns a single contractor by ID within the tenant scope.
func (h *ContractorQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.ContractorResponse, error) {
	contractor, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToContractorResponse(contractor)
	return &resp, nil
}

// List returns a paginated, filtered and searchable list of contractors.
func (h *ContractorQueryHandler) List(ctx context.Context, orgID, appID uuid.UUID, q dto.ListContractorsQuery) (pagination.Result[dto.ContractorResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	sortBy := q.SortBy
	if !constants.IsAllowedContractorSortKey(sortBy) {
		sortBy = constants.ContractorSortKeyCreatedAt
	}

	filter := repository.ContractorFilter{
		OrganizationID: orgID,
		AppID:          appID,
		Search:         q.Search,
		SortBy:         sortBy,
		SortOrder:      constants.NormalizeSortOrder(q.SortOrder),
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}
	if q.Status != nil {
		status := model.ContractorStatus(*q.Status)
		filter.Status = &status
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.ContractorResponse]{}, err
	}
	return pagination.NewResult(dto.ToContractorResponseList(items), params, total), nil
}
