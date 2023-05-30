// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pb

import (
	"time"

	"github.com/stv0g/vand/pkg/owntracks"
)

func (s *GpsState) LocationUpdate() *owntracks.Location {
	loc := &owntracks.Location{
		Type: owntracks.TypeLocation,

		Latitude:         float64(s.Latitude),
		Longitude:        float64(s.Longitude),
		Altitude:         int(s.Altitude),
		CourseOverGround: int(s.Cog),
		Velocity:         int(s.Speed),

		TrackerID:    "BS",
		SSID:         "Wanderlust",
		Connectivity: owntracks.ConnectivityMobile,
	}

	timeNow := time.Now()

	if s.Date != nil {
		timeFix := time.Date(
			int(s.Date.Year), time.Month(s.Date.Month), int(s.Date.Day),
			int(s.Time.Hour), int(s.Time.Minute), int(s.Time.Second), int(s.Time.Thousand*1e6), time.UTC,
		)
		loc.Timestamp = int(timeFix.Unix())

	}

	loc.CreatedAt = int(timeNow.Unix())

	return loc
}
