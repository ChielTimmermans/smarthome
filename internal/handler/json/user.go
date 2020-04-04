package json

import (
	"errors"
	"smarthome-home/internal/domain/user"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type userHandler struct {
	service          user.Servicer
	jsonIteratorPool jsoniter.IteratorPool
	jsonStreamPool   jsoniter.StreamPool
	answer           *Answer
}

func NewUser(service user.Servicer, jip jsoniter.IteratorPool, jsp jsoniter.StreamPool, a *Answer) (user.Handler, error) {
	if service == nil {
		return nil, errors.New("userservice_nil")
	}
	return &userHandler{
		service:          service,
		jsonIteratorPool: jip,
		jsonStreamPool:   jsp,
		answer:           a,
	}, nil
}

func (h *userHandler) Login(ctx *fasthttp.RequestCtx) {
	var u user.LoginJSON
	jsonIterator := h.jsonIteratorPool.BorrowIterator(ctx.PostBody())
	defer h.jsonIteratorPool.ReturnIterator(jsonIterator)
	jsonIterator.ReadVal(&u)

	accessToken, err := h.service.Login(u.ToService())

	if err != nil {
		h.answer.SetAnswer(ctx, err)
		return
	}

	cookie := &fasthttp.Cookie{}
	cookie.SetKey("X-API-KEY")
	cookie.SetValue(accessToken.Token)
	cookie.SetExpire(accessToken.ExpiresAT)
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	ctx.Response.Header.SetCookie(cookie)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
