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
		tz string
	}

	structuredEncoder struct {
		w  StructuredEncoder
		ff []field
		u  userFinder
		tz string
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

func preprocessHeader(hh []field, tz string) []field {
	nhh := make([]field, 0)

	// We need to prepare additional header fields for exporting
	if tz != "" && tz != "UTC" {
		for _, f := range hh {
			switch f.name {
			case "createdAt",
				"updatedAt",
				"deletedAt":

				nhh = append(nhh, f, field{name: f.name + "_date"}, field{name: f.name + "_time"})
				break
			default:
				nhh = append(nhh, f)
				break
			}
		}
	} else {
		return hh
	}
	return nhh
}

func MultiValueField(name string) field {
	return field{name: name, encodeAllMulti: true}
}

func NewFlatWriter(w FlatWriter, header bool, u userFinder, tz string, ff ...field) *flatWriter {
	f := &flatWriter{
		w:  w,
		ff: preprocessHeader(ff, tz),
		u:  u,
		tz: tz,
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

func NewStructuredEncoder(w StructuredEncoder, u userFinder, tz string, ff ...field) *structuredEncoder {
	return &structuredEncoder{
		w: w,
		// No need for additional timezone headers, since the output is structured
		ff: ff,
		u:  u,
		tz: tz,
	}
}

func (enc structuredEncoder) Flush() {
	// noop
}
