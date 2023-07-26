// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
	"google.golang.org/protobuf/proto"
)

func init() {
	rootCmd.AddCommand(loggerCmd)
}

var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "log messages to disk",
	Run:   runLogger,
}

func runLogger(cmd *cobra.Command, args []string) {
	client, err := mqtt.NewClient(&cfg.Broker, "logger", cfg.DataDir, false)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	dbPath := fmt.Sprintf("%s/logger.db", cfg.DataDir)
	dbOpts := badger.DefaultOptions(dbPath)
	dbOpts.ValueLogFileSize = 1 << 20

	db, err := badger.Open(dbOpts)
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}
	defer db.Close()

	topic := fmt.Sprintf("%s/#", cfg.Broker.Topic)

	client.Subscribe(topic, 2, func(c pmqtt.Client, m pmqtt.Message) {
		var sup pb.StateUpdatePoint

		if err := proto.Unmarshal(m.Payload(), &sup); err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			return
		}

		if err := logUpdate(db, &sup); err != nil {
			log.Printf("Failed to log update: %s", err)
			return
		}
	})

	select {}
}

func logUpdate(db *badger.DB, sup *pb.StateUpdatePoint) error {
	return db.Update(func(txn *badger.Txn) error {
		ts := sup.Timestamp.Time().Format(time.RFC3339)
		key := fmt.Sprintf("update/%s", ts)

		value, err := proto.Marshal(sup)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		return txn.Set([]byte(key), value)
	})
}
