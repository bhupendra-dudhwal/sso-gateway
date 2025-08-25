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
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/ingress/response"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	errCodePrefix     string
	config            *models.Config
	repository        ports.Repository
	egressRepository  egress.Repository
	ingressRepository ingress.Repository
}

func NewAuthService(
	config *models.Config,
	repository ports.Repository,
	egressRepository egress.Repository,
	ingressRepository ingress.Repository,
) ingress.AuthServicePorts {
	return &authService{
		errCodePrefix:     "AH-%s-%d",
		config:            config,
		repository:        repository,
		ingressRepository: ingressRepository,
		egressRepository:  egressRepository,
	}
}

// Required service
// 1. Token Service
// 2. Role Repository
func (a *authService) Session(ctx *fasthttp.RequestCtx) {
	reqID := utils.GetField(ctx, constants.CtxRequestID)
	logger := a.repository.Logger.With(zap.String("requestID", reqID))

	response := response.NewResponse(reqID, a.config.App.Server.Compression, logger)

	ctxVal, cancel := withTimeout(ctx, time.Duration(1*time.Minute))
	defer cancel()

	fmt.Printf("\n\n a.egressRepository.Role - %+v\n\n", a.egressRepository.Role)
	role, err := a.egressRepository.Role.GetByID(ctxVal, constants.RoleSessionUser)
	if err != nil {
		if errors.Is(err, utils.ErrDocumentNotFound) {
			logger.Warn("Session role not found", zap.String("role", constants.RoleSessionUser.String()))

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(a.errCodePrefix, "SN", 1),
				Message: "Session role not found",
				Detail:  nil,
			}).SetStatus(false).SetStatusCode(http.StatusNotFound).Send(ctx)
			return
		}

		logger.Error("Failed to fetch session user role", zap.Error(err))
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SN", 2),
			Message: "Unable to process request. Please try again later.",
			Detail:  nil,
		}).SetStatus(false).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	if role.Status != constants.StatusActive {
		logger.Warn("Session role inactive", zap.String("status", role.Status.String()))

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SN", 3),
			Message: "Session role is inactive",
			Detail:  nil,
		}).SetStatus(false).SetStatusCode(http.StatusForbidden).Send(ctx)
		return
	}

	token, err := a.ingressRepository.Token.GenerateToken(role.ID, role.Permissions, nil)
	if err != nil {
		logger.Error("Token generation failed", zap.Error(err))
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SN", 4),
			Message: "Failed to create session token",
			Detail:  nil,
		}).SetStatus(false).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	// Success
	response.SetStatus(true).SetStatusCode(http.StatusOK).SetToken(token).
		SetPermission(role.Permissions).Send(ctx)
}

