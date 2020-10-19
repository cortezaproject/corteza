package decoder

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/types"
)

// DecodeDirectory is a helper to run the decoding process over the entire directory
func DecodeDirectory(ctx context.Context, path string) (types.NodeSet, error) {
	nn := make(types.NodeSet, 0, 100)

	// Decoders
	dy := NewYamlDecoder()
	dc := NewCsvDecoder()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// Decode sources
		fn := info.Name()
		if strings.HasSuffix(fn, ".yaml") {
			ynn, err := dy.Decode(ctx, f)
			if err != nil {
				return nil
			}

			nn = append(nn, ynn...)
		} else if strings.HasSuffix(fn, ".csv") {
			cnn, err := dc.Decode(ctx, f, fn)
			if err != nil {
				return nil
			}

			nn = append(nn, cnn...)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nn, nil
}
