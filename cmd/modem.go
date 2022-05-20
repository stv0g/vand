package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/modem"
)

func init() {
	rootCmd.AddCommand(modemCmd)
}

var modemCmd = &cobra.Command{
	Use:   "modem",
	Short: "Start the modem agent",
	Run:   runModem,
}

func runModem(cmd *cobra.Command, args []string) {
	client, err := newMQTTClient(&cfg.Broker, "vand-modem")
	if err != nil {
		log.Fatal(err)
	}

	topic := fmt.Sprintf("%s/modem", cfg.Broker.Topic)

	modem, err := modem.New(cfg.Modem.Address, cfg.Modem.Username, cfg.Modem.Password)
	if err != nil {
		log.Fatal(err)
	}

	tick := time.NewTicker(cfg.Modem.PollInterval)

	for range tick.C {
		sts, err := modem.GetState()
		if err != nil {
			log.Printf("Failed to get modem state: %s", err)
		}

		pl, err := proto.Marshal(sts)
		if err != nil {
			log.Printf("Failed to marshal modem state: %s", err)
		}

		client.Publish(topic, 2, false, pl)
	}
}
