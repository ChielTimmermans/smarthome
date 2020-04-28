package main

import (
	"fmt"
	"log"
	"os"
	"smarthome-home/internal"
	"smarthome-home/internal/domain/relay"
	"smarthome-home/internal/push/mqtt/paho"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	pahoMQTT "github.com/eclipse/paho.mqtt.golang"
)

type Push struct {
	relay relay.Pusher
}

func initPusher(configMQTT *ConfigMQTT, mqttMode string, sa *internal.ServicesAvailable) (*Push, error) {
	p := &Push{}
	if err := initMQTT(configMQTT, p, mqttMode, sa); err != nil {
		return nil, err
	}
	return p, nil
}

func initMQTT(config *ConfigMQTT, p *Push, mqttMode string, sa *internal.ServicesAvailable) error {
	switch mqttMode {
	case "paho":
		if err := initPaho(config, p, sa); err != nil {
			return err
		}
	default:
		return fmt.Errorf("mqtt mode unknown. \t possible modes %s\t given mode: %s", possibleModes([]string{"paho"}), mqttMode)
	}
	return nil
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func initPaho(config *ConfigMQTT, p *Push, sa *internal.ServicesAvailable) (err error) {

	var client pahoMQTT.Client
	if client, err = newPahoConn(config); err != nil {
		return err
	}
	if p.relay, err = paho.NewRelay(client); err != nil {
		return err
	}
	sa.MQTT = true
	return nil
}

func newPahoConn(config *ConfigMQTT) (pahoMQTT.Client, error) {
	pahoMQTT.ERROR = log.New(os.Stdout, "[ERROR]", 0)
	opts := pahoMQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", config.Protocol, config.Host, config.Port))
	opts.SetClientID(config.UserName).SetPassword(config.Password)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	client := pahoMQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}
