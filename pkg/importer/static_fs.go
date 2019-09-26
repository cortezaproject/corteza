package importer

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/goware/statik/fs"
	"github.com/pkg/errors"
)

// ReadStatic reads static FS and returns slice of io.Readers
func ReadStatic(data string) ([]io.Reader, error) {
	var (
		files = make([]string, 0)

		yamlFilter = func(filename string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				if matched, err := filepath.Match("/*.yaml", filename); matched && err == nil {
					files = append(files, filename)
				}
			}

			return err
		}

		sfs, err = fs.New(data)

		readers []io.Reader
	)

	if err != nil {
		return nil, errors.Wrap(err, "could not read static filesystem")
	}

	if err = fs.Walk(sfs, "/", yamlFilter); err != nil {
		return nil, errors.Wrap(err, "could not filter files")
	}

	if len(files) == 0 {
		return nil, nil
	}

	sort.Strings(files)

	for _, file := range files {
		if bb, err := fs.ReadFile(sfs, file); err != nil {
			return nil, errors.Wrapf(err, "could not read %s", file)
		} else {
			readers = append(readers, bytes.NewBuffer(bb))
		}
	}

	return readers, nil
}
