package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/modem"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
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
	client, err := mqtt.NewClient(&cfg.Broker, "modem", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}

	topic := fmt.Sprintf("%s/update", cfg.Broker.Topic)

	modem, err := modem.New(cfg.Modem.Address, cfg.Modem.Username, cfg.Modem.Password)
	if err != nil {
		log.Fatal(err)
	}

	tick := time.NewTicker(cfg.Modem.PollInterval)
	for {
		sts, err := modem.GetState()
		if err != nil {
			log.Printf("Failed to get modem state: %s", err)
		}

		sup := &pb.StateUpdatePoint{
			Modem: sts,
		}

		client.PublishUpdate(topic, sup)

		<-tick.C
	}
}
