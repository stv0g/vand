// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
)

func init() {
	rootCmd.AddCommand(monitorCmd)
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "log messages to terminal",
	Run:   runMonitor,
}

func runMonitor(cmd *cobra.Command, args []string) {
	client, err := mqtt.NewClient(&cfg.Broker, "monitor", cfg.DataDir, true)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	topic := fmt.Sprintf("%s/#", cfg.Broker.Topic)

	client.Subscribe(topic, 2, func(c pmqtt.Client, m pmqtt.Message) {
		var sup pb.StateUpdatePoint

		if err := proto.Unmarshal(m.Payload(), &sup); err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			return
		}

		sup.Dump(log.Writer())
	})

	select {}
}
