package urbantransform

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/command"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/constants"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/query"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"
	"github.com/masterfabric-go/masterfabric/internal/shared/response"
	"github.com/masterfabric-go/masterfabric/internal/shared/validator"
)

// NotificationHandler exposes CQRS-based HTTP endpoints for notifications.
type NotificationHandler struct {
	cmd *command.NotificationCommandHandler
	qry *query.NotificationQueryHandler
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(cmd *command.NotificationCommandHandler, qry *query.NotificationQueryHandler) *NotificationHandler {
	return &NotificationHandler{cmd: cmd, qry: qry}
}

// Create handles POST /notifications (send a notification).
func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreateNotificationRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Create(r.Context(), orgID, appID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgNotificationCreated, result)
}

// List handles GET /notifications (the current user's notifications).
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	q := dto.ListNotificationsQuery{
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if id := parseUUIDQuery(r, constants.QueryKeyNotifProject); id != nil {
		q.ProjectID = id
	}
	if v := r.URL.Query().Get(constants.QueryKeyNotifType); v != "" {
		q.NotificationType = &v
	}
	if v := r.URL.Query().Get(constants.QueryKeyNotifIsRead); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			q.IsRead = &b
		}
	}
	result, err := h.qry.List(r.Context(), orgID, appID, userID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgNotificationListed, result)
}

// UnreadCount handles GET /notifications/unread-count.
func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	count, err := h.qry.UnreadCount(r.Context(), orgID, appID, userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgNotificationUnread, map[string]int{"unread_count": count})
}

// MarkRead handles PATCH /notifications/{notificationId}/read.
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamNotificationID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidNotificationID})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	if err := h.cmd.MarkRead(r.Context(), orgID, appID, userID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgNotificationRead, nil)
}

// MarkAllRead handles POST /notifications/read-all.
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	count, err := h.cmd.MarkAllRead(r.Context(), orgID, appID, userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgNotificationAllRead, map[string]int{"updated": count})
}
