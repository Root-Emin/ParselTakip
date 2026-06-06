package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// ContractorFilter holds optional filter criteria for querying contractors.
type ContractorFilter struct {
	OrganizationID uuid.UUID
	AppID          uuid.UUID
	Status         *model.ContractorStatus
	Search         string // free-text on company_name/tax_number/authorized_person/email
	SortBy         string
	SortOrder      string
	Offset         int
	Limit          int
}

// ContractorRepository defines persistence operations for contractors.
type ContractorRepository interface {
	Create(ctx context.Context, contractor *model.Contractor) error
	GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Contractor, error)
	GetByTaxNumber(ctx context.Context, orgID uuid.UUID, taxNumber string) (*model.Contractor, error)
	Update(ctx context.Context, contractor *model.Contractor) error
	Delete(ctx context.Context, orgID, appID, id uuid.UUID) error
	List(ctx context.Context, filter ContractorFilter) ([]*model.Contractor, int, error)
}
