package encoder

import (
	syst "github.com/cortezaproject/corteza-server/system/types"
)

type (
	multiple uint

	field struct {
		name           string
		encodeAllMulti bool
	}

	FlatWriter interface {
		Write([]string) error
		Flush()
	}

	StructuredEncoder interface {
		Encode(interface{}) error
	}

	userFinder func(ID uint64) (*syst.User, error)

	flatWriter struct {
		w  FlatWriter
		ff []field
		u  userFinder
	}

	structuredEncoder struct {
		w  StructuredEncoder
		ff []field
		u  userFinder
	}
)

func Field(name string) field {
	return field{name: name}
}

func MakeFields(nn ...string) []field {
	ff := make([]field, len(nn))
	for i := range nn {
		ff[i] = field{name: nn[i]}
	}

	return ff
}

func MultiValueField(name string) field {
	return field{name: name, encodeAllMulti: true}
}

func NewFlatWriter(w FlatWriter, header bool, u userFinder, ff ...field) *flatWriter {
	f := &flatWriter{
		w:  w,
		ff: ff,
		u:  u,
	}

	if header {
		f.writeHeader()
	}

	return f
}

func (enc flatWriter) Flush() {
	enc.w.Flush()
}

func (enc flatWriter) writeHeader() {
	ss := make([]string, len(enc.ff))
	for i := range enc.ff {
		ss[i] = enc.ff[i].name
	}

	_ = enc.w.Write(ss)
}

func NewStructuredEncoder(w StructuredEncoder, u userFinder, ff ...field) *structuredEncoder {
	return &structuredEncoder{
		w:  w,
		ff: ff,
		u:  u,
	}
}

func (enc structuredEncoder) Flush() {
	// noop
}
