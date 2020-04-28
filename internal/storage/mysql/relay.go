package mysql

import (
	"smarthome-home/internal/domain/relay"
)

type relayStorage struct {
	dbs *DBs
}

func NewRelay(dbs *DBs) (relay.Storager, error) {
	return &relayStorage{
		dbs: dbs,
	}, nil
}
