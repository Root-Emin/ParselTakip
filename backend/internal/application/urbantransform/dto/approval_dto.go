package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateApprovalRequest is the input for creating an approval request (command).
type CreateApprovalRequest struct {
	ProjectID    uuid.UUID  `json:"project_id" validate:"required"`
	ApprovalType string     `json:"approval_type" validate:"required,oneof=municipal_initiation owner_consent majority_decision municipal_permit demolition occupancy title_transfer rent_assistance contractor_assignment"`
	ApproverRole string     `json:"approver_role" validate:"required"`
	OwnerID      *uuid.UUID `json:"owner_id,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

// DecideApprovalRequest is the input for deciding an approval (approve/reject).
type DecideApprovalRequest struct {
	Status        string `json:"status" validate:"required,oneof=approved rejected"`
	DecisionNotes string `json:"decision_notes"`
}

// ListApprovalsQuery is the input for listing/filtering approvals (query).
type ListApprovalsQuery struct {
	ProjectID    *uuid.UUID
	ApprovalType *string
	Status       *string
	ApproverID   *uuid.UUID
	OwnerID      *uuid.UUID
	SortBy       string
	SortOrder    string
	Page         int
	PerPage      int
}

// ApprovalResponse is the public representation of an approval (response model).
type ApprovalResponse struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	AppID          uuid.UUID  `json:"app_id"`
	ProjectID      uuid.UUID  `json:"project_id"`
	ApprovalType   string     `json:"approval_type"`
	ApproverID     *uuid.UUID `json:"approver_id,omitempty"`
	ApproverRole   string     `json:"approver_role"`
	OwnerID        *uuid.UUID `json:"owner_id,omitempty"`
	Status         string     `json:"status"`
	DecisionNotes  string     `json:"decision_notes"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	DecidedAt      *time.Time `json:"decided_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ToApprovalResponse maps a domain approval to its response model.
func ToApprovalResponse(a *model.Approval) ApprovalResponse {
	return ApprovalResponse{
		ID:             a.ID,
		OrganizationID: a.OrganizationID,
		AppID:          a.AppID,
		ProjectID:      a.ProjectID,
		ApprovalType:   string(a.ApprovalType),
		ApproverID:     a.ApproverID,
		ApproverRole:   a.ApproverRole,
		OwnerID:        a.OwnerID,
		Status:         string(a.Status),
		DecisionNotes:  a.DecisionNotes,
		ExpiresAt:      a.ExpiresAt,
		DecidedAt:      a.DecidedAt,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

// ToApprovalResponseList maps a slice of approvals to response models.
func ToApprovalResponseList(items []*model.Approval) []ApprovalResponse {
	out := make([]ApprovalResponse, 0, len(items))
	for _, a := range items {
		out = append(out, ToApprovalResponse(a))
	}
	return out
}
