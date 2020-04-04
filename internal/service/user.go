package service

import (
	"errors"
	"smarthome-home/internal"
	"smarthome-home/internal/domain/accesstoken"
	"smarthome-home/internal/domain/user"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	KeepSignedInExpireTime   = time.Hour * 24 * 30
	KeepNotSignedInExpireTim = time.Hour * 24
)

type userService struct {
	storage            user.Storager
	accessTokenService accesstoken.Servicer
	hash               string
}

func NewUser(storage user.Storager, hash string) (user.Servicer, error) {
	if storage == nil {
		return nil, errors.New("userstorage_nil")
	}
	return &userService{
		storage: storage,
		hash:    hash,
	}, nil
}

func (s *userService) AddAccessTokenService(service accesstoken.Servicer) {
	s.accessTokenService = service
}

func (s *userService) Login(ls *user.LoginService) (*accesstoken.Service, error) {
	if err := ls.Validate(); err != nil {
		return nil, err
	}

	var u *user.Service
	var err error
	if u, err = s.storage.GetUserBasedOnEmail(ls.Email); err != nil {
		return nil, err
	}

	if !passwordIsCorrect(u.Password, ls.Password) {
		time.Sleep(time.Second)
		return nil, errors.New("password:incorrect_email_password")
	}

	var token string
	if token, err = internal.GetRandomChars(20); err != nil {
		return nil, err
	}

	var expiresAt time.Time
	if ls.KeepSignedIn {
		expiresAt = time.Now().Add(KeepSignedInExpireTime)
	} else {
		expiresAt = time.Now().Add(KeepNotSignedInExpireTim)
	}

	aT := &accesstoken.Service{
		UserID:    u.ID,
		Token:     token,
		ExpiresAT: expiresAt,
		Role:      u.Role,
	}
	if err := s.accessTokenService.Create(aT); err != nil {
		return nil, err
	}

	return aT, nil
}

func (s *userService) Logout(token string) error {
	return s.accessTokenService.Remove(token)
}

func passwordIsCorrect(hashedPwd, plainPwd string) bool {
	// compare password with db hash
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)) == nil
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *userService) Create(name, email, password, role string) error {
	u := &user.Service{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}

	var err error
	if u.Password, err = hashAndSalt(u.Password); err != nil {
		return err
	}

	if err := s.storage.Create(u); err != nil {
		return err
	}
	return nil
}
