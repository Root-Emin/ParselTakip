package model

import (
	"time"

	"github.com/google/uuid"
)

// ProjectStatus represents the lifecycle status of an urban transformation project.
type ProjectStatus string

const (
	ProjectStatusDraft      ProjectStatus = "draft"
	ProjectStatusInitiated  ProjectStatus = "initiated"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusBlocked    ProjectStatus = "blocked"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusCancelled  ProjectStatus = "cancelled"
)

// IsValidProjectStatus reports whether the given status is a known status.
func IsValidProjectStatus(s ProjectStatus) bool {
	switch s {
	case ProjectStatusDraft, ProjectStatusInitiated, ProjectStatusInProgress,
		ProjectStatusBlocked, ProjectStatusCompleted, ProjectStatusCancelled:
		return true
	default:
		return false
	}
}

// Project is the aggregate root for an urban transformation project.
type Project struct {
	ID                    uuid.UUID     `json:"id"`
	OrganizationID        uuid.UUID     `json:"organization_id"`
	AppID                 uuid.UUID     `json:"app_id"`
	Code                  string        `json:"code"`
	Name                  string        `json:"name"`
	Description           string        `json:"description"`
	Status                ProjectStatus `json:"status"`
	CurrentWorkflowStepID *uuid.UUID    `json:"current_workflow_step_id,omitempty"`
	InitiatedBy           *uuid.UUID    `json:"initiated_by,omitempty"`
	AssignedContractorID  *uuid.UUID    `json:"assigned_contractor_id,omitempty"`
	StartedAt             *time.Time    `json:"started_at,omitempty"`
	TargetCompletionAt    *time.Time    `json:"target_completion_at,omitempty"`
	CompletedAt           *time.Time    `json:"completed_at,omitempty"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

// IsActive reports whether the project is in an active (non-terminal) state.
func (p *Project) IsActive() bool {
	return p.Status != ProjectStatusCompleted && p.Status != ProjectStatusCancelled
}
