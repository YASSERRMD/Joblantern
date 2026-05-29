// Package export emits the recruitment graph in common formats so
// investigators can drop it into Gephi, Cytoscape, or NetworkX.
package export

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Node is the export node type.
type Node struct {
	ID    string
	Label string
}

// Edge is the export edge type.
type Edge struct {
	From, To string
	Weight   float64
}

// GraphML writes a GraphML document.
func GraphML(w io.Writer, nodes []Node, edges []Edge) error {
	type ml struct {
		XMLName xml.Name `xml:"graphml"`
		Body    string   `xml:",innerxml"`
	}
	var inner string
	inner += `<graph edgedefault="undirected">`
	for _, n := range nodes {
		inner += fmt.Sprintf(`<node id=%q><data key="label">%s</data></node>`, n.ID, xmlEscape(n.Label))
	}
	for i, e := range edges {
		inner += fmt.Sprintf(`<edge id=%q source=%q target=%q><data key="weight">%v</data></edge>`,
			fmt.Sprintf("e%d", i), e.From, e.To, e.Weight)
	}
	inner += `</graph>`
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(ml{Body: inner})
}

// Gephi writes a CSV-edge-list that gephi imports cleanly.
func Gephi(w io.Writer, edges []Edge) error {
	if _, err := fmt.Fprintln(w, "Source,Target,Weight"); err != nil {
		return err
	}
	for _, e := range edges {
		if _, err := fmt.Fprintf(w, "%s,%s,%v\n", e.From, e.To, e.Weight); err != nil {
			return err
		}
	}
	return nil
}

func xmlEscape(s string) string {
	var buf [1024]byte
	w := buf[:0]
	for _, r := range s {
		switch r {
		case '<':
			w = append(w, "&lt;"...)
		case '>':
			w = append(w, "&gt;"...)
		case '&':
			w = append(w, "&amp;"...)
		case '"':
			w = append(w, "&quot;"...)
		default:
			w = append(w, string(r)...)
		}
	}
	return string(w)
}
