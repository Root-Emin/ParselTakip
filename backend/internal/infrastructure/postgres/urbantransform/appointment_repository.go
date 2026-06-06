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

var _ repository.AppointmentRepository = (*AppointmentRepo)(nil)

const appointmentColumns = `id, organization_id, app_id, project_id, owner_id, title,
	COALESCE(description, '') AS description, COALESCE(location, '') AS location,
	scheduled_at, duration_minutes, status, created_by, created_at, updated_at`

// AppointmentRepo implements repository.AppointmentRepository using PostgreSQL.
type AppointmentRepo struct {
	db *pgxpool.Pool
}

// NewAppointmentRepo creates a new AppointmentRepo.
func NewAppointmentRepo(db *pgxpool.Pool) *AppointmentRepo {
	return &AppointmentRepo{db: db}
}

func scanAppointment(row pgx.Row) (*model.Appointment, error) {
	var a model.Appointment
	err := row.Scan(&a.ID, &a.OrganizationID, &a.AppID, &a.ProjectID, &a.OwnerID, &a.Title,
		&a.Description, &a.Location, &a.ScheduledAt, &a.DurationMinutes, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Create inserts a new appointment.
func (r *AppointmentRepo) Create(ctx context.Context, a *model.Appointment) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	if a.Status == "" {
		a.Status = model.AppointmentStatusScheduled
	}

	query := `
		INSERT INTO appointments (
			id, organization_id, app_id, project_id, owner_id, title, description, location,
			scheduled_at, duration_minutes, status, created_by, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`

	_, err := r.db.Exec(ctx, query,
		a.ID, a.OrganizationID, a.AppID, a.ProjectID, a.OwnerID, a.Title, a.Description, a.Location,
		a.ScheduledAt, a.DurationMinutes, a.Status, a.CreatedBy, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create appointment", err)
	}
	return nil
}

// GetByID retrieves an appointment by ID scoped to the tenant.
func (r *AppointmentRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Appointment, error) {
	query := fmt.Sprintf(`SELECT %s FROM appointments
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, appointmentColumns)
	a, err := scanAppointment(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "appointment not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get appointment", err)
	}
	return a, nil
}

// Update updates an existing appointment.
func (r *AppointmentRepo) Update(ctx context.Context, a *model.Appointment) error {
	a.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE appointments SET
			project_id = $4, owner_id = $5, title = $6, description = $7, location = $8,
			scheduled_at = $9, duration_minutes = $10, status = $11, updated_at = $12
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		a.ID, a.OrganizationID, a.AppID,
		a.ProjectID, a.OwnerID, a.Title, a.Description, a.Location,
		a.ScheduledAt, a.DurationMinutes, a.Status, a.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update appointment", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "appointment not found", nil)
	}
	return nil
}

// Delete removes an appointment scoped to the tenant.
func (r *AppointmentRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM appointments WHERE id = $1 AND organization_id = $2 AND app_id = $3`,
		id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete appointment", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "appointment not found", nil)
	}
	return nil
}

// List returns filtered, sorted and paginated appointments with a total count.
func (r *AppointmentRepo) List(ctx context.Context, f repository.AppointmentFilter) ([]*model.Appointment, int, error) {
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
	if f.OwnerID != nil {
		add("owner_id = $%d", *f.OwnerID)
	}
	if f.Status != nil {
		add("status = $%d", string(*f.Status))
	}
	if f.From != nil {
		add("scheduled_at >= $%d", *f.From)
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM appointments "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count appointments", err)
	}

	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM appointments %s ORDER BY scheduled_at %s LIMIT $%d OFFSET $%d`,
		appointmentColumns, where, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list appointments", err)
	}
	defer rows.Close()

	var appointments []*model.Appointment
	for rows.Next() {
		a, err := scanAppointment(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan appointment", err)
		}
		appointments = append(appointments, a)
	}
	return appointments, total, rows.Err()
}
