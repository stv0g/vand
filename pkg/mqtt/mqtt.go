package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stv0g/vand/pkg/pb"
	"google.golang.org/protobuf/proto"
)

type BrokerConfig struct {
	Hostname string `yaml:"hostname"`
	Port     uint16 `yaml:"port"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`

	Topic string `yaml:"topic"`
}

func mqttConnectHandler(client pmqtt.Client) {
	log.Println("Connected to MQTT broker")
}

func mqttConnectLostHandler(client pmqtt.Client, err error) {
	log.Printf("Connect lost from MQTT broker: %v", err)
}

func mqttConnectAttemptHandler(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
	log.Printf("Attempting connection to MQTT broker: %s", broker)

	return tlsCfg
}

type Client struct {
	pmqtt.Client

	store *store
}

func NewClient(broker *BrokerConfig, clientID, dataDir string, clean bool) (*Client, error) {
	var err error

	opts := pmqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker.Hostname, broker.Port))

	opts.ClientID = clientID
	opts.Username = broker.Username
	opts.Password = broker.Password
	opts.OnConnect = mqttConnectHandler
	opts.OnConnectAttempt = mqttConnectAttemptHandler
	opts.OnConnectionLost = mqttConnectLostHandler
	opts.CleanSession = clean

	if opts.Store, err = newStore(clientID, dataDir); err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	client := &Client{
		Client: pmqtt.NewClient(opts),
	}
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func (m *Client) Close() error {
	if m.store != nil {
		m.store.Close()
	}

	return nil
}

func (m *Client) PublishUpdate(topic string, sup *pb.StateUpdatePoint) pmqtt.Token {
	pl, err := proto.Marshal(sup)
	if err != nil {
		log.Printf("Failed to marshal modem state: %s", err)
	}

	return m.Publish(topic, 2, false, pl)
}
