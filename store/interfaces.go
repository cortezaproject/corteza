package store

import "io"

type (
	store struct {
		namespace string

		originalFn func(id uint64, ext string) string
		previewFn  func(id uint64, ext string) string
	}

	Store interface {
		Namespace() string

		Original(id uint64, ext string) string
		Preview(id uint64, ext string) string

		Save(filename string, contents io.Reader) error
		Remove(filename string) error
		Open(filename string) (io.Reader, error)
	}
)
