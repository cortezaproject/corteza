package csv

import (
	"bytes"
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
		c      *csv.Reader
		header []string
		count  uint64
	}
)

// Decoder initializes and returns a fresh CSV decoder
func Decoder() *decoder {
	return &decoder{}
}

// CanDecodeFile determines if the file can be determined by this decoder
func (y *decoder) CanDecodeFile(f io.Reader) bool {
	_, ext, err := mimetype.DetectReader(f)
	if err != nil {
		return false
	}

	return y.CanDecodeExt(ext)
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
	cr := &reader{}
	var buff bytes.Buffer

	// So we can reset to the start of the reader
	tr := io.TeeReader(r, &buff)
	err := cr.prepare(tr)
	if err != nil {
		return nil, err
	}

	cr.c = csv.NewReader(&buff)

	// The first one is a header, so let's just get rid of it
	_, err = cr.c.Read()
	if err != nil {
		return nil, err
	}

	return []resource.Interface{resource.NewResourceDataset(do.Name, cr)}, nil
}

func (cr *reader) prepare(r io.Reader) (err error) {
	cReader := csv.NewReader(r)

	// Header
	cr.header, err = cReader.Read()
	if err != nil {
		return err
	}

	// Entry count
	for {
		_, err := cReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		cr.count++
	}

	return nil
}

// Fields returns every available field in this dataset
func (cr *reader) Fields() []string {
	return cr.header
}

// Next returns the field: value mapping for the next row
func (cr *reader) Next() (map[string]string, error) {
	mr := make(map[string]string)
	rr, err := cr.c.Read()
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	for i, h := range cr.header {
		mr[h] = rr[i]
	}

	return mr, nil
}

func (cr *reader) Count() uint64 {
	return cr.count
}
