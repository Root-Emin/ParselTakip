package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreatePropertyOwnerRequest is the input for creating a property owner (command).
type CreatePropertyOwnerRequest struct {
	UnitID           uuid.UUID  `json:"unit_id" validate:"required"`
	UserID           *uuid.UUID `json:"user_id,omitempty"`
	FirstName        string     `json:"first_name" validate:"required"`
	LastName         string     `json:"last_name" validate:"required"`
	IdentityNumber   string     `json:"identity_number" validate:"omitempty,len=11"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email" validate:"omitempty,email"`
	Address          string     `json:"address"`
	IBAN             string     `json:"iban"`
	OwnershipRatio   *float64   `json:"ownership_ratio,omitempty" validate:"omitempty,gt=0,lte=1"`
	IsPrimaryContact bool       `json:"is_primary_contact"`
}

// UpdatePropertyOwnerRequest is the input for updating a property owner (partial command).
type UpdatePropertyOwnerRequest struct {
	FirstName        *string  `json:"first_name,omitempty"`
	LastName         *string  `json:"last_name,omitempty"`
	IdentityNumber   *string  `json:"identity_number,omitempty" validate:"omitempty,len=11"`
	Phone            *string  `json:"phone,omitempty"`
	Email            *string  `json:"email,omitempty" validate:"omitempty,email"`
	Address          *string  `json:"address,omitempty"`
	IBAN             *string  `json:"iban,omitempty"`
	OwnershipRatio   *float64 `json:"ownership_ratio,omitempty" validate:"omitempty,gt=0,lte=1"`
	IsPrimaryContact *bool    `json:"is_primary_contact,omitempty"`
	Status           *string  `json:"status,omitempty" validate:"omitempty,oneof=active objection_filed consent_given archived"`
}

// ListPropertyOwnersQuery is the input for listing/filtering/searching owners (query).
type ListPropertyOwnersQuery struct {
	UnitID    *uuid.UUID
	Status    *string
	Search    string
	SortBy    string
	SortOrder string
	Page      int
	PerPage   int
}

// PropertyOwnerResponse is the public representation of a property owner (response model).
type PropertyOwnerResponse struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	AppID            uuid.UUID  `json:"app_id"`
	UserID           *uuid.UUID `json:"user_id,omitempty"`
	UnitID           uuid.UUID  `json:"unit_id"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	IdentityNumber   string     `json:"identity_number"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email"`
	Address          string     `json:"address"`
	IBAN             string     `json:"iban"`
	OwnershipRatio   float64    `json:"ownership_ratio"`
	IsPrimaryContact bool       `json:"is_primary_contact"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ToPropertyOwnerResponse maps a domain owner to its response model.
func ToPropertyOwnerResponse(o *model.PropertyOwner) PropertyOwnerResponse {
	return PropertyOwnerResponse{
		ID:               o.ID,
		OrganizationID:   o.OrganizationID,
		AppID:            o.AppID,
		UserID:           o.UserID,
		UnitID:           o.UnitID,
		FirstName:        o.FirstName,
		LastName:         o.LastName,
		IdentityNumber:   o.IdentityNumber,
		Phone:            o.Phone,
		Email:            o.Email,
		Address:          o.Address,
		IBAN:             o.IBAN,
		OwnershipRatio:   o.OwnershipRatio,
		IsPrimaryContact: o.IsPrimaryContact,
		Status:           string(o.Status),
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
	}
}

// ToPropertyOwnerResponseList maps a slice of owners to response models.
func ToPropertyOwnerResponseList(items []*model.PropertyOwner) []PropertyOwnerResponse {
	out := make([]PropertyOwnerResponse, 0, len(items))
	for _, o := range items {
		out = append(out, ToPropertyOwnerResponse(o))
	}
	return out
}
