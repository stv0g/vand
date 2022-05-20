package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/pb"
	"google.golang.org/protobuf/proto"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	rootCmd.AddCommand(bridgeCmd)
}

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "MQTT gateway to forward and/or translate protobuf payloads",
	Run:   runBridge,
}

var flatten = false

func bridgeMessageHandler(clientCar mqtt.Client, clientCloud mqtt.Client, msg mqtt.Message) {
	pl := msg.Payload()
	plr := bytes.NewBuffer(pl)

	zr, err := zlib.NewReader(plr)
	if err != nil {
		log.Fatal(err)
	}

	defer zr.Close()

	b, err := ioutil.ReadAll(zr)
	if err != nil {
		log.Fatal(err)
	}

	su := &pb.StateUpdateMessage{}

	proto.Unmarshal(b, su)

	fmt.Printf("Received message from topic %s of size %d:\n%s\n", msg.Topic(), len(msg.Payload()), su.String())

	// if flatten {
	// 	m := Flatten(su, "bus/state", "/")

	// 	for key, value := range m {
	// 		val := reflect.ValueOf(value)

	// 		if val.Kind() == reflect.Uint32 || val.Kind() == reflect.Float32 {
	// 			fmt.Printf("%s = %v\n", key, value)
	// 		} else {
	// 			fmt.Printf("%s = %v\n", key, val.Elem())
	// 		}
	// 	}
	// } else {

	// }
}

func runBridge(cmd *cobra.Command, args []string) {

	clientCar, err := newMQTTClient(&cfg.Broker, "vand-bridge")
	if err != nil {
		log.Fatal(err)
	}

	clientCloud, err := newMQTTClient(&cfg.BrokerCloud, "vand-bridge")
	if err != nil {
		log.Fatal(err)
	}

	topicCar := fmt.Sprintf("%s/#", cfg.Broker.Topic)
	// topicCloud := fmt.Sprintf("%s/#", cfg.BrokerCloud.Topic)

	t := clientCar.Subscribe(topicCar, 1, func(c mqtt.Client, m mqtt.Message) {
		bridgeMessageHandler(clientCar, clientCloud, m)
	})
	t.Wait()
	fmt.Printf("Subscribed to topic %s", topicCar)

	select {}
}
