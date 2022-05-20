package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stv0g/vand/pkg/config"
)

var mqttConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var mqttConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

func newMQTTClient(broker *config.Broker, clientID string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker.Hostname, broker.Port))
	opts.SetClientID(clientID)
	opts.SetUsername(broker.Username)
	opts.SetPassword(broker.Password)

	opts.OnConnect = mqttConnectHandler
	opts.OnConnectionLost = mqttConnectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
