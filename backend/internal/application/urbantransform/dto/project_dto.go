package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateProjectRequest is the input for creating a project (command).
type CreateProjectRequest struct {
	Code                 string     `json:"code" validate:"required,min=2,max=50"`
	Name                 string     `json:"name" validate:"required,min=2,max=255"`
	Description          string     `json:"description"`
	AssignedContractorID *uuid.UUID `json:"assigned_contractor_id,omitempty"`
	TargetCompletionAt   *time.Time `json:"target_completion_at,omitempty"`
}

// UpdateProjectRequest is the input for updating a project (command).
// Pointer fields allow partial updates: nil means "leave unchanged".
type UpdateProjectRequest struct {
	Name                 *string    `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description          *string    `json:"description,omitempty"`
	Status               *string    `json:"status,omitempty" validate:"omitempty,oneof=draft initiated in_progress blocked completed cancelled"`
	AssignedContractorID *uuid.UUID `json:"assigned_contractor_id,omitempty"`
	TargetCompletionAt   *time.Time `json:"target_completion_at,omitempty"`
}

// ListProjectsQuery is the input for listing/filtering/searching projects (query).
type ListProjectsQuery struct {
	Status               *string
	AssignedContractorID *uuid.UUID
	InitiatedBy          *uuid.UUID
	CurrentWorkflowStepID *uuid.UUID
	Search               string
	SortBy               string
	SortOrder            string
	Page                 int
	PerPage              int
}

// ProjectResponse is the public representation of a project (response model).
type ProjectResponse struct {
	ID                    uuid.UUID  `json:"id"`
	OrganizationID        uuid.UUID  `json:"organization_id"`
	AppID                 uuid.UUID  `json:"app_id"`
	Code                  string     `json:"code"`
	Name                  string     `json:"name"`
	Description           string     `json:"description"`
	Status                string     `json:"status"`
	CurrentWorkflowStepID *uuid.UUID `json:"current_workflow_step_id,omitempty"`
	InitiatedBy           *uuid.UUID `json:"initiated_by,omitempty"`
	AssignedContractorID  *uuid.UUID `json:"assigned_contractor_id,omitempty"`
	StartedAt             *time.Time `json:"started_at,omitempty"`
	TargetCompletionAt    *time.Time `json:"target_completion_at,omitempty"`
	CompletedAt           *time.Time `json:"completed_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
