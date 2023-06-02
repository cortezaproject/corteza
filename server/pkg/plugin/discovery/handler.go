package discovery

import (
	"archive/zip"
	"errors"
	"io/fs"
	"path/filepath"
)

type (
	PackageHandler interface {
		Match(fs.DirEntry) bool
		Validate(fs.DirEntry, string) error
	}

	Unboxer interface {
		Unbox() error
	}

	packageHandlerFs  struct{}
	packageHandlerZip struct{}
)

func (phfs packageHandlerFs) Match(d fs.DirEntry) bool {
	return d.IsDir()
}

func (phfs packageHandlerFs) Validate(d fs.DirEntry, path string) (err error) {
	return
}

func (phfs packageHandlerZip) Match(d fs.DirEntry) bool {
	e := filepath.Ext(d.Name())
	return e == ".zip" || e == ".crust"
}

func (phfs packageHandlerZip) Validate(d fs.DirEntry, path string) (err error) {
	fullPath := filepath.Join(path, d.Name())

	zipListing, err := zip.OpenReader(fullPath)

	if err != nil {
		return
	}

	defer zipListing.Close()

	err = func() (err error) {
		for _, file := range zipListing.File {
			if file.Name == "meta.yaml" {
				return
			}
		}

		return errors.New("meta file not found in this bundle")
	}()

	return
}

func (phfs packageHandlerZip) Unbox() error {
	return nil
}
