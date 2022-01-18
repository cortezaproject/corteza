package conv

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/gig"
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
	case gig.DecoderHandleNoop:
		return gig.DecoderNoopParams(w.Params), nil
	case gig.DecoderHandleArchive:
		return gig.DecoderArchiveParams(w.Params), nil
	}

	return nil, fmt.Errorf("unknown decoder: %s", w.Ref)
}

func (conv Gig) WrapDecoderSet(tt gig.DecoderSet) (out ParamWrapSet) {
	for _, t := range tt {
		out = append(out, conv.WrapDecoder(t))
	}

	return
}

func (conv Gig) WrapDecoder(t gig.Decoder) (out ParamWrap) {
	out.Ref = t.Ref()
	out.Params = t.Params()

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
	case gig.PreprocessorHandleAttachmentRemove:
		return gig.PreprocessorAttachmentRemoveParams(w.Params), nil
	case gig.PreprocessorHandleAttachmentTransform:
		return gig.PreprocessorAttachmentTransformParams(w.Params), nil
	case gig.PreprocessorHandleExperimentalExport:
		return gig.PreprocessorExperimentalExportParams(w.Params), nil
	}

	return nil, fmt.Errorf("unknown preprocessor: %s", w.Ref)
}

func (conv Gig) WrapPreprocessorSet(tt gig.PreprocessorSet) (out ParamWrapSet) {
	for _, t := range tt {
		out = append(out, conv.WrapPreprocessor(t))
	}

	return
}

func (conv Gig) WrapPreprocessor(t gig.Preprocessor) (out ParamWrap) {
	out.Ref = t.Ref()
	out.Params = t.Params()

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
	case gig.PostprocessorHandleNoop:
		return gig.PostprocessorNoopParams(w.Params), nil
	case gig.PostprocessorHandleDiscard:
		return gig.PostprocessorDiscardParams(w.Params), nil
	case gig.PostprocessorHandleSave:
		return gig.PostprocessorSaveParams(w.Params), nil
	case gig.PostprocessorHandleArchive:
		return gig.PostprocessorArchiveParams(w.Params), nil
	}

	return nil, fmt.Errorf("unknown postprocessor: %s", w.Ref)
}

func (conv Gig) WrapPostprocessorSet(tt gig.PostprocessorSet) (out ParamWrapSet) {
	for _, t := range tt {
		out = append(out, conv.WrapPostprocessor(t))
	}

	return
}

func (conv Gig) WrapPostprocessor(t gig.Postprocessor) (out ParamWrap) {
	out.Ref = t.Ref()
	out.Params = t.Params()

	return
}
