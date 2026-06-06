package constants

// Building unit query parameter keys.
const (
	QueryKeyUnitBuilding      = "building_id"
	QueryKeyUnitStatus        = "status"
	QueryKeyUnitOwnershipType = "ownership_type"
	QueryKeyUnitSearch        = "q"
)

// Building unit sortable field keys.
const (
	UnitSortKeyCreatedAt = "created_at"
	UnitSortKeyUpdatedAt = "updated_at"
	UnitSortKeyUnitNo    = "unit_no"
	UnitSortKeyStatus    = "status"
)

// AllowedUnitSortKeys whitelists ORDER BY columns for building units.
var AllowedUnitSortKeys = map[string]struct{}{
	UnitSortKeyCreatedAt: {},
	UnitSortKeyUpdatedAt: {},
	UnitSortKeyUnitNo:    {},
	UnitSortKeyStatus:    {},
}

// IsAllowedUnitSortKey reports whether the key is a whitelisted sort column.
func IsAllowedUnitSortKey(key string) bool {
	_, ok := AllowedUnitSortKeys[key]
	return ok
}

// Building unit permission keys (RBAC).
const (
	PermissionUnitCreate = "building_unit.create"
	PermissionUnitRead   = "building_unit.read"
	PermissionUnitUpdate = "building_unit.update"
	PermissionUnitDelete = "building_unit.delete"
	PermissionUnitList   = "building_unit.list"
)

// Building unit path parameter key.
const PathParamUnitID = "unitId"

// Building unit response message keys.
const (
	MsgUnitCreated   = "building_unit.created"
	MsgUnitUpdated   = "building_unit.updated"
	MsgUnitDeleted   = "building_unit.deleted"
	MsgUnitFetched   = "building_unit.fetched"
	MsgUnitListed    = "building_unit.listed"
	MsgInvalidUnitID = "building_unit.invalid_id"
)
