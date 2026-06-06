package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
)

// CreateBuildingUnitRequest is the input for creating a building unit (command).
type CreateBuildingUnitRequest struct {
	BuildingID    uuid.UUID `json:"building_id" validate:"required"`
	UnitNo        string    `json:"unit_no" validate:"required"`
	FloorNo       *int      `json:"floor_no,omitempty"`
	AreaSqm       *float64  `json:"area_sqm,omitempty"`
	RoomCount     string    `json:"room_count"`
	OwnershipType string    `json:"ownership_type" validate:"omitempty,oneof=kat_irtifaki kat_mulkiyeti arsa_payi"`
	TitleDeedNo   string    `json:"title_deed_no"`
}

// UpdateBuildingUnitRequest is the input for updating a building unit (partial command).
type UpdateBuildingUnitRequest struct {
	UnitNo        *string  `json:"unit_no,omitempty"`
	FloorNo       *int     `json:"floor_no,omitempty"`
	AreaSqm       *float64 `json:"area_sqm,omitempty"`
	RoomCount     *string  `json:"room_count,omitempty"`
	OwnershipType *string  `json:"ownership_type,omitempty" validate:"omitempty,oneof=kat_irtifaki kat_mulkiyeti arsa_payi"`
	TitleDeedNo   *string  `json:"title_deed_no,omitempty"`
	Status        *string  `json:"status,omitempty" validate:"omitempty,oneof=active in_transformation transferred archived"`
}

// ListBuildingUnitsQuery is the input for listing/filtering/searching units (query).
type ListBuildingUnitsQuery struct {
	BuildingID    *uuid.UUID
	Status        *string
	OwnershipType *string
	Search        string
	SortBy        string
	SortOrder     string
	Page          int
	PerPage       int
}

// BuildingUnitResponse is the public representation of a building unit (response model).
type BuildingUnitResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AppID          uuid.UUID `json:"app_id"`
	BuildingID     uuid.UUID `json:"building_id"`
	UnitNo         string    `json:"unit_no"`
	FloorNo        *int      `json:"floor_no,omitempty"`
	AreaSqm        *float64  `json:"area_sqm,omitempty"`
	RoomCount      string    `json:"room_count"`
	OwnershipType  string    `json:"ownership_type"`
	TitleDeedNo    string    `json:"title_deed_no"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ToBuildingUnitResponse maps a domain unit to its response model.
func ToBuildingUnitResponse(u *model.BuildingUnit) BuildingUnitResponse {
	return BuildingUnitResponse{
		ID:             u.ID,
		OrganizationID: u.OrganizationID,
		AppID:          u.AppID,
		BuildingID:     u.BuildingID,
		UnitNo:         u.UnitNo,
		FloorNo:        u.FloorNo,
		AreaSqm:        u.AreaSqm,
		RoomCount:      u.RoomCount,
		OwnershipType:  string(u.OwnershipType),
		TitleDeedNo:    u.TitleDeedNo,
		Status:         string(u.Status),
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

// ToBuildingUnitResponseList maps a slice of units to response models.
func ToBuildingUnitResponseList(items []*model.BuildingUnit) []BuildingUnitResponse {
	out := make([]BuildingUnitResponse, 0, len(items))
	for _, u := range items {
		out = append(out, ToBuildingUnitResponse(u))
	}
	return out
}
