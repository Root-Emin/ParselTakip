package constants

// Workflow permission keys (RBAC).
const (
	PermissionWorkflowRead    = "workflow.read"
	PermissionWorkflowAdvance = "workflow.advance"
)

// Workflow response message keys.
const (
	MsgWorkflowStepsListed   = "workflow.steps_listed"
	MsgWorkflowStatesListed  = "workflow.states_listed"
	MsgWorkflowHistoryListed = "workflow.history_listed"
	MsgWorkflowAdvanced      = "workflow.advanced"
)
