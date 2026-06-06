package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// AdvanceWorkflowRequest is the input for advancing/updating a project's step.
type AdvanceWorkflowRequest struct {
	WorkflowStepID uuid.UUID  `json:"workflow_step_id" validate:"required"`
	Status         string     `json:"status" validate:"required,oneof=pending in_progress awaiting_documents awaiting_approval completed blocked skipped"`
	BlockedReason  string     `json:"blocked_reason"`
	Notes          string     `json:"notes"`
	DueAt          *time.Time `json:"due_at,omitempty"`
}

// WorkflowStepResponse is the public representation of a step definition.
type WorkflowStepResponse struct {
	ID              uuid.UUID `json:"id"`
	StepOrder       int       `json:"step_order"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ResponsibleRole string    `json:"responsible_role"`
	SLADays         *int      `json:"sla_days,omitempty"`
}

// ProjectWorkflowStateResponse is the public representation of a project step state.
type ProjectWorkflowStateResponse struct {
	ID             uuid.UUID  `json:"id"`
	ProjectID      uuid.UUID  `json:"project_id"`
	WorkflowStepID uuid.UUID  `json:"workflow_step_id"`
	Status         string     `json:"status"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	DueAt          *time.Time `json:"due_at,omitempty"`
	BlockedReason  string     `json:"blocked_reason"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// WorkflowHistoryResponse is the public representation of a workflow transition.
type WorkflowHistoryResponse struct {
	ID         uuid.UUID  `json:"id"`
	ProjectID  uuid.UUID  `json:"project_id"`
	FromStepID *uuid.UUID `json:"from_step_id,omitempty"`
	ToStepID   uuid.UUID  `json:"to_step_id"`
	Action     string     `json:"action"`
	Notes      string     `json:"notes"`
	ChangedBy  *uuid.UUID `json:"changed_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ToWorkflowStepResponse maps a step definition to its response model.
func ToWorkflowStepResponse(s *model.WorkflowStepDefinition) WorkflowStepResponse {
	return WorkflowStepResponse{
		ID:              s.ID,
		StepOrder:       s.StepOrder,
		Code:            s.Code,
		Name:            s.Name,
		Description:     s.Description,
		ResponsibleRole: s.ResponsibleRole,
		SLADays:         s.SLADays,
	}
}

// ToWorkflowStepResponseList maps a slice of step definitions to response models.
func ToWorkflowStepResponseList(items []*model.WorkflowStepDefinition) []WorkflowStepResponse {
	out := make([]WorkflowStepResponse, 0, len(items))
	for _, s := range items {
		out = append(out, ToWorkflowStepResponse(s))
	}
	return out
}

// ToWorkflowStateResponse maps a project state to its response model.
func ToWorkflowStateResponse(s *model.ProjectWorkflowState) ProjectWorkflowStateResponse {
	return ProjectWorkflowStateResponse{
		ID:             s.ID,
		ProjectID:      s.ProjectID,
		WorkflowStepID: s.WorkflowStepID,
		Status:         string(s.Status),
		StartedAt:      s.StartedAt,
		CompletedAt:    s.CompletedAt,
		DueAt:          s.DueAt,
		BlockedReason:  s.BlockedReason,
		UpdatedBy:      s.UpdatedBy,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

// ToWorkflowStateResponseList maps a slice of project states to response models.
func ToWorkflowStateResponseList(items []*model.ProjectWorkflowState) []ProjectWorkflowStateResponse {
	out := make([]ProjectWorkflowStateResponse, 0, len(items))
	for _, s := range items {
		out = append(out, ToWorkflowStateResponse(s))
	}
	return out
}

// ToWorkflowHistoryResponse maps a history record to its response model.
func ToWorkflowHistoryResponse(h *model.ProjectWorkflowHistory) WorkflowHistoryResponse {
	return WorkflowHistoryResponse{
		ID:         h.ID,
		ProjectID:  h.ProjectID,
		FromStepID: h.FromStepID,
		ToStepID:   h.ToStepID,
		Action:     h.Action,
		Notes:      h.Notes,
		ChangedBy:  h.ChangedBy,
		CreatedAt:  h.CreatedAt,
	}
}

// ToWorkflowHistoryResponseList maps a slice of history records to response models.
func ToWorkflowHistoryResponseList(items []*model.ProjectWorkflowHistory) []WorkflowHistoryResponse {
	out := make([]WorkflowHistoryResponse, 0, len(items))
	for _, h := range items {
		out = append(out, ToWorkflowHistoryResponse(h))
	}
	return out
}
