package accesstoken

type Servicer interface {
	Create(*Service) error
	Remove(token string) (err error)
	GetUser(token string) (ats *Service, err error)

	RemoveAllTokensBasedOnUser(userID int) error
	RemoveAllTokensExceptCurrent(userID int, token string) error
}
