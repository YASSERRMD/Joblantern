// Package federation lets the Joblantern KG federate queries with
// other open KGs (Wikidata, DBpedia, peer Joblantern instances).
package federation

// Peer is a federated SPARQL endpoint.
type Peer struct {
	ID    string
	Name  string
	URL   string
	Vocab string // base vocabulary IRI
}

// Defaults returns the canonical peer list.
func Defaults() []Peer {
	return []Peer{
		{ID: "wikidata", Name: "Wikidata", URL: "https://query.wikidata.org/sparql", Vocab: "http://www.wikidata.org/entity/"},
		{ID: "dbpedia", Name: "DBpedia", URL: "https://dbpedia.org/sparql", Vocab: "http://dbpedia.org/resource/"},
	}
}
