package urbantransform

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

var _ repository.ProjectRepository = (*ProjectRepo)(nil)

const projectColumns = `id, organization_id, app_id, code, name, COALESCE(description, '') AS description, status,
	current_workflow_step_id, initiated_by, assigned_contractor_id,
	started_at, target_completion_at, completed_at, created_at, updated_at`

// ProjectRepo implements repository.ProjectRepository using PostgreSQL.
type ProjectRepo struct {
	db *pgxpool.Pool
}

// NewProjectRepo creates a new ProjectRepo.
func NewProjectRepo(db *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func scanProject(row pgx.Row) (*model.Project, error) {
	var p model.Project
	err := row.Scan(
		&p.ID, &p.OrganizationID, &p.AppID, &p.Code, &p.Name, &p.Description, &p.Status,
		&p.CurrentWorkflowStepID, &p.InitiatedBy, &p.AssignedContractorID,
		&p.StartedAt, &p.TargetCompletionAt, &p.CompletedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create inserts a new project.
func (r *ProjectRepo) Create(ctx context.Context, p *model.Project) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Status == "" {
		p.Status = model.ProjectStatusDraft
	}

	query := `
		INSERT INTO urban_transformation_projects (
			id, organization_id, app_id, code, name, description, status,
			current_workflow_step_id, initiated_by, assigned_contractor_id,
			started_at, target_completion_at, completed_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`

	_, err := r.db.Exec(ctx, query,
		p.ID, p.OrganizationID, p.AppID, p.Code, p.Name, p.Description, p.Status,
		p.CurrentWorkflowStepID, p.InitiatedBy, p.AssignedContractorID,
		p.StartedAt, p.TargetCompletionAt, p.CompletedAt, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create project", err)
	}
	return nil
}

// GetByID retrieves a project by ID scoped to the tenant.
func (r *ProjectRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Project, error) {
	query := fmt.Sprintf(`SELECT %s FROM urban_transformation_projects
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, projectColumns)
	p, err := scanProject(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "project not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get project", err)
	}
	return p, nil
}

// GetByCode retrieves a project by code scoped to the organization.
func (r *ProjectRepo) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*model.Project, error) {
	query := fmt.Sprintf(`SELECT %s FROM urban_transformation_projects
		WHERE organization_id = $1 AND code = $2`, projectColumns)
	p, err := scanProject(r.db.QueryRow(ctx, query, orgID, code))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "project not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get project", err)
	}
	return p, nil
}

// Update updates an existing project.
func (r *ProjectRepo) Update(ctx context.Context, p *model.Project) error {
	p.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE urban_transformation_projects SET
			name = $4, description = $5, status = $6,
			current_workflow_step_id = $7, assigned_contractor_id = $8,
			started_at = $9, target_completion_at = $10, completed_at = $11, updated_at = $12
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		p.ID, p.OrganizationID, p.AppID,
		p.Name, p.Description, p.Status,
		p.CurrentWorkflowStepID, p.AssignedContractorID,
		p.StartedAt, p.TargetCompletionAt, p.CompletedAt, p.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update project", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "project not found", nil)
	}
	return nil
}

// Delete removes a project scoped to the tenant.
func (r *ProjectRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	query := `DELETE FROM urban_transformation_projects
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`
	ct, err := r.db.Exec(ctx, query, id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete project", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "project not found", nil)
	}
	return nil
}

// List returns filtered, searched, sorted and paginated projects with a total count.
func (r *ProjectRepo) List(ctx context.Context, f repository.ProjectFilter) ([]*model.Project, int, error) {
	var conditions []string
	var args []interface{}
	idx := 1

	add := func(cond string, val interface{}) {
		conditions = append(conditions, fmt.Sprintf(cond, idx))
		args = append(args, val)
		idx++
	}

	add("organization_id = $%d", f.OrganizationID)
	add("app_id = $%d", f.AppID)

	if f.Status != nil {
		add("status = $%d", string(*f.Status))
	}
	if f.AssignedContractorID != nil {
		add("assigned_contractor_id = $%d", *f.AssignedContractorID)
	}
	if f.InitiatedBy != nil {
		add("initiated_by = $%d", *f.InitiatedBy)
	}
	if f.CurrentWorkflowStepID != nil {
		add("current_workflow_step_id = $%d", *f.CurrentWorkflowStepID)
	}
	if strings.TrimSpace(f.Search) != "" {
		pattern := "%" + strings.TrimSpace(f.Search) + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(code ILIKE $%d OR name ILIKE $%d OR description ILIKE $%d)", idx, idx, idx))
		args = append(args, pattern)
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	// Count total matching rows.
	countQuery := "SELECT COUNT(*) FROM urban_transformation_projects " + where
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count projects", err)
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM urban_transformation_projects %s
		ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		projectColumns, where, sortBy, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list projects", err)
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan project", err)
		}
		projects = append(projects, p)
	}
	return projects, total, rows.Err()
}
