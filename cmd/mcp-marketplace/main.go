// mcp-marketplace wraps marketplace platforms with open data where
// ToS permits.
package main

import (
	"log/slog"
	"os"
)

// platform is one supported marketplace platform.
type platform struct {
	ID      string
	Name    string
	Country string
	BaseURL string
}

var platforms = []platform{
	{ID: "olx-pl", Name: "OLX (PL)", Country: "PL"},
	{ID: "carousell-sg", Name: "Carousell (SG)", Country: "SG"},
	{ID: "facebook-marketplace", Name: "Facebook Marketplace", Country: ""},
}

func main() {
	slog.New(slog.NewJSONHandler(os.Stderr, nil)).Info("mcp-marketplace starting", "platforms", len(platforms))
}
