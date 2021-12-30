package gig

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	Decoder interface {
		RelSource() uint64
		Clone(uint64) Decoder
		CanDecode(Source) bool
		Decode(context.Context, Source) (SourceSet, error)
	}

	DecoderWrap struct {
		Ref    decoder `json:"ref"`
		Params json.RawMessage
	}
	DecoderWrapSet []DecoderWrap

	decoder string
)

// Noop

func DecoderNoop(rel uint64) Decoder {
	return decoderNoop{Source: rel}
}

func (d decoderNoop) Clone(rel uint64) Decoder {
	return DecoderNoop(rel)
}

func (d decoderNoop) RelSource() uint64 {
	return d.Source
}

func (d decoderNoop) CanDecode(src Source) bool {
	return true
}

func (d decoderNoop) Decode(ctx context.Context, in Source) (out SourceSet, err error) {
	return SourceSet{in}, nil
}

// Archive

func DecoderArchive(rel uint64) Decoder {
	return decoderArchive{Source: rel}
}

func (d decoderArchive) Clone(rel uint64) Decoder {
	return DecoderArchive(rel)
}

func (d decoderArchive) RelSource() uint64 {
	return d.Source
}

func (d decoderArchive) CanDecode(src Source) bool {
	// @todo others...
	return isTarGz(src)
}

func (d decoderArchive) Decode(ctx context.Context, in Source) (out SourceSet, err error) {
	switch {
	case isTarGz(in):
		return extractTarGz(ctx, in)
	}

	err = fmt.Errorf("unknown archive: %s", in.MimeType())
	return
}

func ParseDecoderWrap(ss []string) (out DecoderWrapSet, err error) {
	for _, s := range ss {
		aux := make(DecoderWrapSet, 0, 2)
		err = json.Unmarshal([]byte(s), &aux)
		if err != nil {
			return
		}

		out = append(out, aux...)
	}
	return
}
