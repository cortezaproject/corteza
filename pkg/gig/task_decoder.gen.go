package gig

import (
	"encoding/json"
	"fmt"
)

type (
	decoderNoop struct {
		Source uint64 `json:"source,string"`
	}
	decoderArchive struct {
		Source uint64 `json:"source,string"`
	}
)

var (
	DecoderHandleNoop    decoder = "noop"
	DecoderHandleArchive decoder = "archive"
)

func UnwrapDecoderSet(ww DecoderWrapSet) (out []Decoder, err error) {
	for _, w := range ww {
		var aux Decoder
		aux, err = UnwrapDecoder(w)
		if err != nil {
			return
		}

		out = append(out, aux)
	}
	return
}

func UnwrapDecoder(w DecoderWrap) (Decoder, error) {
	switch w.Ref {
	case DecoderHandleArchive:
		var aux decoderArchive
		if w.Params != nil {
			if err := json.Unmarshal(w.Params, &aux); err != nil {
				return nil, err
			}
		}
		return DecoderArchive(aux.Source), nil
	case DecoderHandleNoop:
		var aux decoderNoop
		if w.Params != nil {
			if err := json.Unmarshal(w.Params, &aux); err != nil {
				return nil, err
			}
		}
		return DecoderNoop(aux.Source), nil
	}

	return nil, fmt.Errorf("unknown decoder: %s", w.Ref)
}

func WrapDecoder(w Decoder) (out DecoderWrap, err error) {
	enc, err := json.Marshal(w)
	if err != nil {
		return
	}

	switch w.(type) {
	case decoderArchive:
		out = DecoderWrap{
			Ref:    DecoderHandleArchive,
			Params: enc,
		}
		return
	case decoderNoop:
		out = DecoderWrap{
			Ref:    DecoderHandleNoop,
			Params: enc,
		}
		return
	default:
		err = fmt.Errorf("unknown decoder: %T", w)
		return
	}
}

func DecoderDefinitions() (out TaskDefSet) {
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
