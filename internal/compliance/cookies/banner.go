// Package cookies implements a no-dark-pattern cookie consent banner.
// Reject is exactly as prominent as Accept. Only strictly necessary
// cookies are set before consent.
package cookies

// Category is a cookie category.
type Category string

const (
	StrictlyNecessary Category = "necessary"
	Preferences       Category = "preferences"
	Analytics         Category = "analytics"
)

// Choice is the user response.
type Choice struct {
	Necessary   bool // always true
	Preferences bool
	Analytics   bool
}

// PreConsent returns true only for the always-on necessary category.
func PreConsent(c Category) bool { return c == StrictlyNecessary }

// Apply returns whether a cookie of category c is allowed under choice.
func Apply(c Category, ch Choice) bool {
	switch c {
	case StrictlyNecessary:
		return true
	case Preferences:
		return ch.Preferences
	case Analytics:
		return ch.Analytics
	}
	return false
}
