package services

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/valyala/fasthttp"
)

type userService struct {
	config *models.Config
	logger ports.Logger
}

func NewUserService(config *models.Config, logger ports.Logger) ingress.UserServicePorts {
	return &userService{
		config: config,
		logger: logger,
	}
}

func (p *userService) List(ctx *fasthttp.RequestCtx) {

}

func (p *userService) Info(ctx *fasthttp.RequestCtx) {

}

func (p *userService) Add(ctx *fasthttp.RequestCtx) {

}

func (p *userService) Update(ctx *fasthttp.RequestCtx) {

}

func (p *userService) Delete(ctx *fasthttp.RequestCtx) {

}
