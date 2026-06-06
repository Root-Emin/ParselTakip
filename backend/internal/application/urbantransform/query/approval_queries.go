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

// ApprovalQueryHandler handles read operations for approvals (CQRS read side).
type ApprovalQueryHandler struct {
	repo repository.ApprovalRepository
}

// NewApprovalQueryHandler creates a new ApprovalQueryHandler.
func NewApprovalQueryHandler(repo repository.ApprovalRepository) *ApprovalQueryHandler {
	return &ApprovalQueryHandler{repo: repo}
}

// Get returns a single approval by ID within the tenant scope.
func (h *ApprovalQueryHandler) Get(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.ApprovalResponse, error) {
	approval, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToApprovalResponse(approval)
	return &resp, nil
}

// List returns a paginated, filtered list of approvals.
func (h *ApprovalQueryHandler) List(ctx context.Context, orgID, appID uuid.UUID, q dto.ListApprovalsQuery) (pagination.Result[dto.ApprovalResponse], error) {
	params := normalizePage(q.Page, q.PerPage)

	sortBy := q.SortBy
	if !constants.IsAllowedApprovalSortKey(sortBy) {
		sortBy = constants.ApprovalSortKeyCreatedAt
	}

	filter := repository.ApprovalFilter{
		OrganizationID: orgID,
		AppID:          appID,
		ProjectID:      q.ProjectID,
		ApproverID:     q.ApproverID,
		OwnerID:        q.OwnerID,
		SortBy:         sortBy,
		SortOrder:      constants.NormalizeSortOrder(q.SortOrder),
		Offset:         params.Offset(),
		Limit:          params.Limit(),
	}
	if q.ApprovalType != nil {
		t := model.ApprovalType(*q.ApprovalType)
		filter.ApprovalType = &t
	}
	if q.Status != nil {
		s := model.ApprovalStatus(*q.Status)
		filter.Status = &s
	}

	items, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return pagination.Result[dto.ApprovalResponse]{}, err
	}
	return pagination.NewResult(dto.ToApprovalResponseList(items), params, total), nil
}
