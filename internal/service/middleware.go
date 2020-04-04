package service

import (
	"errors"
	"smarthome-home/internal/domain/accesstoken"
	"smarthome-home/internal/domain/middleware"
)

type middlewareService struct {
	accessTokenService accesstoken.Servicer
}

func NewMiddleware(accessTokenService accesstoken.Servicer) (middleware.Servicer, error) {
	if accessTokenService == nil {
		return nil, errors.New("accesstokenservice_nil")
	}
	return &middlewareService{
		accessTokenService: accessTokenService,
	}, nil
}

func (s *middlewareService) AccessControl(token string, roles []string) (u *accesstoken.Service, err error) {
	if u, err = s.accessTokenService.GetUser(token); err != nil {
		return nil, err
	}
	for _, v := range roles {
		if v == u.Role {
			return u, nil
		}
	}
	return nil, errors.New("status:status_unauthorized")
}
