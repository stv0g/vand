package obd2

import (
	"fmt"
	"log"

	"github.com/rzetterberg/elmobd"
	"github.com/stv0g/vand/pkg/pb"
)

type Device struct {
	*elmobd.Device

	supportedCommands elmobd.SupportedCommands
}

func New(addr string) (*Device, error) {
	dev, err := elmobd.NewDevice(addr, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create new device: %w", err)
	}

	supported, err := dev.CheckSupportedCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to check supported commands: %w", err)
	}

	return &Device{
		Device:            dev,
		supportedCommands: *supported,
	}, nil
}

func (d *Device) GetState() (*pb.CarOBD2State, error) {

	var state pb.CarOBD2State

	cmds := []elmobd.OBDCommand{
		elmobd.NewThrottlePosition(),
		elmobd.NewVehicleSpeed(),
		elmobd.NewCoolantTemperature(),
		elmobd.NewEngineLoad(),
		elmobd.NewEngineRPM(),
		elmobd.NewFuel(),
		elmobd.NewIntakeAirTemperature(),
		elmobd.NewDistSinceDTCClear(),
		elmobd.NewRuntimeSinceStart(),
	}

	for _, cmd := range d.supportedCommands.FilterSupported(cmds) {
		cmd, err := d.RunOBDCommand(cmd)
		if err != nil {
			return nil, fmt.Errorf("failed run %s command: %w", cmd.Key(), err)
		}

		if err := d.updateStateWithCommand(cmd, &state); err != nil {
			log.Printf("Failed to run OBD2 command: %s", err)
		}
	}

	return &state, nil
}

func (d *Device) updateStateWithCommand(cmd elmobd.OBDCommand, state *pb.CarOBD2State) error {
	switch cmd := cmd.(type) {
	case *elmobd.ThrottlePosition:
		state.ThrottlePercentage = cmd.Value / 100
	case *elmobd.VehicleSpeed:
		state.VehicleSpeed = float32(cmd.Value)
	case *elmobd.CoolantTemperature:
		state.CoolantTemperature = float32(cmd.Value)
	case *elmobd.IntakeAirTemperature:
		state.IntakeAirTemperature = float32(cmd.Value)
	case *elmobd.EngineLoad:
		state.EngineLoad = cmd.Value
	case *elmobd.EngineRPM:
		state.EngineRpm = cmd.Value
	case *elmobd.Fuel:
		state.FuelTankLevel = cmd.Value
	case *elmobd.DistSinceDTCClear:
		state.DistanceSinceDtcClear = float32(cmd.Value)
	case *elmobd.RuntimeSinceStart:

	default:
		return fmt.Errorf("unknown command")
	}

	return nil
}
