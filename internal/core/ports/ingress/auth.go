package ingress

import "github.com/valyala/fasthttp"

type AuthServicePorts interface {
	Session(ctx *fasthttp.RequestCtx)
	Signin(ctx *fasthttp.RequestCtx)
	Signup(ctx *fasthttp.RequestCtx)
	Otp(ctx *fasthttp.RequestCtx)
	Verify(ctx *fasthttp.RequestCtx)
}
