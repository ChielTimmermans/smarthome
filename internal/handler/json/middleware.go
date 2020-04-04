package json

import (
	"smarthome-home/internal/domain/accesstoken"
	"smarthome-home/internal/domain/middleware"

	"github.com/valyala/fasthttp"
)

type middlewareHandler struct {
	answer  *Answer
	service middleware.Servicer
}

func NewMiddleware(service middleware.Servicer, a *Answer) (middleware.Handler, error) {
	return &middlewareHandler{
		service: service,
		answer:  a,
	}, nil
}
func (h *middlewareHandler) AccessControl(next fasthttp.RequestHandler, roles ...string) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token := string(ctx.Request.Header.Peek("X-API-KEY"))
		if token == "" {
			token = string(ctx.Request.Header.Cookie("X-API-KEY"))
		}
		var err error
		var u *accesstoken.Service

		if u, err = h.service.AccessControl(token, roles); err != nil {
			h.answer.SetAnswer(ctx, err)
			return
		}
		ctx.SetUserValue("user", u)
		next(ctx)
	}
}
