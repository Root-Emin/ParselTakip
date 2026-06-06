package usecase

import (
	"context"
	"time"

	"github.com/masterfabric-go/masterfabric/internal/application/iam/dto"
	iamEvent "github.com/masterfabric-go/masterfabric/internal/domain/iam/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/repository"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/service"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
	"github.com/masterfabric-go/masterfabric/internal/shared/pii"
)

// kvkkConsentVersion is the current consent text version recorded with the user.
const kvkkConsentVersion = "1.0"

// RegisterUseCase handles user registration.
type RegisterUseCase struct {
	userRepo repository.UserRepository
	auth     service.AuthService
	eventBus events.EventBus
}

// NewRegisterUseCase creates a new RegisterUseCase.
func NewRegisterUseCase(userRepo repository.UserRepository, auth service.AuthService, eventBus events.EventBus) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo, auth: auth, eventBus: eventBus}
}

// Execute registers a new user.
func (uc *RegisterUseCase) Execute(ctx context.Context, req dto.RegisterRequest) (*dto.UserInfo, error) {
	// Check if user already exists
	existing, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, domainErr.New(domainErr.ErrAlreadyExists, "user with this email already exists", nil)
	}

	// Validate TCKN when provided (KVKK: only collect valid, lawful PII).
	if req.NationalID != "" {
		if !pii.IsValidTCKN(req.NationalID) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid national id (TCKN)", nil)
		}
		if !req.KVKKConsent {
			return nil, domainErr.New(domainErr.ErrValidation, "KVKK consent is required to store national id", nil)
		}
	}

	// Resolve and validate user type (defaults to citizen).
	userType := model.UserType(req.UserType)
	if req.UserType == "" {
		userType = model.UserTypeCitizen
	} else if !model.IsValidUserType(userType) {
		return nil, domainErr.New(domainErr.ErrValidation, "invalid user type", nil)
	}

	// Hash password
	hash, err := uc.auth.HashPassword(req.Password)
	if err != nil {
		return nil, domainErr.New(domainErr.ErrInternal, "failed to hash password", err)
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: hash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Status:       model.UserStatusActive,
		UserType:     userType,
		Phone:        req.Phone,
		NationalID:   req.NationalID,
		Address:      req.Address,
		City:         req.City,
		District:     req.District,
		CompanyID:    req.CompanyID,
	}
	if req.KVKKConsent {
		now := time.Now().UTC()
		user.KVKKConsentAt = &now
		user.KVKKConsentVersion = kvkkConsentVersion
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Publish domain event to Kafka
	_ = uc.eventBus.Publish(ctx, events.TopicIAM, iamEvent.UserRegistered{
		UserID:    user.ID,
		Email:     user.Email,
		Timestamp: time.Now().UTC(),
	})

	return &dto.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    string(user.Status),
		UserType:  string(user.UserType),
		Phone:     user.Phone,
		City:      user.City,
		District:  user.District,
		CreatedAt: user.CreatedAt,
	}, nil
}
