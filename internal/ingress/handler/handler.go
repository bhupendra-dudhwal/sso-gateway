package handler

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/fasthttp/router"
)

type handelr struct {
	route *router.Router
}

func NewHandler() (*router.Router, ingress.HandlerPorts) {
	r := router.New()

	return r, &handelr{
		route: r,
	}
}

func (h *handelr) SetHealthHandler(healthService ingress.HealthServicePorts) {
	healthGroup := h.route.Group("/healthz")
	healthGroup.GET("/readiness", healthService.Readiness)
	healthGroup.GET("/liveness", healthService.Liveness)
}
