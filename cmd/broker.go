package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(brokerCmd)
}

var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "Start the MQTT broker",
	Run:   runBroker,
}

func runBroker(cmd *cobra.Command, args []string) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Create the new MQTT Server.
	server := mqtt.NewServer(nil)

	// Create a TCP listener on a standard port.
	address := fmt.Sprintf(":%d", cfg.Broker.Port)
	tcp := listeners.NewTCP("t1", address)

	// Add the listener to the server with default options (nil).
	if err := server.AddListener(tcp, nil); err != nil {
		log.Fatal(err)
	}

	go server.Serve()
	log.Println("Starting MQTT broker")

	<-done

	server.Close()
}
