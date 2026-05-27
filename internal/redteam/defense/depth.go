// Package defense scores how many distinct layers caught each attack.
// Multiple-layer catches indicate genuine defense-in-depth.
package defense

// Layer is one named defense surface.
type Layer string

const (
	LayerRules     Layer = "rules"
	LayerAgent     Layer = "agent"
	LayerHumanMod  Layer = "moderation"
	LayerRegistry  Layer = "registry"
	LayerNetwork   Layer = "network"
)

// Catch records which layers stopped a given attack.
type Catch struct {
	AttackID string
	Layers   []Layer
}

// Score returns the count of distinct layers that caught the attack.
func (c Catch) Score() int {
	seen := map[Layer]bool{}
	for _, l := range c.Layers {
		seen[l] = true
	}
	return len(seen)
}

// Healthy returns true if the catch reflects at least two layers.
func (c Catch) Healthy() bool { return c.Score() >= 2 }
