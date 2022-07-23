//go:build virtual

package main

import (
	"fmt"
	"log"

	"gioui.org/app"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/display"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/store"
)

func init() {
	rootCmd.AddCommand(displayVirtualCmd)
}

var displayVirtualCmd = &cobra.Command{
	Use:   "display_virtual",
	Short: "Start the virtual display agent",
	Run:   runDisplayVirtual,
}

func runDisplayVirtual(cmd *cobra.Command, args []string) {
	client, err := mqtt.NewClient(&cfg.Broker, "display-virtual", cfg.DataDir, true)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %s", err)
	}

	topic := fmt.Sprintf("%s/#", cfg.Broker.Topic)

	store, err := store.NewStore(client, topic)
	if err != nil {
		log.Fatal(err)
	}

	disp, err := display.NewVirtualDisplay()
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

	app.Main()
}
