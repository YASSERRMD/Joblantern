// Package sla monitors per-tenant SLA and uptime. The metric we
// commit to is verdict latency p99 — under target for the tenant's
// chosen plan.
package sla

import "time"

// Target is the SLA target per plan.
type Target struct {
	PlanID         string
	UptimePct      float64
	VerdictP99Ms   int
	SupportHours   string
}

// Plans is the canonical SLA table.
var Plans = []Target{
	{PlanID: "ngo-free", UptimePct: 99.0, VerdictP99Ms: 15000, SupportHours: "best-effort"},
	{PlanID: "pro", UptimePct: 99.5, VerdictP99Ms: 8000, SupportHours: "business hours"},
	{PlanID: "custom", UptimePct: 99.9, VerdictP99Ms: 5000, SupportHours: "24/5"},
}

// MonthlyReport is the rollup emitted to the tenant dashboard.
type MonthlyReport struct {
	Month         time.Month
	Year          int
	UptimePct     float64
	VerdictP99Ms  int
	BreachesCount int
}
