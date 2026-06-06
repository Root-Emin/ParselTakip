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

var _ repository.ApprovalRepository = (*ApprovalRepo)(nil)

const approvalColumns = `id, organization_id, app_id, project_id, approval_type, approver_id, approver_role,
	owner_id, status, COALESCE(decision_notes, '') AS decision_notes, expires_at, decided_at, created_at, updated_at`

// ApprovalRepo implements repository.ApprovalRepository using PostgreSQL.
type ApprovalRepo struct {
	db *pgxpool.Pool
}

// NewApprovalRepo creates a new ApprovalRepo.
func NewApprovalRepo(db *pgxpool.Pool) *ApprovalRepo {
	return &ApprovalRepo{db: db}
}

func scanApproval(row pgx.Row) (*model.Approval, error) {
	var a model.Approval
	err := row.Scan(
		&a.ID, &a.OrganizationID, &a.AppID, &a.ProjectID, &a.ApprovalType, &a.ApproverID, &a.ApproverRole,
		&a.OwnerID, &a.Status, &a.DecisionNotes, &a.ExpiresAt, &a.DecidedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Create inserts a new approval.
func (r *ApprovalRepo) Create(ctx context.Context, a *model.Approval) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	if a.Status == "" {
		a.Status = model.ApprovalStatusPending
	}

	query := `
		INSERT INTO approvals (
			id, organization_id, app_id, project_id, approval_type, approver_id, approver_role,
			owner_id, status, decision_notes, expires_at, decided_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`

	_, err := r.db.Exec(ctx, query,
		a.ID, a.OrganizationID, a.AppID, a.ProjectID, a.ApprovalType, a.ApproverID, a.ApproverRole,
		a.OwnerID, a.Status, a.DecisionNotes, a.ExpiresAt, a.DecidedAt, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create approval", err)
	}
	return nil
}

// GetByID retrieves an approval by ID scoped to the tenant.
func (r *ApprovalRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Approval, error) {
	query := fmt.Sprintf(`SELECT %s FROM approvals
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, approvalColumns)
	a, err := scanApproval(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "approval not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get approval", err)
	}
	return a, nil
}

// Update updates an existing approval.
func (r *ApprovalRepo) Update(ctx context.Context, a *model.Approval) error {
	a.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE approvals SET
			approver_id = $4, approver_role = $5, status = $6, decision_notes = $7,
			expires_at = $8, decided_at = $9, updated_at = $10
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		a.ID, a.OrganizationID, a.AppID,
		a.ApproverID, a.ApproverRole, a.Status, a.DecisionNotes,
		a.ExpiresAt, a.DecidedAt, a.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update approval", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "approval not found", nil)
	}
	return nil
}

// Delete removes an approval scoped to the tenant.
func (r *ApprovalRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM approvals WHERE id = $1 AND organization_id = $2 AND app_id = $3`,
		id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete approval", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "approval not found", nil)
	}
	return nil
}

// List returns filtered, sorted and paginated approvals with a total count.
func (r *ApprovalRepo) List(ctx context.Context, f repository.ApprovalFilter) ([]*model.Approval, int, error) {
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
	if f.ProjectID != nil {
		add("project_id = $%d", *f.ProjectID)
	}
	if f.ApprovalType != nil {
		add("approval_type = $%d", string(*f.ApprovalType))
	}
	if f.Status != nil {
		add("status = $%d", string(*f.Status))
	}
	if f.ApproverID != nil {
		add("approver_id = $%d", *f.ApproverID)
	}
	if f.OwnerID != nil {
		add("owner_id = $%d", *f.OwnerID)
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM approvals "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count approvals", err)
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM approvals %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		approvalColumns, where, sortBy, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list approvals", err)
	}
	defer rows.Close()

	var approvals []*model.Approval
	for rows.Next() {
		a, err := scanApproval(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan approval", err)
		}
		approvals = append(approvals, a)
	}
	return approvals, total, rows.Err()
}
