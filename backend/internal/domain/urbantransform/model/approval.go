package model

import (
	"time"

	"github.com/google/uuid"
)

// ApprovalType represents the kind of approval in the transformation flow.
type ApprovalType string

const (
	ApprovalTypeMunicipalInitiation  ApprovalType = "municipal_initiation"
	ApprovalTypeOwnerConsent         ApprovalType = "owner_consent"
	ApprovalTypeMajorityDecision     ApprovalType = "majority_decision"
	ApprovalTypeMunicipalPermit      ApprovalType = "municipal_permit"
	ApprovalTypeDemolition           ApprovalType = "demolition"
	ApprovalTypeOccupancy            ApprovalType = "occupancy"
	ApprovalTypeTitleTransfer        ApprovalType = "title_transfer"
	ApprovalTypeRentAssistance       ApprovalType = "rent_assistance"
	ApprovalTypeContractorAssignment ApprovalType = "contractor_assignment"
)

// ApprovalStatus represents the decision state of an approval.
type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
	ApprovalStatusExpired  ApprovalStatus = "expired"
)

// IsValidApprovalType reports whether the given approval type is known.
func IsValidApprovalType(t ApprovalType) bool {
	switch t {
	case ApprovalTypeMunicipalInitiation, ApprovalTypeOwnerConsent, ApprovalTypeMajorityDecision,
		ApprovalTypeMunicipalPermit, ApprovalTypeDemolition, ApprovalTypeOccupancy,
		ApprovalTypeTitleTransfer, ApprovalTypeRentAssistance, ApprovalTypeContractorAssignment:
		return true
	default:
		return false
	}
}

// IsValidApprovalStatus reports whether the given approval status is known.
func IsValidApprovalStatus(s ApprovalStatus) bool {
	switch s {
	case ApprovalStatusPending, ApprovalStatusApproved, ApprovalStatusRejected, ApprovalStatusExpired:
		return true
	default:
		return false
	}
}

// Approval represents a decision gate (e.g. owner consent, municipal permit).
type Approval struct {
	ID             uuid.UUID      `json:"id"`
	OrganizationID uuid.UUID      `json:"organization_id"`
	AppID          uuid.UUID      `json:"app_id"`
	ProjectID      uuid.UUID      `json:"project_id"`
	ApprovalType   ApprovalType   `json:"approval_type"`
	ApproverID     *uuid.UUID     `json:"approver_id,omitempty"`
	ApproverRole   string         `json:"approver_role"`
	OwnerID        *uuid.UUID     `json:"owner_id,omitempty"`
	Status         ApprovalStatus `json:"status"`
	DecisionNotes  string         `json:"decision_notes"`
	ExpiresAt      *time.Time     `json:"expires_at,omitempty"`
	DecidedAt      *time.Time     `json:"decided_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
