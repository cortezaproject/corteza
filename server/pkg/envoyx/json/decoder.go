package json

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type (
	decoder struct {
		ident string

		src    *os.File
		reader *json.Decoder

		header []string
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
	return m == "application/json" || m == "application/jsonlines"
}

func CanDecodeExt(ext string) bool {
	pt := strings.Split(ext, ".")
	ext = strings.TrimSpace(pt[len(pt)-1])
	return ext == "jsonl" || ext == "json" || ext == "ndjson"
}

// Decoder inits a new csv decoder from the given reader
//
// @todo hold small files in mem to avoid needles disc access
func Decoder(r io.Reader, ident string) (out *decoder, err error) {
	out = &decoder{
		ident: ident,
	}

	out.src, err = ioutil.TempFile(os.TempDir(), "*.ndjson")
	if err != nil {
		return
	}

	r, err = out.flushTemp(r)
	defer out.src.Seek(0, 0)
	if err != nil {
		return
	}

	out.reader = json.NewDecoder(r)

	seenHeader := make(map[string]bool)

	var aux map[string]string
	for out.reader.More() {
		err = out.reader.Decode(&aux)
		if err == io.EOF {
			return out, nil
		} else if err != nil {
			return
		}

		for f := range aux {
			if seenHeader[f] {
				continue
			}
			seenHeader[f] = true

			out.header = append(out.header, f)
		}

		// Entry count
		out.count++
	}

	return
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
	return err
}

// Next returns the field: value mapping for the next row
func (d *decoder) Next(_ context.Context, out map[string]string) (more bool, err error) {
	err = d.reader.Decode(&out)
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Empty out missing fields to keep consistent with CSV
	for _, h := range d.header {
		if _, ok := out[h]; !ok {
			out[h] = ""
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
