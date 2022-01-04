package gig

type (
	preprocessorNoop struct{}
)

var (
	PreprocessorHandleNoop = "noop"
)

func preprocessorDefinitions() (out TaskDefSet) {
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

		// Attachment
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
