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

// ListProjectsHandler handles list/filter/search queries (CQRS read side).
type ListProjectsHandler struct {
	repo repository.ProjectRepository
}

// NewListProjectsHandler creates a new ListProjectsHandler.
func NewListProjectsHandler(repo repository.ProjectRepository) *ListProjectsHandler {
	return &ListProjectsHandler{repo: repo}
}

// Execute returns a paginated, filtered and searchable list of projects.
func (h *ListProjectsHandler) Execute(ctx context.Context, orgID, appID uuid.UUID, q dto.ListProjectsQuery) (pagination.Result[dto.ProjectResponse], error) {
	page := q.Page
	if page < 1 {
		page = pagination.DefaultPage
	}
	perPage := q.PerPage
	if perPage < 1 {
		perPage = pagination.DefaultPerPage
	}
	if perPage > pagination.MaxPerPage {
		perPage = pagination.MaxPerPage
	}
	params := pagination.Params{Page: page, PerPage: perPage}

	sortBy := q.SortBy
	if !constants.IsAllowedSortKey(sortBy) {
		sortBy = constants.DefaultSortBy
	}
	sortOrder := constants.NormalizeSortOrder(q.SortOrder)

	filter := repository.ProjectFilter{
		OrganizationID:        orgID,
		AppID:                 appID,
		AssignedContractorID:  q.AssignedContractorID,
		InitiatedBy:           q.InitiatedBy,
		CurrentWorkflowStepID: q.CurrentWorkflowStepID,
		Search:                q.Search,
		SortBy:                sortBy,
		SortOrder:             sortOrder,
		Offset:                params.Offset(),
		Limit:                 params.Limit(),
	}
	if q.Status != nil {
		status := model.ProjectStatus(*q.Status)
		filter.Status = &status
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.ProjectResponse]{}, err
	}

	return pagination.NewResult(dto.ToProjectResponseList(items), params, total), nil
}
