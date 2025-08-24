package services

import (
	"context"
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

type roleService struct {
	errCodePrefix    string
	config           *models.Config
	logger           ports.Logger
	egressRepository egress.Repository
}

func NewRoleService(config *models.Config, logger ports.Logger, egressRepository egress.Repository) ingress.RoleServicePorts {
	return &roleService{
		errCodePrefix:    "RS-%s-%d",
		config:           config,
		logger:           logger,
		egressRepository: egressRepository,
	}
}

func (s *roleService) List(ctx *fasthttp.RequestCtx) {
	var (
		reqID    = utils.GetField(ctx, constants.CtxRequestID)
		logger   = s.logger.With(zap.String("requestID", reqID))
		response = response.NewResponse(reqID, s.config.App.Server.Compression, logger)
	)

	ctxVal, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	roles, err := s.egressRepository.Role.GetRolesWithoutPagination(ctxVal)
	if err != nil {
		logger.Error("Failed to fetch roles", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "LT", 1),
			Message: "Failed to fetch roles",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	msg := "Roles fetched successfully"
	if len(roles) == 0 {
		msg = "No roles found"
		roles = nil
	}

	// Send response
	response.SetStatus(true).SetMessage(msg).SetStatusCode(http.StatusOK).SetPayload(roles).Send(ctx)
}

func (s *roleService) Info(ctx *fasthttp.RequestCtx) {
	var (
		reqID    = utils.GetField(ctx, constants.CtxRequestID)
		logger   = s.logger.With(zap.String("requestID", reqID))
		response = response.NewResponse(reqID, s.config.App.Server.Compression, logger)
	)

	roleID := utils.Sanitize("")

	ctxVal, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	role, err := s.egressRepository.Role.GetByID(ctxVal, constants.Roles(roleID))
	if err != nil {
		if errors.Is(err, utils.ErrDocumentNotFound) {
			logger.Info("Role not found")

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(s.errCodePrefix, "LTIO", 1),
				Message: "Role not found",
				Detail:  err,
			}).SetStatusCode(http.StatusBadRequest).Send(ctx)
			return
		}

		logger.Error("Failed to fetch role info", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "LTIO", 2),
			Message: "Failed to fetch role info",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	// Success
	response.SetStatus(true).SetStatusCode(http.StatusOK).SetMessage("Role found").SetPayload(&role).Send(ctx)
}

func (s *roleService) Add(ctx *fasthttp.RequestCtx) {
	var (
		reqID    = utils.GetField(ctx, constants.CtxRequestID)
		logger   = s.logger.With(zap.String("requestID", reqID))
		response = response.NewResponse(reqID, s.config.App.Server.Compression, logger)
		role     models.Role
	)

	// Decode JSON body
	if err := json.NewDecoder(ctx.RequestBodyStream()).Decode(&role); err != nil {
		logger.Error("Failed to decode login request", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "AD", 1),
			Message: "Invalid request format",
			Detail:  err,
		}).SetStatus(false).SetStatusCode(http.StatusBadRequest).Send(ctx)
		return
	}

	role.Sanitize(constants.Create, 0)
	if err := role.Validate(); err != nil {
		logger.Error("validation failed", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "AD", 2),
			Message: err.Error(),
			Detail:  err,
		}).SetStatusCode(http.StatusBadRequest).Send(ctx)
		return
	}

	seen := map[string]struct{}{}
	uniquePermissions := make([]string, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		if _, exists := seen[p]; !exists {
			seen[p] = struct{}{}
			uniquePermissions = append(uniquePermissions, p)
		}
	}

	// Get permissions from databses
	var permissions []ingressModel.Permission

	ctxVal, cancel := withTimeout(ctx, 1*time.Minute)
	defer cancel()
	permissions, err := s.egressRepository.Permission.GetByIDs(ctxVal, uniquePermissions)
	if err != nil {
		logger.Error("Failed to fetch permission", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "AD", 3),
			Message: "Failed to fetch permission",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	if len(permissions) == 0 {
		logger.Warn("Invalid permissions")

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "AD", 4),
			Message: "Invalid permissions",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	// Filter valid permissions
	var validPermissions = make([]string, 0, len(permissions))
	for _, p := range permissions {
		if p.Status == constants.StatusActive {
			validPermissions = append(validPermissions, p.ID)
		}
	}
	role.Permissions = validPermissions

	// Add role to database
	ctxVal, cancel = withTimeout(ctx, 1*time.Minute)
	defer cancel()
	if err := s.egressRepository.Role.Add(ctx, &role); err != nil {
		if errors.Is(err, utils.ErrDuplicate) {
			logger.Info("Role already exists", zap.Error(err))

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(s.errCodePrefix, "AD", 5),
				Message: "Role already exists",
				Detail:  err,
			}).SetStatusCode(http.StatusBadRequest).Send(ctx)
			return
		}

		logger.Error("Failed to create role", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "AD", 6),
			Message: "Failed to create role",
			Detail:  err,
		}).SetStatusCode(http.StatusInternalServerError).Send(ctx)
	}

	// Success
	response.SetStatus(true).SetMessage("Role created successfully").SetStatusCode(http.StatusOK).SetPayload(&role).Send(ctx)
}

func (s *roleService) Update(ctx *fasthttp.RequestCtx) {

}

func (s *roleService) Delete(ctx *fasthttp.RequestCtx) {
	var (
		reqID    = utils.GetField(ctx, constants.CtxRequestID)
		logger   = s.logger.With(zap.String("requestID", reqID))
		roleID   = utils.Sanitize("")
		response = response.NewResponse(reqID, s.config.App.Server.Compression, logger)
	)

	ctxVal, cancel := withTimeout(ctx, time.Minute)
	defer cancel()

	if err := s.egressRepository.Role.DeleteByID(ctxVal, constants.Roles(roleID)); err != nil {
		if errors.Is(err, utils.ErrDocumentNotFound) {
			msg := fmt.Sprintf("Role '%s' not found", roleID)
			logger.Error(msg)

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(s.errCodePrefix, "DE", 1),
				Message: msg,
				Detail:  err,
			}).SetStatusCode(http.StatusNotFound).Send(ctx)
			return
		}

		logger.Error("Failed to delete role info", zap.Error(err))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(s.errCodePrefix, "DE", 2),
			Message: "Failed to delete role info",
			Detail:  err,
		}).SetStatusCode(http.StatusNotFound).Send(ctx)
		return
	}

	// Success
	response.SetStatus(true).SetMessage(fmt.Sprintf("Role '%s' deleted successfully", roleID)).SetStatusCode(http.StatusOK).Send(ctx)
}
