package json

import (
	"errors"
	"smarthome-home/internal/domain/relay"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type relayHandler struct {
	service          relay.Servicer
	jsonIteratorPool jsoniter.IteratorPool
	jsonStreamPool   jsoniter.StreamPool
	answer           *Answer
}

func NewRelay(service relay.Servicer, jip jsoniter.IteratorPool, jsp jsoniter.StreamPool, a *Answer) (relay.Handler, error) {
	if service == nil {
		return nil, errors.New("relayservice_nil")
	}
	return &relayHandler{
		service:          service,
		jsonIteratorPool: jip,
		jsonStreamPool:   jsp,
		answer:           a,
	}, nil
}

func (h *relayHandler) Enable(ctx *fasthttp.RequestCtx) {
	relayID, err := getID(ctx, "relayID")
	if err != nil {
		h.answer.SetAnswer(ctx, errors.New("relayID:relayID_malformed"))
		return
	}
	itemID, err := getID(ctx, "itemID")
	if err != nil {
		h.answer.SetAnswer(ctx, errors.New("itemID:itemID_malformed"))
		return
	}
	if err := h.service.Enable(relayID, itemID); err != nil {
		h.answer.SetAnswer(ctx, err)
		return
	}
	h.answer.SetAnswer(ctx, nil)
}

func (h *relayHandler) Disable(ctx *fasthttp.RequestCtx) {
	relayID, err := getID(ctx, "relayID")
	if err != nil {
		h.answer.SetAnswer(ctx, errors.New("relayID:relayID_malformed"))
		return
	}
	itemID, err := getID(ctx, "itemID")
	if err != nil {
		h.answer.SetAnswer(ctx, errors.New("itemID:itemID_malformed"))
		return
	}
	if err := h.service.Disable(relayID, itemID); err != nil {
		h.answer.SetAnswer(ctx, err)
		return
	}
	h.answer.SetAnswer(ctx, nil)
}

func (h *relayHandler) Toggle(ctx *fasthttp.RequestCtx) {
	relayID, err := getID(ctx, "relayID")
	if err != nil {
		h.answer.SetAnswer(ctx, errors.New("relayID:relayID_malformed"))
		return
	}
	itemID, err := getID(ctx, "itemID")
	if err != nil {
		h.answer.SetAnswer(ctx, errors.New("itemID:itemID_malformed"))
		return
	}
	if err := h.service.Toggle(relayID, itemID); err != nil {
		h.answer.SetAnswer(ctx, err)
		return
	}
	h.answer.SetAnswer(ctx, nil)
}
