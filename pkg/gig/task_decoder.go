package gig

import (
	"context"
	"fmt"
)

type (
	Decoder interface {
		Source() uint64
		Clone(uint64) Decoder
		CanDecode(Source) bool
		Decode(context.Context, Source) (SourceSet, error)
		Ref() string
		Params() map[string]interface{}
	}
	DecoderSet []Decoder
)

// Noop

func (d decoderNoop) Clone(rel uint64) Decoder {
	return decoderNoop{
		source: rel,
	}
}

func (d decoderNoop) Source() uint64 {
	return d.source
}

func (d decoderNoop) CanDecode(src Source) bool {
	return true
}

func (d decoderNoop) Decode(ctx context.Context, in Source) (out SourceSet, err error) {
	return SourceSet{in}, nil
}

// Archive

func (d decoderArchive) Clone(rel uint64) Decoder {
	return decoderArchive{
		source: rel,
	}
}

func (d decoderArchive) Source() uint64 {
	return d.source
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
