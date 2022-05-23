package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/gps"
	"github.com/stv0g/vand/pkg/pb"
	"github.com/stv0g/vand/pkg/types"
	"github.com/tarm/serial"
)

func init() {
	rootCmd.AddCommand(gpsCmd)
}

var gpsCmd = &cobra.Command{
	Use:   "gps",
	Short: "Start the GPS agent",
	Run:   runGPS,
}

func runGPS(cmd *cobra.Command, args []string) {
	client, err := newMQTTClient(&cfg.Broker, "vand-gps")
	if err != nil {
		log.Fatal(err)
	}

	gps, err := gps.New(&serial.Config{
		Name: cfg.GPS.Port,
		Baud: cfg.GPS.Baudrate,
	})
	if err != nil {
		log.Fatal(err)
	}

	lastPublish := time.Time{}
	lastPosition := types.Coordinate{}

	for {
		state, err := gps.GetState()
		if err != nil {
			log.Printf("Failed to get GPS update: %s", err)
			continue
		}

		newPosition := types.Coordinate{
			Latitude:  state.Latitude,
			Longitude: state.Longitude,
		}

		// Ignore invalid positions
		if newPosition == types.NullIsland {
			continue
		}

		if lastPosition != types.NullIsland && newPosition.DistanceTo(lastPosition) > 1000 {
			continue
		}

		// Owntracks
		now := time.Now()
		intervalElapsed := now.After(lastPublish.Add(cfg.GPS.MinInterval))
		significantChange := lastPosition.DistanceTo(newPosition) > cfg.GPS.MinDistance

		if intervalElapsed || significantChange {
			json.NewEncoder(os.Stdout).Encode(state)

			if err := PublishToOwntracks(client, state); err != nil {
				log.Printf("Error: failed to publish location update to OwnTracks: %s", err)
			}
			lastPublish = now
			lastPosition = newPosition
		}
	}
}

func PublishToOwntracks(client mqtt.Client, s *pb.GpsState) error {
	loc := s.LocationUpdate()
	msg, err := json.Marshal(&loc)
	if err != nil {
		return fmt.Errorf("failed to marshal location update: %w", err)
	}

	t := client.Publish(cfg.GPS.OwnTracks.Topic, 2, false, msg)
	go func() {
		<-t.Done()
		if t.Error() != nil {
			log.Printf("Failed to publish: %s", t.Error())
		}
	}()

	return nil
}
