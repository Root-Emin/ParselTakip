package constants

// Contractor query parameter keys.
const (
	QueryKeyContractorStatus = "status"
	QueryKeyContractorSearch = "q"
)

// Contractor sortable field keys.
const (
	ContractorSortKeyCreatedAt   = "created_at"
	ContractorSortKeyUpdatedAt   = "updated_at"
	ContractorSortKeyCompanyName = "company_name"
	ContractorSortKeyStatus      = "status"
)

// AllowedContractorSortKeys whitelists ORDER BY columns for contractors.
var AllowedContractorSortKeys = map[string]struct{}{
	ContractorSortKeyCreatedAt:   {},
	ContractorSortKeyUpdatedAt:   {},
	ContractorSortKeyCompanyName: {},
	ContractorSortKeyStatus:      {},
}

// IsAllowedContractorSortKey reports whether the key is a whitelisted sort column.
func IsAllowedContractorSortKey(key string) bool {
	_, ok := AllowedContractorSortKeys[key]
	return ok
}

// Contractor permission keys (RBAC).
const (
	PermissionContractorCreate = "contractor.create"
	PermissionContractorRead   = "contractor.read"
	PermissionContractorUpdate = "contractor.update"
	PermissionContractorDelete = "contractor.delete"
	PermissionContractorList   = "contractor.list"
)

// Contractor path parameter key.
const PathParamContractorID = "contractorId"

// Contractor response message keys.
const (
	MsgContractorCreated   = "contractor.created"
	MsgContractorUpdated   = "contractor.updated"
	MsgContractorDeleted   = "contractor.deleted"
	MsgContractorFetched   = "contractor.fetched"
	MsgContractorListed    = "contractor.listed"
	MsgInvalidContractorID = "contractor.invalid_id"
)
