package decoder

import (
	"io"

	"github.com/cortezaproject/corteza-server/pkg/count"
)

type (
	multiple uint

	FlatReader interface {
		Read() ([]string, error)
	}

	StructuredDecoder interface {
		Decode(interface{}) error
		More() bool
	}

	flatReader struct {
		f      io.ReadSeeker
		r      FlatReader
		header []string
		more   bool
	}

	structuredDecoder struct {
		f      io.ReadSeeker
		header []string
		d      StructuredDecoder
		buf    []map[string]interface{}
	}

	// callbacks
	sdCallback func(map[string]interface{}) error
	fdCallback func([]string) error
)

// flat reader
func NewFlatReader(r FlatReader, f io.ReadSeeker) *flatReader {
	return &flatReader{
		f:    f,
		r:    r,
		more: true,
	}
}

func (dec *flatReader) EntryCount() (uint64, error) {
	defer dec.f.Seek(0, 0)

	c, err := count.Lines(dec.f)
	if err != nil {
		return 0, err
	}
	if c <= 0 {
		return 0, nil
	}
	return c - 1, nil
}

func (dec *flatReader) get(fnc fdCallback) error {
	v, err := dec.r.Read()
	if err == io.EOF {
		dec.more = false
		return nil
	} else if err != nil {
		return err
	}

	return fnc(v)
}

func (dec *flatReader) walk(fnc fdCallback) error {
	for dec.more {
		if err := dec.get(fnc); err != nil {
			return err
		}
	}
	return nil
}

func (dec *flatReader) Header() []string {
	if len(dec.header) > 0 {
		return dec.header
	}

	dec.get(func(rtr []string) error {
		dec.header = rtr
		return nil
	})

	return dec.header
}

// structured decoder
func NewStructuredDecoder(d StructuredDecoder, f io.ReadSeeker) *structuredDecoder {
	return &structuredDecoder{
		f: f,
		d: d,
	}
}

func (dec *structuredDecoder) EntryCount() (uint64, error) {
	defer dec.f.Seek(0, 0)
	return count.Lines(dec.f)
}

func (dec *structuredDecoder) get(fnc sdCallback) error {
	if !dec.d.More() {
		return nil
	}

	var tmp map[string]interface{}
	err := dec.d.Decode(&tmp)
	if err != nil {
		return err
	}

	return fnc(tmp)
}

func (dec *structuredDecoder) exhaustBuffer(fnc sdCallback) error {
	if dec.buf != nil {
		for _, b := range dec.buf {
			fnc(b)
		}
		dec.buf = nil
	}
	return nil
}

func (dec *structuredDecoder) walk(fnc sdCallback) error {
	if err := dec.exhaustBuffer(fnc); err != nil {
		return err
	}

	for dec.d.More() {
		if err := dec.get(fnc); err != nil {
			return err
		}
	}

	return nil
}

func (dec *structuredDecoder) Header() []string {
	if len(dec.header) > 0 {
		return dec.header
	}

	var tmp []string
	dec.get(func(rtr map[string]interface{}) error {
		// buffer first row or else it will be lost
		dec.buf = append(dec.buf, rtr)

		tmp = make([]string, len(rtr))
		i := 0
		for k := range rtr {
			tmp[i] = k
			i++
		}

		return nil
	})

	dec.header = tmp
	return tmp
}
