package constants

// Building query parameter keys.
const (
	QueryKeyBuildingStatus     = "status"
	QueryKeyBuildingRiskStatus = "risk_status"
	QueryKeyBuildingType       = "building_type"
	QueryKeyBuildingCity       = "city"
	QueryKeyBuildingDistrict   = "district"
	QueryKeyBuildingSearch     = "q"
)

// Building sortable field keys.
const (
	BuildingSortKeyCreatedAt = "created_at"
	BuildingSortKeyUpdatedAt = "updated_at"
	BuildingSortKeyName      = "name"
	BuildingSortKeyCity      = "city"
	BuildingSortKeyStatus    = "status"
)

// AllowedBuildingSortKeys whitelists ORDER BY columns for buildings.
var AllowedBuildingSortKeys = map[string]struct{}{
	BuildingSortKeyCreatedAt: {},
	BuildingSortKeyUpdatedAt: {},
	BuildingSortKeyName:      {},
	BuildingSortKeyCity:      {},
	BuildingSortKeyStatus:    {},
}

// IsAllowedBuildingSortKey reports whether the key is a whitelisted sort column.
func IsAllowedBuildingSortKey(key string) bool {
	_, ok := AllowedBuildingSortKeys[key]
	return ok
}

// Building permission keys (RBAC).
const (
	PermissionBuildingCreate = "building.create"
	PermissionBuildingRead   = "building.read"
	PermissionBuildingUpdate = "building.update"
	PermissionBuildingDelete = "building.delete"
	PermissionBuildingList   = "building.list"
)

// Building path parameter key.
const PathParamBuildingID = "buildingId"

// Building response message keys.
const (
	MsgBuildingCreated   = "building.created"
	MsgBuildingUpdated   = "building.updated"
	MsgBuildingDeleted   = "building.deleted"
	MsgBuildingFetched   = "building.fetched"
	MsgBuildingListed    = "building.listed"
	MsgInvalidBuildingID = "building.invalid_id"
)
