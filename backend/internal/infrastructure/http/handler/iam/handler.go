package iam

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/iam/dto"
	"github.com/masterfabric-go/masterfabric/internal/application/iam/usecase"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"
	"github.com/masterfabric-go/masterfabric/internal/shared/pagination"
	"github.com/masterfabric-go/masterfabric/internal/shared/response"
	"github.com/masterfabric-go/masterfabric/internal/shared/validator"

	"github.com/masterfabric-go/masterfabric/internal/domain/iam/model"
	iamRepo "github.com/masterfabric-go/masterfabric/internal/domain/iam/repository"
)

// Handler provides IAM HTTP handlers.
type Handler struct {
	registerUC    *usecase.RegisterUseCase
	loginUC       *usecase.LoginUseCase
	assignRoleUC  *usecase.AssignRoleUseCase
	manageUsersUC *usecase.ManageUsersUseCase
	manageRolesUC *usecase.ManageRolesUseCase
	userRepo      iamRepo.UserRepository
}

// NewHandler creates a new IAM handler.
func NewHandler(
	registerUC *usecase.RegisterUseCase,
	loginUC *usecase.LoginUseCase,
	assignRoleUC *usecase.AssignRoleUseCase,
	manageUsersUC *usecase.ManageUsersUseCase,
	manageRolesUC *usecase.ManageRolesUseCase,
	userRepo iamRepo.UserRepository,
) *Handler {
	return &Handler{
		registerUC:    registerUC,
		loginUC:       loginUC,
		assignRoleUC:  assignRoleUC,
		manageUsersUC: manageUsersUC,
		manageRolesUC: manageRolesUC,
		userRepo:      userRepo,
	}
}

// Register handles user registration.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.registerUC.Execute(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, user)
}

// Login handles user authentication.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.loginUC.Execute(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// AssignRole handles role assignment.
func (h *Handler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignRoleRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.assignRoleUC.Execute(r.Context(), req); err != nil {
		response.Error(w, err)
		return
	}

	response.NoContent(w)
}

// GetMe returns the current authenticated user.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	orgID, _ := middleware.OrgIDFromContext(r.Context())
	roles, _ := middleware.RolesFromContext(r.Context())
	permissions, _ := middleware.PermissionsFromContext(r.Context())

	response.JSON(w, http.StatusOK, dto.MeResponse{
		User: dto.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Status:    string(user.Status),
			UserType:  string(user.UserType),
			Phone:     user.Phone,
			City:      user.City,
			District:  user.District,
			CreatedAt: user.CreatedAt,
		},
		OrganizationID: orgID,
		Roles:          roles,
		Permissions:    permissions,
	})
}

// ExportMyData handles GET /me/kvkk/export and returns the authenticated user's
// own personal data (KVKK right of access). PII is returned UNMASKED because the
// requester is the data owner. The access is audit-logged.
func (h *Handler) ExportMyData(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
		return
	}
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	slog.Info("kvkk pii access", "action", "export", "user_id", userID.String())
	response.JSON(w, http.StatusOK, dto.ToKVKKExportResponse(user))
}

// RecordConsent handles POST /me/kvkk/consent and stamps the user's KVKK
// consent (açık rıza) with the policy version they accepted.
func (h *Handler) RecordConsent(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
		return
	}
	var req dto.KVKKConsentRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	now := time.Now().UTC()
	user.KVKKConsentAt = &now
	user.KVKKConsentVersion = req.Version
	if err := h.userRepo.Update(r.Context(), user); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.KVKKConsentResponse{
		Message:   "consent recorded",
		Version:   req.Version,
		ConsentAt: now,
	})
}

// EraseMyData handles DELETE /me/kvkk/erase and anonymizes the data subject's
// PII (KVKK right to erasure / unutulma hakkı) while retaining the row for audit
// and referential integrity. The account is deactivated.
func (h *Handler) EraseMyData(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
		return
	}
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	user.Email = fmt.Sprintf("deleted-%s@anonymized.local", user.ID)
	user.FirstName = "Silinmis"
	user.LastName = "Kullanici"
	user.Phone = ""
	user.NationalID = ""
	user.Address = ""
	user.City = ""
	user.District = ""
	user.Status = model.UserStatusInactive
	if err := h.userRepo.Update(r.Context(), user); err != nil {
		response.Error(w, err)
		return
	}
	slog.Info("kvkk pii access", "action", "erase", "user_id", userID.String())
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "personal data anonymized and account deactivated",
	})
}

// GetUser returns a user by ID.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ToAdminUserResponse(user))
}

// ListUsers returns a paginated list of users.
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	params := pagination.FromRequest(r)

	users, total, err := h.userRepo.List(r.Context(), params.Offset(), params.Limit())
	if err != nil {
		response.Error(w, err)
		return
	}

	infos := make([]*dto.AdminUserResponse, 0, len(users))
	for _, u := range users {
		infos = append(infos, dto.ToAdminUserResponse(u))
	}

	response.JSON(w, http.StatusOK, pagination.NewResult(infos, params, total))
}
