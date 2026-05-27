// Package integration wires marketplace verdicts into Phase 42
// (caller-ID) and Phase 43 (graph analysis) so a scam network across
// jobs, rentals, and marketplaces collapses into one component.
package integration

// Bridge wraps the small set of hooks the marketplace module needs.
type Bridge interface {
	PushPhoneHash(phoneHash, band string) error
	PushEntity(kind, id string) error
	PushEdge(a, b string) error
}

// Wire performs the per-verdict propagation using the bridge.
func Wire(b Bridge, phones []string, sellerEntity string, marketEntity string) error {
	for _, p := range phones {
		if err := b.PushPhoneHash(p, "red"); err != nil {
			return err
		}
	}
	if err := b.PushEntity("seller", sellerEntity); err != nil {
		return err
	}
	if marketEntity != "" {
		if err := b.PushEdge(sellerEntity, marketEntity); err != nil {
			return err
		}
	}
	return nil
}
