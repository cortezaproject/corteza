package filter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type (
	// Paging is a helper struct that should be embedded in filter types
	// to help with the paging
	Paging struct {
		// How many items per fetch do we want
		Limit uint `json:"limit,omitempty"`

		// PeakNextPage bool `json:"peakNextPage,omitempty"`

		// Opaque pointer to 1st item on page
		// read-only
		PageCursor *PagingCursor `json:"cursor,omitempty"`

		// Opaque cursor that points to the first item on the next page
		// read-only
		NextPage *PagingCursor `json:"nextPage,omitempty"`

		// Opaque cursor that points to the first item on the previous page
		// value of {cursor} will be copied here
		// read-only
		PrevPage *PagingCursor `json:"prevPage,omitempty"`
	}

	PagingCursor struct {
		keys    []string
		values  []interface{}
		desc    []bool
		Reverse bool
	}
)

func NewPaging(limit uint, cursor string) (p Paging, err error) {
	p = Paging{Limit: limit}

	if p.PageCursor, err = parseCursor(cursor); err != nil {
		return
	}

	return
}

func (p *PagingCursor) Set(k string, v interface{}, d bool) {
	for i, key := range p.keys {
		if key == k {
			p.values[i] = v
			return
		}
	}

	p.keys = append(p.keys, k)
	p.values = append(p.values, v)
	p.desc = append(p.desc, d)
}

func (p *PagingCursor) Keys() []string {
	return p.keys
}

func (p *PagingCursor) Values() []interface{} {
	return p.values
}

func (p *PagingCursor) String() string {
	var o = "<"

	for i, key := range p.keys {
		o += fmt.Sprintf("%s: %v, ", key, p.values[i])
	}

	if p.Reverse {
		o += "reverse"
	} else {
		o += "forward"
	}

	return o + ">"
}

func (p *PagingCursor) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(struct {
		K []string
		V []interface{}
		D []bool
		R bool
	}{
		p.keys,
		p.values,
		p.desc,
		p.Reverse,
	})

	if err != nil {
		return nil, err
	}

	std := base64.StdEncoding
	dbuf := make([]byte, std.EncodedLen(len(buf)))
	std.Encode(dbuf, buf)

	return append([]byte{'"'}, append(dbuf, '"')...), nil
}

func (p *PagingCursor) Encode() string {
	b, _ := p.MarshalJSON()
	return string(b)
}

func (p *PagingCursor) UnmarshalJSON(in []byte) error {
	var (
		aux struct {
			K []string
			V []interface{}
			D []bool
			R bool
		}
	)

	if err := json.Unmarshal(in, &aux); err != nil {
		return err
	}

	p.keys = aux.K
	p.values = aux.V
	p.desc = aux.D
	p.Reverse = aux.R

	return nil
}

func (p *PagingCursor) Decode(cursor string) error {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	err = p.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	return err
}

func parseCursor(in string) (p *PagingCursor, err error) {
	if len(in) == 0 {
		return nil, nil
	}

	var buf []byte
	if buf, err = base64.StdEncoding.DecodeString(in); err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	p = &PagingCursor{}
	if err = p.UnmarshalJSON(buf); err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	return p, nil
}
