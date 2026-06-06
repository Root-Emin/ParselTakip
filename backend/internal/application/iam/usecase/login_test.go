package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/iam/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/service"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubUserRepo struct {
	user *model.User
}

func (s *stubUserRepo) Create(context.Context, *model.User) error { return nil }
func (s *stubUserRepo) GetByID(context.Context, uuid.UUID) (*model.User, error) {
	return s.user, nil
}
func (s *stubUserRepo) GetByEmail(context.Context, string) (*model.User, error) {
	return s.user, nil
}
func (s *stubUserRepo) Update(context.Context, *model.User) error { return nil }
func (s *stubUserRepo) Delete(context.Context, uuid.UUID) error   { return nil }
func (s *stubUserRepo) List(context.Context, int, int) ([]*model.User, int, error) {
	return nil, 0, nil
}

type stubOrgUserRepo struct {
	member bool
}

func (s *stubOrgUserRepo) Add(context.Context, *model.OrganizationUser) error { return nil }
func (s *stubOrgUserRepo) Remove(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (s *stubOrgUserRepo) GetByOrgAndUser(context.Context, uuid.UUID, uuid.UUID) (*model.OrganizationUser, error) {
	if !s.member {
		return nil, domainErr.New(domainErr.ErrNotFound, "organization user not found", nil)
	}
	return &model.OrganizationUser{Status: model.OrgUserStatusActive}, nil
}
func (s *stubOrgUserRepo) ListByOrg(context.Context, uuid.UUID, int, int) ([]*model.OrganizationUser, int, error) {
	return nil, 0, nil
}
func (s *stubOrgUserRepo) ListByUser(context.Context, uuid.UUID) ([]*model.OrganizationUser, error) {
	return nil, nil
}

type stubRoleRepo struct {
	roles       []string
	permissions []string
}

func (s *stubRoleRepo) Create(context.Context, *model.Role) error { return nil }
func (s *stubRoleRepo) GetByID(context.Context, uuid.UUID) (*model.Role, error) {
	return nil, nil
}
func (s *stubRoleRepo) ListByScope(context.Context, model.ScopeType, uuid.UUID) ([]*model.Role, error) {
	return nil, nil
}
func (s *stubRoleRepo) Update(context.Context, *model.Role) error { return nil }
func (s *stubRoleRepo) Delete(context.Context, uuid.UUID) error   { return nil }
func (s *stubRoleRepo) AddPermission(context.Context, uuid.UUID, string) error {
	return nil
}
func (s *stubRoleRepo) RemovePermission(context.Context, uuid.UUID, string) error { return nil }
func (s *stubRoleRepo) GetPermissions(context.Context, uuid.UUID) ([]string, error) {
	return nil, nil
}
func (s *stubRoleRepo) AssignRoleToUser(context.Context, *model.UserRole) error { return nil }
func (s *stubRoleRepo) RemoveRoleFromUser(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}
func (s *stubRoleRepo) GetUserRoles(context.Context, uuid.UUID, uuid.UUID) ([]*model.UserRole, error) {
	return nil, nil
}
func (s *stubRoleRepo) GetUserPermissions(context.Context, uuid.UUID, uuid.UUID) ([]string, error) {
	return s.permissions, nil
}
func (s *stubRoleRepo) GetUserRoleNames(context.Context, uuid.UUID, uuid.UUID) ([]string, error) {
	return s.roles, nil
}
func (s *stubRoleRepo) GetByNameAndScope(context.Context, model.ScopeType, uuid.UUID, string) (*model.Role, error) {
	return nil, domainErr.New(domainErr.ErrNotFound, "role not found", nil)
}

type stubRBAC struct {
	permissions []string
}

func (s *stubRBAC) HasPermission(context.Context, uuid.UUID, uuid.UUID, string) (bool, error) {
	return false, nil
}
func (s *stubRBAC) HasAnyPermission(context.Context, uuid.UUID, uuid.UUID, []string) (bool, error) {
	return false, nil
}
func (s *stubRBAC) GetUserPermissions(context.Context, uuid.UUID, uuid.UUID) ([]string, error) {
	return s.permissions, nil
}
func (s *stubRBAC) InvalidateCache(context.Context, uuid.UUID, uuid.UUID) error { return nil }

type stubAuth struct {
	hash string
}

func (s *stubAuth) HashPassword(password string) (string, error) { return s.hash, nil }
func (s *stubAuth) VerifyPassword(hashedPassword, password string) error {
	if hashedPassword != password {
		return domainErr.New(domainErr.ErrUnauthorized, "invalid credentials", nil)
	}
	return nil
}
func (s *stubAuth) GenerateToken(_ context.Context, claims service.TokenClaims) (string, error) {
	if len(claims.Roles) > 0 {
		return "token-with-rbac", nil
	}
	return "token-basic", nil
}
func (s *stubAuth) ValidateToken(context.Context, string) (*service.TokenClaims, error) {
	return nil, nil
}

func TestLoginUseCase_BasicTokenWithoutOrganization(t *testing.T) {
	userID := uuid.New()
	uc := NewLoginUseCase(
		&stubUserRepo{user: &model.User{ID: userID, Email: "user@test.com", PasswordHash: "secret", Status: model.UserStatusActive}},
		&stubOrgUserRepo{},
		&stubRoleRepo{},
		&stubRBAC{},
		&stubAuth{},
		nil,
		24,
	)

	resp, err := uc.Execute(context.Background(), dto.LoginRequest{
		Email:    "user@test.com",
		Password: "secret",
	})
	require.NoError(t, err)
	assert.Equal(t, "token-basic", resp.Token)
	assert.Equal(t, "Bearer", resp.TokenType)
	assert.Empty(t, resp.Roles)
}

func TestLoginUseCase_OrgScopedTokenIncludesRBACClaims(t *testing.T) {
	userID := uuid.New()
	orgID := uuid.New()
	uc := NewLoginUseCase(
		&stubUserRepo{user: &model.User{ID: userID, Email: "admin@test.com", PasswordHash: "secret", Status: model.UserStatusActive}},
		&stubOrgUserRepo{member: true},
		&stubRoleRepo{roles: []string{"municipality_admin"}, permissions: []string{"project.create"}},
		&stubRBAC{permissions: []string{"project.create"}},
		&stubAuth{},
		nil,
		24,
	)

	resp, err := uc.Execute(context.Background(), dto.LoginRequest{
		Email:          "admin@test.com",
		Password:       "secret",
		OrganizationID: &orgID,
	})
	require.NoError(t, err)
	assert.Equal(t, "token-with-rbac", resp.Token)
	assert.Equal(t, orgID, resp.OrganizationID)
	assert.Equal(t, []string{"municipality_admin"}, resp.Roles)
	assert.Equal(t, []string{"project.create"}, resp.Permissions)
}
