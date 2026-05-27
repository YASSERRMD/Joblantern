package graphql

import (
	"errors"
	"strings"
)

// CostLimits defines per-tier complexity ceilings.
type CostLimits struct {
	MaxDepth      int
	MaxComplexity int
	MaxAliases    int
}

// Default returns the limits applied per tier.
func Default(tier string) CostLimits {
	switch tier {
	case "academic":
		return CostLimits{MaxDepth: 10, MaxComplexity: 5000, MaxAliases: 25}
	case "journalist":
		return CostLimits{MaxDepth: 8, MaxComplexity: 3000, MaxAliases: 15}
	case "regulator":
		return CostLimits{MaxDepth: 12, MaxComplexity: 10000, MaxAliases: 50}
	default:
		return CostLimits{MaxDepth: 5, MaxComplexity: 500, MaxAliases: 5}
	}
}

// Analyse estimates the cost of a query using simple structural
// heuristics. It is deliberately conservative — production runs the
// gqlparser AST through complexity rules from gqlgen.
func Analyse(query string, l CostLimits) error {
	depth := strings.Count(query, "{")
	if depth > l.MaxDepth {
		return errors.New("query depth exceeds tier limit")
	}
	aliases := strings.Count(query, ":")
	if aliases > l.MaxAliases {
		return errors.New("alias count exceeds tier limit")
	}
	complexity := depth*100 + aliases*10
	if complexity > l.MaxComplexity {
		return errors.New("query complexity exceeds tier limit")
	}
	return nil
}
