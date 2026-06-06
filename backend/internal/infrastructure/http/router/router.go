package router

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	// Handlers
	apimgmtHandler "github.com/masterfabric-go/masterfabric/internal/infrastructure/http/handler/apimanagement"
	auditHandler "github.com/masterfabric-go/masterfabric/internal/infrastructure/http/handler/audit"
	"github.com/masterfabric-go/masterfabric/internal/infrastructure/http/handler/health"
	iamHandler "github.com/masterfabric-go/masterfabric/internal/infrastructure/http/handler/iam"
	statsHandler "github.com/masterfabric-go/masterfabric/internal/infrastructure/http/handler/stats"
	tenantHandler "github.com/masterfabric-go/masterfabric/internal/infrastructure/http/handler/tenant"
	urbanHandler "github.com/masterfabric-go/masterfabric/internal/infrastructure/http/handler/urbantransform"

	// Services & middleware
	iamService "github.com/masterfabric-go/masterfabric/internal/domain/iam/service"
	"github.com/masterfabric-go/masterfabric/internal/gateway"
	"github.com/masterfabric-go/masterfabric/internal/shared/middleware"

	// Repositories (for tenant resolver middleware)
	tenantRepo "github.com/masterfabric-go/masterfabric/internal/domain/tenant/repository"
)

// Dependencies holds all injected dependencies for the router.
type Dependencies struct {
	Logger *slog.Logger
	DB     *pgxpool.Pool
	Redis  *redis.Client

	// Services
	AuthService iamService.AuthService
	RBACService iamService.RBACService

	// Handlers
	IAMHandler    *iamHandler.Handler
	TenantHandler *tenantHandler.Handler
	APIMgmtHandler *apimgmtHandler.Handler
	AuditHandler  *auditHandler.Handler
	ProjectHandler      *urbanHandler.ProjectHandler
	ContractorHandler   *urbanHandler.ContractorHandler
	BuildingHandler     *urbanHandler.BuildingHandler
	BuildingUnitHandler *urbanHandler.BuildingUnitHandler
	PropertyOwnerHandler *urbanHandler.PropertyOwnerHandler
	DocumentHandler     *urbanHandler.DocumentHandler
	ApprovalHandler     *urbanHandler.ApprovalHandler
	NotificationHandler *urbanHandler.NotificationHandler
	WorkflowHandler     *urbanHandler.WorkflowHandler
	MessageHandler      *urbanHandler.MessageHandler
	AppointmentHandler  *urbanHandler.AppointmentHandler
	StatsHandler        *statsHandler.Handler

	// Gateway
	GatewayPipeline *gateway.Pipeline

	// Repos needed for middleware
	OrgRepo        tenantRepo.OrgRepository
	WorkspaceRepo  tenantRepo.WorkspaceRepository
}

// registerCRUD wires the standard CQRS CRUD + list/filter/search routes for a
// resource. When rbac is non-nil each route is guarded by a "<resource>.<action>"
// permission; otherwise the routes are registered without permission checks.
func registerCRUD(
	r chi.Router,
	rbac iamService.RBACService,
	resource string,
	idParam string,
	list, create, get, update, del http.HandlerFunc,
) {
	itemPath := "/{" + idParam + "}"
	if rbac != nil {
		r.With(middleware.RequirePermission(rbac, resource+".list")).Get("/", list)
		r.With(middleware.RequirePermission(rbac, resource+".create")).Post("/", create)
		r.With(middleware.RequirePermission(rbac, resource+".read")).Get(itemPath, get)
		r.With(middleware.RequirePermission(rbac, resource+".update")).Patch(itemPath, update)
		r.With(middleware.RequirePermission(rbac, resource+".delete")).Delete(itemPath, del)
		return
	}
	r.Get("/", list)
	r.Post("/", create)
	r.Get(itemPath, get)
	r.Patch(itemPath, update)
	r.Delete(itemPath, del)
}

