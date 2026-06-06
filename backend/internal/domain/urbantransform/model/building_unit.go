package model

import (
	"time"

	"github.com/google/uuid"
)

// OwnershipType represents the title/ownership type of a building unit.
type OwnershipType string

const (
	OwnershipTypeKatIrtifaki  OwnershipType = "kat_irtifaki"
	OwnershipTypeKatMulkiyeti OwnershipType = "kat_mulkiyeti"
	OwnershipTypeArsaPayi     OwnershipType = "arsa_payi"
)

// UnitStatus represents the lifecycle status of a building unit.
type UnitStatus string

const (
	UnitStatusActive           UnitStatus = "active"
	UnitStatusInTransformation UnitStatus = "in_transformation"
	UnitStatusTransferred      UnitStatus = "transferred"
	UnitStatusArchived         UnitStatus = "archived"
)

// IsValidOwnershipType reports whether the given ownership type is known.
func IsValidOwnershipType(t OwnershipType) bool {
	switch t {
	case OwnershipTypeKatIrtifaki, OwnershipTypeKatMulkiyeti, OwnershipTypeArsaPayi:
		return true
	default:
		return false
	}
}

// IsValidUnitStatus reports whether the given status is known.
func IsValidUnitStatus(s UnitStatus) bool {
	switch s {
	case UnitStatusActive, UnitStatusInTransformation, UnitStatusTransferred, UnitStatusArchived:
		return true
	default:
		return false
	}
}

// BuildingUnit is an independent section (bağımsız bölüm) of a building.
type BuildingUnit struct {
	ID             uuid.UUID     `json:"id"`
	OrganizationID uuid.UUID     `json:"organization_id"`
	AppID          uuid.UUID     `json:"app_id"`
	BuildingID     uuid.UUID     `json:"building_id"`
	UnitNo         string        `json:"unit_no"`
	FloorNo        *int          `json:"floor_no,omitempty"`
	AreaSqm        *float64      `json:"area_sqm,omitempty"`
	RoomCount      string        `json:"room_count"`
	OwnershipType  OwnershipType `json:"ownership_type"`
	TitleDeedNo    string        `json:"title_deed_no"`
	Status         UnitStatus    `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
