package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// BuildingFilter holds optional filter criteria for querying buildings.
type BuildingFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	Status         *model.BuildingStatus
	RiskStatus     *model.RiskStatus
	BuildingType   *model.BuildingType
	City           string
	District       string
	Search         string // free-text on name/address/block_no/parcel_no
	SortBy         string
	SortOrder      string
	Offset         int
	Limit          int
}

// BuildingRepository defines persistence operations for buildings.
type BuildingRepository interface {
	Create(ctx context.Context, building *model.Building) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Building, error)
	Update(ctx context.Context, building *model.Building) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter BuildingFilter) ([]*model.Building, int, error)
}
