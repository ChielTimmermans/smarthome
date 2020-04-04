package user

type Storager interface {
	GetUserBasedOnEmail(email string) (*Service, error)
	Create(*Service) error
}
