package decoder

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"regexp"

	"github.com/cortezaproject/corteza-server/pkg/envoy/types"
)

type (
	CsvDecoder struct{}
)

var (
	ErrorNoCsvHeader        = errors.New("csv decoder: no header")
	ErrorCsvHeaderMalformed = errors.New("csv decoder: header malformed")

	// This strict regexp for field names will do for now.
	// Later we can add support for matching over field labels as well.
	headerRegexp, _ = regexp.Compile("^[A-Za-z][0-9A-Za-z_]*[A-Za-z0-9]$")
)

func NewCsvDecoder() *CsvDecoder {
	return &CsvDecoder{}
}

// A quick header field validator
//
// @note should we complicate it any further?
func (c *CsvDecoder) validateHeader(header []string) error {
	for _, h := range header {
		if !headerRegexp.MatchString(h) {
			return ErrorCsvHeaderMalformed
		}
	}

	return nil
}

func (c *CsvDecoder) Decode(ctx context.Context, r io.Reader, filename string) ([]types.Node, error) {
	n := &types.ComposeRecordNode{}

	// Determine base module for dependency resolution
	// -4 is to remove .csv ext
	//
	// @todo tweak this a bit
	modRes := filename[0 : len(filename)-4]
	mod := &types.ComposeModule{}
	mod.Handle = modRes
	mod.Name = modRes
	n.Mod = mod

	// Prepare reader
	//
	// For optimization we reuse allocated memory; keep this in mind!
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	// Get header
	hh, err := cr.Read()
	if err == io.EOF {
		return nil, ErrorNoCsvHeader
	} else if err != nil {
		return nil, err
	}

	header := make([]string, 0, len(hh))
	for _, h := range hh {
		header = append(header, h)
	}

	err = c.validateHeader(header)
	if err != nil {
		return nil, err
	}

	// Iterator function for providing records to be imported.
	// This doesn't do any validation; that should be handled by other layers.
	n.Walk = func(f func(*types.ComposeRecord) error) error {
		for {
			record, err := cr.Read()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			rvs := make(types.ComposeRecordValueSet, 0)
			for i, h := range header {
				v := &types.ComposeRecordValue{}
				v.Name = h
				v.Value = record[i]

				rvs = append(rvs, v)
			}

			rec := &types.ComposeRecord{}
			rec.Values = rvs

			err = f(rec)
			if err != nil {
				return err
			}
		}
	}

	return []types.Node{n}, nil
}
