package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateBuildingRequest is the input for creating a building (command).
type CreateBuildingRequest struct {
	Name             string `json:"name"`
	Address          string `json:"address" validate:"required"`
	City             string `json:"city" validate:"required"`
	District         string `json:"district" validate:"required"`
	Neighborhood     string `json:"neighborhood"`
	BlockNo          string `json:"block_no"`
	ParcelNo         string `json:"parcel_no"`
	IslandNo         string `json:"island_no"`
	FloorCount       *int   `json:"floor_count,omitempty"`
	UnitCount        int    `json:"unit_count" validate:"min=1"`
	ConstructionYear *int   `json:"construction_year,omitempty"`
	BuildingType     string `json:"building_type" validate:"omitempty,oneof=residential commercial mixed"`
	RiskStatus       string `json:"risk_status" validate:"omitempty,oneof=unknown risky not_risky under_assessment"`
}

// UpdateBuildingRequest is the input for updating a building (partial command).
type UpdateBuildingRequest struct {
	Name             *string `json:"name,omitempty"`
	Address          *string `json:"address,omitempty"`
	City             *string `json:"city,omitempty"`
	District         *string `json:"district,omitempty"`
	Neighborhood     *string `json:"neighborhood,omitempty"`
	BlockNo          *string `json:"block_no,omitempty"`
	ParcelNo         *string `json:"parcel_no,omitempty"`
	IslandNo         *string `json:"island_no,omitempty"`
	FloorCount       *int    `json:"floor_count,omitempty"`
	UnitCount        *int    `json:"unit_count,omitempty" validate:"omitempty,min=1"`
	ConstructionYear *int    `json:"construction_year,omitempty"`
	BuildingType     *string `json:"building_type,omitempty" validate:"omitempty,oneof=residential commercial mixed"`
	RiskStatus       *string `json:"risk_status,omitempty" validate:"omitempty,oneof=unknown risky not_risky under_assessment"`
	Status           *string `json:"status,omitempty" validate:"omitempty,oneof=active in_transformation demolished rebuilt archived"`
}

// ListBuildingsQuery is the input for listing/filtering/searching buildings (query).
type ListBuildingsQuery struct {
	Status       *string
	RiskStatus   *string
	BuildingType *string
	City         string
	District     string
	Search       string
	SortBy       string
	SortOrder    string
	Page         int
	PerPage      int
}

// BuildingResponse is the public representation of a building (response model).
type BuildingResponse struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	AppID            uuid.UUID  `json:"app_id"`
	Name             string     `json:"name"`
	Address          string     `json:"address"`
	City             string     `json:"city"`
	District         string     `json:"district"`
	Neighborhood     string     `json:"neighborhood"`
	BlockNo          string     `json:"block_no"`
	ParcelNo         string     `json:"parcel_no"`
	IslandNo         string     `json:"island_no"`
	FloorCount       *int       `json:"floor_count,omitempty"`
	UnitCount        int        `json:"unit_count"`
	ConstructionYear *int       `json:"construction_year,omitempty"`
	BuildingType     string     `json:"building_type"`
	RiskStatus       string     `json:"risk_status"`
	Status           string     `json:"status"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ToBuildingResponse maps a domain building to its response model.
func ToBuildingResponse(b *model.Building) BuildingResponse {
	return BuildingResponse{
		ID:               b.ID,
		OrganizationID:   b.OrganizationID,
		AppID:            b.AppID,
		Name:             b.Name,
		Address:          b.Address,
		City:             b.City,
		District:         b.District,
		Neighborhood:     b.Neighborhood,
		BlockNo:          b.BlockNo,
		ParcelNo:         b.ParcelNo,
		IslandNo:         b.IslandNo,
		FloorCount:       b.FloorCount,
		UnitCount:        b.UnitCount,
		ConstructionYear: b.ConstructionYear,
		BuildingType:     string(b.BuildingType),
		RiskStatus:       string(b.RiskStatus),
		Status:           string(b.Status),
		CreatedBy:        b.CreatedBy,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
	}
}

// ToBuildingResponseList maps a slice of buildings to response models.
func ToBuildingResponseList(items []*model.Building) []BuildingResponse {
	out := make([]BuildingResponse, 0, len(items))
	for _, b := range items {
		out = append(out, ToBuildingResponse(b))
	}
	return out
}
