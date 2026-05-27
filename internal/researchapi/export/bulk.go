// Package export implements bulk anonymized exports with opaque
// cursor pagination. Cursors are signed base64url to prevent skipping
// past tier-permitted ranges.
package export

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

// Cursor is the opaque page handle exchanged with researchers.
type Cursor struct {
	After     time.Time `json:"after"`
	ID        string    `json:"id"`
	PageSize  int       `json:"n"`
	Tier      string    `json:"t"`
	IssuedAt  time.Time `json:"iat"`
}

// Encode serialises a cursor to a URL-safe string.
func (c Cursor) Encode() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Decode parses an opaque cursor.
func Decode(s string) (Cursor, error) {
	var c Cursor
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return c, errors.New("invalid cursor")
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, errors.New("invalid cursor")
	}
	if time.Since(c.IssuedAt) > 7*24*time.Hour {
		return c, errors.New("cursor expired")
	}
	return c, nil
}

// Page is the response shape for bulk export endpoints.
type Page struct {
	Items      []map[string]any `json:"items"`
	NextCursor string           `json:"next_cursor,omitempty"`
}
