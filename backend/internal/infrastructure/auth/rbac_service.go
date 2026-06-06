package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/repository"
	"github.com/redis/go-redis/v9"
)

const permissionCacheTTL = 15 * time.Minute

// roleSystemAdmin holds full access; its permissions are augmented with the
// wildcard so it is authorized for every action without enumerating each grant.
const roleSystemAdmin = "system_admin"

// permissionWildcard, when present in a user's effective permissions, authorizes
// any requested permission (super-admin).
const permissionWildcard = "*"

// RBACServiceImpl implements service.RBACService with Redis caching.
type RBACServiceImpl struct {
	roleRepo repository.RoleRepository
	redis    *redis.Client
}

// NewRBACService creates a new RBACServiceImpl.
func NewRBACService(roleRepo repository.RoleRepository, redisClient *redis.Client) *RBACServiceImpl {
	return &RBACServiceImpl{roleRepo: roleRepo, redis: redisClient}
}

func (s *RBACServiceImpl) cacheKey(userID, orgID uuid.UUID) string {
	return fmt.Sprintf("user:%s:org:%s:permissions", userID, orgID)
}

func (s *RBACServiceImpl) GetUserPermissions(ctx context.Context, userID, orgID uuid.UUID) ([]string, error) {
	// Try cache first
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, s.cacheKey(userID, orgID)).Result()
		if err == nil {
			var perms []string
			if json.Unmarshal([]byte(cached), &perms) == nil {
				return perms, nil
			}
		}
	}

	// Fetch from DB
	perms, err := s.roleRepo.GetUserPermissions(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	// Super-admin (system_admin role) is granted a wildcard so it is authorized
	// for every endpoint regardless of the granular permission catalog.
	if roles, rerr := s.roleRepo.GetUserRoleNames(ctx, userID, orgID); rerr == nil {
		for _, r := range roles {
			if r == roleSystemAdmin {
				perms = append(perms, permissionWildcard)
				break
			}
		}
	}

	// Cache result
	if s.redis != nil && len(perms) > 0 {
		data, _ := json.Marshal(perms)
		_ = s.redis.Set(ctx, s.cacheKey(userID, orgID), data, permissionCacheTTL).Err()
	}

	return perms, nil
}

func (s *RBACServiceImpl) HasPermission(ctx context.Context, userID, orgID uuid.UUID, permission string) (bool, error) {
	perms, err := s.GetUserPermissions(ctx, userID, orgID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		if p == permission || p == permissionWildcard {
			return true, nil
		}
	}
	return false, nil
}

func (s *RBACServiceImpl) HasAnyPermission(ctx context.Context, userID, orgID uuid.UUID, permissions []string) (bool, error) {
	perms, err := s.GetUserPermissions(ctx, userID, orgID)
	if err != nil {
		return false, err
	}
	permSet := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		permSet[p] = struct{}{}
	}
	if _, ok := permSet[permissionWildcard]; ok {
		return true, nil
	}
	for _, required := range permissions {
		if _, ok := permSet[required]; ok {
			return true, nil
		}
	}
	return false, nil
}

func (s *RBACServiceImpl) InvalidateCache(ctx context.Context, userID, orgID uuid.UUID) error {
	if s.redis != nil {
		return s.redis.Del(ctx, s.cacheKey(userID, orgID)).Err()
	}
	return nil
}
