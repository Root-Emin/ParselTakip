package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/stats/model"
)

// StatsRepository defines read-only aggregate queries for dashboards.
type StatsRepository interface {
	ProjectDashboard(ctx context.Context, orgID, appID uuid.UUID) (*model.ProjectDashboardStats, error)
	AdminDashboard(ctx context.Context) (*model.AdminDashboardStats, error)
}
