package user

import "github.com/valyala/fasthttp"

type Handler interface {
	Login(ctx *fasthttp.RequestCtx)
}
