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

// MessageHandler exposes CQRS-based HTTP endpoints for messages.
type MessageHandler struct {
	cmd *command.MessageCommandHandler
	qry *query.MessageQueryHandler
}

// NewMessageHandler creates a new MessageHandler.
func NewMessageHandler(cmd *command.MessageCommandHandler, qry *query.MessageQueryHandler) *MessageHandler {
	return &MessageHandler{cmd: cmd, qry: qry}
}

// Create handles POST /messages (send).
func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreateMessageRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	senderID, _ := middleware.UserIDFromContext(r.Context())
	result, err := h.cmd.Create(r.Context(), orgID, appID, senderID, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgMessageCreated, result)
}

// List handles GET /messages (?box=inbox|sent).
func (h *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	q := dto.ListMessagesQuery{
		Box:     r.URL.Query().Get(constants.QueryKeyMsgBox),
		Page:    queryInt(r, constants.QueryKeyPage),
		PerPage: queryInt(r, constants.QueryKeyPerPage),
	}
	if id := parseUUIDQuery(r, constants.QueryKeyMsgProject); id != nil {
		q.ProjectID = id
	}
	if v := r.URL.Query().Get(constants.QueryKeyMsgIsRead); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			q.IsRead = &b
		}
	}
	result, err := h.qry.List(r.Context(), orgID, appID, userID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgMessageListed, result)
}

// UnreadCount handles GET /messages/unread-count.
func (h *MessageHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
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
	response.Success(w, constants.MsgMessageUnread, map[string]int{"unread_count": count})
}

// Get handles GET /messages/{messageId}.
func (h *MessageHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamMessageID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidMessageID})
		return
	}
	result, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgMessageFetched, result)
}

// MarkRead handles PATCH /messages/{messageId}/read.
func (h *MessageHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamMessageID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidMessageID})
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	if err := h.cmd.MarkRead(r.Context(), orgID, appID, userID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgMessageRead, nil)
}
