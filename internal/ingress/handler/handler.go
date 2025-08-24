package handler

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/fasthttp/router"
)

type handler struct {
	route *router.Router
}

func NewHandler() (*router.Router, ingress.HandlerPorts) {
	r := router.New()

	return r, &handler{
		route: r,
	}
}

func (h *handler) SetHealthHandler(healthService ingress.HealthServicePorts) {
	healthGroup := h.route.Group("/healthz")
	healthGroup.GET("/readiness", healthService.Readiness)
	healthGroup.GET("/liveness", healthService.Liveness)
}

func (h *handler) SetAuthHandler(authService ingress.AuthServicePorts) {
	api := h.route.Group("/api/v1")
	userGroup := api.Group("/auth")
	userGroup.POST("/signin", authService.Signin)
	userGroup.POST("/signup", authService.Signin)
	userGroup.POST("/otp", authService.Otp)
	userGroup.POST("/otp/Verify", authService.Verify)
	userGroup.POST("/session", authService.Session)
}

func (h *handler) SetUserHandler(userService ingress.UserServicePorts) {
	api := h.route.Group("/api/v1")
	userGroup := api.Group("/users")
	userGroup.GET("", userService.List)           // List
	userGroup.GET("/{id}", userService.Info)      // Info
	userGroup.POST("", userService.Add)           // Add
	userGroup.PUT("/{id}", userService.Update)    // Update
	userGroup.DELETE("/{id}", userService.Delete) // Delete
}

func (h *handler) SetRoleHandler(roleService ingress.RoleServicePorts) {
	api := h.route.Group("/api/v1")
	roleGroup := api.Group("/roles")
	roleGroup.GET("", roleService.List)           // List
	roleGroup.GET("/{id}", roleService.Info)      // Info
	roleGroup.POST("", roleService.Add)           // Add
	roleGroup.PUT("/{id}", roleService.Update)    // Update
	roleGroup.DELETE("/{id}", roleService.Delete) // Delete
}

func (h *handler) SetPermissionHandler(permissionsService ingress.PermissionServicePorts) {
	api := h.route.Group("/api/v1")
	permissionGroup := api.Group("/permissions")
	permissionGroup.GET("", permissionsService.List)           // List
	permissionGroup.GET("/{id}", permissionsService.Info)      // Info
	permissionGroup.POST("", permissionsService.Add)           // Add
	permissionGroup.PUT("/{id}", permissionsService.Update)    // Update
	permissionGroup.DELETE("/{id}", permissionsService.Delete) // Delete
}
