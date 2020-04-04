package mysql

import (
	"database/sql"
	"errors"
	"smarthome-home/internal/domain/user"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type userStorage struct {
	dbs *DBs
}

func NewUserStorage(dbs *DBs) (user.Storager, error) {
	return &userStorage{
		dbs: dbs,
	}, nil
}

func (s *userStorage) GetUserBasedOnEmail(email string) (*user.Service, error) {
	uP := &user.MySQL{}
	switch err := s.dbs.Master.QueryRow(`SELECT id, name, password, role, created_at 
		FROM users WHERE LOWER(email) = LOWER(?) AND deleted_at IS NULL`,
		email).Scan(&uP.ID, &uP.Name, &uP.Password, &uP.Role, &uP.CreatedAt); err {
	case sql.ErrNoRows:
		return nil, errors.New("password:incorrect_email_password")
	case nil:
		uS := uP.ToService()
		uS.Email = email
		return uS, nil
	default:
		return nil, err
	}
}

func (s *userStorage) Create(u *user.Service) (err error) {
	um := u.ToMySQL()
	if _, err := s.dbs.Master.Exec(`INSERT INTO users (name, email, role, password) 
		VALUES (?, ?, ?, ?)`,
		um.Name, um.Email, um.Role, um.Password); err != nil {
		if sqlErr, ok := err.(*mysql.MySQLError); ok {
			if sqlErr.Number == 1062 {
				if strings.Contains(sqlErr.Message, "email") {
					return errors.New("email:email_duplicated")
				}
			}
		}
		return err
	}
	return nil
}
