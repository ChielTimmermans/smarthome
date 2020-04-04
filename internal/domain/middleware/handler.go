package middleware

import "github.com/valyala/fasthttp"

type Handler interface {
	AccessControl(next fasthttp.RequestHandler, roles ...string) fasthttp.RequestHandler
}
