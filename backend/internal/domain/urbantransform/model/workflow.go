package model

import (
	"time"

	"github.com/google/uuid"
)

// WorkflowStepDefinition is the global, ordered definition of a workflow step.
type WorkflowStepDefinition struct {
	ID              uuid.UUID `json:"id"`
	StepOrder       int       `json:"step_order"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ResponsibleRole string    `json:"responsible_role"`
	SLADays         *int      `json:"sla_days,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// WorkflowStateStatus represents the status of a project's step.
type WorkflowStateStatus string

const (
	WorkflowStatusPending           WorkflowStateStatus = "pending"
	WorkflowStatusInProgress        WorkflowStateStatus = "in_progress"
	WorkflowStatusAwaitingDocuments WorkflowStateStatus = "awaiting_documents"
	WorkflowStatusAwaitingApproval  WorkflowStateStatus = "awaiting_approval"
	WorkflowStatusCompleted         WorkflowStateStatus = "completed"
	WorkflowStatusBlocked           WorkflowStateStatus = "blocked"
	WorkflowStatusSkipped           WorkflowStateStatus = "skipped"
)

// IsValidWorkflowStatus reports whether the given status is known.
func IsValidWorkflowStatus(s WorkflowStateStatus) bool {
	switch s {
	case WorkflowStatusPending, WorkflowStatusInProgress, WorkflowStatusAwaitingDocuments,
		WorkflowStatusAwaitingApproval, WorkflowStatusCompleted, WorkflowStatusBlocked, WorkflowStatusSkipped:
		return true
	default:
		return false
	}
}

// ProjectWorkflowState is the current state of a step for a given project.
type ProjectWorkflowState struct {
	ID             uuid.UUID           `json:"id"`
	OrganizationID uuid.UUID           `json:"organization_id"`
	AppID          uuid.UUID           `json:"app_id"`
	ProjectID      uuid.UUID           `json:"project_id"`
	WorkflowStepID uuid.UUID           `json:"workflow_step_id"`
	Status         WorkflowStateStatus `json:"status"`
	StartedAt      *time.Time          `json:"started_at,omitempty"`
	CompletedAt    *time.Time          `json:"completed_at,omitempty"`
	DueAt          *time.Time          `json:"due_at,omitempty"`
	BlockedReason  string              `json:"blocked_reason"`
	UpdatedBy      *uuid.UUID          `json:"updated_by,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

// ProjectWorkflowHistory records a transition in a project's workflow.
type ProjectWorkflowHistory struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	AppID          uuid.UUID  `json:"app_id"`
	ProjectID      uuid.UUID  `json:"project_id"`
	FromStepID     *uuid.UUID `json:"from_step_id,omitempty"`
	ToStepID       uuid.UUID  `json:"to_step_id"`
	Action         string     `json:"action"`
	Notes          string     `json:"notes"`
	ChangedBy      *uuid.UUID `json:"changed_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
