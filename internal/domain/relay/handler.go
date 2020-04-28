package relay

import "github.com/valyala/fasthttp"

type Handler interface {
	Enable(ctx *fasthttp.RequestCtx)
	Disable(ctx *fasthttp.RequestCtx)
	Toggle(ctx *fasthttp.RequestCtx)
}
