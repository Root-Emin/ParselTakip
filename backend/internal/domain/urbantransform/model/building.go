package model

import (
	"time"

	"github.com/google/uuid"
)

// RiskStatus represents the seismic/structural risk status of a building.
type RiskStatus string

const (
	RiskStatusUnknown         RiskStatus = "unknown"
	RiskStatusRisky           RiskStatus = "risky"
	RiskStatusNotRisky        RiskStatus = "not_risky"
	RiskStatusUnderAssessment RiskStatus = "under_assessment"
)

// BuildingStatus represents the lifecycle status of a building.
type BuildingStatus string

const (
	BuildingStatusActive           BuildingStatus = "active"
	BuildingStatusInTransformation BuildingStatus = "in_transformation"
	BuildingStatusDemolished       BuildingStatus = "demolished"
	BuildingStatusRebuilt          BuildingStatus = "rebuilt"
	BuildingStatusArchived         BuildingStatus = "archived"
)

// BuildingType represents the usage type of a building.
type BuildingType string

const (
	BuildingTypeResidential BuildingType = "residential"
	BuildingTypeCommercial  BuildingType = "commercial"
	BuildingTypeMixed       BuildingType = "mixed"
)

// IsValidRiskStatus reports whether the given risk status is known.
func IsValidRiskStatus(s RiskStatus) bool {
	switch s {
	case RiskStatusUnknown, RiskStatusRisky, RiskStatusNotRisky, RiskStatusUnderAssessment:
		return true
	default:
		return false
	}
}

// IsValidBuildingStatus reports whether the given status is known.
func IsValidBuildingStatus(s BuildingStatus) bool {
	switch s {
	case BuildingStatusActive, BuildingStatusInTransformation, BuildingStatusDemolished,
		BuildingStatusRebuilt, BuildingStatusArchived:
		return true
	default:
		return false
	}
}

// IsValidBuildingType reports whether the given building type is known.
func IsValidBuildingType(t BuildingType) bool {
	switch t {
	case BuildingTypeResidential, BuildingTypeCommercial, BuildingTypeMixed:
		return true
	default:
		return false
	}
}

// Building represents a physical building subject to urban transformation.
type Building struct {
	ID               uuid.UUID      `json:"id"`
	OrganizationID   uuid.UUID      `json:"organization_id"`
	AppID            uuid.UUID      `json:"app_id"`
	Name             string         `json:"name"`
	Address          string         `json:"address"`
	City             string         `json:"city"`
	District         string         `json:"district"`
	Neighborhood     string         `json:"neighborhood"`
	BlockNo          string         `json:"block_no"`
	ParcelNo         string         `json:"parcel_no"`
	IslandNo         string         `json:"island_no"`
	FloorCount       *int           `json:"floor_count,omitempty"`
	UnitCount        int            `json:"unit_count"`
	ConstructionYear *int           `json:"construction_year,omitempty"`
	BuildingType     BuildingType   `json:"building_type"`
	RiskStatus       RiskStatus     `json:"risk_status"`
	Status           BuildingStatus `json:"status"`
	CreatedBy        *uuid.UUID     `json:"created_by,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
