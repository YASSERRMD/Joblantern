// Package detection emits the rolling detection-rate dashboard.
package detection

import "time"

// Point is one daily measurement.
type Point struct {
	Day             time.Time
	DetectionRate   float64 // 0..1
	AdversarialRate float64
}

// History is a sliding window.
type History struct {
	points []Point
	max    int
}

// New returns a history with capacity max.
func New(max int) *History { return &History{max: max} }

// Add appends a point and trims to capacity.
func (h *History) Add(p Point) {
	h.points = append(h.points, p)
	if len(h.points) > h.max {
		h.points = h.points[len(h.points)-h.max:]
	}
}

// Series returns a copy of the points.
func (h *History) Series() []Point {
	out := make([]Point, len(h.points))
	copy(out, h.points)
	return out
}
