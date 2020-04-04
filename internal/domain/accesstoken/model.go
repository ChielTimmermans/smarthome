package accesstoken

//nolint file not goimports -ed
import (
	"time"

	"github.com/go-sql-driver/mysql"
)

type Service struct {
	UserID    int
	Token     string
	Role      string
	ExpiresAT time.Time
}

type MySQL struct {
	UserID    int
	Token     string
	Role      string
	ExpiresAT mysql.NullTime
}

func (m *MySQL) ToService() *Service {
	return &Service{
		UserID:    m.UserID,
		Token:     m.Token,
		Role:      m.Role,
		ExpiresAT: m.ExpiresAT.Time,
	}
}

func (s *Service) ToMySQL() *MySQL {
	return &MySQL{
		UserID: s.UserID,
		Token:  s.Token,
		Role:   s.Role,
		ExpiresAT: mysql.NullTime{
			Time:  s.ExpiresAT,
			Valid: true,
		},
	}
}
