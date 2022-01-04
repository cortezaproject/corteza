package conv

import (
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/gig"
)

type (
	Gig struct{}

	ParamWrap struct {
		Ref    string
		Params map[string]interface{}
	}
	ParamWrapSet []ParamWrap
)

func (conv Gig) UnwrapDecoderSet(wraps ParamWrapSet) (out gig.DecoderSet, err error) {
	for _, w := range wraps {
		var aux gig.Decoder
		aux, err = conv.UnwrapDecoder(w)
		if err != nil {
			return
		}

		out = append(out, aux)
	}
	return
}

func (conv Gig) UnwrapDecoder(w ParamWrap) (out gig.Decoder, err error) {
	switch w.Ref {
	case gig.DecoderHandleArchive:
		return gig.DecoderArchiveParams(w.Params), nil
	case gig.DecoderHandleNoop:
		return gig.DecoderNoopParams(w.Params), nil
	}

	return nil, fmt.Errorf("unknown decoder: %s", w.Ref)
}

func (conv Gig) WrapDecoder(w gig.Decoder) (out ParamWrap, err error) {
	out.Ref = w.Ref()
	out.Params = w.Params()

	return
}

func (conv Gig) UnwrapPreprocessorSet(wraps ParamWrapSet) (out gig.PreprocessorSet, err error) {
	for _, w := range wraps {
		var aux gig.Preprocessor
		aux, err = conv.UnwrapPreprocessor(w)
		if err != nil {
			return
		}

		out = append(out, aux)
	}
	return
}

func (conv Gig) UnwrapPreprocessor(w ParamWrap) (out gig.Preprocessor, err error) {
	switch w.Ref {
	case gig.PreprocessorHandleNoop:
		return gig.PreprocessorNoopParams(w.Params), nil

	case gig.PreprocessorHandleResourceRemove:
		return gig.PreprocessorResourceRemoveParams(w.Params), nil

	case gig.PreprocessorHandleResourceLoad:
		return gig.PreprocessorResourceLoadParams(w.Params), nil

	case gig.PreprocessorHandleNamespaceLoad:
		return gig.PreprocessorNamespaceLoadParams(w.Params), nil
	}

	return nil, fmt.Errorf("unknown preprocessor: %s", w.Ref)
}

func (conv Gig) WrapPreprocessor(w gig.Preprocessor) (out ParamWrap, err error) {
	out.Ref = w.Ref()
	out.Params = w.Params()

	return
}

func (conv Gig) UnwrapPostprocessorSet(wraps ParamWrapSet) (out gig.PostprocessorSet, err error) {
	for _, w := range wraps {
		var aux gig.Postprocessor
		aux, err = conv.UnwrapPostprocessor(w)
		if err != nil {
			return
		}

		out = append(out, aux)
	}
	return
}

func (conv Gig) UnwrapPostprocessor(w ParamWrap) (out gig.Postprocessor, err error) {
	switch w.Ref {
	case gig.PostprocessorHandleArchive:
		return gig.PostprocessorArchiveParams(w.Params), nil

	case gig.PostprocessorHandleDiscard:
		return gig.PostprocessorDiscardParams(w.Params), nil

	case gig.PostprocessorHandleNoop:
		return gig.PostprocessorNoopParams(w.Params), nil

	case gig.PostprocessorHandleSave:
		return gig.PostprocessorSaveParams(w.Params), nil
	}

	return nil, fmt.Errorf("unknown postprocessor: %s", w.Ref)
}

func (conv Gig) WrapPostprocessor(w gig.Postprocessor) (out ParamWrap, err error) {
	out.Ref = w.Ref()
	out.Params = w.Params()

	return
}

// REST parsers

func ParseParamWrap(ss []string) (out ParamWrapSet, err error) {
	for _, s := range ss {
		aux := make(ParamWrapSet, 0, 2)
		err = json.Unmarshal([]byte(s), &aux)
		if err != nil {
			return
		}

		out = append(out, aux...)
	}
	return
}
