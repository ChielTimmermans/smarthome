package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"smarthome-home/internal"
	"smarthome-home/internal/domain/accesstoken"
	"smarthome-home/internal/domain/relay"
	"smarthome-home/internal/domain/user"
	"smarthome-home/internal/storage/mysql"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	user        user.Storager
	accessToken accesstoken.Storager
	relay       relay.Storager
}

const (
	IntervalDBCheck = 1 * time.Second
)

func initStorage(config *ConfigStorageMySQL, databaseStorageMode string, servicesAvailable *internal.ServicesAvailable) (*Storage, error) {
	s := &Storage{}

	if err := initDatabase(config, s, databaseStorageMode, servicesAvailable); err != nil {
		return nil, err
	}
	return s, nil
}

func initDatabase(config *ConfigStorageMySQL, s *Storage, databaseMode string, servicesAvailable *internal.ServicesAvailable) (err error) {
	log.Printf("Connecting to %s", databaseMode)
	switch databaseMode {
	case "mysql":
		dbs := &mysql.DBs{}
		errs := make(chan error)
		go NewMySQLConn(config.Hostname, config.User, config.Password, config.Database, config.Port, dbs, servicesAvailable, errs)

		if s.user, err = mysql.NewUserStorage(dbs); err != nil {
			return err
		}
		if s.accessToken, err = mysql.NewAccessToken(dbs); err != nil {
			return err
		}
		if s.relay, err = mysql.NewRelay(dbs); err != nil {
			return err
		}

	default:
		return fmt.Errorf("storage mode unknown. \t possible modes: %s\t given mode: %s", possibleModes([]string{"vitess"}), databaseMode)
	}
	return nil
}

func possibleModes(possibleDBModes []string) string {
	var data string
	for _, v := range possibleDBModes {
		data = internal.Concat(data, " ['", v, "']")
	}
	data = internal.Concat(data, " ")
	return data
}

//nolint : to complex
func NewMySQLConn(hostname, user, password, database string, port int, dbs *mysql.DBs,
	servicesAvailable *internal.ServicesAvailable, errs chan error) {

	mySQLInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		user, password, hostname, port, database)

	db, err := sql.Open("mysql", mySQLInfo)
	if err != nil {
		errs <- err
	}
	err = db.Ping()
	if err != nil {
		errs <- err
	}
	dbs.Master = db
	servicesAvailable.DB = true

	go func() {
		//nolint, must be infinite loop
		for {
			select {
			case <-time.After(IntervalDBCheck):
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
				defer cancel()

				if servicesAvailable.DB {
					if _, err := dbs.Master.ExecContext(ctx, "SELECT 1 + 1;"); err != nil {
						dbs.Master.Close()
						servicesAvailable.DB = false
					}
				}

				if !servicesAvailable.DB {
					db, err := sql.Open("mysql", mySQLInfo)
					if err != nil {
						errs <- err
					}
					err = db.Ping()
					if err != nil {
						errs <- err
					}
					dbs.Master = db
					servicesAvailable.DB = true
				}
			}
		}
	}()

	go func() {
		for {
			err := <-errs
			log.Println(err)
		}
	}()
}
