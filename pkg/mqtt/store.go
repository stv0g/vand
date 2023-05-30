// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"bytes"
	"fmt"
	"log"
	"path"

	"github.com/dgraph-io/badger/v3"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

type store struct {
	*badger.DB

	path string
}

func newStore(clientID, dataDir string) (*store, error) {
	return &store{
		path: path.Join(dataDir, "mqtt", fmt.Sprintf("%s.db", clientID)),
	}, nil
}

func (s *store) Open() {
	var err error

	opts := badger.DefaultOptions(s.path)
	opts.ValueLogFileSize = 1 << 20

	if s.DB, err = badger.Open(opts); err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}
}

func (s *store) Put(key string, msg packets.ControlPacket) {
	wr := bytes.NewBuffer(nil)

	if err := msg.Write(wr); err != nil {
		log.Fatalf("Failed to write: %s", err)
	}

	if err := s.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), wr.Bytes())
	}); err != nil {
		log.Fatalf("Failed to list all keys: %s", err)
	}
}

func (s *store) Get(key string) packets.ControlPacket {
	buf := bytes.NewBuffer(nil)

	if err := s.DB.View(func(txn *badger.Txn) error {
		if it, err := txn.Get([]byte(key)); err == nil {
			return it.Value(func(val []byte) error {
				_, err := buf.Write(val)
				return err
			})
		} else {
			return err
		}
	}); err != nil {
		log.Fatalf("Failed to list all keys: %s", err)
	}

	pkt, _ := packets.ReadPacket(buf)
	return pkt
}

func (s *store) All() []string {
	all := []string{}

	if err := s.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			all = append(all, string(it.Item().Key()))
		}
		return nil
	}); err != nil {
		log.Fatalf("Failed to list all keys: %s", err)
	}

	return all
}

func (s *store) Del(key string) {
	if err := s.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	}); err != nil {
		log.Fatalf("Failed to list all keys: %s", err)
	}
}

func (s *store) Close() {
	if err := s.DB.Close(); err != nil {
		log.Fatalf("Failed to close database: %s", err)
	}
}

func (s *store) Reset() {
	if err := s.DB.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			if err := txn.Delete(it.Item().Key()); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatalf("failed to reset store: %s", err)
	}
}
