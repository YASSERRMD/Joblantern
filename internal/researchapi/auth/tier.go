// Package auth gates the research API behind a tiered token model.
// Researchers must hold a signed data-use agreement (DUA) and a valid
// API token whose claims include their assigned tier.
package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Tier identifies the access level granted to a token holder.
type Tier string

const (
	TierPublic     Tier = "public"     // anonymous, hard-capped
	TierAcademic   Tier = "academic"   // requires DUA
	TierJournalist Tier = "journalist" // requires DUA + editor sign-off
	TierRegulator  Tier = "regulator"  // requires Phase 37 verification
)

// Token is the verified representation of a researcher credential.
type Token struct {
	Subject      string
	Tier         Tier
	AgreementID  string
	Organization string
}

// Verifier resolves a bearer token into a Token. Implementations are
// typically backed by ed25519 signature verification.
type Verifier interface {
	Verify(bearer string) (*Token, error)
}

// ErrUnauthorized is returned when a token is missing or invalid.
var ErrUnauthorized = errors.New("researcher token invalid")

type ctxKey struct{}

// Middleware extracts and verifies the bearer token. Missing or
// invalid tokens fall back to TierPublic so unauthenticated callers
// still get the rate-limited public surface.
func Middleware(v Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok := &Token{Tier: TierPublic}
			authz := r.Header.Get("Authorization")
			if strings.HasPrefix(authz, "Bearer ") {
				if got, err := v.Verify(strings.TrimPrefix(authz, "Bearer ")); err == nil {
					tok = got
				}
			}
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(withToken(ctx, tok)))
		})
	}
}
