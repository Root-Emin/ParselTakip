package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
)

// GetProjectHandler handles the get-project query (CQRS read side).
type GetProjectHandler struct {
	repo repository.ProjectRepository
}

// NewGetProjectHandler creates a new GetProjectHandler.
func NewGetProjectHandler(repo repository.ProjectRepository) *GetProjectHandler {
	return &GetProjectHandler{repo: repo}
}

// Execute returns a single project by ID within the tenant scope.
func (h *GetProjectHandler) Execute(ctx context.Context, orgID, appID, id uuid.UUID) (*dto.ProjectResponse, error) {
	project, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToProjectResponse(project)
	return &resp, nil
}
