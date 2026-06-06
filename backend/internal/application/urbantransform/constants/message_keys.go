package constants

// Message query parameter keys.
const (
	QueryKeyMsgBox     = "box"
	QueryKeyMsgProject = "project_id"
	QueryKeyMsgIsRead  = "is_read"
)

// Message permission keys (RBAC).
const (
	PermissionMessageCreate = "message.create"
	PermissionMessageList   = "message.list"
	PermissionMessageRead   = "message.read"
)

// Message path parameter key.
const PathParamMessageID = "messageId"

// Message response message keys.
const (
	MsgMessageCreated   = "message.created"
	MsgMessageListed    = "message.listed"
	MsgMessageFetched   = "message.fetched"
	MsgMessageRead      = "message.read"
	MsgMessageUnread    = "message.unread_count"
	MsgInvalidMessageID = "message.invalid_id"
)
