//go:build mockup

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
	rootCmd.AddCommand(displayCmd)
}

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Start the display agent",
	Run:   runDisplay,
}

func runDisplay(cmd *cobra.Command, args []string) {
	client, err := mqtt.NewClient(&cfg.Broker, "display-mockup", cfg.DataDir, true)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %s", err)
	}

	topic := fmt.Sprintf("%s/#", cfg.Broker.Topic)

	store, err := store.NewStore(client, topic)
	if err != nil {
		log.Fatal(err)
	}

	disp, err := display.NewMockupDisplay()
	if err != nil {
		log.Fatal(err)
	}

	pages := map[string]*display.Page{}
	for id, page := range cfg.Display.Pages {
		pages[id] = &display.Page{DisplayPage: page}
	}

	for id, page := range pages {
		if page.Next != "" {
			page.next = pages[page.Next]
		}

		if page.Over != "" {
			page.over = pages[page.Over]
		}
	}

	log.Printf("Pages: %+#v", pages)

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
