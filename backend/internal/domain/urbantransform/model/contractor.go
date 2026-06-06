package model

import (
	"time"

	"github.com/google/uuid"
)

// ContractorStatus represents the status of a contractor company.
type ContractorStatus string

const (
	ContractorStatusActive      ContractorStatus = "active"
	ContractorStatusSuspended   ContractorStatus = "suspended"
	ContractorStatusBlacklisted ContractorStatus = "blacklisted"
)

// IsValidContractorStatus reports whether the given status is known.
func IsValidContractorStatus(s ContractorStatus) bool {
	switch s {
	case ContractorStatusActive, ContractorStatusSuspended, ContractorStatusBlacklisted:
		return true
	default:
		return false
	}
}

// Contractor is a contractor company participating in transformation projects.
type Contractor struct {
	ID               uuid.UUID        `json:"id"`
	OrganizationID   uuid.UUID        `json:"organization_id"`
	AppID            uuid.UUID        `json:"app_id"`
	UserID           *uuid.UUID       `json:"user_id,omitempty"`
	CompanyName      string           `json:"company_name"`
	TaxNumber        string           `json:"tax_number"`
	TradeRegistryNo  string           `json:"trade_registry_no"`
	AuthorizedPerson string           `json:"authorized_person"`
	Phone            string           `json:"phone"`
	Email            string           `json:"email"`
	Address          string           `json:"address"`
	Status           ContractorStatus `json:"status"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}
