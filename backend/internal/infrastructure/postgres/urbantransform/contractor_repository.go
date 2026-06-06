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

var _ repository.ContractorRepository = (*ContractorRepo)(nil)

const contractorColumns = `id, organization_id, app_id, user_id, company_name, tax_number,
	COALESCE(trade_registry_no, '') AS trade_registry_no, COALESCE(authorized_person, '') AS authorized_person,
	COALESCE(phone, '') AS phone, COALESCE(email, '') AS email, COALESCE(address, '') AS address,
	status, created_at, updated_at`

// ContractorRepo implements repository.ContractorRepository using PostgreSQL.
type ContractorRepo struct {
	db *pgxpool.Pool
}

// NewContractorRepo creates a new ContractorRepo.
func NewContractorRepo(db *pgxpool.Pool) *ContractorRepo {
	return &ContractorRepo{db: db}
}

func scanContractor(row pgx.Row) (*model.Contractor, error) {
	var c model.Contractor
	err := row.Scan(
		&c.ID, &c.OrganizationID, &c.AppID, &c.UserID, &c.CompanyName, &c.TaxNumber,
		&c.TradeRegistryNo, &c.AuthorizedPerson, &c.Phone, &c.Email, &c.Address,
		&c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Create inserts a new contractor.
func (r *ContractorRepo) Create(ctx context.Context, c *model.Contractor) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Status == "" {
		c.Status = model.ContractorStatusActive
	}

	query := `
		INSERT INTO contractor_companies (
			id, organization_id, app_id, user_id, company_name, tax_number,
			trade_registry_no, authorized_person, phone, email, address, status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`

	_, err := r.db.Exec(ctx, query,
		c.ID, c.OrganizationID, c.AppID, c.UserID, c.CompanyName, c.TaxNumber,
		c.TradeRegistryNo, c.AuthorizedPerson, c.Phone, c.Email, c.Address, c.Status, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create contractor", err)
	}
	return nil
}

// GetByID retrieves a contractor by ID scoped to the tenant.
func (r *ContractorRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Contractor, error) {
	query := fmt.Sprintf(`SELECT %s FROM contractor_companies
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, contractorColumns)
	c, err := scanContractor(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "contractor not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get contractor", err)
	}
	return c, nil
}

// GetByTaxNumber retrieves a contractor by tax number scoped to the organization.
func (r *ContractorRepo) GetByTaxNumber(ctx context.Context, orgID uuid.UUID, taxNumber string) (*model.Contractor, error) {
	query := fmt.Sprintf(`SELECT %s FROM contractor_companies
		WHERE organization_id = $1 AND tax_number = $2`, contractorColumns)
	c, err := scanContractor(r.db.QueryRow(ctx, query, orgID, taxNumber))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "contractor not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get contractor", err)
	}
	return c, nil
}

// Update updates an existing contractor.
func (r *ContractorRepo) Update(ctx context.Context, c *model.Contractor) error {
	c.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE contractor_companies SET
			user_id = $4, company_name = $5, trade_registry_no = $6, authorized_person = $7,
			phone = $8, email = $9, address = $10, status = $11, updated_at = $12
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		c.ID, c.OrganizationID, c.AppID,
		c.UserID, c.CompanyName, c.TradeRegistryNo, c.AuthorizedPerson,
		c.Phone, c.Email, c.Address, c.Status, c.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update contractor", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "contractor not found", nil)
	}
	return nil
}

// Delete removes a contractor scoped to the tenant.
func (r *ContractorRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM contractor_companies WHERE id = $1 AND organization_id = $2 AND app_id = $3`,
		id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete contractor", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "contractor not found", nil)
	}
	return nil
}

// List returns filtered, searched, sorted and paginated contractors with a total count.
func (r *ContractorRepo) List(ctx context.Context, f repository.ContractorFilter) ([]*model.Contractor, int, error) {
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
	if strings.TrimSpace(f.Search) != "" {
		pattern := "%" + strings.TrimSpace(f.Search) + "%"
		conditions = append(conditions, fmt.Sprintf(
			"(company_name ILIKE $%d OR tax_number ILIKE $%d OR authorized_person ILIKE $%d OR email ILIKE $%d)",
			idx, idx, idx, idx))
		args = append(args, pattern)
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM contractor_companies "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count contractors", err)
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM contractor_companies %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		contractorColumns, where, sortBy, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list contractors", err)
	}
	defer rows.Close()

	var contractors []*model.Contractor
	for rows.Next() {
		c, err := scanContractor(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan contractor", err)
		}
		contractors = append(contractors, c)
	}
	return contractors, total, rows.Err()
}
