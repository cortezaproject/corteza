package directory

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"io"
	"os"
	"path/filepath"
)

type (
	decoder interface {
		CanDecodeFile(os.FileInfo) bool
		Decode(context.Context, io.Reader, os.FileInfo) ([]envoy.Node, error)
	}
)

// DecodeDirectory is a helper to run the decoding process over the entire directory
func Decode(ctx context.Context, path string, decoders ...decoder) ([]envoy.Node, error) {
	var (
		f *os.File

		d decoder

		// decoded nodes
		dnn envoy.NodeSet

		// agregated nodes
		nn = make(envoy.NodeSet, 0, 100)
	)

	if len(decoders) == 0 {
		return nil, fmt.Errorf("no decoders provided")
	}

	return nn, filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		for _, d = range decoders {
			// find compatible decoder
			if d.CanDecodeFile(info) {
				break
			}
		}

		if d == nil {
			// no decoder found
			return nil
		}

		if f, err = os.Open(path); err != nil {
			return err
		}
		defer f.Close()

		if dnn, err = d.Decode(ctx, f, info); err != nil {
			return err
		}

		nn = append(nn, dnn...)
		return nil
	})
}
