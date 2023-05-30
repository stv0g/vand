// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package gps

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/adrianmo/go-nmea"
	"github.com/stv0g/vand/pkg/pb"
	"github.com/tarm/serial"
)

const (
	KnotsInKPH = 1.852
)

type Device struct {
	port *serial.Port

	state pb.GpsState
}

func New(cfg *serial.Config) (*Device, error) {
	var err error

	d := &Device{}

	if d.port, err = serial.OpenPort(cfg); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Device) GetSentence() (nmea.Sentence, error) {
	scanner := bufio.NewScanner(d.port)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read from port: %w", err)
		}

		return nil, nil
	}

	t := scanner.Text()
	s, err := nmea.Parse(t)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NMEA sentence: %s: %w", t, err)
	}

	return s, nil
}

func (d *Device) GetState() (*pb.GpsState, error) {
	s, err := d.GetSentence()
	if err != nil {
		return nil, err
	}

	switch s := s.(type) {
	case nmea.GGA:
		d.setTime(&s.Time)

		d.state.Latitude = s.Latitude
		d.state.Longitude = s.Longitude
		d.state.Altitude = s.Altitude

		if q, err := strconv.Atoi(s.FixQuality); err == nil {
			d.state.Fix = pb.GpsFix(q)
		}

		d.state.SatsInUse = uint32(s.NumSatellites)

	case nmea.GSA:
		switch s.FixType {
		case nmea.Fix2D:
			d.state.FixMode = pb.GpsFixMode_FIX_MODE_2D
		case nmea.Fix3D:
			d.state.FixMode = pb.GpsFixMode_FIX_MODE_3D
		case nmea.FixNone:
			d.state.FixMode = pb.GpsFixMode_FIX_MODE_INVALID
		}

		d.state.SatsIdInUse = []int64{}
		for _, sv := range s.SV {
			if n, err := strconv.Atoi(sv); err == nil {
				d.state.SatsIdInUse = append(d.state.SatsIdInUse, int64(n))
			}
		}

		d.state.DopH = s.HDOP
		d.state.DopP = s.PDOP
		d.state.DopV = s.VDOP

	case nmea.RMC:
		d.setTime(&s.Time)
		d.setDate(&s.Date)

		d.state.Latitude = s.Latitude
		d.state.Longitude = s.Longitude

		d.state.Speed = s.Speed * KnotsInKPH
		d.state.Cog = s.Course

		d.state.Valid = s.Validity == "A"

	case nmea.ZDA:
		d.setTime(&s.Time)
		d.setDate(&nmea.Date{
			DD: int(s.Day),
			MM: int(s.Month),
			YY: int(s.Year),
		})

	case nmea.GSV:
		d.state.SatsInView = uint32(s.NumberSVsInView)

		if s.MessageNumber == 1 {
			d.state.SatsDescInView = []*pb.GpsSatellite{}
		}

		for _, sv := range s.Info {
			d.state.SatsDescInView = append(d.state.SatsDescInView, &pb.GpsSatellite{
				Num:       sv.SVPRNNumber,
				Azimuth:   uint32(sv.Azimuth),
				Elevation: uint32(sv.Elevation),
				Snr:       uint32(sv.SNR),
			})
		}
	}

	return &d.state, nil
}

func (d *Device) setTime(t *nmea.Time) {
	d.state.Time = &pb.GpsTime{
		Hour:     uint32(t.Hour),
		Minute:   uint32(t.Minute),
		Second:   uint32(t.Second),
		Thousand: uint32(t.Millisecond),
	}
}

func (d *Device) setDate(t *nmea.Date) {
	d.state.Date = &pb.GpsDate{
		Year:  2000 + uint32(t.YY),
		Month: uint32(t.MM),
		Day:   uint32(t.DD),
	}
}
