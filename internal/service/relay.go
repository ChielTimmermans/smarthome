package service

import (
	"errors"
	"smarthome-home/internal/domain/relay"
)

type relayService struct {
	storage relay.Storager
	mqtt    relay.Pusher
	status  map[int]map[int]int
}

func NewRelay(storage relay.Storager, push relay.Pusher) (relay.Servicer, error) {
	if storage == nil {
		return nil, errors.New("relaystorage_nil")
	}
	return &relayService{
		storage: storage,
		mqtt:    push,
		status:  map[int]map[int]int{},
	}, nil
}

func (s *relayService) Enable(relayID, itemID int) error {
	s.Register(relayID, itemID)
	if err := s.mqtt.Enable(relayID, itemID); err != nil {
		s.status[relayID][itemID] = -1
		return err
	}
	s.status[relayID][itemID] = 1
	return nil
}

func (s *relayService) Disable(relayID, itemID int) error {
	s.Register(relayID, itemID)
	if err := s.mqtt.Disable(relayID, itemID); err != nil {
		s.status[relayID][itemID] = -1
		return err
	}
	s.status[relayID][itemID] = 0
	return nil
}

func (s *relayService) Toggle(relayID, itemID int) error {
	s.Register(relayID, itemID)
	if s.status[relayID][itemID] == 0 {
		if err := s.Enable(relayID, itemID); err != nil {
			return err
		}
	} else if s.status[relayID][itemID] == 1 {
		if err := s.Disable(relayID, itemID); err != nil {
			return err
		}
	} else {
		if err := s.Enable(relayID, itemID); err != nil {
			return err
		}
	}
	return nil
}

func (s *relayService) Register(relayID, itemID int) {
	if s.status[relayID] == nil {
		s.status[relayID] = map[int]int{}
		s.status[relayID][itemID] = -1
	}
}
