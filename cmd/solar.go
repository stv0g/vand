package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/solar/renogy"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
)

func init() {
	rootCmd.AddCommand(solarCmd)
}

var solarCmd = &cobra.Command{
	Use:   "solar",
	Short: "Start the photo-voltaic agent",
	Run:   runSolar,
}

func runSolar(cmd *cobra.Command, args []string) {
	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGINT, syscall.SIGTERM)

	updateSig := make(chan os.Signal, 1)
	signal.Notify(updateSig, syscall.SIGUSR1)

	client, err := mqtt.NewClient(&cfg.Broker, "solar", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}

	topic := fmt.Sprintf("%s/update", cfg.Broker.Topic)

	d, err := renogy.NewDevice(cfg.Solar.Address)
	if err != nil {
		log.Fatalf("Failed to initialize device: %s", err)
	}

	tick := time.NewTicker(cfg.Solar.PollInterval)
out:
	for {
		sts, err := d.GetState()
		if err != nil {
			log.Printf("Failed to get state: %s", err)
			continue
		}

		sup := &pb.StateUpdatePoint{
			Solar: sts,
		}

		client.PublishUpdate(topic, sup)

		select {
		case <-termSig:
			break out

		case <-updateSig:
		case <-tick.C:
		}
	}
}