func (a *authService) Signin(ctx *fasthttp.RequestCtx) {
	reqID := utils.GetField(ctx, constants.CtxRequestID)
	logger := a.repository.Logger.With(zap.String("requestID", reqID))

	response := response.NewResponse(reqID, a.config.App.Server.Compression, logger)
	var loginPayload = models.LoginRequest{}
	if err := json.NewDecoder(ctx.RequestBodyStream()).Decode(&loginPayload); err != nil {
		logger.Warn("Failed to decode login request", zap.Error(err))
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 1),
			Message: "Invalid request format",
			Detail:  err,
		}).SetStatus(false).SetStatusCode(http.StatusBadRequest).Send(ctx)
		return
	}

	loginPayload.Sanitize()

	if err := loginPayload.Validate(); err != nil {
		logger.Warn("Login request validation failed", zap.Error(err))
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 2),
			Message: err.Error(),
			Detail:  err,
		}).SetStatus(false).SetStatusCode(http.StatusBadRequest).Send(ctx)
		return
	}

	ctxVal, cancel := withTimeout(ctx, time.Duration(1*time.Minute))
	defer cancel()

	user, err := a.egressRepository.User.GetByEmail(ctxVal, loginPayload.Email)
	if err != nil {
		if errors.Is(err, utils.ErrDocumentNotFound) {
			logger.Warn("invalid credentials", zap.String("email", loginPayload.Email))

			response.SetError(&models.Error{
				Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 3),
				Message: "invalid credentials",
				Detail:  nil,
			}).SetStatus(false).SetStatusCode(http.StatusNotFound).Send(ctx)
			return
		}

		logger.Error("Something went wrong! Please try after sometime", zap.Error(err))
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 4),
			Message: "Something went wrong! Please try after sometime",
			Detail:  nil,
		}).SetStatus(false).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	if user.Status != constants.StatusActive {
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 5),
			Message: fmt.Sprintf("your account is %s", user.Status),
			Detail:  nil,
		}).SetStatus(false).SetStatusCode(http.StatusUnauthorized).Send(ctx)
		return
	}

	if !user.LockoutUntil.IsZero() && user.LockoutUntil.After(time.Now()) {
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 6),
			Message: fmt.Sprintf("Too many failed attempts. Try again at %s", user.LockoutUntil.Format(time.RFC1123)),
			Detail:  nil,
		}).SetStatus(false).SetStatusCode(http.StatusForbidden).Send(ctx)
		return
	}
	fail := a.countRecentFailures(ctx, user.ID)
	if err := isValidPassword(user.Password, loginPayload.Password); err != nil {
		fail++

		//add fail login history :: in go routine
		go func(user *models.User, fail int) {
			a.handleFailCounts(context.Background(), user, fail)
		}(user, fail)

		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 7),
			Message: "invalid credentials",
			Detail:  err,
		}).SetStatus(false).SetStatusCode(http.StatusForbidden).Send(ctx)
		return
	}

	// create a token
	token, err := a.ingressRepository.Token.GenerateToken(user.Role, user.Permissions, user)
	if err != nil {
		response.SetError(&models.Error{
			Code:    fmt.Sprintf(a.errCodePrefix, "SIN", 8),
			Message: "Something went wrong! Please try after sometime",
			Detail:  err,
		}).SetStatus(false).SetStatusCode(http.StatusInternalServerError).Send(ctx)
		return
	}

	// :: in go routine
	go func(user *models.User, fail int) {
		a.handleFailCounts(context.Background(), user, fail)
	}(user, fail)

	go func(user *models.User) {
		a.egressRepository.LoginHistory.Add(context.Background(), &models.LoginHistory{
			UserID:     user.ID,
			Status:     constants.StatusSuccess,
			Token:      token,
			Permission: user.Permissions,
		})
	}(user)

	response.SetStatus(true).SetMessage("success").SetPayload(user).SetToken(token).SetPermission(user.Permissions)
}

func (a *authService) Signup(ctx *fasthttp.RequestCtx) {
	// reqID := utils.GetField(ctx, constants.CtxRequestID)
	// logger := a.logger.With(zap.String("requestID", reqID))

	// response := response.NewResponse(reqID, a.config.App.Server.Compression, logger)
}

func (a *authService) Otp(ctx *fasthttp.RequestCtx) {

}

func (a *authService) Verify(ctx *fasthttp.RequestCtx) {

}

// `hashedPassword` is fetched from DB (as []byte or string)
// `inputPassword` is the plain password the user typed
func isValidPassword(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func (a *authService) handleFailCounts(ctx context.Context, user *models.User, failCount int) {
	if failCount >= a.config.App.Login.MaxFailedAttempts {
		a.egressRepository.User.LockByID(ctx, user.ID, time.Now().Add(a.config.App.Login.LockoutDurationMinutes))
	}
}

func (a *authService) countRecentFailures(ctx context.Context, userID int) int {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	before := time.Now().Add(-a.config.App.Login.LockoutDurationMinutes)
	history, _ := a.egressRepository.LoginHistory.GetByIDAndLoginAt(ctx, userID, before)
	fail := 0
	for _, login := range history {
		if login.Status == constants.StatusSuccess {
			break
		}
		fail++
	}
	return fail
}
