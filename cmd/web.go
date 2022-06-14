package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/store"
	"github.com/stv0g/vand/pkg/web"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Run:   runWeb,
}

func runWeb(cmd *cobra.Command, args []string) {
	client, err := mqtt.NewClient(&cfg.Broker, "web", cfg.DataDir, true)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %s", err)
	}

	topic := fmt.Sprintf("%s/#", cfg.Broker.Topic)

	store, err := store.NewStore(client, topic)
	if err != nil {
		log.Fatal(err)
	}

	web.Run(cfg, store, version, commit, date)
}
