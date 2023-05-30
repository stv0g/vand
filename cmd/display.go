// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/display"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/store"
)

func init() {
	rootCmd.AddCommand(displayCmd)
}

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Start the display agent",
	Run:   runDisplay,
}

func runDisplay(cmd *cobra.Command, args []string) {
	client, err := mqtt.NewClient(&cfg.Broker, "display", cfg.DataDir, true)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %s", err)
	}

	topic := fmt.Sprintf("%s/#", cfg.Broker.Topic)

	store, err := store.NewStore(client, topic)
	if err != nil {
		log.Fatal(err)
	}

	disp, err := display.NewDisplay(&cfg.Display)
	if err != nil {
		log.Fatal(err)
	}

	pages, err := display.LoadPages(cfg.Display.Pages)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := disp.Play(pages, store); err != nil {
			log.Fatalf("failed to playback page: %s", err)
		}

		// showCanvas(dev)
		// time.Sleep(1 * time.Second)

		// playGif(dev)
	}()

	select {}
}
