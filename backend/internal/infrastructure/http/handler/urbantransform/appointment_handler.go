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

// AppointmentHandler exposes CQRS-based HTTP endpoints for appointments.
type AppointmentHandler struct {
	cmd *command.AppointmentCommandHandler
	qry *query.AppointmentQueryHandler
}

// NewAppointmentHandler creates a new AppointmentHandler.
func NewAppointmentHandler(cmd *command.AppointmentCommandHandler, qry *query.AppointmentQueryHandler) *AppointmentHandler {
	return &AppointmentHandler{cmd: cmd, qry: qry}
}

// Create handles POST /appointments.
func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	var req dto.CreateAppointmentRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	createdBy, _ := middleware.UserIDFromContext(r.Context())
	result, err := h.cmd.Create(r.Context(), orgID, appID, createdBy, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.CreatedEnvelope(w, constants.MsgAppointmentCreated, result)
}

// Update handles PATCH /appointments/{appointmentId}.
func (h *AppointmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamAppointmentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidAppointmentID})
		return
	}
	var req dto.UpdateAppointmentRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := h.cmd.Update(r.Context(), orgID, appID, id, req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgAppointmentUpdated, result)
}

// Delete handles DELETE /appointments/{appointmentId}.
func (h *AppointmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamAppointmentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidAppointmentID})
		return
	}
	if err := h.cmd.Delete(r.Context(), orgID, appID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgAppointmentDeleted, nil)
}

// Get handles GET /appointments/{appointmentId}.
func (h *AppointmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, constants.PathParamAppointmentID))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": constants.MsgInvalidAppointmentID})
		return
	}
	result, err := h.qry.Get(r.Context(), orgID, appID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgAppointmentFetched, result)
}

// List handles GET /appointments (filter project_id, owner_id, status, upcoming).
func (h *AppointmentHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, appID, ok := resolveTenant(r)
	if !ok {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "organization and app context required"})
		return
	}
	q := dto.ListAppointmentsQuery{
		SortOrder: r.URL.Query().Get(constants.QueryKeySortOrder),
		Page:      queryInt(r, constants.QueryKeyPage),
		PerPage:   queryInt(r, constants.QueryKeyPerPage),
	}
	if id := parseUUIDQuery(r, constants.QueryKeyApptProject); id != nil {
		q.ProjectID = id
	}
	if id := parseUUIDQuery(r, constants.QueryKeyApptOwner); id != nil {
		q.OwnerID = id
	}
	if v := r.URL.Query().Get(constants.QueryKeyApptStatus); v != "" {
		q.Status = &v
	}
	if v := r.URL.Query().Get(constants.QueryKeyApptUpcoming); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			q.Upcoming = b
		}
	}
	result, err := h.qry.List(r.Context(), orgID, appID, q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, constants.MsgAppointmentListed, result)
}
