package constants

// Property owner query parameter keys.
const (
	QueryKeyOwnerUnit   = "unit_id"
	QueryKeyOwnerStatus = "status"
	QueryKeyOwnerSearch = "q"
)

// Property owner sortable field keys.
const (
	OwnerSortKeyCreatedAt = "created_at"
	OwnerSortKeyUpdatedAt = "updated_at"
	OwnerSortKeyLastName  = "last_name"
	OwnerSortKeyStatus    = "status"
)

// AllowedOwnerSortKeys whitelists ORDER BY columns for property owners.
var AllowedOwnerSortKeys = map[string]struct{}{
	OwnerSortKeyCreatedAt: {},
	OwnerSortKeyUpdatedAt: {},
	OwnerSortKeyLastName:  {},
	OwnerSortKeyStatus:    {},
}

// IsAllowedOwnerSortKey reports whether the key is a whitelisted sort column.
func IsAllowedOwnerSortKey(key string) bool {
	_, ok := AllowedOwnerSortKeys[key]
	return ok
}

// Property owner permission keys (RBAC).
const (
	PermissionOwnerCreate = "property_owner.create"
	PermissionOwnerRead   = "property_owner.read"
	PermissionOwnerUpdate = "property_owner.update"
	PermissionOwnerDelete = "property_owner.delete"
	PermissionOwnerList   = "property_owner.list"
)

// Property owner path parameter key.
const PathParamOwnerID = "ownerId"

// Property owner response message keys.
const (
	MsgOwnerCreated   = "property_owner.created"
	MsgOwnerUpdated   = "property_owner.updated"
	MsgOwnerDeleted   = "property_owner.deleted"
	MsgOwnerFetched   = "property_owner.fetched"
	MsgOwnerListed    = "property_owner.listed"
	MsgInvalidOwnerID = "property_owner.invalid_id"
)
