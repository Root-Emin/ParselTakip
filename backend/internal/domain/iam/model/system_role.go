package model

import (
	"time"

	"github.com/google/uuid"
)

// SystemRoleDefinition is a template role for urban transformation workflows.
type SystemRoleDefinition struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}
