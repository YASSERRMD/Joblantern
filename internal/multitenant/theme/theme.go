// Package theme stores per-tenant theming overrides — colour, logo,
// locale defaults. Values are loaded once at tenant boot and cached.
package theme

// Theme is the per-tenant override bundle.
type Theme struct {
	TenantID    string
	BrandName   string
	LogoURL     string
	Primary     string // hex
	Accent      string // hex
	DefaultLang string // BCP-47
}

// Default returns the Joblantern default.
func Default() Theme {
	return Theme{BrandName: "Joblantern", Primary: "#0e1a2b", Accent: "#3b6", DefaultLang: "en"}
}

// Apply overlays non-empty fields of override on top of base.
func Apply(base, override Theme) Theme {
	if override.BrandName != "" {
		base.BrandName = override.BrandName
	}
	if override.LogoURL != "" {
		base.LogoURL = override.LogoURL
	}
	if override.Primary != "" {
		base.Primary = override.Primary
	}
	if override.Accent != "" {
		base.Accent = override.Accent
	}
	if override.DefaultLang != "" {
		base.DefaultLang = override.DefaultLang
	}
	return base
}
