package services

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/valyala/fasthttp"
)

type roleService struct {
	config         *models.Config
	logger         ports.Logger
	roleRepository egress.RoleRepositoryPorts
}

func NewRoleService(config *models.Config, logger ports.Logger, roleRepository egress.RoleRepositoryPorts) ingress.RoleServicePorts {
	return &roleService{
		config:         config,
		logger:         logger,
		roleRepository: roleRepository,
	}
}

func (p *roleService) List(ctx *fasthttp.RequestCtx) {

}

func (p *roleService) Info(ctx *fasthttp.RequestCtx) {

}

func (p *roleService) Add(ctx *fasthttp.RequestCtx) {

}

func (p *roleService) Update(ctx *fasthttp.RequestCtx) {

}

func (p *roleService) Delete(ctx *fasthttp.RequestCtx) {

}
