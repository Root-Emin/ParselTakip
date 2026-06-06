package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
)

// KVKKConsentRequest records/updates the data subject's consent to processing
// of their personal data (KVKK açık rıza).
type KVKKConsentRequest struct {
	Version string `json:"version" validate:"required"`
}

// KVKKConsentResponse confirms a recorded consent.
type KVKKConsentResponse struct {
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	ConsentAt time.Time `json:"consent_at"`
}

// KVKKExportResponse is the full personal-data export for the data subject
// themselves (KVKK Madde 11 - bilgi talep etme / veri taşınabilirliği). PII is
// returned UNMASKED here because the requester IS the data owner, authenticated
// as themselves.
type KVKKExportResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	Status             string     `json:"status"`
	UserType           string     `json:"user_type"`
	Phone              string     `json:"phone,omitempty"`
	NationalID         string     `json:"national_id,omitempty"`
	Address            string     `json:"address,omitempty"`
	City               string     `json:"city,omitempty"`
	District           string     `json:"district,omitempty"`
	CompanyID          *uuid.UUID `json:"company_id,omitempty"`
	KVKKConsentAt      *time.Time `json:"kvkk_consent_at,omitempty"`
	KVKKConsentVersion string     `json:"kvkk_consent_version,omitempty"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ExportedAt         time.Time  `json:"exported_at"`
}

// ToKVKKExportResponse maps the domain user to a full self-service export.
func ToKVKKExportResponse(u *model.User) *KVKKExportResponse {
	return &KVKKExportResponse{
		ID:                 u.ID,
		Email:              u.Email,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Status:             string(u.Status),
		UserType:           string(u.UserType),
		Phone:              u.Phone,
		NationalID:         u.NationalID,
		Address:            u.Address,
		City:               u.City,
		District:           u.District,
		CompanyID:          u.CompanyID,
		KVKKConsentAt:      u.KVKKConsentAt,
		KVKKConsentVersion: u.KVKKConsentVersion,
		LastLoginAt:        u.LastLoginAt,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		ExportedAt:         time.Now().UTC(),
	}
}
