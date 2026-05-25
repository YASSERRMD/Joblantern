package nominatim

import "errors"

// ErrRateLimited signals that the upstream Nominatim instance returned
// HTTP 429. Callers should fall back to cache or surface a
// RATE_LIMITED error to the MCP consumer.
var ErrRateLimited = errors.New("nominatim: rate limited")
