// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

// Vivaro B 2017 Supported PIDs
//
// PID: 0x1		Monitor status since DTCs cleared. (Includes malfunction indicator lamp (MIL), status and number of DTCs, components tests, DTC readiness checks)
// PID: 0x4	ok	Calculated engine load
// PID: 0x5	ok	Engine coolant temperature
// PID: 0xc	ok	Engine speed
// PID: 0xd	ok	Vehicle speed
// PID: 0xf	ok	Intake air temperature
// PID: 0x10		Mass air flow sensor (MAF) air flow rate
// PID: 0x11	do	Throttle position
// PID: 0x1c	ign	OBD standards this vehicle conforms to
// PID: 0x1f	ok	Run time since engine start
// PID: 0x20	ign	PIDs supported [21 - 40]
// PID: 0x21		Distance traveled with malfunction indicator lamp (MIL) on
// PID: 0x23		Fuel Rail Gauge Pressure (diesel, or gasoline direct injection)
// PID: 0x24		Oxygen Sensor 1
// PID: 0x2c		Commanded EGR
// PID: 0x2d		EGR Error
// PID: 0x30		Warm-ups since codes cleared
// PID: 0x31	ok	Distance traveled since codes cleared
// PID: 0x33	ok	Absolute Barometric Pressure
// PID: 0x3c		Catalyst Temperature: Bank 1, Sensor 1
// PID: 0x3e		Catalyst Temperature: Bank 1, Sensor 2
// PID: 0x40	ign	PIDs supported [41 - 60]
// PID: 0x41		Monitor status this drive cycle
// PID: 0x42	ok	Control module voltage
// PID: 0x45		Relative throttle position
// PID: 0x46	ok	Ambient air temperature
// PID: 0x49		Accelerator pedal position D
// PID: 0x4a		Accelerator pedal position E
// PID: 0x4c		Commanded throttle actuator
// PID: 0x4f		Maximum value for Fuel–Air equivalence ratio, oxygen sensor voltage, oxygen sensor current, and intake manifold absolute pressure
// PID: 0x5c	ok	Engine oil temperature
// PID: 0x60	ign	PIDs supported [61 - 80]
// PID: 0x80	ign	PIDs supported [81 - A0]
// PID: 0x88		SCR Induce System

package obd2

import (
	"fmt"
	"log"

	"github.com/rzetterberg/elmobd"
	"github.com/stv0g/vand/pkg/pb"
)

type Device struct {
	*elmobd.Device

	supportedCommands *elmobd.SupportedCommands
}

func New(addr string, debug bool) (*Device, error) {
	d, err := elmobd.NewDevice(addr, debug)
	if err != nil {
		return nil, fmt.Errorf("failed to create new device: %w", err)
	}

	v, err := d.GetVoltage()
	if err != nil {
		log.Fatalf("Failed to get voltage: %s", err)
	}

	log.Printf("ELM327 system voltage: %f", v)

	return &Device{
		Device: d,
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
		elmobd.NewControlModuleVoltage(),
		elmobd.NewAmbientTemperature(),
		elmobd.NewAbsoluteBarometricPressure(),
		elmobd.NewEngineOilTemperature(),
	}

	s, err := d.GetIgnitionState()
	if err != nil {
		return nil, err
	}

	log.Printf("Ignition state: %v", s)

	if s == false {
		return nil, fmt.Errorf("ignition is off")
	}

	if d.supportedCommands == nil {
		if d.supportedCommands, err = d.CheckSupportedCommands(); err != nil {
			return nil, fmt.Errorf("failed to check supported commands: %w", err)
		}

		for pid := elmobd.OBDParameterID(0); pid != 0xff; pid++ {
			part, err := d.supportedCommands.GetPartByPID(pid)
			if err != nil {
				break
			}

			if part.SupportsPID(pid) {
				log.Printf("PID: %#x", pid)
			}
		}

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
	case *elmobd.AmbientTemperature:
		state.AmbientTemperature = float32(cmd.Value)
	case *elmobd.EngineLoad:
		state.EngineLoad = cmd.Value
	case *elmobd.EngineRPM:
		state.EngineRpm = cmd.Value
	case *elmobd.Fuel:
		state.FuelTankLevel = cmd.Value
	case *elmobd.DistSinceDTCClear:
		state.DistanceSinceDtcClear = float32(cmd.Value)
	case *elmobd.RuntimeSinceStart:
		state.RuntimeSinceStart = cmd.Value
	case *elmobd.ControlModuleVoltage:
		state.ControlModuleVoltage = cmd.Value
	case *elmobd.AbsoluteBarometricPressure:
		state.AbsoluteBarometricPressure = float32(cmd.Value)
	case *elmobd.EngineOilTemperature:
		state.EngineOilTemperature = float32(cmd.Value)
	default:
		return fmt.Errorf("unknown command")
	}

	return nil
}
