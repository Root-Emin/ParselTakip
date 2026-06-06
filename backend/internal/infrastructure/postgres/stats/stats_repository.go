package stats

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric-go/masterfabric/internal/domain/stats/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/stats/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

var _ repository.StatsRepository = (*StatsRepo)(nil)

// StatsRepo implements repository.StatsRepository using PostgreSQL.
type StatsRepo struct {
	db *pgxpool.Pool
}

// NewStatsRepo creates a new StatsRepo.
func NewStatsRepo(db *pgxpool.Pool) *StatsRepo {
	return &StatsRepo{db: db}
}

func (r *StatsRepo) countByStatus(ctx context.Context, query string, args ...interface{}) (map[string]int, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]int{}
	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		result[key] = count
	}
	return result, rows.Err()
}

// ProjectDashboard computes municipality dashboard metrics scoped to org+app.
func (r *StatsRepo) ProjectDashboard(ctx context.Context, orgID, appID uuid.UUID) (*model.ProjectDashboardStats, error) {
	stats := &model.ProjectDashboardStats{ProjectsByStatus: map[string]int{}}

	byStatus, err := r.countByStatus(ctx,
		`SELECT status, COUNT(*) FROM urban_transformation_projects
		 WHERE organization_id = $1 AND app_id = $2 GROUP BY status`, orgID, appID)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to count projects by status", err)
	}
	stats.ProjectsByStatus = byStatus
	for status, count := range byStatus {
		stats.TotalProjects += count
		switch status {
		case "in_progress", "initiated":
			stats.OngoingProjects += count
		case "completed":
			stats.CompletedProjects += count
		}
	}

	queries := []struct {
		dest  *int
		query string
	}{
		{&stats.PendingApprovals, `SELECT COUNT(*) FROM approvals WHERE organization_id = $1 AND app_id = $2 AND status = 'pending'`},
		{&stats.MissingDocuments, `SELECT COUNT(*) FROM documents WHERE organization_id = $1 AND app_id = $2 AND status = 'missing'`},
		{&stats.PendingDocumentReviews, `SELECT COUNT(*) FROM documents WHERE organization_id = $1 AND app_id = $2 AND status IN ('submitted','under_review')`},
		{&stats.TotalBuildings, `SELECT COUNT(*) FROM buildings WHERE organization_id = $1 AND app_id = $2`},
		{&stats.TotalPropertyOwners, `SELECT COUNT(*) FROM property_owners WHERE organization_id = $1 AND app_id = $2`},
	}
	for _, q := range queries {
		if err := r.db.QueryRow(ctx, q.query, orgID, appID).Scan(q.dest); err != nil {
			return nil, domainErr.New(domainErr.ErrInternal, "failed to compute project dashboard metric", err)
		}
	}

	return stats, nil
}

// AdminDashboard computes platform-wide system admin metrics.
func (r *StatsRepo) AdminDashboard(ctx context.Context) (*model.AdminDashboardStats, error) {
	stats := &model.AdminDashboardStats{
		OrganizationsByStatus: map[string]int{},
		UsersByStatus:         map[string]int{},
		RoleDistribution:      map[string]int{},
	}

	orgsByStatus, err := r.countByStatus(ctx, `SELECT status, COUNT(*) FROM organizations GROUP BY status`)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to count organizations", err)
	}
	stats.OrganizationsByStatus = orgsByStatus
	for _, c := range orgsByStatus {
		stats.TotalOrganizations += c
	}

	usersByStatus, err := r.countByStatus(ctx, `SELECT status, COUNT(*) FROM users GROUP BY status`)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to count users", err)
	}
	stats.UsersByStatus = usersByStatus
	for status, c := range usersByStatus {
		stats.TotalUsers += c
		if status == "active" {
			stats.ActiveUsers += c
		}
	}

	roleDist, err := r.countByStatus(ctx,
		`SELECT r.name, COUNT(*) FROM user_roles ur JOIN roles r ON r.id = ur.role_id GROUP BY r.name`)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to compute role distribution", err)
	}
	stats.RoleDistribution = roleDist

	return stats, nil
}
