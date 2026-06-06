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

var _ repository.DocumentRepository = (*DocumentRepo)(nil)

const documentColumns = `id, organization_id, app_id, project_id, document_type_id, building_id, unit_id, owner_id,
	file_name, file_path, file_size, COALESCE(mime_type, '') AS mime_type, status, is_notarized,
	notary_date, expiry_date, uploaded_by, COALESCE(uploaded_by_role, '') AS uploaded_by_role,
	version, metadata, created_at, updated_at`

// DocumentRepo implements repository.DocumentRepository using PostgreSQL.
type DocumentRepo struct {
	db *pgxpool.Pool
}

// NewDocumentRepo creates a new DocumentRepo.
func NewDocumentRepo(db *pgxpool.Pool) *DocumentRepo {
	return &DocumentRepo{db: db}
}

func scanDocument(row pgx.Row) (*model.Document, error) {
	var d model.Document
	err := row.Scan(
		&d.ID, &d.OrganizationID, &d.AppID, &d.ProjectID, &d.DocumentTypeID, &d.BuildingID, &d.UnitID, &d.OwnerID,
		&d.FileName, &d.FilePath, &d.FileSize, &d.MimeType, &d.Status, &d.IsNotarized,
		&d.NotaryDate, &d.ExpiryDate, &d.UploadedBy, &d.UploadedByRole,
		&d.Version, &d.Metadata, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Create inserts a new document.
func (r *DocumentRepo) Create(ctx context.Context, d *model.Document) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Status == "" {
		d.Status = model.DocumentStatusDraft
	}
	if d.Version < 1 {
		d.Version = 1
	}

	query := `
		INSERT INTO documents (
			id, organization_id, app_id, project_id, document_type_id, building_id, unit_id, owner_id,
			file_name, file_path, file_size, mime_type, status, is_notarized,
			notary_date, expiry_date, uploaded_by, uploaded_by_role, version, metadata, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)`

	_, err := r.db.Exec(ctx, query,
		d.ID, d.OrganizationID, d.AppID, d.ProjectID, d.DocumentTypeID, d.BuildingID, d.UnitID, d.OwnerID,
		d.FileName, d.FilePath, d.FileSize, d.MimeType, d.Status, d.IsNotarized,
		d.NotaryDate, d.ExpiryDate, d.UploadedBy, d.UploadedByRole, d.Version, d.Metadata, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create document", err)
	}
	return nil
}

// GetByID retrieves a document by ID scoped to the tenant.
func (r *DocumentRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Document, error) {
	query := fmt.Sprintf(`SELECT %s FROM documents
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, documentColumns)
	d, err := scanDocument(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "document not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get document", err)
	}
	return d, nil
}

// Update updates an existing document.
func (r *DocumentRepo) Update(ctx context.Context, d *model.Document) error {
	d.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE documents SET
			file_name = $4, file_path = $5, file_size = $6, mime_type = $7, status = $8,
			is_notarized = $9, notary_date = $10, expiry_date = $11, version = $12, metadata = $13, updated_at = $14
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`

	ct, err := r.db.Exec(ctx, query,
		d.ID, d.OrganizationID, d.AppID,
		d.FileName, d.FilePath, d.FileSize, d.MimeType, d.Status,
		d.IsNotarized, d.NotaryDate, d.ExpiryDate, d.Version, d.Metadata, d.UpdatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to update document", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "document not found", nil)
	}
	return nil
}

// Delete removes a document scoped to the tenant.
func (r *DocumentRepo) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`DELETE FROM documents WHERE id = $1 AND organization_id = $2 AND app_id = $3`,
		id, orgID, appID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to delete document", err)
	}
	if ct.RowsAffected() == 0 {
		return domainErr.New(domainErr.ErrNotFound, "document not found", nil)
	}
	return nil
}

// List returns filtered, searched, sorted and paginated documents with a total count.
func (r *DocumentRepo) List(ctx context.Context, f repository.DocumentFilter) ([]*model.Document, int, error) {
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
	if f.BuildingID != nil {
		add("building_id = $%d", *f.BuildingID)
	}
	if f.OwnerID != nil {
		add("owner_id = $%d", *f.OwnerID)
	}
	if f.DocumentTypeID != nil {
		add("document_type_id = $%d", *f.DocumentTypeID)
	}
	if f.Status != nil {
		add("status = $%d", string(*f.Status))
	}
	if strings.TrimSpace(f.Search) != "" {
		pattern := "%" + strings.TrimSpace(f.Search) + "%"
		conditions = append(conditions, fmt.Sprintf("(file_name ILIKE $%d)", idx))
		args = append(args, pattern)
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM documents "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count documents", err)
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM documents %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		documentColumns, where, sortBy, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list documents", err)
	}
	defer rows.Close()

	var docs []*model.Document
	for rows.Next() {
		d, err := scanDocument(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan document", err)
		}
		docs = append(docs, d)
	}
	return docs, total, rows.Err()
}

// CreateReview inserts a new document review.
func (r *DocumentRepo) CreateReview(ctx context.Context, rv *model.DocumentReview) error {
	if rv.ID == uuid.Nil {
		rv.ID = uuid.New()
	}
	now := time.Now().UTC()
	rv.ReviewedAt = now
	rv.CreatedAt = now

	query := `
		INSERT INTO document_reviews (
			id, organization_id, app_id, document_id, reviewer_id, status, missing_items, review_notes, reviewed_at, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := r.db.Exec(ctx, query,
		rv.ID, rv.OrganizationID, rv.AppID, rv.DocumentID, rv.ReviewerID, rv.Status, rv.MissingItems, rv.ReviewNotes, rv.ReviewedAt, rv.CreatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create document review", err)
	}
	return nil
}

// ListReviews returns review history for a document, newest first.
func (r *DocumentRepo) ListReviews(ctx context.Context, orgID, appID, documentID uuid.UUID) ([]*model.DocumentReview, error) {
	query := `
		SELECT id, organization_id, app_id, document_id, reviewer_id, status,
			missing_items, COALESCE(review_notes, '') AS review_notes, reviewed_at, created_at
		FROM document_reviews
		WHERE organization_id = $1 AND app_id = $2 AND document_id = $3
		ORDER BY reviewed_at DESC`

	rows, err := r.db.Query(ctx, query, orgID, appID, documentID)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to list document reviews", err)
	}
	defer rows.Close()

	var reviews []*model.DocumentReview
	for rows.Next() {
		var rv model.DocumentReview
		if err := rows.Scan(
			&rv.ID, &rv.OrganizationID, &rv.AppID, &rv.DocumentID, &rv.ReviewerID, &rv.Status,
			&rv.MissingItems, &rv.ReviewNotes, &rv.ReviewedAt, &rv.CreatedAt,
		); err != nil {
			return nil, domainErr.New(domainErr.ErrInternal, "failed to scan document review", err)
		}
		reviews = append(reviews, &rv)
	}
	return reviews, rows.Err()
}
