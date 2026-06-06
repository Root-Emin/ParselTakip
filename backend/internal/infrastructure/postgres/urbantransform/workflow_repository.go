package urbantransform

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

var _ repository.WorkflowRepository = (*WorkflowRepo)(nil)

// WorkflowRepo implements repository.WorkflowRepository using PostgreSQL.
type WorkflowRepo struct {
	db *pgxpool.Pool
}

// NewWorkflowRepo creates a new WorkflowRepo.
func NewWorkflowRepo(db *pgxpool.Pool) *WorkflowRepo {
	return &WorkflowRepo{db: db}
}

// ListStepDefinitions returns all step definitions ordered by step_order.
func (r *WorkflowRepo) ListStepDefinitions(ctx context.Context) ([]*model.WorkflowStepDefinition, error) {
	query := `
		SELECT id, step_order, code, name, COALESCE(description, '') AS description,
			responsible_role, sla_days, created_at
		FROM workflow_step_definitions
		ORDER BY step_order`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to list workflow steps", err)
	}
	defer rows.Close()

	var steps []*model.WorkflowStepDefinition
	for rows.Next() {
		var s model.WorkflowStepDefinition
		if err := rows.Scan(&s.ID, &s.StepOrder, &s.Code, &s.Name, &s.Description,
			&s.ResponsibleRole, &s.SLADays, &s.CreatedAt); err != nil {
			return nil, domainErr.New(domainErr.ErrInternal, "failed to scan workflow step", err)
		}
		steps = append(steps, &s)
	}
	return steps, rows.Err()
}

// GetStepDefinition returns a single step definition by ID.
func (r *WorkflowRepo) GetStepDefinition(ctx context.Context, id uuid.UUID) (*model.WorkflowStepDefinition, error) {
	query := `
		SELECT id, step_order, code, name, COALESCE(description, '') AS description,
			responsible_role, sla_days, created_at
		FROM workflow_step_definitions WHERE id = $1`
	var s model.WorkflowStepDefinition
	err := r.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.StepOrder, &s.Code, &s.Name, &s.Description,
		&s.ResponsibleRole, &s.SLADays, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "workflow step not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get workflow step", err)
	}
	return &s, nil
}

const wfStateColumns = `id, organization_id, app_id, project_id, workflow_step_id, status,
	started_at, completed_at, due_at, COALESCE(blocked_reason, '') AS blocked_reason, updated_by, created_at, updated_at`

func scanState(row pgx.Row) (*model.ProjectWorkflowState, error) {
	var s model.ProjectWorkflowState
	err := row.Scan(&s.ID, &s.OrganizationID, &s.AppID, &s.ProjectID, &s.WorkflowStepID, &s.Status,
		&s.StartedAt, &s.CompletedAt, &s.DueAt, &s.BlockedReason, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// ListStates returns all workflow states for a project ordered by step order.
func (r *WorkflowRepo) ListStates(ctx context.Context, orgID, appID, projectID uuid.UUID) ([]*model.ProjectWorkflowState, error) {
	query := `
		SELECT s.id, s.organization_id, s.app_id, s.project_id, s.workflow_step_id, s.status,
			s.started_at, s.completed_at, s.due_at, COALESCE(s.blocked_reason, '') AS blocked_reason,
			s.updated_by, s.created_at, s.updated_at
		FROM project_workflow_states s
		JOIN workflow_step_definitions d ON d.id = s.workflow_step_id
		WHERE s.organization_id = $1 AND s.app_id = $2 AND s.project_id = $3
		ORDER BY d.step_order`
	rows, err := r.db.Query(ctx, query, orgID, appID, projectID)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to list workflow states", err)
	}
	defer rows.Close()

	var states []*model.ProjectWorkflowState
	for rows.Next() {
		s, err := scanState(rows)
		if err != nil {
			return nil, domainErr.New(domainErr.ErrInternal, "failed to scan workflow state", err)
		}
		states = append(states, s)
	}
	return states, rows.Err()
}

// GetState returns a single state for a (project, step) pair.
func (r *WorkflowRepo) GetState(ctx context.Context, orgID, appID, projectID, stepID uuid.UUID) (*model.ProjectWorkflowState, error) {
	query := `SELECT ` + wfStateColumns + ` FROM project_workflow_states
		WHERE organization_id = $1 AND app_id = $2 AND project_id = $3 AND workflow_step_id = $4`
	s, err := scanState(r.db.QueryRow(ctx, query, orgID, appID, projectID, stepID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "workflow state not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get workflow state", err)
	}
	return s, nil
}

// UpsertState inserts or updates the state for a (project, step) pair.
func (r *WorkflowRepo) UpsertState(ctx context.Context, s *model.ProjectWorkflowState) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now

	query := `
		INSERT INTO project_workflow_states (
			id, organization_id, app_id, project_id, workflow_step_id, status,
			started_at, completed_at, due_at, blocked_reason, updated_by, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		ON CONFLICT (project_id, workflow_step_id) DO UPDATE SET
			status = EXCLUDED.status,
			started_at = COALESCE(project_workflow_states.started_at, EXCLUDED.started_at),
			completed_at = EXCLUDED.completed_at,
			due_at = EXCLUDED.due_at,
			blocked_reason = EXCLUDED.blocked_reason,
			updated_by = EXCLUDED.updated_by,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.Exec(ctx, query,
		s.ID, s.OrganizationID, s.AppID, s.ProjectID, s.WorkflowStepID, s.Status,
		s.StartedAt, s.CompletedAt, s.DueAt, s.BlockedReason, s.UpdatedBy, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to upsert workflow state", err)
	}
	return nil
}

// AddHistory inserts a workflow transition record.
func (r *WorkflowRepo) AddHistory(ctx context.Context, h *model.ProjectWorkflowHistory) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	h.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO project_workflow_history (
			id, organization_id, app_id, project_id, from_step_id, to_step_id, action, notes, changed_by, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := r.db.Exec(ctx, query,
		h.ID, h.OrganizationID, h.AppID, h.ProjectID, h.FromStepID, h.ToStepID, h.Action, h.Notes, h.ChangedBy, h.CreatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to add workflow history", err)
	}
	return nil
}

// ListHistory returns the workflow transition history for a project, newest first.
func (r *WorkflowRepo) ListHistory(ctx context.Context, orgID, appID, projectID uuid.UUID) ([]*model.ProjectWorkflowHistory, error) {
	query := `
		SELECT id, organization_id, app_id, project_id, from_step_id, to_step_id,
			action, COALESCE(notes, '') AS notes, changed_by, created_at
		FROM project_workflow_history
		WHERE organization_id = $1 AND app_id = $2 AND project_id = $3
		ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, orgID, appID, projectID)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to list workflow history", err)
	}
	defer rows.Close()

	var history []*model.ProjectWorkflowHistory
	for rows.Next() {
		var h model.ProjectWorkflowHistory
		if err := rows.Scan(&h.ID, &h.OrganizationID, &h.AppID, &h.ProjectID, &h.FromStepID, &h.ToStepID,
			&h.Action, &h.Notes, &h.ChangedBy, &h.CreatedAt); err != nil {
			return nil, domainErr.New(domainErr.ErrInternal, "failed to scan workflow history", err)
		}
		history = append(history, &h)
	}
	return history, rows.Err()
}
