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

var _ repository.BuildingRepository = (*BuildingRepo)(nil)

const buildingColumns = `id, organization_id, app_id, COALESCE(name, '') AS name, address, city, district,
	COALESCE(neighborhood, '') AS neighborhood, COALESCE(block_no, '') AS block_no,
	COALESCE(parcel_no, '') AS parcel_no, COALESCE(island_no, '') AS island_no,
	floor_count, unit_count, construction_year,
	COALESCE(building_type, 'residential') AS building_type, risk_status, status,
	created_by, created_at, updated_at`

// BuildingRepo implements repository.BuildingRepository using PostgreSQL.
type BuildingRepo struct {
	db *pgxpool.Pool
}

// NewBuildingRepo creates a new BuildingRepo.
func NewBuildingRepo(db *pgxpool.Pool) *BuildingRepo {
	return &BuildingRepo{db: db}
}

func scanBuilding(row pgx.Row) (*model.Building, error) {
	var b model.Building
	err := row.Scan(
		&b.ID, &b.OrganizationID, &b.AppID, &b.Name, &b.Address, &b.City, &b.District,
		&b.Neighborhood, &b.BlockNo, &b.ParcelNo, &b.IslandNo,
		&b.FloorCount, &b.UnitCount, &b.ConstructionYear,
		&b.BuildingType, &b.RiskStatus, &b.Status,
		&b.CreatedBy, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Create inserts a new building.
func (r *BuildingRepo) Create(ctx context.Context, b *model.Building) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now
	if b.Status == "" {
		b.Status = model.BuildingStatusActive
	}
	if b.RiskStatus == "" {
		b.RiskStatus = model.RiskStatusUnknown
	}
	if b.BuildingType == "" {
		b.BuildingType = model.BuildingTypeResidential
	}
	if b.UnitCount < 1 {
		b.UnitCount = 1
	}

	query := `
		INSERT INTO buildings (
			id, organization_id, app_id, name, address, city, district, neighborhood,
			block_no, parcel_no, island_no, floor_count, unit_count, construction_year,
			building_type, risk_status, status, created_by, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`

	_, err := r.db.Exec(ctx, query,
		b.ID, b.OrganizationID, b.AppID, b.Name, b.Address, b.City, b.District, b.Neighborhood,
		b.BlockNo, b.ParcelNo, b.IslandNo, b.FloorCount, b.UnitCount, b.ConstructionYear,
		b.BuildingType, b.RiskStatus, b.Status, b.CreatedBy, b.CreatedAt, b.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create building", err)
	}
	return nil
}

// GetByID retrieves a building by ID scoped to the tenant.
func (r *BuildingRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Building, error) {
	query := fmt.Sprintf(`SELECT %s FROM buildings
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, buildingColumns)
	b, err := scanBuilding(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "building not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get building", err)
	}
	return b, nil
}

// Update updates an existing building.
func (r *BuildingRepo) Update(ctx context.Context, b *model.Building) error {
	b.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE buildings SET
			name = $4, address = $5, city = $6, district = $7, neighborhood = $8,
			block_no = $9, parcel_no = $10, island_no = $11, floor_count = $12, unit_count = $13,
			construction_year = $14, building_type = $15, risk_status = $16, status = $17, updated_at = $18
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		b.ID, b.OrganizationID, b.AppID,
		b.Name, b.Address, b.City, b.District, b.Neighborhood,
		b.BlockNo, b.ParcelNo, b.IslandNo, b.FloorCount, b.UnitCount,
		b.ConstructionYear, b.BuildingType, b.RiskStatus, b.Status, b.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update building", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "building not found", nil)
	}
	return nil
}

// Delete removes a building scoped to the tenant.
func (r *BuildingRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM buildings WHERE id = $1 AND organization_id = $2 AND app_id = $3`,
		id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete building", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "building not found", nil)
	}
	return nil
}

// List returns filtered, searched, sorted and paginated buildings with a total count.
func (r *BuildingRepo) List(ctx context.Context, f repository.BuildingFilter) ([]*model.Building, int, error) {
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
	if f.RiskStatus != nil {
		add("risk_status = $%d", string(*f.RiskStatus))
	}
	if f.BuildingType != nil {
		add("building_type = $%d", string(*f.BuildingType))
	}
	if strings.TrimSpace(f.City) != "" {
		add("city ILIKE $%d", strings.TrimSpace(f.City))
	}
	if strings.TrimSpace(f.District) != "" {
		add("district ILIKE $%d", strings.TrimSpace(f.District))
	}
	if strings.TrimSpace(f.Search) != "" {
		pattern := "%" + strings.TrimSpace(f.Search) + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(name ILIKE $%d OR address ILIKE $%d OR block_no ILIKE $%d OR parcel_no ILIKE $%d)",
			idx, idx, idx, idx))
		args = append(args, pattern)
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM buildings "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count buildings", err)
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM buildings %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		buildingColumns, where, sortBy, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list buildings", err)
	}
	defer rows.Close()

	var buildings []*model.Building
	for rows.Next() {
		b, err := scanBuilding(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan building", err)
		}
		buildings = append(buildings, b)
	}
	return buildings, total, rows.Err()
}
