package ports

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/valyala/fasthttp"
)

type Response interface {
	SetStatusCode(code int) Response
	SetStatus(status bool) Response
	SetMessage(msg string) Response
	SetPayload(payload any) Response
	SetToken(token string) Response
	SetPermission(permissions []string) Response
	SetErrorCode(code string) Response
	SetErrorMessage(msg string) Response
	SetErrorDetails(details any) Response

	SetError(err *models.Error) Response // Set error in one go

	Send(ctx *fasthttp.RequestCtx) // Send the response
}
