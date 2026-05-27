// mcp-cheki wraps open course/program databases used to corroborate
// that a claimed program exists in published catalogues. Sources
// include IEEE Xplore (subject metadata only), DAAD's Hochschulkompass,
// UCAS, and similar public catalogues — used per source ToS.
package main

import (
	"log/slog"
	"os"
)

// catalogue is a source name we know how to query.
type catalogue struct {
	ID      string
	Name    string
	Country string
	BaseURL string
}

var catalogues = []catalogue{
	{ID: "daad", Name: "DAAD Hochschulkompass", Country: "DE", BaseURL: "https://www.hochschulkompass.de"},
	{ID: "ucas", Name: "UCAS", Country: "UK", BaseURL: "https://www.ucas.com"},
	{ID: "ieee", Name: "IEEE Xplore (metadata)", Country: "", BaseURL: "https://ieeexplore.ieee.org"},
}

func main() {
	slog.New(slog.NewJSONHandler(os.Stderr, nil)).Info("mcp-cheki starting", "sources", len(catalogues))
}
