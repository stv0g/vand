// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/obd2"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
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
	client, err := mqtt.NewClient(&cfg.Broker, "car", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}

	d, err := obd2.New(cfg.Car.Address, cfg.Debug)
	if err != nil {
		log.Fatalf("Failed to create device: %s", err)
	}

	topic := fmt.Sprintf("%s/update", cfg.Broker.Topic)

	tick := time.NewTicker(cfg.Car.PollInterval)
	for {
		sts, err := d.GetState()
		if err != nil {
			log.Printf("Failed to get car state: %s", err)
			continue
		}

		sup := &pb.StateUpdatePoint{
			Car: sts,
		}

		client.PublishUpdate(topic, sup)

		<-tick.C
	}
}
