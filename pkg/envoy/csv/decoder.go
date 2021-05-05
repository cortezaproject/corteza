package csv

import (
	"context"
	"encoding/csv"
	"io"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/gabriel-vasile/mimetype"
)

type (
	// wrapper struct for csv related methods
	decoder struct{}

	// csv decoder wrapper for additional bits
	reader struct {
		header []string
		count  uint64

		cache      []map[string]string
		cacheIndex int
	}
)

// Decoder initializes and returns a fresh CSV decoder
func Decoder() *decoder {
	return &decoder{}
}

// CanDecodeFile determines if the file can be determined by this decoder
func (y *decoder) CanDecodeFile(f io.Reader) bool {
	m, err := mimetype.DetectReader(f)
	if err != nil {
		return false
	}

	return y.CanDecodeExt(m.Extension())
}

func (y *decoder) CanDecodeMime(m string) bool {
	return m == "text/csv"
}

func (y *decoder) CanDecodeExt(ext string) bool {
	pt := strings.Split(ext, ".")
	return strings.TrimSpace(pt[len(pt)-1]) == "csv"
}

// Decode decodes the given io.Reader into a generic resource dataset
func (c *decoder) Decode(ctx context.Context, r io.Reader, do *envoy.DecoderOpts) ([]resource.Interface, error) {
	cr := &reader{
		cacheIndex: 0,
		cache:      make([]map[string]string, 0, 1000),
	}

	err := cr.prepare(r)
	if err != nil {
		return nil, err
	}

	return []resource.Interface{resource.NewResourceDataset(do.Name, cr)}, nil
}

// The prepare step caches all of the rows into a cache array that will be used
// to access already visited CSV rows.
//
// @todo implement FS cache file for larger files
func (cr *reader) prepare(r io.Reader) (err error) {
	cReader := csv.NewReader(r)

	// Header
	cr.header, err = cReader.Read()
	if err != nil {
		return err
	}

	for {
		aux := make(map[string]string)
		rr, err := cReader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		// Entry count
		cr.count++

		// Cache entry
		for i, h := range cr.header {
			aux[h] = rr[i]
		}
		cr.cache = append(cr.cache, aux)
	}
}

// Fields returns every available field in this dataset
func (cr *reader) Fields() []string {
	return cr.header
}

func (cr *reader) Reset() error {
	cr.cacheIndex = 0
	return nil
}

// Next returns the field: value mapping for the next row
func (cr *reader) Next() (map[string]string, error) {
	if cr.cacheIndex >= len(cr.cache) {
		return nil, nil
	}

	mr := cr.cache[cr.cacheIndex]
	cr.cacheIndex++
	return mr, nil
}

func (cr *reader) Count() uint64 {
	return cr.count
}
