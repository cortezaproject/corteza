package csv

import (
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type (
	decoder struct {
		ident string

		src *os.File

		reader   *csv.Reader
		skipHead bool

		header []string
		row    map[string]string
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
func Decoder(r io.Reader, ident string) (out *decoder, err error) {
	out = &decoder{
		ident: ident,
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

	// Header
	aux, err := out.reader.Read()
	out.header = append(out.header, aux...)
	if err != nil {
		return
	}

	out.row = make(map[string]string, len(out.header))

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
func (d *decoder) Next(_ context.Context, out map[string]string) (more bool, err error) {
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
		out[h] = aux[i]
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
