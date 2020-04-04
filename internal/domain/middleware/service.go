package middleware

import "smarthome-home/internal/domain/accesstoken"

type Servicer interface {
	AccessControl(token string, roles []string) (u *accesstoken.Service, err error)
}
