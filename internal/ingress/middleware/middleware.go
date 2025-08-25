package middleware

import (
	"context"
	"net/http"
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

func NewMiddleWare(config *models.Config, logger ports.Logger, tokenService ingress.TokenServicePorts) *middleware {
	return &middleware{
		config:       config,
		logger:       logger,
		tokenService: tokenService,
	}
}

func (m *middleware) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new UUID for this request
		reqID := uuid.NewString()

		// Add the request ID to the context
		ctx := context.WithValue(r.Context(), constants.CtxRequestID, reqID)

		// Add it to the response headers
		w.Header().Set("X-Request-ID", reqID)

		// Pass the new context to the request and call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

// PanicRecover recovers from panics and logs the error
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
