package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/iam/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
)

// ManageUsersUseCase handles user updates and (de)activation (system admin).
type ManageUsersUseCase struct {
	userRepo repository.UserRepository
}

// NewManageUsersUseCase creates a new ManageUsersUseCase.
func NewManageUsersUseCase(userRepo repository.UserRepository) *ManageUsersUseCase {
	return &ManageUsersUseCase{userRepo: userRepo}
}

func toUserInfo(u *model.User) *dto.UserInfo {
	return &dto.UserInfo{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Status:    string(u.Status),
		CreatedAt: u.CreatedAt,
	}
}

// Update applies a partial update to a user (name and/or status).
func (uc *ManageUsersUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserInfo, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Status != "" {
		status := model.UserStatus(req.Status)
		switch status {
		case model.UserStatusActive, model.UserStatusInactive, model.UserStatusSuspended:
			user.Status = status
		default:
			return nil, domainErr.New(domainErr.ErrValidation, "invalid user status", nil)
		}
	}
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return toUserInfo(user), nil
}

// SetStatus activates, deactivates or suspends a user.
func (uc *ManageUsersUseCase) SetStatus(ctx context.Context, id uuid.UUID, status model.UserStatus) (*dto.UserInfo, error) {
	switch status {
	case model.UserStatusActive, model.UserStatusInactive, model.UserStatusSuspended:
	default:
		return nil, domainErr.New(domainErr.ErrValidation, "invalid user status", nil)
	}
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Status = status
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return toUserInfo(user), nil
}
