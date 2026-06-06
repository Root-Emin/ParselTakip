package model

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account.
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// UserType classifies the kind of platform user, aligned with the system roles.
type UserType string

const (
	UserTypeCitizen        UserType = "citizen"
	UserTypeContractor     UserType = "contractor"
	UserTypeMunicipalStaff UserType = "municipality_staff"
	UserTypeMunicipalAdmin UserType = "municipality_admin"
	UserTypeSystemAdmin    UserType = "system_admin"
)

// IsValidUserType reports whether the given user type is known.
func IsValidUserType(t UserType) bool {
	switch t {
	case UserTypeCitizen, UserTypeContractor, UserTypeMunicipalStaff, UserTypeMunicipalAdmin, UserTypeSystemAdmin:
		return true
	default:
		return false
	}
}

// User represents a platform user entity. Sensitive fields (NationalID) are held
// in plaintext in memory but persisted encrypted (AES-256-GCM) at rest.
type User struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"-"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	Status             UserStatus `json:"status"`
	UserType           UserType   `json:"user_type"`
	Phone              string     `json:"phone,omitempty"`
	NationalID         string     `json:"-"` // TCKN, plaintext in memory, encrypted at rest
	Address            string     `json:"address,omitempty"`
	City               string     `json:"city,omitempty"`
	District           string     `json:"district,omitempty"`
	CompanyID          *uuid.UUID `json:"company_id,omitempty"`
	KVKKConsentAt      *time.Time `json:"kvkk_consent_at,omitempty"`
	KVKKConsentVersion string     `json:"kvkk_consent_version,omitempty"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	FailedLoginCount   int        `json:"-"`
	LockedUntil        *time.Time `json:"-"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// IsLocked reports whether the account is currently locked due to failed logins.
func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now().UTC())
}

// FullName returns the user's full name.
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Email
	}
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user account is active.
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}
