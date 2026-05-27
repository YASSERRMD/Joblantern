// Package objstore wraps S3-compatible storage for evidence
// artifacts. Self-host deployments point this at MinIO; cloud
// deployments at AWS S3, Cloudflare R2, or Backblaze B2.
package objstore

import (
	"context"
	"errors"
	"io"
)

// Object is a stored artifact.
type Object struct {
	Bucket  string
	Key     string
	Size    int64
	SHA256  string
	Reader  io.ReadCloser
}

// Store is the abstraction.
type Store interface {
	Put(ctx context.Context, o Object) error
	Get(ctx context.Context, bucket, key string) (Object, error)
	Delete(ctx context.Context, bucket, key string) error
}

// ErrNotFound is returned for missing keys.
var ErrNotFound = errors.New("object not found")
