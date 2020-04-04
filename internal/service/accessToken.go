package service

import (
	"errors"
	"smarthome-home/internal/domain/accesstoken"
)

type accessTokenService struct {
	storage accesstoken.Storager
}

func NewAccessToken(storage accesstoken.Storager) (accesstoken.Servicer, error) {
	if storage == nil {
		return nil, errors.New("userstorage_nil")
	}
	return &accessTokenService{
		storage: storage,
	}, nil
}

func (s *accessTokenService) Create(at *accesstoken.Service) error {
	if err := s.storage.Create(at); err != nil {
		return err
	}

	return nil
}

func (s *accessTokenService) Remove(token string) (err error) {
	return s.storage.Remove(token)
}

func (s *accessTokenService) GetUser(token string) (accessToken *accesstoken.Service, err error) {
	return s.storage.GetUser(token)
}

func (s *accessTokenService) RemoveAllTokensBasedOnUser(userID int) error {
	return s.storage.RemoveAllTokensBasedOnUser(userID)
}

func (s *accessTokenService) RemoveAllTokensExceptCurrent(userID int, token string) error {
	return s.storage.RemoveAllTokensExceptCurrent(userID, token)
}
