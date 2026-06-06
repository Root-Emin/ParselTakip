package constants

// Appointment query parameter keys.
const (
	QueryKeyApptProject  = "project_id"
	QueryKeyApptOwner    = "owner_id"
	QueryKeyApptStatus   = "status"
	QueryKeyApptUpcoming = "upcoming"
)

// Appointment permission keys (RBAC).
const (
	PermissionAppointmentCreate = "appointment.create"
	PermissionAppointmentRead   = "appointment.read"
	PermissionAppointmentUpdate = "appointment.update"
	PermissionAppointmentDelete = "appointment.delete"
	PermissionAppointmentList   = "appointment.list"
)

// Appointment path parameter key.
const PathParamAppointmentID = "appointmentId"

// Appointment response message keys.
const (
	MsgAppointmentCreated   = "appointment.created"
	MsgAppointmentUpdated   = "appointment.updated"
	MsgAppointmentDeleted   = "appointment.deleted"
	MsgAppointmentFetched   = "appointment.fetched"
	MsgAppointmentListed    = "appointment.listed"
	MsgInvalidAppointmentID = "appointment.invalid_id"
)
