package constants

// Notification query parameter keys.
const (
	QueryKeyNotifProject = "project_id"
	QueryKeyNotifType    = "notification_type"
	QueryKeyNotifIsRead  = "is_read"
)

// Notification permission keys (RBAC).
const (
	PermissionNotificationCreate = "notification.create"
	PermissionNotificationList   = "notification.list"
	PermissionNotificationRead   = "notification.read"
)

// Notification path parameter key.
const PathParamNotificationID = "notificationId"

// Notification response message keys.
const (
	MsgNotificationCreated   = "notification.created"
	MsgNotificationListed    = "notification.listed"
	MsgNotificationRead      = "notification.read"
	MsgNotificationAllRead   = "notification.all_read"
	MsgNotificationUnread    = "notification.unread_count"
	MsgInvalidNotificationID = "notification.invalid_id"
)
