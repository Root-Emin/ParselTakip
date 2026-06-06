package constants

// Approval query parameter keys.
const (
	QueryKeyApprovalProject  = "project_id"
	QueryKeyApprovalType     = "approval_type"
	QueryKeyApprovalStatus   = "status"
	QueryKeyApprovalApprover = "approver_id"
	QueryKeyApprovalOwner    = "owner_id"
)

// Approval sortable field keys.
const (
	ApprovalSortKeyCreatedAt = "created_at"
	ApprovalSortKeyUpdatedAt = "updated_at"
	ApprovalSortKeyStatus    = "status"
)

// AllowedApprovalSortKeys whitelists ORDER BY columns for approvals.
var AllowedApprovalSortKeys = map[string]struct{}{
	ApprovalSortKeyCreatedAt: {},
	ApprovalSortKeyUpdatedAt: {},
	ApprovalSortKeyStatus:    {},
}

// IsAllowedApprovalSortKey reports whether the key is a whitelisted sort column.
func IsAllowedApprovalSortKey(key string) bool {
	_, ok := AllowedApprovalSortKeys[key]
	return ok
}

// Approval permission keys (RBAC).
const (
	PermissionApprovalCreate = "approval.create"
	PermissionApprovalRead   = "approval.read"
	PermissionApprovalDelete = "approval.delete"
	PermissionApprovalList   = "approval.list"
	PermissionApprovalDecide = "approval.decide"
)

// Approval path parameter key.
const PathParamApprovalID = "approvalId"

// Approval response message keys.
const (
	MsgApprovalCreated   = "approval.created"
	MsgApprovalDecided   = "approval.decided"
	MsgApprovalDeleted   = "approval.deleted"
	MsgApprovalFetched   = "approval.fetched"
	MsgApprovalListed    = "approval.listed"
	MsgInvalidApprovalID = "approval.invalid_id"
)
