package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	ingressModel "github.com/bhupendra-dudhwal/sso-gateway/internal/core/models/ingress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/ingress/response"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type permissionService struct {
	errCodePrefix    string
	config           *models.Config
	logger           ports.Logger
	egressRepository egress.Repository
}

func NewPermissionService(config *models.Config, logger ports.Logger, egressRepository egress.Repository) ingress.PermissionServicePorts {
	return &permissionService{
		errCodePrefix:    "PN-%s-%d",
		config:           config,
		logger:           logger,
		egressRepository: egressRepository,
	}
}

func (s *permissionService) List(ctx *fasthttp.RequestCtx) {
	var (
		reqID    = utils.GetField(&fasthttp.RequestCtx{}, constants.CtxRequestID)
		logger   = s.logger.With(zap.String("requestID", reqID))
		response = response.NewResponse(reqID, s.config.App.Server.Compression, logger)
	)

	ctxVal, cancel := withTimeout(ctx, 1*time.Minute)
	defer cancel()

	permissions, err := s.egressRepository.Permission.GetPermissionWithoutPagination(ctxVal)
	if err != nil {
		logger.Error("Failed to fetch permissions", zap.Error(err))
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "LT", 1),
			Message: "Failed to fetch permissions",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)

		return
	}

	msg := "Permissions fetched successfully"
	if len(permissions) == 0 {
		msg = "No permissions found"
		permissions = nil
	}

	// Send response
	response.SetStatus(true).SetStatusCode(http.StatusOK).SetMessage(msg).SetPayload(permissions).Send(ctx)
}

func (s *permissionService) Info(ctx *fasthttp.RequestCtx) {
	var (
		reqID        = utils.GetField(ctx, constants.CtxRequestID)
		logger       = s.logger.With(zap.String("requestID", reqID))
		permissionID = utils.Sanitize("")
		response     = response.NewResponse(reqID, s.config.App.Server.Compression, logger)
	)

	ctxVal, cancel := withTimeout(ctx, 1*time.Minute)
	defer cancel()

	permission, err := s.egressRepository.Permission.GetByID(ctxVal, permissionID)
	if err != nil {
		if errors.Is(err, utils.ErrDocumentNotFound) {
			logger.Info("Permission not found")

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(s.errCodePrefix, "IO", 1),
				Message: "Permission not found",
				Detail:  err,
			}).SetStatusCode(http.StatusNotFound).Send(ctx)
			return
		}

		logger.Error("Failed to fetch permission info", zap.String("requestID", reqID), zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "IO", 2),
			Message: "Failed to fetch permission info",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	// Success
	response.SetStatus(true).SetStatusCode(http.StatusOK).SetMessage("Permission found").SetPayload(&permission).Send(ctx)
}

func (p *permissionService) Add(ctx *fasthttp.RequestCtx) {
	var (
		reqID      = utils.GetField(ctx, constants.CtxRequestID)
		logger     = p.logger.With(zap.String("requestID", reqID))
		response   = response.NewResponse(reqID, p.config.App.Server.Compression, logger)
		permission ingressModel.Permission
	)

	// Decode JSON body
	if err := json.NewDecoder(ctx.RequestBodyStream()).Decode(&permission); err != nil {
		logger.Error("Failed to decode login request", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(p.errCodePrefix, "AD", 1),
			Message: "Invalid request format",
			Detail:  err,
		}).SetStatusCode(http.StatusBadRequest).Send(ctx)
		return
	}

	permission.Sanitize(constants.Create, 0)
	if err := permission.Validate(); err != nil {
		logger.Error("validation failed", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(p.errCodePrefix, "AD", 2),
			Message: err.Error(),
			Detail:  err,
		}).SetStatusCode(http.StatusBadRequest).Send(ctx)
		return
	}

	ctxVal, cancel := withTimeout(ctx, 1*time.Minute)
	defer cancel()

	if err := p.egressRepository.Permission.Add(ctxVal, &permission); err != nil {
		if errors.Is(err, utils.ErrDuplicate) {
			logger.Info("Permission already exists")

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(p.errCodePrefix, "AD", 3),
				Message: "Permission already exists",
				Detail:  err,
			}).SetStatusCode(http.StatusBadRequest).Send(ctx)
			return
		}

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(p.errCodePrefix, "AD", 4),
			Message: "Failed to create permission",
			Detail:  err,
		}).SetStatusCode(http.StatusBadRequest).Send(ctx)
		return
	}

	// Success
	response.SetStatus(true).SetMessage("Permission created successfully").SetStatusCode(http.StatusOK).SetPayload(&permission).Send(ctx)
}

func (p *permissionService) Update(ctx *fasthttp.RequestCtx) {

}

func (s *permissionService) Delete(ctx *fasthttp.RequestCtx) {
	var (
		reqID        = utils.GetField(ctx, constants.CtxRequestID)
		logger       = s.logger.With(zap.String("requestID", reqID))
		permissionID = utils.Sanitize("")
		response     = response.NewResponse(reqID, s.config.App.Server.Compression, logger)
	)

	ctxVal, cancel := withTimeout(ctx, 1*time.Minute)
	defer cancel()

	if err := s.egressRepository.Permission.DeleteByID(ctxVal, permissionID); err != nil {
		if errors.Is(err, utils.ErrDocumentNotFound) {
			msg := fmt.Sprintf("Permission '%s' not found", permissionID)
			logger.Error(msg)

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(s.errCodePrefix, "DE", 1),
				Message: msg,
				Detail:  err,
			}).SetStatusCode(http.StatusNotFound).Send(ctx)
			return
		}

		logger.Error("Failed to delete permission info", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "DE", 2),
			Message: "Failed to delete permission info",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	// Success
	response.SetStatus(true).SetStatusCode(http.StatusOK).
		SetMessage(fmt.Sprintf("Permission '%s' deleted successfully", permissionID)).Send(ctx)
}
