package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/obd2"
)

func init() {
	rootCmd.AddCommand(carCmd)
}

var carCmd = &cobra.Command{
	Use:   "car",
	Short: "Start the car agent",
	Run:   runCar,
}

func runCar(cmd *cobra.Command, args []string) {
	client, err := newMQTTClient(&cfg.Broker, "vand-car")
	if err != nil {
		log.Fatal(err)
	}

	d, err := obd2.New(cfg.Car.Address)
	if err != nil {
		log.Fatalf("Failed to create device: %s", err)
	}

	topic := fmt.Sprintf("%s/car", cfg.Broker.Topic)

	for {
		sts, err := d.GetState()
		if err != nil {
			log.Printf("Failed to get car state: %s", err)
			continue
		}

		pl, err := proto.Marshal(sts)
		if err != nil {
			log.Printf("Failed to marshal status update: %s", err)
			continue
		}

		client.Publish(topic, 2, false, pl)
	}
}
