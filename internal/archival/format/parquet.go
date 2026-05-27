// Package format declares the archival file formats. We pick formats
// that survive: Parquet for the columnar data, JSON-LD for the
// schema and the linked-data sidecar.
package format

// File is one archive file specification.
type File struct {
	Name     string
	Format   string // "parquet" | "json-ld" | "csv" | "manifest"
	Required bool
}

// Defaults returns the canonical file set per archive.
func Defaults() []File {
	return []File{
		{Name: "verdicts.parquet", Format: "parquet", Required: true},
		{Name: "evidence.parquet", Format: "parquet", Required: true},
		{Name: "schema.json-ld", Format: "json-ld", Required: true},
		{Name: "manifest.json", Format: "manifest", Required: true},
		{Name: "README.md", Format: "markdown", Required: true},
	}
}
