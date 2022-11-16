package directory

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

type (
	Decoder interface {
		CanDecodeExt(string) bool
		CanDecodeFile(io.Reader) bool
		Decode(context.Context, io.Reader, *envoy.DecoderOpts) ([]resource.Interface, error)
	}
)

// DecodeDirectory is a helper to run the decoding process over the entire directory
func Decode(ctx context.Context, p string, decoders ...Decoder) ([]resource.Interface, error) {
	var (
		f *os.File

		// decoded nodes
		dnn []resource.Interface

		// agregated resources
		nn = make([]resource.Interface, 0, 100)
	)

	if len(decoders) == 0 {
		return nil, fmt.Errorf("no decoders provided")
	}

	return nn, filepath.Walk(p, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if f, err = os.Open(p); err != nil {
			return err
		}
		defer f.Close()

		for _, d := range decoders {
			if !d.CanDecodeFile(f) {
				// decoder cannot handle this file
				// Make sure to reset it, as the above check consumes the reader
				if _, err = f.Seek(0, 0); err != nil {
					return err
				}

				continue
			}

			if !d.CanDecodeExt(info.Name()) {
				// this decoder cannot handle this extension
				continue
			}

			if _, err = f.Seek(0, 0); err != nil {
				return err
			}

			dir, fn := path.Split(p)
			do := &envoy.DecoderOpts{
				Name: fn,
				Path: dir,
			}

			if dnn, err = d.Decode(ctx, f, do); err != nil {
				return fmt.Errorf("failed to decode %s: %w", info.Name(), err)
			}

			nn = append(nn, dnn...)

			// found compatible decoder
			return nil
		}

		return nil
	})
}
