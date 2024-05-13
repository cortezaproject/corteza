package plain

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/spf13/afero"
)

type (
	store struct {
		fs afero.Fs

		namespace string

		originalFn func(id uint64, ext string) string
		previewFn  func(id uint64, ext string) string
	}
)

var (
	defPreviewFn = func(id uint64, ext string) string {
		return fmt.Sprintf("%d_preview.%s", id, ext)
	}

	defOriginalFn = func(id uint64, ext string) string {
		return fmt.Sprintf("%d.%s", id, ext)
	}
)

func New(namespace string) (*store, error) {
	return NewWithAfero(afero.NewOsFs(), namespace)
}

func NewWithAfero(fs afero.Fs, namespace string) (*store, error) {
	return &store{
		fs:        fs,
		namespace: namespace,

		originalFn: defOriginalFn,
		previewFn:  defPreviewFn,
	}, nil
}

func (s *store) check(filename string) error {
	if len(filename) == 0 {
		return fmt.Errorf("Invalid filename when trying to store file: '%s' (for %s)", filename, s.namespace)
	}

	if filename[:len(s.namespace)+1] != s.namespace+"/" {
		return fmt.Errorf("Invalid namespace when trying to store file: '%s' (for %s)", filename, s.namespace)
	}

	return nil
}

func (s *store) Original(id uint64, ext string) string {
	return path.Join(s.namespace, s.originalFn(id, ext))
}

func (s *store) Preview(id uint64, ext string) string {
	return path.Join(s.namespace, s.previewFn(id, ext))
}

func (s *store) Save(filename string, contents io.Reader) (err error) {
	// check filename for validity
	if err = s.check(filename); err != nil {
		return
	}

	folder := path.Dir(filename)

	if err = s.fs.MkdirAll(folder, 0755); err != nil {
		return
	}

	return afero.WriteReader(s.fs, filename, contents)
}

func (s *store) Remove(filename string) error {
	// check filename for validity
	if err := s.check(filename); err != nil {
		return err
	}

	return s.fs.Remove(filename)
}

func (s *store) Open(filename string) (out io.ReadSeekCloser, err error) {
	// check filename for validity
	if err := s.check(filename); err != nil {
		return nil, err
	}

	for i := 1; i <= 3; i++ {
		out, err = s.fs.Open(filename)
		if err != nil {
			time.Sleep(time.Millisecond * 500 * time.Duration(i))

			continue
		}

		return
	}

	return
}

func (s *store) Healthcheck(ctx context.Context) error {
	var (
		fname = fmt.Sprintf("%s/.healthcheck_%d", s.namespace, time.Now().UnixNano())
		buf   = &bytes.Buffer{}
	)

	if s == nil {
		return fmt.Errorf("uninitialized")
	}

	buf.Write([]byte("healthcheck"))

	if err := s.Save(fname, buf); err != nil {
		return err
	}

	if _, err := s.Open(fname); err != nil {
		return err
	}

	if err := s.Remove(fname); err != nil {
		return err
	}

	return nil
}
