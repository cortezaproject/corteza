package gig

import (
	"io"
	"io/ioutil"
	"os"
)

func createTempFile(src io.Reader) (*os.File, error) {
	tmpF, err := ioutil.TempFile(os.TempDir(), "gig")
	if err != nil {
		return nil, err
	}

	if src != nil {
		_, err = io.Copy(tmpF, src)
		if err != nil {
			return nil, err
		}
	}

	return tmpF, nil
}

func getTempFile(path string) (*os.File, error) {
	return os.Open(path)
}

func deleteTempFile(path string) error {
	return os.Remove(path)
}
