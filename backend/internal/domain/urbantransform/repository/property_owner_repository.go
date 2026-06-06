package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// PropertyOwnerFilter holds optional filter criteria for querying property owners.
type PropertyOwnerFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	UnitID         *uuid.UUID
	Status         *model.OwnerStatus
	Search         string // free-text on first_name/last_name/identity_number/email
	SortBy         string
	SortOrder      string
	Offset         int
	Limit          int
}

// PropertyOwnerRepository defines persistence operations for property owners.
type PropertyOwnerRepository interface {
	Create(ctx context.Context, owner *model.PropertyOwner) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.PropertyOwner, error)
	Update(ctx context.Context, owner *model.PropertyOwner) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter PropertyOwnerFilter) ([]*model.PropertyOwner, int, error)
}
