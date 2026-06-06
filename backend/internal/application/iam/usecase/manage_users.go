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

// Update applies a partial update to a user (profile, type and/or status).
func (uc *ManageUsersUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.AdminUserResponse, error) {
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
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.District != "" {
		user.District = req.District
	}
	if req.UserType != "" {
		ut := model.UserType(req.UserType)
		if !model.IsValidUserType(ut) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid user type", nil)
		}
		user.UserType = ut
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
	return dto.ToAdminUserResponse(user), nil
}

// SetStatus activates, deactivates or suspends a user.
func (uc *ManageUsersUseCase) SetStatus(ctx context.Context, id uuid.UUID, status model.UserStatus) (*dto.AdminUserResponse, error) {
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
	return dto.ToAdminUserResponse(user), nil
}
