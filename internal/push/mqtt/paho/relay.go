package paho

import (
	"errors"
	"smarthome-home/internal/domain/relay"
	"strconv"

	pahoMQTT "github.com/eclipse/paho.mqtt.golang"
)

type relayPusher struct {
	client pahoMQTT.Client
}

func NewRelay(c pahoMQTT.Client) (relay.Pusher, error) {
	return &relayPusher{
		client: c,
	}, nil
}

func (p *relayPusher) Enable(relayID, itemID int) error {
	if relayID == 0 {
		token := p.client.Publish("ChielT/Relay/Enable", 0, false, strconv.Itoa(itemID))
		if token.Wait() {
			return token.Error()
		}
	} else {
		return errors.New("relay_not_found")
	}
	return nil
}

func (p *relayPusher) Disable(relayID, itemID int) error {
	if relayID == 0 {
		token := p.client.Publish("ChielT/Relay/Disable", 0, false, strconv.Itoa(itemID))
		if token.Wait() {
			return token.Error()
		}
	} else {
		return errors.New("relay_not_found")
	}
	return nil
}
