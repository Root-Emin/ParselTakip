package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// BuildingUnitFilter holds optional filter criteria for querying building units.
type BuildingUnitFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	BuildingID     *uuid.UUID
	Status         *model.UnitStatus
	OwnershipType  *model.OwnershipType
	Search         string // free-text on unit_no/title_deed_no
	SortBy         string
	SortOrder      string
	Offset         int
	Limit          int
}

// BuildingUnitRepository defines persistence operations for building units.
type BuildingUnitRepository interface {
	Create(ctx context.Context, unit *model.BuildingUnit) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.BuildingUnit, error)
	Update(ctx context.Context, unit *model.BuildingUnit) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter BuildingUnitFilter) ([]*model.BuildingUnit, int, error)
}
