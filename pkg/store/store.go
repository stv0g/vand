// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"log"
	"sync"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/imdario/mergo"
	"github.com/stv0g/vand/pkg/mqtt"
	"github.com/stv0g/vand/pkg/pb"
	"github.com/stv0g/vand/pkg/types"
)

type Store struct {
	client *mqtt.Client

	subs     map[chan *pb.StateUpdatePoint]struct{}
	subsLock sync.RWMutex

	State pb.StateUpdatePoint
}

func NewStore(client *mqtt.Client, topic string) (*Store, error) {
	s := &Store{
		subs: map[chan *pb.StateUpdatePoint]struct{}{},
	}

	client.Subscribe(topic, 2, s.messageHandler)

	return s, nil
}

func (s *Store) Update(sup *pb.StateUpdatePoint) error {
	s.Notify(sup)

	return mergo.Merge(&s.State, sup, mergo.WithOverride)
}

func (s *Store) messageHandler(client pmqtt.Client, msg pmqtt.Message) {
	var sup pb.StateUpdatePoint

	if err := proto.Unmarshal(msg.Payload(), &sup); err != nil {
		log.Printf("Failed to unmarshal payload: %s", err)
	}

	if err := s.Update(&sup); err != nil {
		log.Printf("Failed to update store: %s", err)
	}
}

func (s *Store) Flatten(sep string) map[string]interface{} {
	return types.Flatten(&s.State, sep)
}

func (s *Store) Notify(sup *pb.StateUpdatePoint) {
	s.subsLock.RLock()
	defer s.subsLock.RUnlock()

	for ch := range s.subs {
		ch <- sup
	}
}

func (s *Store) Subscribe() (chan *pb.StateUpdatePoint, error) {
	ch := make(chan *pb.StateUpdatePoint)

	s.subsLock.Lock()
	defer s.subsLock.Unlock()

	s.subs[ch] = struct{}{}

	return ch, nil
}

func (s *Store) Unsubscribe(ch chan *pb.StateUpdatePoint) {
	s.subsLock.Lock()
	defer s.subsLock.Unlock()

	delete(s.subs, ch)

	close(ch)
}
