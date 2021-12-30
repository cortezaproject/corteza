package gig

import (
	"encoding/json"
	"fmt"
)

type (
	postprocessorNoop    struct{}
	postprocessorDiscard struct{}
	postprocessorSave    struct{}
	postprocessorArchive struct {
		Encoding archive `json:"encoding,string"`
		Name     string  `json:"name,string"`
	}
)

var (
	PostprocessorHandleArchive postprocessor = "archive"
	PostprocessorHandleDiscard postprocessor = "discard"
	PostprocessorHandleNoop    postprocessor = "noop"
	PostprocessorHandleSave    postprocessor = "save"
)

func UnwrapPostprocessorSet(ww PostprocessorWrapSet) (out []Postprocessor, err error) {
	for _, w := range ww {
		var aux Postprocessor
		aux, err = UnwrapPostprocessor(w)
		if err != nil {
			return
		}

		out = append(out, aux)
	}
	return
}

func UnwrapPostprocessor(w PostprocessorWrap) (Postprocessor, error) {
	switch w.Ref {
	case PostprocessorHandleArchive:
		var aux postprocessorArchive
		if w.Params != nil {
			if err := json.Unmarshal(w.Params, &aux); err != nil {
				return nil, err
			}
		}
		return PostprocessorArchive(aux.Encoding, aux.Name)

	case PostprocessorHandleDiscard:
		return PostprocessorDiscard()

	case PostprocessorHandleNoop:
		return PostprocessorNoop()

	case PostprocessorHandleSave:
		return PostprocessorSave()
	}

	return nil, fmt.Errorf("unknown postprocessor: %s", w.Ref)
}

func WrapPostprocessTask(w Postprocessor) (out PostprocessorWrap, err error) {
	enc, err := json.Marshal(w)
	if err != nil {
		return
	}

	switch w.(type) {
	case postprocessorDiscard:
		out = PostprocessorWrap{
			Ref:    PostprocessorHandleDiscard,
			Params: enc,
		}
		return
	case postprocessorNoop:
		out = PostprocessorWrap{
			Ref:    PostprocessorHandleNoop,
			Params: enc,
		}
		return
	case postprocessorArchive:
		out = PostprocessorWrap{
			Ref:    PostprocessorHandleArchive,
			Params: enc,
		}
		return

	default:
		err = fmt.Errorf("unknown postprocessor: %T", w)
		return
	}
}

func PostprocessorDefinitions() TaskDefSet {
	return TaskDefSet{
		{
			Ref:  string(PostprocessorHandleDiscard),
			Kind: TaskDecoder,
		},
		{
			Ref:  string(PostprocessorHandleNoop),
			Kind: TaskDecoder,
		},
		{
			Ref:  string(PostprocessorHandleSave),
			Kind: TaskDecoder,
		},
		{
			Ref:  string(PostprocessorHandleArchive),
			Kind: TaskDecoder,
			Params: []taskDefParam{
				{
					Name:        "format",
					Kind:        "String",
					Required:    false,
					Description: "What archive format to encode into",
				},
				{
					Name:        "name",
					Kind:        "String",
					Required:    false,
					Description: "The output name of the archive",
				},
			},
		},
	}
}
