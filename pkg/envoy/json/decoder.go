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
		d      *json.Decoder
		header []string
		count  uint64

		buff bytes.Buffer
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

func (d *decoder) CanDecodeExt(ext string) bool {
	pt := strings.Split(ext, ".")
	ext = strings.TrimSpace(pt[len(pt)-1])
	return ext == "jsonl" || ext == "json" || ext == "ndjson"
}

// Decode decodes the given io.Reader into a generic resource dataset
func (d *decoder) Decode(ctx context.Context, r io.Reader, do *envoy.DecoderOpts) ([]resource.Interface, error) {
	jr := &reader{}
	var buff bytes.Buffer

	// So we can reset to the start of the reader
	err := jr.prepare(io.TeeReader(r, &buff))
	if err != nil {
		return nil, err
	}

	jr.d = json.NewDecoder(io.TeeReader(&buff, &jr.buff))

	return []resource.Interface{resource.NewResourceDataset(do.Name, jr)}, nil
}

func (jr *reader) prepare(r io.Reader) (err error) {
	jReader := json.NewDecoder(r)

	// JSON can omit empty values, so we can't be 100% that all of the headers
	// were read in the first go.
	// We'll determine all of the headers below where we count the entries
	hx := make(map[string]bool)
	jr.header = make([]string, 0, 100)

	aux := make(map[string]interface{})
	for jReader.More() {
		err = jReader.Decode(&aux)
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		// Get all the header fields
		// @todo do we want to preserve order?
		for h := range aux {
			if !hx[h] {
				jr.header = append(jr.header, h)
				hx[h] = true
			}
		}

		jr.count++
	}

	return nil
}

func (jr *reader) Reset() error {
	if len(jr.buff.Bytes()) == 0 {
		return nil
	}

	buff := jr.buff
	jr.buff = bytes.Buffer{}

	jr.d = json.NewDecoder(io.TeeReader(&buff, &jr.buff))

	return nil
}

// Fields returns every available field in this dataset
func (jr *reader) Fields() []string {
	return jr.header
}

// Next returns the field: value mapping for the next row
func (jr *reader) Next() (map[string]string, error) {
	// It's over
	if !jr.d.More() {
		return nil, nil
	}

	mr := make(map[string]string)
	err := jr.d.Decode(&mr)
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return mr, nil
}

func (cr *reader) Count() uint64 {
	return cr.count
}
