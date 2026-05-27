package auth

import "context"

func withToken(ctx context.Context, t *Token) context.Context {
	return context.WithValue(ctx, ctxKey{}, t)
}

// FromContext returns the verified token attached to the request
// context. It always returns a non-nil token; unauthenticated callers
// get a TierPublic token.
func FromContext(ctx context.Context) *Token {
	if t, ok := ctx.Value(ctxKey{}).(*Token); ok && t != nil {
		return t
	}
	return &Token{Tier: TierPublic}
}
