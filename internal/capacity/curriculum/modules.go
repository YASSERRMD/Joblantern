// Package curriculum lays out the NGO partner onboarding curriculum.
// The curriculum is short by design: any operator should reach first
// verdict within one working day.
package curriculum

// Module is one curriculum unit.
type Module struct {
	ID       string
	Title    string
	Minutes  int
	Required bool
}

// Defaults returns the canonical curriculum.
func Defaults() []Module {
	return []Module{
		{ID: "overview", Title: "What Joblantern is (and isn't)", Minutes: 20, Required: true},
		{ID: "deploy", Title: "Deploy your instance (deploy-in-a-box)", Minutes: 30, Required: true},
		{ID: "first-verdict", Title: "Run your first verdict", Minutes: 20, Required: true},
		{ID: "intake", Title: "Build your intake form", Minutes: 30, Required: true},
		{ID: "appeals", Title: "How appeals work", Minutes: 30, Required: false},
		{ID: "privacy", Title: "Privacy posture for your jurisdiction", Minutes: 40, Required: true},
		{ID: "regulator", Title: "Engaging your country's regulator", Minutes: 30, Required: false},
	}
}
