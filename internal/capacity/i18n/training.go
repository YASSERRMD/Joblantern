// Package i18n indexes the per-language training materials shipped
// for the curriculum.
package i18n

// Material is one training artifact in one language.
type Material struct {
	ModuleID string
	Lang     string
	Title    string
	Path     string // path inside the docs/ngo-training tree
}

// Catalog returns the published catalogue. Production loads this from
// docs/ngo-training/index.yaml so translators can update without code
// changes.
func Catalog() []Material {
	langs := []string{"en", "hi", "bn", "tl", "id", "ar", "ur", "ne"}
	mods := []string{"overview", "deploy", "first-verdict", "intake", "appeals", "privacy", "regulator"}
	var out []Material
	for _, m := range mods {
		for _, l := range langs {
			out = append(out, Material{
				ModuleID: m,
				Lang:     l,
				Title:    m,
				Path:     "docs/ngo-training/" + l + "/" + m + ".md",
			})
		}
	}
	return out
}
