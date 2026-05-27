// Package grid wraps the Electricity Maps API to fetch real-time
// carbon intensity for a region.
//
// The package only exposes the small subset Joblantern needs: a
// region -> gCO2e/kWh lookup and a "low-carbon hours" predicate.
package grid

import (
	"errors"
	"time"
)

// Intensity is gCO2-equivalent per kWh.
type Intensity float64

// Sample is one measurement.
type Sample struct {
	Region   string
	At       time.Time
	Value    Intensity
}

// Client is the abstraction.
type Client interface {
	Now(region string) (Sample, error)
	Forecast(region string, hours int) ([]Sample, error)
}

// LowCarbonThreshold is the intensity below which a region counts as
// "low-carbon hour" for scheduling purposes. The number is region-
// dependent in production; the default is a conservative global mean.
const LowCarbonThreshold Intensity = 250

// IsLowCarbon reports whether the sample is below the threshold.
func IsLowCarbon(s Sample) bool { return s.Value <= LowCarbonThreshold }

// Average computes the mean of a series of samples.
func Average(samples []Sample) (Intensity, error) {
	if len(samples) == 0 {
		return 0, errors.New("no samples")
	}
	sum := 0.0
	for _, s := range samples {
		sum += float64(s.Value)
	}
	return Intensity(sum / float64(len(samples))), nil
}
