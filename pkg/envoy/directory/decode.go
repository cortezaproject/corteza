package directory

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	decoder interface {
		CanDecodeFile(os.FileInfo) bool
		Decode(context.Context, io.Reader, *envoy.DecoderOpts) ([]resource.Interface, error)
	}
)

// DecodeDirectory is a helper to run the decoding process over the entire directory
func Decode(ctx context.Context, p string, decoders ...decoder) ([]resource.Interface, error) {
	var (
		f *os.File

		d decoder

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

		if f, err = os.Open(p); err != nil {
			return err
		}
		defer f.Close()

		dir, fn := path.Split(p)
		do := &envoy.DecoderOpts{
			Name: fn,
			Path: dir,
		}

		if dnn, err = d.Decode(ctx, f, do); err != nil {
			return fmt.Errorf("failed to decode %s: %w", info.Name(), err)
		}

		nn = append(nn, dnn...)
		return nil
	})
}
