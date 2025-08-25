package ingress

import (
	"github.com/valyala/fasthttp"
)

type MiddlewarePorts interface {
	RequestID(next fasthttp.RequestHandler) fasthttp.RequestHandler
	PanicRecover(next fasthttp.RequestHandler) fasthttp.RequestHandler
	Authorization(requiredPermission string) func(next fasthttp.RequestHandler) fasthttp.RequestHandler
}
