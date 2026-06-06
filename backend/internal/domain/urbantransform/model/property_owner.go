package model

import (
	"time"

	"github.com/google/uuid"
)

// OwnerStatus represents the status of a property owner in the transformation flow.
type OwnerStatus string

const (
	OwnerStatusActive        OwnerStatus = "active"
	OwnerStatusObjectionFiled OwnerStatus = "objection_filed"
	OwnerStatusConsentGiven  OwnerStatus = "consent_given"
	OwnerStatusArchived      OwnerStatus = "archived"
)

// IsValidOwnerStatus reports whether the given status is known.
func IsValidOwnerStatus(s OwnerStatus) bool {
	switch s {
	case OwnerStatusActive, OwnerStatusObjectionFiled, OwnerStatusConsentGiven, OwnerStatusArchived:
		return true
	default:
		return false
	}
}

// PropertyOwner represents an owner (hak sahibi) of a building unit.
type PropertyOwner struct {
	ID               uuid.UUID   `json:"id"`
	OrganizationID   uuid.UUID   `json:"organization_id"`
	AppID            uuid.UUID   `json:"app_id"`
	UserID           *uuid.UUID  `json:"user_id,omitempty"`
	UnitID           uuid.UUID   `json:"unit_id"`
	FirstName        string      `json:"first_name"`
	LastName         string      `json:"last_name"`
	IdentityNumber   string      `json:"identity_number"`
	Phone            string      `json:"phone"`
	Email            string      `json:"email"`
	Address          string      `json:"address"`
	IBAN             string      `json:"iban"`
	OwnershipRatio   float64     `json:"ownership_ratio"`
	IsPrimaryContact bool        `json:"is_primary_contact"`
	Status           OwnerStatus `json:"status"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}
