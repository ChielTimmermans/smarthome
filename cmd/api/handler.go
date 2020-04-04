package main

import (
	"smarthome-home/internal/domain/middleware"
	"smarthome-home/internal/domain/user"
	"smarthome-home/internal/handler/json"

	jsoniter "github.com/json-iterator/go"
)

type Handler struct {
	user       user.Handler
	middleware middleware.Handler
}

func initHandler(s *Service) (h *Handler, err error) {
	jsonIteratorPool, jsonStreamPool := initJSON()
	errors := json.InitAnswer(jsonStreamPool)
	h = &Handler{}

	if h.user, err = json.NewUser(s.user, jsonIteratorPool, jsonStreamPool, errors); err != nil {
		return nil, err
	}
	if h.middleware, err = json.NewMiddleware(s.middleware, errors); err != nil {
		return nil, err
	}
	return
}

func initJSON() (jsonIteratorPool jsoniter.IteratorPool, jsonStreamPool jsoniter.StreamPool) {
	jsonIteratorPool = jsoniter.NewIterator(jsoniter.ConfigFastest).Pool()

	// TODO not sure about the 100 buffer size, will lookup in the future
	jsonStreamPool = jsoniter.NewStream(jsoniter.ConfigFastest, nil, 100).Pool()
	return jsonIteratorPool, jsonStreamPool
}
