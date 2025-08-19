package ingress

import (
	"github.com/valyala/fasthttp"
)

type HealthServicePorts interface {
	Liveness(ctx *fasthttp.RequestCtx)
	Readiness(ctx *fasthttp.RequestCtx)
}
