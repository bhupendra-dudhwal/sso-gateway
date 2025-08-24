package response

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/valyala/fasthttp"
)

type response struct {
	logger      ports.Logger
	compression bool
	payload     models.Response
}

func NewResponse(requestID string, compression bool, logger ports.Logger) ports.Response {
	return &response{
		logger:      logger,
		compression: compression,
		payload: models.Response{
			RequestID: requestID,
		},
	}
}

func (r *response) SetStatusCode(code int) ports.Response {
	r.payload.StatusCode = code
	return r
}

func (r *response) SetMessage(msg string) ports.Response {
	r.payload.Message = msg
	return r
}

func (r *response) SetStatus(status bool) ports.Response {
	r.payload.Status = status
	return r
}

func (r *response) SetToken(token string) ports.Response {
	r.payload.Token = token
	return r
}

func (r *response) SetPermission(permissions []string) ports.Response {
	r.payload.Permissions = permissions
	return r
}

func (r *response) SetPayload(payload any) ports.Response {
	r.payload.Payload = payload
	return r
}

func (r *response) SetErrorCode(code string) ports.Response {
	if r.payload.Error == nil {
		r.payload.Error = &models.Error{}
	}
	r.payload.Error.Code = code

	return r
}

func (r *response) SetErrorMessage(msg string) ports.Response {
	if r.payload.Error == nil {
		r.payload.Error = &models.Error{}
	}
	r.payload.Error.Message = msg

	return r
}

func (r *response) SetErrorDetails(details any) ports.Response {
	if r.payload.Error == nil {
		r.payload.Error = &models.Error{}
	}
	r.payload.Error.Detail = details

	return r
}

func (r *response) SetError(err *models.Error) ports.Response {
	r.payload.Error = err

	return r
}

func (r *response) Send(ctx *fasthttp.RequestCtx) {
	// Set common headers
	ctx.Response.Header.Set(constants.ContentType.String(), constants.Json.String())

	// Enable gzip encoding if requested
	if r.compression {
		ctx.Response.Header.Set(constants.ContentEncoding.String(), constants.Gzip.String())
		ctx.SetStatusCode(r.payload.StatusCode)
		ctx.SetBodyStreamWriter(func(w *bufio.Writer) {

			gw := gzip.NewWriter(w)
			defer gw.Close()

			encoder := json.NewEncoder(gw)
			if err := encoder.Encode(r.payload); err != nil {
				r.logger.ErrorCtx(ctx, err.Error())
			}
		})
		return
	}

	// No encoding
	ctx.SetStatusCode(r.payload.StatusCode)

	body, err := json.Marshal(r.payload)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"error":"Failed to encode response"}`)
		return
	}

	fmt.Println(string(body))

	ctx.SetBody(body)
}
