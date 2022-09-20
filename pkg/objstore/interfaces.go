package objstore

import (
	"context"
	"io"
)

type Store interface {
	// Original returns URL to the original file
	Original(id uint64, ext string) string

	// Preview returns URL to the preview (of the original) file
	Preview(id uint64, ext string) string

	// Save stores the file
	Save(filename string, f io.Reader) error

	// Remove deletes the file
	Remove(filename string) error

	// Open returns file handle
	Open(filename string) (io.ReadSeekCloser, error)

	// Healthcheck checks health status of the store
	Healthcheck(ctx context.Context) error
}
