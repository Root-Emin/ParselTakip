package repository

import (
	"context"

	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
)

// SystemRoleRepository defines persistence for system role templates.
type SystemRoleRepository interface {
	ListAll(ctx context.Context) ([]*model.SystemRoleDefinition, error)
}
