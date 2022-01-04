package gig

type (
	decoderNoop struct {
		source uint64
	}
	decoderArchive struct {
		source uint64
	}
)

var (
	DecoderHandleNoop    = "noop"
	DecoderHandleArchive = "archive"
)

func decoderDefinitions() (out TaskDefSet) {
	return TaskDefSet{{
		Ref:  string(DecoderHandleArchive),
		Kind: TaskDecoder,
		Params: []taskDefParam{{
			Name:        "rel",
			Kind:        "ID",
			Required:    true,
			Description: "The related source this decoder applies to",
		}},
	}, {
		Ref:    string(DecoderHandleNoop),
		Kind:   TaskDecoder,
		Params: []taskDefParam{},
	}}
}
