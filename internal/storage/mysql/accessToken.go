package mysql

import (
	"errors"
	"smarthome-home/internal/domain/accesstoken"
)

type accessTokenStorage struct {
	dbs *DBs
}

func NewAccessToken(dbs *DBs) (accesstoken.Storager, error) {
	return &accessTokenStorage{
		dbs: dbs,
	}, nil
}

func (s *accessTokenStorage) Create(at *accesstoken.Service) error {
	atp := at.ToMySQL()
	if _, err := s.dbs.Master.Exec(`INSERT INTO access_tokens (user_id, token, role, expires_at) VALUES (?,?,?,?,?)`,
		atp.UserID, atp.Token, atp.Role, atp.ExpiresAT); err != nil {
		return err
	}
	return nil
}

func (s *accessTokenStorage) Remove(token string) (err error) {
	if _, err := s.dbs.Master.Exec("DELETE FROM access_tokens where token = ?", token); err != nil {
		return err
	}
	return nil
}

func (s *accessTokenStorage) GetUser(token string) (*accesstoken.Service, error) {
	rows, err := s.dbs.Master.Query(`SELECT user_id, token, role, expires_at 
		FROM access_tokens where token = ? AND expires_at >= NOW() `, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessTokens []accesstoken.MySQL
	for rows.Next() {
		var at accesstoken.MySQL
		if err := rows.Scan(&at.UserID, &at.Token, &at.Role, &at.ExpiresAT); err != nil {
			return nil, err
		}
		accessTokens = append(accessTokens, at)
	}
	lengthAT := len(accessTokens)

	if lengthAT == 1 {
		return accessTokens[0].ToService(), nil
	} else if lengthAT > 1 {
		if err := s.Remove(token); err != nil {
			return nil, err
		}
		return nil, errors.New("token:duplicated_token")
	}
	return nil, errors.New("status:status_unauthorized")
}

func (s *accessTokenStorage) RemoveAllTokensBasedOnUser(userID int) error {
	if _, err := s.dbs.Master.Exec("DELETE FROM access_tokens where user_id = ?", userID); err != nil {
		return err
	}
	return nil
}

func (s *accessTokenStorage) RemoveAllTokensExceptCurrent(userID int, token string) error {
	if _, err := s.dbs.Master.Exec("DELETE FROM access_tokens where user_id = ? AND token != ?", userID, token); err != nil {
		return err
	}
	return nil
}
