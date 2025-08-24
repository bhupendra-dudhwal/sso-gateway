package services

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/valyala/fasthttp"
)

type permissionService struct {
	config *models.Config
	logger ports.Logger
}

func NewPermissionService(config *models.Config, logger ports.Logger) ingress.PermissionServicePorts {
	return &permissionService{
		config: config,
		logger: logger,
	}
}

func (p *permissionService) List(ctx *fasthttp.RequestCtx) {

}

func (p *permissionService) Info(ctx *fasthttp.RequestCtx) {

}

func (p *permissionService) Add(ctx *fasthttp.RequestCtx) {

}

func (p *permissionService) Update(ctx *fasthttp.RequestCtx) {

}

func (p *permissionService) Delete(ctx *fasthttp.RequestCtx) {

}
