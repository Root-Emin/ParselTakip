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

var _ repository.BuildingUnitRepository = (*BuildingUnitRepo)(nil)

const buildingUnitColumns = `id, organization_id, app_id, building_id, unit_no, floor_no, area_sqm,
	COALESCE(room_count, '') AS room_count, ownership_type,
	COALESCE(title_deed_no, '') AS title_deed_no, status, created_at, updated_at`

// BuildingUnitRepo implements repository.BuildingUnitRepository using PostgreSQL.
type BuildingUnitRepo struct {
	db *pgxpool.Pool
}

// NewBuildingUnitRepo creates a new BuildingUnitRepo.
func NewBuildingUnitRepo(db *pgxpool.Pool) *BuildingUnitRepo {
	return &BuildingUnitRepo{db: db}
}

func scanBuildingUnit(row pgx.Row) (*model.BuildingUnit, error) {
	var u model.BuildingUnit
	err := row.Scan(
		&u.ID, &u.OrganizationID, &u.AppID, &u.BuildingID, &u.UnitNo, &u.FloorNo, &u.AreaSqm,
		&u.RoomCount, &u.OwnershipType, &u.TitleDeedNo, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create inserts a new building unit.
func (r *BuildingUnitRepo) Create(ctx context.Context, u *model.BuildingUnit) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	if u.Status == "" {
		u.Status = model.UnitStatusActive
	}
	if u.OwnershipType == "" {
		u.OwnershipType = model.OwnershipTypeKatMulkiyeti
	}

	query := `
		INSERT INTO building_units (
			id, organization_id, app_id, building_id, unit_no, floor_no, area_sqm,
			room_count, ownership_type, title_deed_no, status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`

	_, err := r.db.Exec(ctx, query,
		u.ID, u.OrganizationID, u.AppID, u.BuildingID, u.UnitNo, u.FloorNo, u.AreaSqm,
		u.RoomCount, u.OwnershipType, u.TitleDeedNo, u.Status, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create building unit", err)
	}
	return nil
}

// GetByID retrieves a building unit by ID scoped to the tenant.
func (r *BuildingUnitRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.BuildingUnit, error) {
	query := fmt.Sprintf(`SELECT %s FROM building_units
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, buildingUnitColumns)
	u, err := scanBuildingUnit(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "building unit not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get building unit", err)
	}
	return u, nil
}

// Update updates an existing building unit.
func (r *BuildingUnitRepo) Update(ctx context.Context, u *model.BuildingUnit) error {
	u.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE building_units SET
			unit_no = $4, floor_no = $5, area_sqm = $6, room_count = $7,
			ownership_type = $8, title_deed_no = $9, status = $10, updated_at = $11
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		u.ID, u.OrganizationID, u.AppID,
		u.UnitNo, u.FloorNo, u.AreaSqm, u.RoomCount,
		u.OwnershipType, u.TitleDeedNo, u.Status, u.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update building unit", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "building unit not found", nil)
	}
	return nil
}

// Delete removes a building unit scoped to the tenant.
func (r *BuildingUnitRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM building_units WHERE id = $1 AND organization_id = $2 AND app_id = $3`,
		id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete building unit", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "building unit not found", nil)
	}
	return nil
}

// List returns filtered, searched, sorted and paginated building units with a total count.
func (r *BuildingUnitRepo) List(ctx context.Context, f repository.BuildingUnitFilter) ([]*model.BuildingUnit, int, error) {
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
	if f.BuildingID != nil {
		add("building_id = $%d", *f.BuildingID)
	}
	if f.Status != nil {
		add("status = $%d", string(*f.Status))
	}
	if f.OwnershipType != nil {
		add("ownership_type = $%d", string(*f.OwnershipType))
	}
	if strings.TrimSpace(f.Search) != "" {
		pattern := "%" + strings.TrimSpace(f.Search) + "%"
		conditions = append(conditions, fmt.Sprintf("(unit_no ILIKE $%d OR title_deed_no ILIKE $%d)", idx, idx))
		args = append(args, pattern)
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM building_units "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count building units", err)
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM building_units %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		buildingUnitColumns, where, sortBy, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list building units", err)
	}
	defer rows.Close()

	var units []*model.BuildingUnit
	for rows.Next() {
		u, err := scanBuildingUnit(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan building unit", err)
		}
		units = append(units, u)
	}
	return units, total, rows.Err()
}