// New creates the root Chi router with all middleware and routes.
func New(deps Dependencies) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(deps.Logger))
	r.Use(middleware.Recoverer(deps.Logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-Organization-ID", "X-App-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health endpoints
	healthHandler := health.NewHandler(deps.DB, deps.Redis)
	r.Get("/health/live", healthHandler.Liveness)
	r.Get("/health/ready", healthHandler.Readiness)

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes (no JWT required)
		r.Route("/auth", func(r chi.Router) {
			if deps.IAMHandler != nil {
				r.Post("/register", deps.IAMHandler.Register)
				r.Post("/login", deps.IAMHandler.Login)
			}
		})

		// Protected routes (require JWT)
		r.Group(func(r chi.Router) {
			if deps.AuthService != nil {
				r.Use(middleware.JWTAuth(deps.AuthService))
			}

			// Tenant resolution middleware (with workspace support)
			if deps.OrgRepo != nil {
				// Note: WorkspaceRepo can be nil - workspace resolution is optional
				r.Use(middleware.TenantResolverWithWorkspace(deps.OrgRepo, deps.WorkspaceRepo))
			}

			// Gateway pipeline (rate limiting, permission enforcement for managed endpoints)
			// Must be applied before specific routes so it can handle dynamic endpoints
			if deps.GatewayPipeline != nil {
				r.Use(deps.GatewayPipeline.Enforce)
			}

			// User routes
			if deps.IAMHandler != nil {
				r.Get("/me", deps.IAMHandler.GetMe)
				r.Route("/users", func(r chi.Router) {
					r.Get("/", deps.IAMHandler.ListUsers)
					r.Get("/{id}", deps.IAMHandler.GetUser)
					if deps.RBACService != nil {
						r.With(middleware.RequirePermission(deps.RBACService, "user.manage")).
							Patch("/{id}", deps.IAMHandler.UpdateUser)
						r.With(middleware.RequirePermission(deps.RBACService, "user.manage")).
							Post("/{id}/deactivate", deps.IAMHandler.DeactivateUser)
						r.With(middleware.RequirePermission(deps.RBACService, "user.manage")).
							Post("/{id}/activate", deps.IAMHandler.ActivateUser)
					} else {
						r.Patch("/{id}", deps.IAMHandler.UpdateUser)
						r.Post("/{id}/deactivate", deps.IAMHandler.DeactivateUser)
						r.Post("/{id}/activate", deps.IAMHandler.ActivateUser)
					}
				})

				// Role administration (system admin: role CRUD + permission editing)
				r.Route("/roles", func(r chi.Router) {
					if deps.RBACService != nil {
						r.Use(middleware.RequirePermission(deps.RBACService, "role.manage"))
					}
					r.Post("/assign", deps.IAMHandler.AssignRole)
					r.Post("/revoke", deps.IAMHandler.RevokeRole)
					r.Get("/", deps.IAMHandler.ListRoles)
					r.Post("/", deps.IAMHandler.CreateRole)
					r.Get("/{roleId}", deps.IAMHandler.GetRole)
					r.Patch("/{roleId}", deps.IAMHandler.UpdateRole)
					r.Delete("/{roleId}", deps.IAMHandler.DeleteRole)
					r.Put("/{roleId}/permissions", deps.IAMHandler.SetRolePermissions)
				})
			}

			// Organization routes
			if deps.TenantHandler != nil {
				r.Route("/organizations", func(r chi.Router) {
					r.Post("/", deps.TenantHandler.CreateOrg)
					r.Get("/", deps.TenantHandler.ListOrgs)
					r.Route("/{orgId}", func(r chi.Router) {
						r.Get("/", deps.TenantHandler.GetOrg)
						if deps.RBACService != nil {
							r.With(middleware.RequirePermission(deps.RBACService, "organization.manage")).
								Patch("/", deps.TenantHandler.UpdateOrg)
							r.With(middleware.RequirePermission(deps.RBACService, "organization.manage")).
								Delete("/", deps.TenantHandler.DeleteOrg)
						} else {
							r.Patch("/", deps.TenantHandler.UpdateOrg)
							r.Delete("/", deps.TenantHandler.DeleteOrg)
						}

						// Apps under organization
						r.Route("/apps", func(r chi.Router) {
							r.Post("/", deps.TenantHandler.CreateApp)
							r.Get("/", deps.TenantHandler.ListApps)
							r.Route("/{appId}", func(r chi.Router) {
								r.Get("/", deps.TenantHandler.GetApp)

								// API keys under app
								r.Route("/keys", func(r chi.Router) {
									r.Post("/", deps.TenantHandler.CreateAPIKey)
									r.Get("/", deps.TenantHandler.ListAPIKeys)
									r.Delete("/{keyId}", deps.TenantHandler.RevokeAPIKey)
								})

								// Endpoints under app
								if deps.APIMgmtHandler != nil {
									r.Route("/endpoints", func(r chi.Router) {
										r.Post("/", deps.APIMgmtHandler.DefineEndpoint)
										r.Get("/", deps.APIMgmtHandler.ListEndpoints)
										r.Route("/{endpointId}", func(r chi.Router) {
											r.Get("/", deps.APIMgmtHandler.GetEndpoint)
											r.Post("/retire", deps.APIMgmtHandler.RetireEndpoint)
											r.Post("/activate", deps.APIMgmtHandler.ActivateEndpoint)
											r.Put("/policy", deps.APIMgmtHandler.UpdatePolicy)
											r.Get("/policy", deps.APIMgmtHandler.GetPolicy)
										})
									})
								}
							})
						})

						// Workspaces under organization
						r.Route("/workspaces", func(r chi.Router) {
							r.Post("/", deps.TenantHandler.CreateWorkspace)
							r.Get("/", deps.TenantHandler.ListWorkspaces)
							r.Route("/{workspaceId}", func(r chi.Router) {
								r.Put("/", deps.TenantHandler.UpdateWorkspace)
							})
						})

						// Audit logs under organization
						if deps.AuditHandler != nil {
							r.Get("/audit-logs", deps.AuditHandler.ListByOrg)
						}
					})
				})
			}

			// Audit logs by user
			if deps.AuditHandler != nil {
				r.Get("/users/{userId}/audit-logs", deps.AuditHandler.ListByUser)
			}

			// Urban transformation: projects (CQRS CRUD + filter/search)
			if deps.ProjectHandler != nil {
				r.Route("/projects", func(r chi.Router) {
					if deps.RBACService != nil {
						r.With(middleware.RequirePermission(deps.RBACService, "project.list")).
							Get("/", deps.ProjectHandler.List)
						r.With(middleware.RequirePermission(deps.RBACService, "project.create")).
							Post("/", deps.ProjectHandler.Create)
						r.With(middleware.RequirePermission(deps.RBACService, "project.read")).
							Get("/{projectId}", deps.ProjectHandler.Get)
						r.With(middleware.RequirePermission(deps.RBACService, "project.update")).
							Patch("/{projectId}", deps.ProjectHandler.Update)
						r.With(middleware.RequirePermission(deps.RBACService, "project.delete")).
							Delete("/{projectId}", deps.ProjectHandler.Delete)
					} else {
						r.Get("/", deps.ProjectHandler.List)
						r.Post("/", deps.ProjectHandler.Create)
						r.Get("/{projectId}", deps.ProjectHandler.Get)
						r.Patch("/{projectId}", deps.ProjectHandler.Update)
						r.Delete("/{projectId}", deps.ProjectHandler.Delete)
					}

					// Per-project workflow states/history + advance (nested under projects)
					if deps.WorkflowHandler != nil {
						r.Route("/{projectId}/workflow", func(r chi.Router) {
							if deps.RBACService != nil {
								r.With(middleware.RequirePermission(deps.RBACService, "workflow.read")).
									Get("/", deps.WorkflowHandler.ListProjectWorkflow)
								r.With(middleware.RequirePermission(deps.RBACService, "workflow.read")).
									Get("/history", deps.WorkflowHandler.ListProjectHistory)
								r.With(middleware.RequirePermission(deps.RBACService, "workflow.advance")).
									Post("/advance", deps.WorkflowHandler.Advance)
							} else {
								r.Get("/", deps.WorkflowHandler.ListProjectWorkflow)
								r.Get("/history", deps.WorkflowHandler.ListProjectHistory)
								r.Post("/advance", deps.WorkflowHandler.Advance)
							}
						})
					}
				})
			}

			// Urban transformation: contractors (CQRS CRUD + filter/search)
			if deps.ContractorHandler != nil {
				r.Route("/contractors", func(r chi.Router) {
					registerCRUD(r, deps.RBACService, "contractor", "contractorId",
						deps.ContractorHandler.List, deps.ContractorHandler.Create,
						deps.ContractorHandler.Get, deps.ContractorHandler.Update, deps.ContractorHandler.Delete)
				})
			}

			// Urban transformation: buildings (CQRS CRUD + filter/search)
			if deps.BuildingHandler != nil {
				r.Route("/buildings", func(r chi.Router) {
					registerCRUD(r, deps.RBACService, "building", "buildingId",
						deps.BuildingHandler.List, deps.BuildingHandler.Create,
						deps.BuildingHandler.Get, deps.BuildingHandler.Update, deps.BuildingHandler.Delete)
				})
			}

			// Urban transformation: building units (CQRS CRUD + filter/search)
			if deps.BuildingUnitHandler != nil {
				r.Route("/building-units", func(r chi.Router) {
					registerCRUD(r, deps.RBACService, "building_unit", "unitId",
						deps.BuildingUnitHandler.List, deps.BuildingUnitHandler.Create,
						deps.BuildingUnitHandler.Get, deps.BuildingUnitHandler.Update, deps.BuildingUnitHandler.Delete)
				})
			}

			// Urban transformation: property owners (CQRS CRUD + filter/search)
			if deps.PropertyOwnerHandler != nil {
				r.Route("/property-owners", func(r chi.Router) {
					registerCRUD(r, deps.RBACService, "property_owner", "ownerId",
						deps.PropertyOwnerHandler.List, deps.PropertyOwnerHandler.Create,
						deps.PropertyOwnerHandler.Get, deps.PropertyOwnerHandler.Update, deps.PropertyOwnerHandler.Delete)
				})
			}

			// Urban transformation: documents (CQRS CRUD + filter/search + reviews)
			if deps.DocumentHandler != nil {
				r.Route("/documents", func(r chi.Router) {
					registerCRUD(r, deps.RBACService, "document", "documentId",
						deps.DocumentHandler.List, deps.DocumentHandler.Create,
						deps.DocumentHandler.Get, deps.DocumentHandler.Update, deps.DocumentHandler.Delete)
					if deps.RBACService != nil {
						r.With(middleware.RequirePermission(deps.RBACService, "document.review")).
							Post("/{documentId}/reviews", deps.DocumentHandler.Review)
						r.With(middleware.RequirePermission(deps.RBACService, "document.read")).
							Get("/{documentId}/reviews", deps.DocumentHandler.ListReviews)
					} else {
						r.Post("/{documentId}/reviews", deps.DocumentHandler.Review)
						r.Get("/{documentId}/reviews", deps.DocumentHandler.ListReviews)
					}
				})
				// Document types master data (read-only)
				r.Get("/document-types", deps.DocumentHandler.ListTypes)
			}

			// Urban transformation: approvals (pending approvals + decide)
			if deps.ApprovalHandler != nil {
				r.Route("/approvals", func(r chi.Router) {
					if deps.RBACService != nil {
						r.With(middleware.RequirePermission(deps.RBACService, "approval.list")).
							Get("/", deps.ApprovalHandler.List)
						r.With(middleware.RequirePermission(deps.RBACService, "approval.create")).
							Post("/", deps.ApprovalHandler.Create)
						r.With(middleware.RequirePermission(deps.RBACService, "approval.read")).
							Get("/{approvalId}", deps.ApprovalHandler.Get)
						r.With(middleware.RequirePermission(deps.RBACService, "approval.decide")).
							Patch("/{approvalId}/decision", deps.ApprovalHandler.Decide)
						r.With(middleware.RequirePermission(deps.RBACService, "approval.delete")).
							Delete("/{approvalId}", deps.ApprovalHandler.Delete)
					} else {
						r.Get("/", deps.ApprovalHandler.List)
						r.Post("/", deps.ApprovalHandler.Create)
						r.Get("/{approvalId}", deps.ApprovalHandler.Get)
						r.Patch("/{approvalId}/decision", deps.ApprovalHandler.Decide)
						r.Delete("/{approvalId}", deps.ApprovalHandler.Delete)
					}
				})
			}

			// Urban transformation: notifications (current user's inbox)
			if deps.NotificationHandler != nil {
				r.Route("/notifications", func(r chi.Router) {
					r.Get("/", deps.NotificationHandler.List)
					r.Get("/unread-count", deps.NotificationHandler.UnreadCount)
					r.Post("/read-all", deps.NotificationHandler.MarkAllRead)
					r.Patch("/{notificationId}/read", deps.NotificationHandler.MarkRead)
					if deps.RBACService != nil {
						r.With(middleware.RequirePermission(deps.RBACService, "notification.create")).
							Post("/", deps.NotificationHandler.Create)
					} else {
						r.Post("/", deps.NotificationHandler.Create)
					}
				})
			}

			// Urban transformation: workflow step definitions (master data).
			// Per-project workflow routes are nested under /projects/{projectId}/workflow.
			if deps.WorkflowHandler != nil {
				r.Get("/workflow-steps", deps.WorkflowHandler.ListSteps)
			}

			// Urban transformation: messages (inbox/sent + send + mark read)
			if deps.MessageHandler != nil {
				r.Route("/messages", func(r chi.Router) {
					r.Get("/", deps.MessageHandler.List)
					r.Get("/unread-count", deps.MessageHandler.UnreadCount)
					r.Get("/{messageId}", deps.MessageHandler.Get)
					r.Patch("/{messageId}/read", deps.MessageHandler.MarkRead)
					if deps.RBACService != nil {
						r.With(middleware.RequirePermission(deps.RBACService, "message.create")).
							Post("/", deps.MessageHandler.Create)
					} else {
						r.Post("/", deps.MessageHandler.Create)
					}
				})
			}

			// Urban transformation: appointments (CQRS CRUD + filter)
			if deps.AppointmentHandler != nil {
				r.Route("/appointments", func(r chi.Router) {
					registerCRUD(r, deps.RBACService, "appointment", "appointmentId",
						deps.AppointmentHandler.List, deps.AppointmentHandler.Create,
						deps.AppointmentHandler.Get, deps.AppointmentHandler.Update, deps.AppointmentHandler.Delete)
				})
			}

			// Dashboard statistics
			if deps.StatsHandler != nil {
				r.Get("/dashboard/stats", deps.StatsHandler.ProjectDashboard)
				if deps.RBACService != nil {
					r.With(middleware.RequirePermission(deps.RBACService, "system.configure")).
						Get("/admin/stats", deps.StatsHandler.AdminDashboard)
				} else {
					r.Get("/admin/stats", deps.StatsHandler.AdminDashboard)
				}
			}

			// Catch-all handler for managed endpoints (must be last in the group)
			// This allows the gateway pipeline to handle dynamic endpoints like /api/v1/products
			// The gateway middleware will validate and return responses for managed endpoints
			r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
				// Gateway middleware should have already handled this if it's a managed endpoint
				// If we reach here, it means no endpoint was found, return 404
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":"endpoint not found","code":404,"message":"No endpoint registered for this path. Define the endpoint first using POST /api/v1/organizations/{orgId}/apps/{appId}/endpoints"}`))
			})
		})
	})

	// Catch-all handler for managed endpoints (must be after all specific routes)
	// This allows the gateway pipeline to handle dynamic endpoints like /api/v1/products
	// The gateway middleware will validate and return responses for managed endpoints
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// If this is an API v1 path, let the gateway handle it (if it hasn't already)
		// Otherwise return 404
		if !strings.HasPrefix(r.URL.Path, "/api/v1") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"not found","code":404}`))
			return
		}
		
		// For /api/v1 paths, check if gateway pipeline already handled it
		// If not, return 404 (gateway would have returned response if endpoint existed)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"endpoint not found","code":404,"message":"No endpoint registered for this path. Define the endpoint first."}`))
	})

	return r
}
