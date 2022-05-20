package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/faker/v3"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	pb "github.com/stv0g/vand/pkg/pb"
)

func init() {
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish test messages to broker",
	Run:   runPublish,
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

func getTestData() *pb.StateUpdateMessage {
	// return &pb.StateUpdate{
	// 	Gps: &pb.GpsState{
	// 		Latitude: 12,
	// 		Longitude: 13
	// 	},
	// }

	su := &pb.StateUpdateMessage{}

	_ = faker.SetRandomMapAndSliceSize(10)

	err := faker.FakeData(su)
	if err != nil {
		log.Fatalf("Failed to generate fake data: %s", err)
	}

	return su
}

func runPublish(cmd *cobra.Command, args []string) {
	client, err := newMQTTClient(&cfg.Broker, "vand-publisher")
	if err != nil {
		log.Fatal(err)
	}

	topic := "bus"
	su := getTestData()

	pl, err := proto.Marshal(su)
	if err != nil {
		log.Fatal(err)
	}

	var plc bytes.Buffer
	w := zlib.NewWriter(&plc)
	w.Write(pl)
	w.Close()

	log.Printf("Sending message to topic %s of size %d (uncompressed %d):\n%s\n", topic, plc.Len(), len(pl), su.String())

	for {
		token := client.Publish(topic, 0, false, plc)
		token.Wait()
		time.Sleep(time.Second)
	}
}
