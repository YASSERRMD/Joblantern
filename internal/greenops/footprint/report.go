// Package footprint computes a monthly CO2-equivalent footprint per
// verdict.
package footprint

import "time"

// Sample is one verdict's resource use.
type Sample struct {
	VerdictID      string
	At             time.Time
	CPUSeconds     float64
	LLMtokens      int
	BytesEgress    int64
	Region         string
	GridGCO2PerKWh float64
}

// PerKwhCPU is the assumed conversion from CPU-seconds to kWh (very
// conservative; replace with measured values per fleet).
const PerKwhCPU = 0.000035

// PerToken1k is the inference cost per 1k tokens in kWh.
const PerToken1k = 0.0006

// PerGBEgress is kWh per GB egress (global average).
const PerGBEgress = 0.06

// GramsCO2 returns the grams CO2-eq for one verdict.
func (s Sample) GramsCO2() float64 {
	kwh := s.CPUSeconds*PerKwhCPU + float64(s.LLMtokens)/1000*PerToken1k + float64(s.BytesEgress)/(1024*1024*1024)*PerGBEgress
	return kwh * s.GridGCO2PerKWh
}

// MonthlySummary is the aggregate.
type MonthlySummary struct {
	Month           time.Month
	Year            int
	Verdicts        int
	TotalGramsCO2   float64
	PerVerdictGrams float64
}
