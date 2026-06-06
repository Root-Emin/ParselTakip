package stats

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/stats/repository"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"
	"github.com/masterfabric-go/masterfabric/internal/shared/response"
)

// Handler exposes dashboard statistics endpoints.
type Handler struct {
	repo repository.StatsRepository
}

// NewHandler creates a new stats Handler.
func NewHandler(repo repository.StatsRepository) *Handler {
	return &Handler{repo: repo}
}

// ProjectDashboard handles GET /dashboard/stats (municipality manager/staff).
func (h *Handler) ProjectDashboard(w http.ResponseWriter, r *http.Request) {
	orgID, found := middleware.OrgIDFromContext(r.Context())
	if !found || orgID == uuid.Nil {
		if tid, tok := middleware.TenantIDFromContext(r.Context()); tok {
			orgID = tid
		}
	}
	appID, _ := uuid.Parse(r.Header.Get("X-App-ID"))
	if orgID == uuid.Nil || appID == uuid.Nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required (X-App-ID)"})
		return
	}

	stats, err := h.repo.ProjectDashboard(r.Context(), orgID, appID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, stats)
}

// AdminDashboard handles GET /admin/stats (system admin).
func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.AdminDashboard(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, stats)
}
