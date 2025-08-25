package handler

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// middlewareFunc type for fasthttp middleware functions
type middlewareFunc func(fasthttp.RequestHandler) fasthttp.RequestHandler

type handler struct {
	route           *router.Router
	middlewarePorts ingress.MiddlewarePorts
}

// ChainMiddleware applies multiple middlewares to a handler
func chainMiddleware(base fasthttp.RequestHandler, middlewares ...middlewareFunc) fasthttp.RequestHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		base = middlewares[i](base)
	}
	return base
}

func NewHandler(middlewarePorts ingress.MiddlewarePorts) (fasthttp.RequestHandler, ingress.HandlerPorts) {
	r := router.New()

	newHandler := chainMiddleware(r.Handler,
		middlewarePorts.RequestID,
		middlewarePorts.PanicRecover,
	)

	return newHandler, &handler{
		route:           r,
		middlewarePorts: middlewarePorts,
	}
}

func (h *handler) SetHealthHandler(healthService ingress.HealthServicePorts) {
	healthGroup := h.route.Group("/healthz")
	healthGroup.GET("/readiness", healthService.Readiness)
	healthGroup.GET("/liveness", healthService.Liveness)
}

func (h *handler) SetAuthHandler(authService ingress.AuthServicePorts) {
	userGroup := h.route.Group("/api/v1/auth")
	userGroup.GET("/session", authService.Session)
	userGroup.POST("/signin", h.middlewarePorts.Authorization(constants.PrmSignin)(authService.Signin))
	userGroup.POST("/signup", h.middlewarePorts.Authorization(constants.PrmSignup)(authService.Signin))
}

func (h *handler) SetUserHandler(userService ingress.UserServicePorts) {
	userGroup := h.route.Group("/api/v1/users")
	userGroup.GET("/", h.middlewarePorts.Authorization(constants.PrmListUser)(userService.List))            // List
	userGroup.GET("/{id}", h.middlewarePorts.Authorization(constants.PrmInfoUser)(userService.Info))        // Info
	userGroup.POST("/", h.middlewarePorts.Authorization(constants.PrmAdduser)(userService.Add))             // Add
	userGroup.PUT("/{id}", h.middlewarePorts.Authorization(constants.PrmEditUser)(userService.Update))      // Update
	userGroup.DELETE("/{id}", h.middlewarePorts.Authorization(constants.PrmDeleteUser)(userService.Delete)) // Delete
}

func (h *handler) SetRoleHandler(roleService ingress.RoleServicePorts) {
	roleGroup := h.route.Group("/api/v1/roles")
	roleGroup.GET("/", h.middlewarePorts.Authorization(constants.PrmListRoles)(roleService.List))            // List
	roleGroup.GET("/{id}", h.middlewarePorts.Authorization(constants.PrmInfoRole)(roleService.Info))         // Info
	roleGroup.POST("/", h.middlewarePorts.Authorization(constants.PrmAddRoles)(roleService.Add))             // Add
	roleGroup.PUT("/{id}", h.middlewarePorts.Authorization(constants.PrmEditRoles)(roleService.Update))      // Update
	roleGroup.DELETE("/{id}", h.middlewarePorts.Authorization(constants.PrmDeleteRoles)(roleService.Delete)) // Delete
}

func (h *handler) SetPermissionHandler(permissionsService ingress.PermissionServicePorts) {
	permissionGroup := h.route.Group("/api/v1/permissions")
	permissionGroup.GET("/", h.middlewarePorts.Authorization(constants.PrmListPermissions)(permissionsService.List))            // List
	permissionGroup.GET("/{id}", h.middlewarePorts.Authorization(constants.PrmInfoPermission)(permissionsService.Info))         // Info
	permissionGroup.POST("/", h.middlewarePorts.Authorization(constants.PrmAddPermissions)(permissionsService.Add))             // Add
	permissionGroup.PUT("/{id}", h.middlewarePorts.Authorization(constants.PrmEditPermissions)(permissionsService.Update))      // Update
	permissionGroup.DELETE("/{id}", h.middlewarePorts.Authorization(constants.PrmDeletePermissions)(permissionsService.Delete)) // Delete
}
