// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/bms/jbd"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
)

func init() {
	rootCmd.AddCommand(batteryCmd)
}

var batteryCmd = &cobra.Command{
	Use:   "battery",
	Short: "Start the Battery agent",
	Run:   runBattery,
}

func runBattery(cmd *cobra.Command, args []string) {
	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGINT, syscall.SIGTERM)

	updateSig := make(chan os.Signal, 1)
	signal.Notify(updateSig, syscall.SIGUSR1)

	client, err := mqtt.NewClient(&cfg.Broker, "battery", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}

	topic := fmt.Sprintf("%s/update", cfg.Broker.Topic)

	d, err := jbd.NewDevice(cfg.Battery.Address)
	if err != nil {
		log.Fatalf("Failed to initialize device: %s", err)
	}

	// d.SetFET(false, true)

	tick := time.NewTicker(cfg.Battery.PollInterval)
out:
	for {
		sts, err := d.GetState()
		if err != nil {
			log.Printf("Failed to get state: %s", err)
			continue
		}

		log.Printf("State: %+#v", sts)

		sup := &pb.StateUpdatePoint{
			Bat: sts,
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
