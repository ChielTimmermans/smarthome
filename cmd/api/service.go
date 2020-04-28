package main

import (
	"smarthome-home/internal/domain/accesstoken"
	"smarthome-home/internal/domain/middleware"
	"smarthome-home/internal/domain/relay"
	"smarthome-home/internal/domain/user"
	"smarthome-home/internal/service"
)

type Service struct {
	user        user.Servicer
	accessToken accesstoken.Servicer
	middleware  middleware.Servicer
	relay       relay.Servicer
}

func initService(storage *Storage, push *Push, hash string) (s *Service, err error) {
	s = &Service{}

	if s.user, err = service.NewUser(storage.user, hash); err != nil {
		return nil, err
	}
	if s.accessToken, err = service.NewAccessToken(storage.accessToken); err != nil {
		return nil, err
	}
	if s.relay, err = service.NewRelay(storage.relay, push.relay); err != nil {
		return nil, err
	}

	if s.middleware, err = service.NewMiddleware(s.accessToken); err != nil {
		return nil, err
	}

	s.user.AddAccessTokenService(s.accessToken)

	return s, nil
}
