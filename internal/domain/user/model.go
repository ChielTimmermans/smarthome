package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-sql-driver/mysql"
)

const (
	NameLengthMin     = 3
	NameLengthMax     = 50
	EmailLengthMin    = 3
	EmailLengthMax    = 255
	PasswordLengthMin = 8
	PasswordLengthMax = 50
	HashLengthMin     = 30
	HashLengthMax     = 30

	USER  = "user"
	ADMIN = "admin"
)

type JSON struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Service struct {
	ID        int
	Name      string
	Password  string
	Email     string
	Role      string
	CreatedAt time.Time
}

type MySQL struct {
	ID        int
	Name      string
	Password  string
	Email     string
	Role      string
	CreatedAt mysql.NullTime
}

func (j *JSON) ToService() *Service {
	return &Service{
		ID:        j.ID,
		Name:      j.Name,
		Password:  j.Password,
		Email:     j.Email,
		Role:      j.Role,
		CreatedAt: j.CreatedAt,
	}
}

func (s *Service) ToJSON() *JSON {
	return &JSON{
		ID:        s.ID,
		Name:      s.Name,
		Password:  s.Password,
		Email:     s.Email,
		Role:      s.Role,
		CreatedAt: s.CreatedAt,
	}
}

func (p *MySQL) ToService() *Service {
	s := &Service{
		ID:       p.ID,
		Name:     p.Name,
		Password: p.Password,
		Email:    p.Email,
		Role:     p.Role,
	}
	if p.CreatedAt.Valid {
		s.CreatedAt = p.CreatedAt.Time
	}
	return s
}

func (s *Service) ToMySQL() *MySQL {
	return &MySQL{
		ID:       s.ID,
		Name:     s.Name,
		Password: s.Password,
		Email:    s.Email,
		Role:     s.Role,
		CreatedAt: mysql.NullTime{
			Time:  s.CreatedAt,
			Valid: true,
		},
	}
}

func (s *Service) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Email,
			validation.Required.Error("Email_required"),
			validation.Length(3, 255).Error("Email_length")),
	)
}

type LoginJSON struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	KeepSignedIn bool   `json:"keepSignedIn"`
}

type LoginService struct {
	Email        string
	Password     string
	KeepSignedIn bool
}

type LoginMySQL struct {
	Email        string
	Password     string
	KeepSignedIn bool
}

func (j *LoginJSON) ToService() *LoginService {
	return &LoginService{
		Password:     j.Password,
		Email:        j.Email,
		KeepSignedIn: j.KeepSignedIn,
	}
}

func (s *LoginService) ToJSON() *LoginJSON {
	return &LoginJSON{
		Password:     s.Password,
		Email:        s.Email,
		KeepSignedIn: s.KeepSignedIn,
	}
}

func (m *LoginMySQL) ToService() *LoginService {
	return &LoginService{
		Password:     m.Password,
		Email:        m.Email,
		KeepSignedIn: m.KeepSignedIn,
	}
}

func (s *LoginService) ToMySQL() *LoginMySQL {
	return &LoginMySQL{
		Password:     s.Password,
		Email:        s.Email,
		KeepSignedIn: s.KeepSignedIn,
	}
}

func (s *LoginService) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Email,
			validation.Required.Error("email_required"),
			validation.Length(EmailLengthMin, EmailLengthMax).Error("email_length"),
			is.Email.Error("email_must_be_email")),

		validation.Field(&s.Password,
			validation.Required.Error("password_required"),
			validation.Length(PasswordLengthMin, PasswordLengthMax).Error("password_length")),
	)
}
