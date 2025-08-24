package ingress

import "github.com/valyala/fasthttp"

type RoleServicePorts interface {
	List(ctx *fasthttp.RequestCtx)
	Info(ctx *fasthttp.RequestCtx)
	Add(ctx *fasthttp.RequestCtx)
	Update(ctx *fasthttp.RequestCtx)
	Delete(ctx *fasthttp.RequestCtx)
}
