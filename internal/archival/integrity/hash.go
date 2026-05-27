// Package integrity computes per-file and per-archive integrity
// hashes.
package integrity

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// Manifest pairs file name with its hash.
type Manifest struct {
	Files []FileEntry
}

// FileEntry is one entry in a manifest.
type FileEntry struct {
	Name   string
	Bytes  int64
	SHA256 string
}

// HashReader streams through r and returns the hex digest plus byte count.
func HashReader(r io.Reader) (string, int64, error) {
	h := sha256.New()
	n, err := io.Copy(h, r)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}
