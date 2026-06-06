// Package constants holds dedicated key definitions for the urban transformation
// module: query parameter keys, sortable field keys, permission keys and
// response message keys. Handlers and use cases reference these keys instead of
// hard-coded strings so the contract stays consistent across layers.
package constants

// Query parameter keys (read from URL query string for list/search/filter).
const (
	QueryKeyStatus     = "status"
	QueryKeyContractor = "contractor_id"
	QueryKeyInitiator  = "initiated_by"
	QueryKeyStep       = "workflow_step_id"
	QueryKeySearch     = "q"
	QueryKeySortBy     = "sort_by"
	QueryKeySortOrder  = "sort_order"
	QueryKeyPage       = "page"
	QueryKeyPerPage    = "per_page"
)

// Sortable field keys (whitelist of columns allowed in ORDER BY).
const (
	SortKeyCreatedAt = "created_at"
	SortKeyUpdatedAt = "updated_at"
	SortKeyName      = "name"
	SortKeyCode      = "code"
	SortKeyStatus    = "status"
)

// Sort order values.
const (
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// DefaultSortBy and DefaultSortOrder define the fallback ordering.
const (
	DefaultSortBy    = SortKeyCreatedAt
	DefaultSortOrder = SortOrderDesc
)

// AllowedSortKeys is the set of columns permitted in ORDER BY clauses.
var AllowedSortKeys = map[string]struct{}{
	SortKeyCreatedAt: {},
	SortKeyUpdatedAt: {},
	SortKeyName:      {},
	SortKeyCode:      {},
	SortKeyStatus:    {},
}

// IsAllowedSortKey reports whether the given key is a whitelisted sort column.
func IsAllowedSortKey(key string) bool {
	_, ok := AllowedSortKeys[key]
	return ok
}

// NormalizeSortOrder returns a safe, lowercase sort order, defaulting to desc.
func NormalizeSortOrder(order string) string {
	if order == SortOrderAsc {
		return SortOrderAsc
	}
	return SortOrderDesc
}

// Permission keys (RBAC) for project operations.
const (
	PermissionProjectCreate = "project.create"
	PermissionProjectRead   = "project.read"
	PermissionProjectUpdate = "project.update"
	PermissionProjectDelete = "project.delete"
	PermissionProjectList   = "project.list"
)

// URL path parameter keys (chi route params).
const (
	PathParamProjectID = "projectId"
)

// Response message keys (stable, client-facing message identifiers).
const (
	MsgProjectCreated   = "project.created"
	MsgProjectUpdated   = "project.updated"
	MsgProjectDeleted   = "project.deleted"
	MsgProjectFetched   = "project.fetched"
	MsgProjectListed    = "project.listed"
	MsgProjectNotFound  = "project.not_found"
	MsgProjectExists    = "project.already_exists"
	MsgInvalidProjectID = "project.invalid_id"
)
