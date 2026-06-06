package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	"github.com/masterfabric-go/masterfabric/internal/shared/pii"
)

// RegisterRequest is the input for user registration.
type RegisterRequest struct {
	Email       string     `json:"email" validate:"required,email"`
	Password    string     `json:"password" validate:"required,min=8"`
	FirstName   string     `json:"first_name" validate:"required"`
	LastName    string     `json:"last_name" validate:"required"`
	Phone       string     `json:"phone,omitempty"`
	NationalID  string     `json:"national_id,omitempty"` // TCKN; validated when present
	UserType    string     `json:"user_type,omitempty"`   // defaults to "citizen"
	Address     string     `json:"address,omitempty"`
	City        string     `json:"city,omitempty"`
	District    string     `json:"district,omitempty"`
	CompanyID   *uuid.UUID `json:"company_id,omitempty"`
	KVKKConsent bool       `json:"kvkk_consent,omitempty"` // required to persist sensitive PII
}

// LoginRequest is the input for user login.
type LoginRequest struct {
	Email          string     `json:"email" validate:"required,email"`
	Password       string     `json:"password" validate:"required"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
}

// LoginResponse is the output for successful login.
type LoginResponse struct {
	Token            string    `json:"token"`
	TokenType        string    `json:"token_type"`
	ExpiresInHours   int       `json:"expires_in_hours"`
	User             UserInfo  `json:"user"`
	OrganizationID   uuid.UUID `json:"organization_id,omitempty"`
	Roles            []string  `json:"roles,omitempty"`
	Permissions      []string  `json:"permissions,omitempty"`
}

// MeResponse is the authenticated user profile with JWT context.
type MeResponse struct {
	User           UserInfo  `json:"user"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	Roles          []string  `json:"roles,omitempty"`
	Permissions    []string  `json:"permissions,omitempty"`
}

// UserInfo is a self-facing user representation (the caller's own profile).
type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    string    `json:"status"`
	UserType  string    `json:"user_type,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	City      string    `json:"city,omitempty"`
	District  string    `json:"district,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// AdminUserResponse is returned by admin endpoints listing other users; sensitive
// PII (TCKN, phone) is masked for KVKK data-minimization.
type AdminUserResponse struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Status         string     `json:"status"`
	UserType       string     `json:"user_type"`
	PhoneMasked    string     `json:"phone_masked,omitempty"`
	NationalMasked string     `json:"national_id_masked,omitempty"`
	City           string     `json:"city,omitempty"`
	District       string     `json:"district,omitempty"`
	CompanyID      *uuid.UUID `json:"company_id,omitempty"`
	KVKKConsentAt  *time.Time `json:"kvkk_consent_at,omitempty"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ToAdminUserResponse maps a user to a masked admin representation.
func ToAdminUserResponse(u *model.User) *AdminUserResponse {
	resp := &AdminUserResponse{
		ID:            u.ID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Status:        string(u.Status),
		UserType:      string(u.UserType),
		City:          u.City,
		District:      u.District,
		CompanyID:     u.CompanyID,
		KVKKConsentAt: u.KVKKConsentAt,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
	}
	if u.Phone != "" {
		resp.PhoneMasked = pii.MaskPhone(u.Phone)
	}
	if u.NationalID != "" {
		resp.NationalMasked = pii.MaskTCKN(u.NationalID)
	}
	return resp
}

// AssignRoleRequest is the input for assigning a role to a user.
type AssignRoleRequest struct {
	UserID         uuid.UUID  `json:"user_id" validate:"required"`
	RoleID         uuid.UUID  `json:"role_id" validate:"required"`
	OrganizationID uuid.UUID  `json:"organization_id" validate:"required"`
	AppID          *uuid.UUID `json:"app_id,omitempty"`
}

// UpdateUserRequest is the input for updating a user (admin or self profile).
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	District  string `json:"district"`
	UserType  string `json:"user_type"`
}