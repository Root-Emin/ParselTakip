package constants

// Document query parameter keys.
const (
	QueryKeyDocProject = "project_id"
	QueryKeyDocBuilding = "building_id"
	QueryKeyDocOwner   = "owner_id"
	QueryKeyDocType    = "document_type_id"
	QueryKeyDocStatus  = "status"
	QueryKeyDocSearch  = "q"
	QueryKeyDocCategory = "category"
)

// Document sortable field keys.
const (
	DocSortKeyCreatedAt = "created_at"
	DocSortKeyUpdatedAt = "updated_at"
	DocSortKeyFileName  = "file_name"
	DocSortKeyStatus    = "status"
)

// AllowedDocSortKeys whitelists ORDER BY columns for documents.
var AllowedDocSortKeys = map[string]struct{}{
	DocSortKeyCreatedAt: {},
	DocSortKeyUpdatedAt: {},
	DocSortKeyFileName:  {},
	DocSortKeyStatus:    {},
}

// IsAllowedDocSortKey reports whether the key is a whitelisted sort column.
func IsAllowedDocSortKey(key string) bool {
	_, ok := AllowedDocSortKeys[key]
	return ok
}

// Document permission keys (RBAC).
const (
	PermissionDocumentCreate = "document.create"
	PermissionDocumentRead   = "document.read"
	PermissionDocumentUpdate = "document.update"
	PermissionDocumentDelete = "document.delete"
	PermissionDocumentList   = "document.list"
	PermissionDocumentReview = "document.review"
)

// Document path parameter key.
const PathParamDocumentID = "documentId"

// Document response message keys.
const (
	MsgDocumentCreated   = "document.created"
	MsgDocumentUpdated   = "document.updated"
	MsgDocumentDeleted   = "document.deleted"
	MsgDocumentFetched   = "document.fetched"
	MsgDocumentListed    = "document.listed"
	MsgDocumentReviewed  = "document.reviewed"
	MsgInvalidDocumentID = "document.invalid_id"
	MsgDocumentTypesListed = "document_type.listed"
	MsgReviewsListed     = "document_review.listed"
)
