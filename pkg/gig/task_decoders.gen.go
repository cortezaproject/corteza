package gig

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/spf13/cast"
)

type (
	decoderNoop struct {
		source uint64
	}
	decoderArchive struct {
		source uint64
	}
)

const (
	DecoderHandleNoop    = "noop"
	DecoderHandleArchive = "archive"
)

// ------------------------------------------------------------------------
// Constructors and utils

// DecoderNoopParams returns a new decoderNoop from the params
func DecoderNoopParams(params map[string]interface{}) Decoder {
	out := decoderNoop{
		source: cast.ToUint64(params["source"]),
	}
	return out
}

func DecoderNoop(source uint64) Decoder {
	out := decoderNoop{
		source: source,
	}

	return out
}

func (t decoderNoop) Ref() string {
	return DecoderHandleNoop
}

func (t decoderNoop) Params() map[string]interface{} {
	return map[string]interface{}{
		"source": t.source,
	}
}

// DecoderArchiveParams returns a new decoderArchive from the params
func DecoderArchiveParams(params map[string]interface{}) Decoder {
	out := decoderArchive{
		source: cast.ToUint64(params["source"]),
	}
	return out
}

func DecoderArchive(source uint64) Decoder {
	out := decoderArchive{
		source: source,
	}

	return out
}

func (t decoderArchive) Ref() string {
	return DecoderHandleArchive
}

func (t decoderArchive) Params() map[string]interface{} {
	return map[string]interface{}{
		"source": t.source,
	}
}

// ------------------------------------------------------------------------
// Task registry

func decoderDefinitions() TaskDefSet {
	return TaskDefSet{
		{
			Ref:         DecoderHandleNoop,
			Kind:        TaskDecoder,
			Description: "Noop does nothing.",
			Params: []taskDefParam{
				{
					Name:     "source",
					Kind:     "String",
					Required: true,
				},
			},
		},
		{
			Ref:         DecoderHandleArchive,
			Kind:        TaskDecoder,
			Description: "Extracts the contents of the archive into sepparate sources; extraction is not recursive.",
			Params: []taskDefParam{
				{
					Name:     "source",
					Kind:     "String",
					Required: true,
				},
			},
		},
	}
}
