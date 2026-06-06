package event

import (
	"time"

	"github.com/google/uuid"
)

// EntityEvent is a generic lifecycle event for urban transformation entities.
// EntityType identifies the resource (e.g. "contractor", "building", "building_unit",
// "property_owner") and Action the operation ("created", "updated", "deleted").
type EntityEvent struct {
	EntityType     string    `json:"entity_type"`
	Action         string    `json:"action"`
	EntityID       uuid.UUID `json:"entity_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AppID          uuid.UUID `json:"app_id"`
	Timestamp      time.Time `json:"timestamp"`
}

// Entity type constants for EntityEvent.EntityType.
const (
	EntityTypeContractor    = "contractor"
	EntityTypeBuilding      = "building"
	EntityTypeBuildingUnit  = "building_unit"
	EntityTypePropertyOwner = "property_owner"
	EntityTypeDocument      = "document"
	EntityTypeApproval      = "approval"
	EntityTypeNotification  = "notification"
)

// Action constants for EntityEvent.Action.
const (
	ActionCreated  = "created"
	ActionUpdated  = "updated"
	ActionDeleted  = "deleted"
	ActionReviewed = "reviewed"
	ActionDecided  = "decided"
	ActionRead     = "read"
)

// NewEntityEvent builds an EntityEvent stamped with the current UTC time.
func NewEntityEvent(entityType, action string, entityID, orgID, appID uuid.UUID) EntityEvent {
	return EntityEvent{
		EntityType:     entityType,
		Action:         action,
		EntityID:       entityID,
		OrganizationID: orgID,
		AppID:          appID,
		Timestamp:      time.Now().UTC(),
	}
}
