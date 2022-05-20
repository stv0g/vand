package owntracks

type Type string

const (
	TypeLocation Type = "location"

	BatteryUnknown   = iota
	BatteryUnplugged = iota
	BatteryCharging  = iota
	BatteryFull      = iota

	TriggerPing     = "p" // Ping issued randomly by background task (iOS,Android)
	TriggerCircular = "c" // Circular region enter/leave event (iOS,Android)
	TriggerBeacon   = "b" // Beacon region enter/leave event (iOS)
	TriggerResponse = "r" // Response to a reportLocation cmd message (iOS,Android)
	TriggerManual   = "u" // Manual publish requested by the user (iOS,Android)
	TriggerTimer    = "t" // Timer based publish in move move (iOS)
	TriggerUpdated  = "v" // Updated by Settings/Privacy/Locations Services/System Services/Frequent Locations monitoring (iOS)

	ConnectivityWifi    = "w" // Phone is connected to a WiFi connection (iOS,Android)
	ConnectivityOffline = "o" // Phone is offline (iOS,Android)
	ConnectivityMobile  = "m" // Mobile data (iOS,Android)

	MonitorModeSignificant = 1
	MonitorModeMove        = 2
)

type Location struct {
	Type Type `json:"_type"`

	Accuracy         int `json:"acc,omitempty"`  // Accuracy of the reported location in meters without unit (iOS,Android/integer/meters/optional)
	Altitude         int `json:"alt,omitempty"`  // Altitude measured above sea level (iOS,Android/integer/meters/optional)
	Battery          int `json:"batt,omitempty"` // Device battery level (iOS,Android/integer/percent/optional)
	BatteryStatus    int `json:"bs,omitempty"`   // Battery Status 0=unknown, 1=unplugged, 2=charging, 3=full (iOS, Android)
	CourseOverGround int `json:"cog,omitempty"`  // Course over ground (iOS/integer/degree/optional)

	Latitude  float64 `json:"lat"` // (iOS,Android/float/degree/required)
	Longitude float64 `json:"lon"` // (iOS,Android/float/degree/required)

	Radius  int    `json:"rad,omitempty"` // Radius around the region when entering/leaving (iOS/integer/meters/optional)
	Trigger string `json:"t,omitempty"`   // Trigger for the location report (iOS,Android/string/optional)

	TrackerID        string  `json:"tid,omitempty"`  // Tracker ID used to display the initials of a user (iOS,Android/string/optional) required for http mode
	Timestamp        int     `json:"tst"`            // UNIX epoch timestamp in seconds of the location fix (iOS,Android/integer/epoch/required)
	VerticalAccuracy int     `json:"vac,omitempty"`  // Vertical accuracy of the alt element (iOS/integer/meters/optional)
	Velocity         int     `json:"vel,omitempty"`  // Velocity (iOS,Android/integer/kmh/optional)
	Pressure         float64 `json:"p,omitempty"`    // Barometric pressure (iOS/float/kPa/optional/extended data)
	Connectivity     string  `json:"conn,omitempty"` // Internet connectivity status (route to host) when the message is created (iOS,Android/string/optional/extended data)

	Topic string `json:"topic,omitempty"` // (only in HTTP payloads) contains the original publish topic (e.g. owntracks/jane/phone). (iOS,Android >= 2.4,string)

	InRegions   []string `json:"inregions,omitempty"` // Contains a list of regions the device is currently in (e.g. ["Home","Garage"]). Might be empty. (iOS,Android/list of strings/optional)
	InRIDs      []string `json:"inrids,omitempty"`    // Contains a list of region IDs the device is currently in (e.g. ["6da9cf","3defa7"]). Might be empty. (iOS,Android/list of strings/optional)
	SSID        string   `json:"SSID,omitempty"`      // If available, is the unique name of the WLAN. (iOS,string/optional)
	BSSID       string   `json:"BSSID,omitempty"`     // If available, identifies the access point. (iOS,string/optional)
	CreatedAt   int      `json:"created_at"`          // Identifies the time at which the message is constructed (vs. tst which is the timestamp of the GPS fix) (iOS,Android)
	MonitorMode int      `json:"m,omitempty"`         // Identifies the monitoring mode at which the message is constructed (significant=1, move=2) (iOS/integer/optional)
}
