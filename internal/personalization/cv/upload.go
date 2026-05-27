// Package cv handles the optional CV upload at submission time. The
// raw bytes never leave RAM unless the user explicitly opts in.
package cv

import (
	"errors"
	"io"
	"mime/multipart"
)

// MaxBytes caps in-memory CV size.
const MaxBytes = 5 * 1024 * 1024 // 5 MiB

// Upload is the in-memory CV payload.
type Upload struct {
	Filename    string
	ContentType string
	Bytes       []byte
}

// ReadOptional accepts an optional multipart form file. If the field
// is absent the caller proceeds without personalization.
func ReadOptional(form *multipart.Form, field string) (*Upload, error) {
	if form == nil || form.File == nil {
		return nil, nil
	}
	files := form.File[field]
	if len(files) == 0 {
		return nil, nil
	}
	fh := files[0]
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lim := io.LimitReader(f, MaxBytes+1)
	body, err := io.ReadAll(lim)
	if err != nil {
		return nil, err
	}
	if len(body) > MaxBytes {
		return nil, errors.New("cv exceeds 5 MiB cap")
	}
	return &Upload{Filename: fh.Filename, ContentType: fh.Header.Get("Content-Type"), Bytes: body}, nil
}
