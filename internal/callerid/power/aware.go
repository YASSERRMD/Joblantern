// Package power adjusts call-screening behaviour to battery state.
// Low battery prefers the on-device cache and skips the round-trip
// to the server.
package power

// Mode selects the lookup behaviour.
type Mode string

const (
	ModeNormal       Mode = "normal"
	ModeBatterySaver Mode = "battery_saver"
)

// FromBattery returns the recommended Mode for a battery level (0..1)
// and whether the device is currently charging.
func FromBattery(level float64, charging bool) Mode {
	if charging {
		return ModeNormal
	}
	if level < 0.15 {
		return ModeBatterySaver
	}
	return ModeNormal
}

// AllowsNetworkLookup reports whether the mode permits a network call.
func AllowsNetworkLookup(m Mode) bool { return m == ModeNormal }
