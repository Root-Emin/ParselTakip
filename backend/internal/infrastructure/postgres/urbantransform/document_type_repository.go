package urbantransform

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

var _ repository.DocumentTypeRepository = (*DocumentTypeRepo)(nil)

// DocumentTypeRepo implements repository.DocumentTypeRepository using PostgreSQL.
type DocumentTypeRepo struct {
	db *pgxpool.Pool
}

// NewDocumentTypeRepo creates a new DocumentTypeRepo.
func NewDocumentTypeRepo(db *pgxpool.Pool) *DocumentTypeRepo {
	return &DocumentTypeRepo{db: db}
}

// ListAll returns all document types, optionally filtered by category.
func (r *DocumentTypeRepo) ListAll(ctx context.Context, category string) ([]*model.DocumentType, error) {
	query := `
		SELECT id, code, name, COALESCE(description, '') AS description, category, is_mandatory,
			requires_notary, requires_municipal_stamp, is_valid_without_notary,
			COALESCE(invalid_reason, '') AS invalid_reason, created_at, updated_at
		FROM document_types`
	args := []interface{}{}
	if category != "" {
		query += ` WHERE category = $1`
		args = append(args, category)
	}
	query += ` ORDER BY name`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to list document types", err)
	}
	defer rows.Close()

	var types []*model.DocumentType
	for rows.Next() {
		var t model.DocumentType
		if err := rows.Scan(
			&t.ID, &t.Code, &t.Name, &t.Description, &t.Category, &t.IsMandatory,
			&t.RequiresNotary, &t.RequiresMunicipalStamp, &t.IsValidWithoutNotary,
			&t.InvalidReason, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, domainErr.New(domainErr.ErrInternal, "failed to scan document type", err)
		}
		types = append(types, &t)
	}
	return types, rows.Err()
}
