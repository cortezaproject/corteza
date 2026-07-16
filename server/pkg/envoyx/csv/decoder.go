package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx/datasource"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cast"
)

type (
	decoder struct {
		ident string

		src *os.File

		reader   *csv.Reader
		skipHead bool

		delimiter           string
		multiValueDelimiter string
		multiValueBrackets  bool

		header []string
		row    datasource.RawRecord
		count  uint64
	}
)

func CanDecodeFile(f io.Reader) bool {
	m, err := mimetype.DetectReader(f)
	if err != nil {
		return false
	}

	return CanDecodeExt(m.Extension())
}

func CanDecodeMime(m string) bool {
	return m == "text/csv"
}

func CanDecodeExt(ext string) bool {
	pt := strings.Split(ext, ".")
	return strings.TrimSpace(pt[len(pt)-1]) == "csv"
}

// Decoder inits a new csv decoder from the given reader
//
// @todo hold small files in mem to avoid needles disc access
func Decoder(r io.Reader, ident string, config map[string]any) (out *decoder, err error) {
	out = &decoder{
		ident: ident,
	}

	err = out.parseConfig(config)
	if err != nil {
		return
	}

	out.src, err = ioutil.TempFile(os.TempDir(), "*.csv")
	if err != nil {
		return
	}

	r, err = out.flushTemp(r)
	defer out.Reset(nil)
	if err != nil {
		return
	}

	out.reader = csv.NewReader(r)
	out.reader.ReuseRecord = true
	if out.delimiter != "" {
		out.reader.Comma = []rune(out.delimiter)[0]
	}

	// Header
	aux, err := out.reader.Read()
	out.header = append(out.header, aux...)
	if err != nil {
		return
	}

	out.row = make(datasource.RawRecord, len(out.header))

	for {
		_, err = out.reader.Read()
		if err == io.EOF {
			return out, nil
		} else if err != nil {
			return
		}

		// Entry count
		out.count++
	}
}

func (d *decoder) SetConfigs(config map[string]any) (err error) {
	return d.parseConfig(config)
}

// Cleanup should be called before we stop using the decoder
func (d *decoder) Cleanup() error {
	return os.Remove(d.src.Name())
}

// SetIdent overwrites the system defined identifier
func (d *decoder) SetIdent(ident string) {
	d.ident = ident
}

// Ident returns the assigned identifier
func (d *decoder) Ident() string {
	return d.ident
}

// Fields returns every available field in this dataset
func (d *decoder) Fields() []string {
	return d.header
}

// Reset resets the decoder to the start
func (d *decoder) Reset(_ context.Context) error {
	_, err := d.src.Seek(0, 0)
	d.skipHead = true
	return err
}

// Next returns the field: value mapping for the next row
func (d *decoder) Next(_ context.Context, out datasource.RawRecord) (more bool, err error) {
	if d.skipHead {
		_, err = d.reader.Read()
		if err != nil {
			return
		}
		d.skipHead = false
	}

	aux, err := d.reader.Read()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		return false, err
	}
	for i, h := range d.header {
		for j, v := range d.parseComplexCSVCell(aux[i]) {
			err = out.SetValue(h, uint(j), v)
			if err != nil {
				return
			}
		}
	}

	return true, nil
}

// Count returns the total number of rows in the dataset
func (d *decoder) Count() uint64 {
	return d.count
}

func (d *decoder) flushTemp(r io.Reader) (_ io.Reader, err error) {
	_, err = io.Copy(d.src, r)
	if err != nil {
		return
	}

	d.src.Seek(0, 0)
	return d.src, nil
}

// parseComplexCSVCell helps us properly multi value and complex JSON definitions
func (d *decoder) parseComplexCSVCell(cell string) []string {
	// Cover multi value fields
	//
	// Multi value fields are "mini csvs" encoded in the csv so we need to
	// clean up the value and then CSV decode again.
	if d.multiValueBrackets {
		if strings.HasPrefix(cell, "[") && strings.HasSuffix(cell, "]") {
			cell = cell[1 : len(cell)-1]
		}
	}

	// Cover complex json definitions (geometry field)
	//
	// Multi value complex fields are delimiter-joined JSON objects; the JSON
	// itself may contain the delimiter so split on complete objects instead.
	if strings.HasPrefix(cell, "{") && strings.HasSuffix(cell, "}") {
		return splitConcatenatedJSON(cell)
	}

	if len(cell) == 0 {
		return nil
	}

	r := csv.NewReader(strings.NewReader(cell))
	if d.multiValueDelimiter != "" {
		r.Comma = []rune(d.multiValueDelimiter)[0]
	}

	records, err := r.ReadAll()
	if err != nil {
		return nil
	}

	return records[0]
}

// splitConcatenatedJSON splits a cell of delimiter-joined JSON objects into
// the individual objects by tracking brace depth and string literals, so
// delimiters and braces inside the JSON never produce a false split.
//
// A cell holding a single (possibly nested) object comes back as one value;
// a malformed cell is returned whole so the value fails validation upstream
// instead of being silently corrupted.
func splitConcatenatedJSON(cell string) (out []string) {
	var (
		depth int
		inStr bool
		esc   bool
		start = -1
	)

	for i := 0; i < len(cell); i++ {
		c := cell[i]

		if inStr {
			switch {
			case esc:
				esc = false
			case c == '\\':
				esc = true
			case c == '"':
				inStr = false
			}
			continue
		}

		switch c {
		case '"':
			inStr = true
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 && start >= 0 {
				out = append(out, cell[start:i+1])
				start = -1
			}
		}
	}

	if depth != 0 || inStr || len(out) == 0 {
		return []string{cell}
	}

	return out
}

func (d *decoder) parseConfig(cfg map[string]any) error {
	if cfg == nil {
		return nil
	}

	for k := range cfg {
		switch k {
		case "delimiter":
			d.delimiter = cast.ToString(cfg[k])
		case "multiValueDelimiter":
			switch cast.ToString(cast.ToString(cfg[k])) {
			case ",":
				d.multiValueDelimiter = ","
				d.multiValueBrackets = false
			case ";":
				d.multiValueDelimiter = ";"
				d.multiValueBrackets = false
			case "|":
				d.multiValueDelimiter = "|"
				d.multiValueBrackets = false
			case "[,]":
				d.multiValueDelimiter = ","
				d.multiValueBrackets = true
			case "[;]":
				d.multiValueDelimiter = ";"
				d.multiValueBrackets = true
			case "[|]":
				d.multiValueDelimiter = "|"
				d.multiValueBrackets = true
			}

		default:
			return fmt.Errorf("unknown parameter %s", k)
		}
	}

	return nil
}
