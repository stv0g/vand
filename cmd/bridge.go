// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
	"github.com/stv0g/vand/pkg/types"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	rootCmd.AddCommand(bridgeCmd)
}

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "MQTT gateway to forward and/or translate protobuf payloads",
	Run:   runBridge,
}

var flatten = true

func runBridge(cmd *cobra.Command, args []string) {
	clientCar, err := mqtt.NewClient(&cfg.Broker, "bridge-car", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}

	clientCloud, err := mqtt.NewClient(&cfg.BrokerCloud, "bridge-cloud", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}

	topicCar := fmt.Sprintf("%s/update", cfg.Broker.Topic)

	t := clientCar.Subscribe(topicCar, 0, func(clientCar pmqtt.Client, msg pmqtt.Message) {
		bridgeMessageHandler(clientCar, clientCloud, msg)
	})
	t.Wait()

	log.Printf("Subscribed to topic %s", topicCar)

	select {}
}

func bridgeMessageHandler(clientCar pmqtt.Client, clientCloud pmqtt.Client, msg pmqtt.Message) {
	var sup pb.StateUpdatePoint
	if err := proto.Unmarshal(msg.Payload(), &sup); err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}

	log.Printf("Forwarding: %+#v", &sup)

	if flatten {
		m := types.Flatten(&sup, "/")

		for key, value := range m {
			topic := fmt.Sprintf("%s/flat/%s", cfg.BrokerCloud.Topic, key)
			val, err := json.Marshal(value)
			if err != nil {
				log.Printf("Failed to marshal value: %s", err)
			}

			clientCloud.Publish(topic, 2, false, val)
		}
	} else {
		clientCloud.Publish(cfg.BrokerCloud.Topic, msg.Qos(), msg.Retained(), msg.Payload())
	}
}
