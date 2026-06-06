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

var _ repository.PropertyOwnerRepository = (*PropertyOwnerRepo)(nil)

const propertyOwnerColumns = `id, organization_id, app_id, user_id, unit_id, first_name, last_name,
	COALESCE(identity_number, '') AS identity_number, COALESCE(phone, '') AS phone,
	COALESCE(email, '') AS email, COALESCE(address, '') AS address, COALESCE(iban, '') AS iban,
	ownership_ratio, is_primary_contact, status, created_at, updated_at`

// PropertyOwnerRepo implements repository.PropertyOwnerRepository using PostgreSQL.
type PropertyOwnerRepo struct {
	db *pgxpool.Pool
}

// NewPropertyOwnerRepo creates a new PropertyOwnerRepo.
func NewPropertyOwnerRepo(db *pgxpool.Pool) *PropertyOwnerRepo {
	return &PropertyOwnerRepo{db: db}
}

func scanPropertyOwner(row pgx.Row) (*model.PropertyOwner, error) {
	var o model.PropertyOwner
	err := row.Scan(
		&o.ID, &o.OrganizationID, &o.AppID, &o.UserID, &o.UnitID, &o.FirstName, &o.LastName,
		&o.IdentityNumber, &o.Phone, &o.Email, &o.Address, &o.IBAN,
		&o.OwnershipRatio, &o.IsPrimaryContact, &o.Status, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// Create inserts a new property owner.
func (r *PropertyOwnerRepo) Create(ctx context.Context, o *model.PropertyOwner) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	if o.Status == "" {
		o.Status = model.OwnerStatusActive
	}
	if o.OwnershipRatio <= 0 {
		o.OwnershipRatio = 1.0
	}

	query := `
		INSERT INTO property_owners (
			id, organization_id, app_id, user_id, unit_id, first_name, last_name,
			identity_number, phone, email, address, iban, ownership_ratio, is_primary_contact,
			status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`

	_, err := r.db.Exec(ctx, query,
		o.ID, o.OrganizationID, o.AppID, o.UserID, o.UnitID, o.FirstName, o.LastName,
		o.IdentityNumber, o.Phone, o.Email, o.Address, o.IBAN, o.OwnershipRatio, o.IsPrimaryContact,
		o.Status, o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create property owner", err)
	}
	return nil
}

// GetByID retrieves a property owner by ID scoped to the tenant.
func (r *PropertyOwnerRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.PropertyOwner, error) {
	query := fmt.Sprintf(`SELECT %s FROM property_owners
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, propertyOwnerColumns)
	o, err := scanPropertyOwner(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "property owner not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get property owner", err)
	}
	return o, nil
}

// Update updates an existing property owner.
func (r *PropertyOwnerRepo) Update(ctx context.Context, o *model.PropertyOwner) error {
	o.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE property_owners SET
			first_name = $4, last_name = $5, identity_number = $6, phone = $7, email = $8,
			address = $9, iban = $10, ownership_ratio = $11, is_primary_contact = $12,
			status = $13, updated_at = $14
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		o.ID, o.OrganizationID, o.AppID,
		o.FirstName, o.LastName, o.IdentityNumber, o.Phone, o.Email,
		o.Address, o.IBAN, o.OwnershipRatio, o.IsPrimaryContact, o.Status, o.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update property owner", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "property owner not found", nil)
	}
	return nil
}

// Delete removes a property owner scoped to the tenant.
func (r *PropertyOwnerRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM property_owners WHERE id = $1 AND organization_id = $2 AND app_id = $3`,
		id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete property owner", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "property owner not found", nil)
	}
	return nil
}

// List returns filtered, searched, sorted and paginated property owners with a total count.
func (r *PropertyOwnerRepo) List(ctx context.Context, f repository.PropertyOwnerFilter) ([]*model.PropertyOwner, int, error) {
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
	if f.UnitID != nil {
		add("unit_id = $%d", *f.UnitID)
	}
	if f.Status != nil {
		add("status = $%d", string(*f.Status))
	}
	if strings.TrimSpace(f.Search) != "" {
		pattern := "%" + strings.TrimSpace(f.Search) + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(first_name ILIKE $%d OR last_name ILIKE $%d OR identity_number ILIKE $%d OR email ILIKE $%d)",
			idx, idx, idx, idx))
		args = append(args, pattern)
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM property_owners "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count property owners", err)
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM property_owners %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		propertyOwnerColumns, where, sortBy, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list property owners", err)
	}
	defer rows.Close()

	var owners []*model.PropertyOwner
	for rows.Next() {
		o, err := scanPropertyOwner(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan property owner", err)
		}
		owners = append(owners, o)
	}
	return owners, total, rows.Err()
}
