package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// ApprovalFilter holds optional filter criteria for querying approvals.
type ApprovalFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	ProjectID      *uuid.UUID
	ApprovalType   *model.ApprovalType
	Status         *model.ApprovalStatus
	ApproverID     *uuid.UUID
	OwnerID        *uuid.UUID
	SortBy         string
	SortOrder      string
	Offset         int
	Limit          int
}

// ApprovalRepository defines persistence operations for approvals.
type ApprovalRepository interface {
	Create(ctx context.Context, approval *model.Approval) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Approval, error)
	Update(ctx context.Context, approval *model.Approval) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter ApprovalFilter) ([]*model.Approval, int, error)
}
