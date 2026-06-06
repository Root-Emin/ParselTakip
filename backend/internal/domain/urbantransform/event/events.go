package event

import (
	"time"

	"github.com/google/uuid"
)

// ProjectCreated is published after a project is created.
type ProjectCreated struct {
	ProjectID      uuid.UUID `json:"project_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AppID          uuid.UUID `json:"app_id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Timestamp      time.Time `json:"timestamp"`
}

// ProjectUpdated is published after a project is updated.
type ProjectUpdated struct {
	ProjectID      uuid.UUID `json:"project_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AppID          uuid.UUID `json:"app_id"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
}

// ProjectDeleted is published after a project is deleted (cancelled).
type ProjectDeleted struct {
	ProjectID      uuid.UUID `json:"project_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AppID          uuid.UUID `json:"app_id"`
	Timestamp      time.Time `json:"timestamp"`
}
