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

var _ repository.NotificationRepository = (*NotificationRepo)(nil)

const notificationColumns = `id, organization_id, app_id, project_id, user_id, notification_type,
	title, message, channel, is_read, read_at, metadata, created_at`

// NotificationRepo implements repository.NotificationRepository using PostgreSQL.
type NotificationRepo struct {
	db *pgxpool.Pool
}

// NewNotificationRepo creates a new NotificationRepo.
func NewNotificationRepo(db *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func scanNotification(row pgx.Row) (*model.Notification, error) {
	var n model.Notification
	err := row.Scan(
		&n.ID, &n.OrganizationID, &n.AppID, &n.ProjectID, &n.UserID, &n.NotificationType,
		&n.Title, &n.Message, &n.Channel, &n.IsRead, &n.ReadAt, &n.Metadata, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// Create inserts a new notification.
func (r *NotificationRepo) Create(ctx context.Context, n *model.Notification) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	n.CreatedAt = time.Now().UTC()
	if n.Channel == "" {
		n.Channel = model.NotificationChannelInApp
	}

	query := `
		INSERT INTO project_notifications (
			id, organization_id, app_id, project_id, user_id, notification_type,
			title, message, channel, is_read, read_at, metadata, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`

	_, err := r.db.Exec(ctx, query,
		n.ID, n.OrganizationID, n.AppID, n.ProjectID, n.UserID, n.NotificationType,
		n.Title, n.Message, n.Channel, n.IsRead, n.ReadAt, n.Metadata, n.CreatedAt,
	)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to create notification", err)
	}
	return nil
}

// GetByID retrieves a notification by ID scoped to the tenant.
func (r *NotificationRepo) GetByID(ctx context.Context, orgID, appID, id uuid.UUID) (*model.Notification, error) {
	query := fmt.Sprintf(`SELECT %s FROM project_notifications
		WHERE id = $1 AND organization_id = $2 AND app_id = $3`, notificationColumns)
	n, err := scanNotification(r.db.QueryRow(ctx, query, id, orgID, appID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.New(domainErr.ErrNotFound, "notification not found", nil)
		}
		return nil, domainErr.New(domainErr.ErrInternal, "failed to get notification", err)
	}
	return n, nil
}

// List returns the recipient's notifications, filtered, sorted and paginated.
func (r *NotificationRepo) List(ctx context.Context, f repository.NotificationFilter) ([]*model.Notification, int, error) {
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
	add("user_id = $%d", f.UserID)
	if f.ProjectID != nil {
		add("project_id = $%d", *f.ProjectID)
	}
	if f.NotificationType != nil {
		add("notification_type = $%d", *f.NotificationType)
	}
	if f.IsRead != nil {
		add("is_read = $%d", *f.IsRead)
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM project_notifications "+where, args...).Scan(&total); err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to count notifications", err)
	}

	sortOrder := strings.ToUpper(f.SortOrder)
	if sortOrder != "ASC" {
		sortOrder = "DESC"
	}

	listQuery := fmt.Sprintf(`SELECT %s FROM project_notifications %s ORDER BY created_at %s LIMIT $%d OFFSET $%d`,
		notificationColumns, where, sortOrder, idx, idx+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to list notifications", err)
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, 0, domainErr.New(domainErr.ErrInternal, "failed to scan notification", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, total, rows.Err()
}

// MarkRead marks a single notification owned by the user as read.
func (r *NotificationRepo) MarkRead(ctx context.Context, orgID, appID, userID, id uuid.UUID) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE project_notifications SET is_read = TRUE, read_at = $1
		 WHERE id = $2 AND organization_id = $3 AND app_id = $4 AND user_id = $5 AND is_read = FALSE`,
		time.Now().UTC(), id, orgID, appID, userID)
	if err != nil {
		return domainErr.New(domainErr.ErrInternal, "failed to mark notification read", err)
	}
	if ct.RowsAffected() == 0 {
		// Either not found or already read; verify existence for a clearer error.
		if _, gErr := r.GetByID(ctx, orgID, appID, id); gErr != nil {
			return gErr
		}
	}
	return nil
}

// MarkAllRead marks all of the user's unread notifications as read.
func (r *NotificationRepo) MarkAllRead(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error) {
	ct, err := r.db.Exec(ctx,
		`UPDATE project_notifications SET is_read = TRUE, read_at = $1
		 WHERE organization_id = $2 AND app_id = $3 AND user_id = $4 AND is_read = FALSE`,
		time.Now().UTC(), orgID, appID, userID)
	if err != nil {
		return 0, domainErr.New(domainErr.ErrInternal, "failed to mark all notifications read", err)
	}
	return int(ct.RowsAffected()), nil
}

// CountUnread returns the number of unread notifications for the user.
func (r *NotificationRepo) CountUnread(ctx context.Context, orgID, appID, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM project_notifications
		 WHERE organization_id = $1 AND app_id = $2 AND user_id = $3 AND is_read = FALSE`,
		orgID, appID, userID).Scan(&count)
	if err != nil {
		return 0, domainErr.New(domainErr.ErrInternal, "failed to count unread notifications", err)
	}
	return count, nil
}
