package middleware

import (
	"runtime/debug"
	"strings"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/ingress/response"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type middleware struct {
	config       *models.Config
	logger       ports.Logger
	tokenService ingress.TokenServicePorts
}

func NewMiddleWare(config *models.Config, logger ports.Logger, tokenService ingress.TokenServicePorts) ingress.MiddlewarePorts {
	return &middleware{
		config:       config,
		logger:       logger,
		tokenService: tokenService,
	}
}

func (m *middleware) RequestID(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if ctx.UserValue(constants.CtxRequestID) == nil {
			reqID := uuid.NewString()
			ctx.SetUserValue(constants.CtxRequestID, reqID)
			ctx.Response.Header.Set("X-Request-ID", reqID)
		}
		next(ctx)
	}
}

// Authorization returns a fasthttp middleware for auth + permission check
func (m *middleware) Authorization(requiredPermission string) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {

			authHeader := string(ctx.Request.Header.Peek(constants.Authorization))
			if authHeader == "" || !strings.HasPrefix(authHeader, constants.AuthType) {
				reqID := utils.GetField(ctx, constants.CtxRequestID)
				m.logger.Info("missing or invalid authorization header", zap.String("requestID", reqID))

				response.NewResponse(reqID, m.config.App.Server.Compression, m.logger).
					SetStatusCode(fasthttp.StatusUnauthorized).
					SetError(&models.Error{
						Code:    "ME-AN-1",
						Message: "Unauthorized access",
					}).Send(ctx)
				return
			}

			token := strings.TrimPrefix(authHeader, constants.AuthType)
			if !m.tokenService.HavePermission(token, requiredPermission) {
				reqID := utils.GetField(ctx, constants.CtxRequestID)
				m.logger.Info("permission denied", zap.String("requestID", reqID), zap.String("requiredPermission", requiredPermission))

				response.NewResponse(reqID, m.config.App.Server.Compression, m.logger).
					SetStatusCode(fasthttp.StatusForbidden).
					SetError(&models.Error{
						Code:    "ME-AN-2",
						Message: "Permission denied",
					}).Send(ctx)
				return
			}

			// All good — proceed to next handler
			next(ctx)
		}
	}
}

// PanicRecover handles panics and responds with a standard error message
func (m *middleware) PanicRecover(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rec := recover(); rec != nil {
				reqID := utils.GetField(ctx, constants.CtxRequestID)

				m.logger.Error("Panic occurred",
					zap.Any("error", rec),
					zap.ByteString("stackTrace", debug.Stack()),
					zap.String("requestID", reqID),
					zap.String("path", string(ctx.Path())),
					zap.String("method", string(ctx.Method())),
				)

				response.NewResponse(reqID, m.config.App.Server.Compression, m.logger).
					SetStatusCode(fasthttp.StatusInternalServerError).
					SetError(&models.Error{
						Code:    "ME-PR-1",
						Message: "Something went wrong! Please try again later.",
					}).Send(ctx)
			}
		}()

		// Proceed with next handler
		next(ctx)
	}
}
