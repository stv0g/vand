package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/devices/gps"
	"github.com/stv0g/vand/pkg/mqtt"
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
	client, err := mqtt.NewClient(&cfg.Broker, "gps", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	topic := fmt.Sprintf("%s/update", cfg.Broker.Topic)

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
		sts, err := gps.GetState()
		if err != nil {
			log.Printf("Failed to get GPS update: %s", err)
			continue
		}

		newPosition := types.Coordinate{
			Latitude:  sts.Latitude,
			Longitude: sts.Longitude,
		}

		// Ignore invalid positions
		if (sts.Latitude == 0 && sts.Longitude == 0) || sts.Time == nil || sts.Date == nil || int(sts.Date.Year) != time.Now().Year() {
			continue
		}

		if lastPosition != types.NullIsland && newPosition.DistanceTo(lastPosition) > 1000 {
			lastPosition = newPosition
			continue
		}

		now := time.Now()
		intervalElapsed := now.After(lastPublish.Add(cfg.GPS.MinInterval))
		significantChange := lastPosition.DistanceTo(newPosition) > cfg.GPS.MinDistance

		if intervalElapsed || significantChange {
			// Owntracks
			if err := PublishToOwntracks(client, sts); err != nil {
				log.Printf("Error: failed to publish location update to OwnTracks: %s", err)
			}

			// Normal Update
			sup := &pb.StateUpdatePoint{
				Gps: sts,
			}

			client.PublishUpdate(topic, sup)

			lastPublish = now
			lastPosition = newPosition
		}
	}
}

func PublishToOwntracks(client *mqtt.Client, s *pb.GpsState) error {
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
