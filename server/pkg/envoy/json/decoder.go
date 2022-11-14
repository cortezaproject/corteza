package json

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/gabriel-vasile/mimetype"
)

type (
	// wrapper struct for json related methods
	decoder struct{}

	// json decoder wrapper for additional bits
	reader struct {
		header []string
		count  uint64

		cache      []map[string]string
		cacheIndex int
	}
)

// Decoder initializes and returns a fresh JSON decoder
//
// We'll only do jsonl for now -- for record importing.
// We'll expand this to work with other resources later on.
func Decoder() *decoder {
	return &decoder{}
}

// CanDecodeFile determines if the file can be decoded by this decoder
func (d *decoder) CanDecodeFile(f io.Reader) bool {
	var buff bytes.Buffer
	tr := io.TeeReader(f, &buff)

	m, err := mimetype.DetectReader(tr)
	if err != nil {
		return false
	}
	ext := m.Extension()

	return d.CanDecodeExt(ext)
}

func (d *decoder) CanDecodeMime(m string) bool {
	return m == "application/json" || m == "application/jsonlines"
}

func (d *decoder) CanDecodeExt(ext string) bool {
	pt := strings.Split(ext, ".")
	ext = strings.TrimSpace(pt[len(pt)-1])
	return ext == "jsonl" || ext == "json" || ext == "ndjson"
}

// Decode decodes the given io.Reader into a generic resource dataset
func (d *decoder) Decode(ctx context.Context, r io.Reader, do *envoy.DecoderOpts) ([]resource.Interface, error) {
	jr := &reader{
		cacheIndex: 0,
		cache:      make([]map[string]string, 0, 1000),
	}

	err := jr.prepare(r)
	if err != nil {
		return nil, err
	}

	return []resource.Interface{resource.NewResourceDataset(do.Name, jr)}, nil
}

// The prepare step caches all of the rows into a cache array that will be used
// to access already visited json rows.
//
// @todo implement FS cache file for larger files
func (jr *reader) prepare(r io.Reader) (err error) {
	jReader := json.NewDecoder(r)

	// JSON can omit empty values, so we can't be 100% that all of the headers
	// were read in the first go.
	// We'll determine all of the headers below where we count the entries
	hx := make(map[string]bool)
	jr.header = make([]string, 0, 100)

	for jReader.More() {
		aux := make(map[string]string)
		err = jReader.Decode(&aux)
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		// Get all the header fields
		for h := range aux {
			if !hx[h] {
				jr.header = append(jr.header, h)
				hx[h] = true
			}
		}

		// Entry count
		jr.count++

		// Cache entry
		jr.cache = append(jr.cache, aux)
	}

	return nil
}

func (jr *reader) Reset() error {
	jr.cacheIndex = 0
	return nil
}

// Fields returns every available field in this dataset
func (jr *reader) Fields() []string {
	return jr.header
}

// Next returns the field: value mapping for the next row
func (jr *reader) Next() (map[string]string, error) {
	if jr.cacheIndex >= len(jr.cache) {
		return nil, nil
	}

	mr := jr.cache[jr.cacheIndex]
	jr.cacheIndex++
	return mr, nil
}

func (cr *reader) Count() uint64 {
	return cr.count
}
