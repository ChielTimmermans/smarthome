package accesstoken

type Storager interface {
	Create(*Service) error
	Remove(token string) (err error)
	GetUser(token string) (u *Service, err error)

	RemoveAllTokensBasedOnUser(userID int) error
	RemoveAllTokensExceptCurrent(userID int, token string) error
}
