package user

import "smarthome-home/internal/domain/accesstoken"

type Servicer interface {
	AddAccessTokenService(accesstoken.Servicer)
	Login(*LoginService) (*accesstoken.Service, error)
	Logout(hash string) error
	Create(name, email, password, role string) error
}
