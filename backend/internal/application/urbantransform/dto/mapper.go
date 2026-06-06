package dto

import "github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"

// ToProjectResponse maps a domain project to its response model.
func ToProjectResponse(p *model.Project) ProjectResponse {
	return ProjectResponse{
		ID:                    p.ID,
		OrganizationID:        p.OrganizationID,
		AppID:                 p.AppID,
		Code:                  p.Code,
		Name:                  p.Name,
		Description:           p.Description,
		Status:                string(p.Status),
		CurrentWorkflowStepID: p.CurrentWorkflowStepID,
		InitiatedBy:           p.InitiatedBy,
		AssignedContractorID:  p.AssignedContractorID,
		StartedAt:             p.StartedAt,
		TargetCompletionAt:    p.TargetCompletionAt,
		CompletedAt:           p.CompletedAt,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}

// ToProjectResponseList maps a slice of domain projects to response models.
func ToProjectResponseList(items []*model.Project) []ProjectResponse {
	out := make([]ProjectResponse, 0, len(items))
	for _, p := range items {
		out = append(out, ToProjectResponse(p))
	}
	return out
}
