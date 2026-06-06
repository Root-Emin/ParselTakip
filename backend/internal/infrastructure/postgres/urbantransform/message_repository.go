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

var _ repository.MessageRepository = (*MessageRepo)(nil)

const messageColumns = `id, organization_id, app_id, project_id, sender_id, recipient_id, parent_id,
	COALESCE(subject, '') AS subject, body, is_read, read_at, created_at`

// MessageRepo implements repository.MessageRepository using PostgreSQL.
type MessageRepo struct {
	db *pgxpool.Pool
}

// NewMessageRepo creates a new MessageRepo.
func NewMessageRepo(db *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{db: db}
}

func scanMessage(row pgx.Row) (*model.Message, error) {
	var m model.Message
	err := row.Scan(&m.ID, &m.OrganizationID, &m.AppID, &m.ProjectID, &m.SenderID, &m.RecipientID, &m.ParentID,
		&m.Subject, &m.Body, &m.IsRead, &m.ReadAt, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Create inserts a new message.
func (r *MessageRepo) Create(ctx context.Context, m *model.Message) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	m.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO messages (
			id, organization_id, app_id, project_id, sender_id, recipient_id, parent_id, subject, body, is_read, read_at, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	_, err := r.db.Exec(ctx, query,
		m.ID, m.OrganizationID, m.AppID, m.ProjectID, m.SenderID, m.RecipientID, m.ParentID, m.Subject, m.Body, m.IsRead, m.ReadAt, m.CreatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create message", err)
	}
	return nil
}

// GetByID retrieves a message by ID scoped to the tenant.
func (r *MessageRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Message, error) {
	query := fmt.Sprintf(`SELECT %s FROM messages
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, messageColumns)
	m, err := scanMessage(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "message not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get message", err)
	}
	return m, nil
}

// List returns the user's inbox or sent messages, filtered and paginated.
func (r *MessageRepo) List(ctx context.Context, f repository.MessageFilter) ([]*model.Message, int, error) {
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
	if f.Box == "sent" {
		add("sender_id = $%d", f.UserID)
	} else {
		add("recipient_id = $%d", f.UserID)
	}
	if f.ProjectID != nil {
		add("project_id = $%d", *f.ProjectID)
	}
	if f.IsRead != nil {
		add("is_read = $%d", *f.IsRead)
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM messages "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count messages", err)
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM messages %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		messageColumns, where, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list messages", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m, err := scanMessage(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan message", err)
		}
		messages = append(messages, m)
	}
	return messages, total, rows.Err()
}

// MarkRead marks a message owned by the user (as recipient) as read.
func (r *MessageRepo) MarkRead(ctx context.Context, orgID, appID, userID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE messages SET is_read = TRUE, read_at = $1
		 WHERE id = $2 AND organization_id = $3 AND app_id = $4 AND recipient_id = $5 AND is_read = FALSE`,
		time.Now().UTC(), id, orgID, appID, userID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to mark message read", err)
	}
	if ct.RowsAffected() == 0 {
		if _, gErr := r.GetByID(ctx, orgID, appID, id); gErr != nil {
			return gErr
		}
	}
	return nil
}

// CountUnread returns the number of unread inbox messages for the user.
func (r *MessageRepo) CountUnread(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM messages
		 WHERE organization_id = $1 AND app_id = $2 AND recipient_id = $3 AND is_read = FALSE`,
		orgID, appID, userID).Scan(&count)
	if err != nil {
		return 0, domainErr.New(domainErr.ErrInternal, "failed to count unread messages", err)
	}
	return count, nil
}
