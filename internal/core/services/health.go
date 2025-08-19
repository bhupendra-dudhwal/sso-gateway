package services

import (
	"net/http"

	"github.com/bhupendra-dudhwal/go-hexagonal/internal/constants"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/core/models"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/core/ports"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/core/ports/ingress"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/ingress/response"
	"github.com/bhupendra-dudhwal/go-hexagonal/internal/utils"
	"github.com/valyala/fasthttp"
)

type healthService struct {
	config *models.Config
	logger ports.Logger
}

func NewHealthService(config *models.Config, logger ports.Logger) ingress.HealthServicePorts {
	return &healthService{
		config: config,
		logger: logger,
	}
}

func (h *healthService) Readiness(ctx *fasthttp.RequestCtx) {
	requestID := utils.GetField(ctx, constants.CtxRequestID)
	response.NewResponse(requestID, h.config.App.Server.Compression, h.logger).SetStatus(true).
		SetMessage("server is ready to server").
		SetStatusCode(http.StatusOK).
		Send(ctx)
}

func (h *healthService) Liveness(ctx *fasthttp.RequestCtx) {
	requestID := utils.GetField(ctx, constants.CtxRequestID)
	response.NewResponse(requestID, h.config.App.Server.Compression, h.logger).SetStatus(true).
		SetMessage("server is live").
		SetStatusCode(http.StatusOK).
		Send(ctx)
}
