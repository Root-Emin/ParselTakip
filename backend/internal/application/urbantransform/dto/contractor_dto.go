package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateContractorRequest is the input for creating a contractor (command).
type CreateContractorRequest struct {
	CompanyName      string     `json:"company_name" validate:"required,min=2,max=255"`
	TaxNumber        string     `json:"tax_number" validate:"required,min=10,max=20"`
	TradeRegistryNo  string     `json:"trade_registry_no"`
	AuthorizedPerson string     `json:"authorized_person"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email" validate:"omitempty,email"`
	Address          string     `json:"address"`
	UserID           *uuid.UUID `json:"user_id,omitempty"`
}

// UpdateContractorRequest is the input for updating a contractor (partial command).
type UpdateContractorRequest struct {
	CompanyName      *string `json:"company_name,omitempty" validate:"omitempty,min=2,max=255"`
	TradeRegistryNo  *string `json:"trade_registry_no,omitempty"`
	AuthorizedPerson *string `json:"authorized_person,omitempty"`
	Phone            *string `json:"phone,omitempty"`
	Email            *string `json:"email,omitempty" validate:"omitempty,email"`
	Address          *string `json:"address,omitempty"`
	Status           *string `json:"status,omitempty" validate:"omitempty,oneof=active suspended blacklisted"`
}

// ListContractorsQuery is the input for listing/filtering/searching contractors (query).
type ListContractorsQuery struct {
	Status    *string
	Search    string
	SortBy    string
	SortOrder string
	Page      int
	PerPage   int
}

// ContractorResponse is the public representation of a contractor (response model).
type ContractorResponse struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	AppID            uuid.UUID  `json:"app_id"`
	UserID           *uuid.UUID `json:"user_id,omitempty"`
	CompanyName      string     `json:"company_name"`
	TaxNumber        string     `json:"tax_number"`
	TradeRegistryNo  string     `json:"trade_registry_no"`
	AuthorizedPerson string     `json:"authorized_person"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email"`
	Address          string     `json:"address"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ToContractorResponse maps a domain contractor to its response model.
func ToContractorResponse(c *model.Contractor) ContractorResponse {
	return ContractorResponse{
		ID:               c.ID,
		OrganizationID:   c.OrganizationID,
		AppID:            c.AppID,
		UserID:           c.UserID,
		CompanyName:      c.CompanyName,
		TaxNumber:        c.TaxNumber,
		TradeRegistryNo:  c.TradeRegistryNo,
		AuthorizedPerson: c.AuthorizedPerson,
		Phone:            c.Phone,
		Email:            c.Email,
		Address:          c.Address,
		Status:           string(c.Status),
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

// ToContractorResponseList maps a slice of contractors to response models.
func ToContractorResponseList(items []*model.Contractor) []ContractorResponse {
	out := make([]ContractorResponse, 0, len(items))
	for _, c := range items {
		out = append(out, ToContractorResponse(c))
	}
	return out
}
