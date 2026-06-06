package iam

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

// SystemRoleRepo implements repository.SystemRoleRepository with PostgreSQL.
type SystemRoleRepo struct {
	db *pgxpool.Pool
}

// NewSystemRoleRepo creates a new SystemRoleRepo.
func NewSystemRoleRepo(db *pgxpool.Pool) *SystemRoleRepo {
	return &SystemRoleRepo{db: db}
}

func (r *SystemRoleRepo) ListAll(ctx context.Context) ([]*model.SystemRoleDefinition, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, code, name, description, permissions, created_at
		 FROM system_role_definitions ORDER BY code`,
	)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to list system roles", err)
	}
	defer rows.Close()

	var roles []*model.SystemRoleDefinition
	for rows.Next() {
		var role model.SystemRoleDefinition
		var permsJSON []byte
		if err := rows.Scan(&role.ID, &role.Code, &role.Name, &role.Description, &permsJSON, &role.CreatedAt); err != nil {
			return nil, domainErr.New(domainErr.ErrInternal, "failed to scan system role", err)
		}
		if len(permsJSON) > 0 {
			_ = json.Unmarshal(permsJSON, &role.Permissions)
		}
		roles = append(roles, &role)
	}
	return roles, nil
}
