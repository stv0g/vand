package store

import "github.com/stv0g/vand/pkg/pb"

type State struct {
	Solar       *pb.SolarState
	Battery     *pb.BatteryState
	Gps         *pb.GpsState
	Car         *pb.CarOBD2State
	Environment *pb.EnvironmentState
	Modem       *pb.ModemState
}

type Store struct {
	State State
}

func (s *Store) Update(msg *pb.StateUpdateMessage) error {
	return nil
}

func (s *Store) GetValues() map[string]string {
	return map[string]string{}
}
