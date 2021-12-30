package gig

import (
	"encoding/json"
	"fmt"
)

type (
	preprocessorNoop struct{}
)

var (
	PreprocessorHandleNoop preprocessor = "noop"
)

func UnwrapPreprocessorSet(ww PreprocessorWrapSet) (out []Preprocessor, err error) {
	for _, w := range ww {
		var aux Preprocessor
		aux, err = UnwrapPreprocessor(w)
		if err != nil {
			return
		}

		out = append(out, aux)
	}
	return
}

func UnwrapPreprocessor(w PreprocessorWrap) (Preprocessor, error) {
	switch w.Ref {
	case PreprocessorHandleNoop:
		return PreprocessorNoop()

	case PreprocessorHandleResourceRemove:
		var aux preprocessorResourceRemove
		if w.Params != nil {
			if err := json.Unmarshal(w.Params, &aux); err != nil {
				return nil, err
			}
		}
		return PreprocessorResourceRemove(aux.Resource, aux.Identifier), nil

	case PreprocessorHandleResourceLoad:
		var aux preprocessorResourceLoad
		if w.Params != nil {
			if err := json.Unmarshal(w.Params, &aux); err != nil {
				return nil, err
			}
		}
		if aux.ID != 0 {
			return PreprocessorResourceLoadByID(aux.Resource, aux.ID), nil
		} else if aux.Handle != "" {
			return PreprocessorResourceLoadByHandle(aux.Resource, aux.Handle), nil
		} else if aux.Query != "" {
			return PreprocessorResourceLoadByQuery(aux.Resource, aux.Query), nil
		}

		return nil, fmt.Errorf("invalid arguments to preprocessor: %s: expecting ID, handle, or query", PreprocessorHandleResourceLoad)

	case PreprocessorHandleNamespaceLoad:
		var aux preprocessorNamespaceLoad
		if w.Params != nil {
			if err := json.Unmarshal(w.Params, &aux); err != nil {
				return nil, err
			}
		}
		if aux.ID != 0 {
			return PreprocessorNamespaceLoadByID(aux.ID), nil
		} else if aux.Handle != "" {
			return PreprocessorNamespaceLoadByHandle(aux.Handle), nil
		}

		return nil, fmt.Errorf("invalid arguments to preprocessor: %s: expecting ID, handle, or query", PreprocessorHandleNamespaceLoad)
	}

	return nil, fmt.Errorf("unknown preprocessor: %s", w.Ref)
}

func WrapPreprocessor(w Preprocessor) (out PreprocessorWrap, err error) {
	enc, err := json.Marshal(w)
	if err != nil {
		return
	}

	switch w.(type) {
	case preprocessorNoop:
		out = PreprocessorWrap{
			Ref:    PreprocessorHandleNoop,
			Params: enc,
		}
		return

	case preprocessorResourceRemove:
		out = PreprocessorWrap{
			Ref:    PreprocessorHandleResourceRemove,
			Params: enc,
		}
		return

	case preprocessorResourceLoad:
		out = PreprocessorWrap{
			Ref:    PreprocessorHandleResourceLoad,
			Params: enc,
		}
		return

	case preprocessorNamespaceLoad:
		out = PreprocessorWrap{
			Ref:    PreprocessorHandleNamespaceLoad,
			Params: enc,
		}
		return

	default:
		err = fmt.Errorf("unknown preprocessor: %T", w)
		return
	}
}

func PreprocessorDefinitions() (out TaskDefSet) {
	return TaskDefSet{
		// root
		{
			Ref:  string(PreprocessorHandleNoop),
			Kind: TaskPreprocessor,
		}, {
			Ref:    string(DecoderHandleNoop),
			Kind:   TaskDecoder,
			Params: []taskDefParam{},
		},

		// envoy
		{
			Ref:  string(PreprocessorHandleResourceRemove),
			Kind: TaskPreprocessor,
			Params: []taskDefParam{
				{
					Name:        "resource",
					Kind:        "String",
					Required:    true,
					Description: "The resource type to remove",
				},
				{
					Name:        "identifier",
					Kind:        "String",
					Required:    false,
					Description: "The resource identifier to remove",
				},
			},
		},
		{
			Ref:  string(PreprocessorHandleResourceLoad),
			Kind: TaskPreprocessor,
			Params: []taskDefParam{
				{
					Name:        "resource",
					Kind:        "String",
					Required:    true,
					Description: "The resource type to load",
				},
				{
					Name:        "id",
					Kind:        "ID",
					Required:    false,
					Description: "Load the resource with this ID",
				},
				{
					Name:        "handle",
					Kind:        "Handle",
					Required:    false,
					Description: "Load the resource with this handle/slug",
				},
				{
					Name:        "query",
					Kind:        "String",
					Required:    false,
					Description: "Load the resource matching this query",
				},
			},
		},
		{
			Ref:  string(PreprocessorHandleNamespaceLoad),
			Kind: TaskPreprocessor,
			Params: []taskDefParam{
				{
					Name:        "id",
					Kind:        "ID",
					Required:    false,
					Description: "Load the namespace by this ID",
				},
				{
					Name:        "handle",
					Kind:        "Handle",
					Required:    false,
					Description: "Load the namespace with this handle",
				},
			},
		},

		{
			Ref:  string(PreprocessorHandleAttachmentRemove),
			Kind: TaskPreprocessor,
			Params: []taskDefParam{
				{
					Name:        "mimeType",
					Kind:        "String",
					Required:    true,
					Description: "File mime types to remove",
				},
			},
		},
		{
			Ref:  string(PreprocessorHandleAttachmentTransform),
			Kind: TaskPreprocessor,
			Params: []taskDefParam{
				{
					Name:        "width",
					Kind:        "Number",
					Required:    false,
					Description: "Width",
				},
				{
					Name:        "height",
					Kind:        "Number",
					Required:    false,
					Description: "Height",
				},
			},
		},
	}
}
