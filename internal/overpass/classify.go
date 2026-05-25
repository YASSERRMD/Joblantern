package overpass

import "strings"

// BuildingClassification is the derived land-use signal.
type BuildingClassification struct {
	PrimaryType   string         `json:"primary_type"`
	IsCommercial  bool           `json:"is_commercial"`
	IsResidential bool           `json:"is_residential"`
	IsMixed       bool           `json:"is_mixed"`
	TopTags       map[string]int `json:"top_tags"`
}

// commercialBuildingValues marks building=* values that are clearly
// commercial.
var commercialBuildingValues = map[string]bool{
	"commercial": true, "retail": true, "office": true,
	"warehouse": true, "supermarket": true, "industrial": true,
	"hotel": true, "kiosk": true,
}

// residentialBuildingValues marks building=* values that are clearly
// residential.
var residentialBuildingValues = map[string]bool{
	"residential": true, "house": true, "apartments": true,
	"detached": true, "dormitory": true, "terrace": true,
	"bungalow": true, "semi": true, "semidetached_house": true,
}

// Classify derives a BuildingClassification from Overpass elements.
func Classify(r *Response) BuildingClassification {
	c := BuildingClassification{TopTags: map[string]int{}}
	if r == nil {
		return c
	}
	for _, el := range r.Elements {
		for k, v := range el.Tags {
			key := k
			if v != "yes" && v != "" {
				key = k + "=" + v
			}
			c.TopTags[key]++

			switch k {
			case "building":
				if commercialBuildingValues[strings.ToLower(v)] {
					c.IsCommercial = true
				}
				if residentialBuildingValues[strings.ToLower(v)] {
					c.IsResidential = true
				}
			case "landuse":
				switch strings.ToLower(v) {
				case "residential":
					c.IsResidential = true
				case "commercial", "retail", "industrial":
					c.IsCommercial = true
				}
			case "office", "shop":
				c.IsCommercial = true
			}
		}
	}
	if c.IsCommercial && c.IsResidential {
		c.IsMixed = true
	}
	switch {
	case c.IsMixed:
		c.PrimaryType = "mixed"
	case c.IsCommercial:
		c.PrimaryType = "commercial"
	case c.IsResidential:
		c.PrimaryType = "residential"
	default:
		c.PrimaryType = "unknown"
	}
	return c
}
