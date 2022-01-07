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
	preprocessorNoop             struct{}
	preprocessorAttachmentRemove struct {
		mimeType string
	}
	preprocessorAttachmentTransform struct {
		height int
		width  int
	}
	preprocessorResourceRemove struct {
		identifier string
		resource   string
	}
	preprocessorResourceLoad struct {
		handle   string
		id       uint64
		query    string
		resource string
	}
	preprocessorNamespaceLoad struct {
		handle string
		id     uint64
	}
)

const (
	PreprocessorHandleNoop                = "noop"
	PreprocessorHandleAttachmentRemove    = "attachmentRemove"
	PreprocessorHandleAttachmentTransform = "attachmentTransform"
	PreprocessorHandleResourceRemove      = "resourceRemove"
	PreprocessorHandleResourceLoad        = "resourceLoad"
	PreprocessorHandleNamespaceLoad       = "namespaceLoad"
)

// ------------------------------------------------------------------------
// Constructors and utils

// PreprocessorNoopParams returns a new preprocessorNoop from the params
func PreprocessorNoopParams(params map[string]interface{}) Preprocessor {
	out := preprocessorNoop{}
	return out
}

func (t preprocessorNoop) Ref() string {
	return PreprocessorHandleNoop
}

func (t preprocessorNoop) Params() map[string]interface{} {
	return nil
}

// PreprocessorAttachmentRemoveParams returns a new preprocessorAttachmentRemove from the params
func PreprocessorAttachmentRemoveParams(params map[string]interface{}) Preprocessor {
	out := preprocessorAttachmentRemove{
		mimeType: cast.ToString(params["mimeType"]),
	}
	return out
}

func PreprocessorAttachmentRemove(mimeType string) Preprocessor {
	out := preprocessorAttachmentRemove{
		mimeType: mimeType,
	}

	return out
}

func (t preprocessorAttachmentRemove) Ref() string {
	return PreprocessorHandleAttachmentRemove
}

func (t preprocessorAttachmentRemove) Params() map[string]interface{} {
	return map[string]interface{}{
		"mimeType": t.mimeType,
	}
}

// PreprocessorAttachmentTransformParams returns a new preprocessorAttachmentTransform from the params
func PreprocessorAttachmentTransformParams(params map[string]interface{}) Preprocessor {
	out := preprocessorAttachmentTransform{
		height: cast.ToInt(params["height"]),

		width: cast.ToInt(params["width"]),
	}
	return out
}

// PreprocessorAttachmentTransformHeight returns a new preprocessorAttachmentTransform from the required fields and height
func PreprocessorAttachmentTransformHeight(height int) Preprocessor {
	out := preprocessorAttachmentTransform{
		height: height,
	}

	return out
}

// PreprocessorAttachmentTransformWidth returns a new preprocessorAttachmentTransform from the required fields and width
func PreprocessorAttachmentTransformWidth(width int) Preprocessor {
	out := preprocessorAttachmentTransform{
		width: width,
	}

	return out
}

func (t preprocessorAttachmentTransform) Ref() string {
	return PreprocessorHandleAttachmentTransform
}

func (t preprocessorAttachmentTransform) Params() map[string]interface{} {
	return map[string]interface{}{
		"height": t.height,
		"width":  t.width,
	}
}

// PreprocessorResourceRemoveParams returns a new preprocessorResourceRemove from the params
func PreprocessorResourceRemoveParams(params map[string]interface{}) Preprocessor {
	out := preprocessorResourceRemove{
		identifier: cast.ToString(params["identifier"]),

		resource: cast.ToString(params["resource"]),
	}
	out = preprocessorResourceRemoveTransformer(out)
	return out
}

// PreprocessorResourceRemoveIdentifier returns a new preprocessorResourceRemove from the required fields and identifier
func PreprocessorResourceRemoveIdentifier(resource string, identifier string) Preprocessor {
	out := preprocessorResourceRemove{
		resource:   resource,
		identifier: identifier,
	}
	out = preprocessorResourceRemoveTransformer(out)

	return out
}
func PreprocessorResourceRemove(resource string) Preprocessor {
	out := preprocessorResourceRemove{
		resource: resource,
	}
	out = preprocessorResourceRemoveTransformer(out)

	return out
}

func (t preprocessorResourceRemove) Ref() string {
	return PreprocessorHandleResourceRemove
}

func (t preprocessorResourceRemove) Params() map[string]interface{} {
	return map[string]interface{}{
		"identifier": t.identifier,
		"resource":   t.resource,
	}
}

// PreprocessorResourceLoadParams returns a new preprocessorResourceLoad from the params
func PreprocessorResourceLoadParams(params map[string]interface{}) Preprocessor {
	out := preprocessorResourceLoad{
		handle: cast.ToString(params["handle"]),

		id: cast.ToUint64(params["id"]),

		query: cast.ToString(params["query"]),

		resource: cast.ToString(params["resource"]),
	}
	return out
}

// PreprocessorResourceLoadHandle returns a new preprocessorResourceLoad from the required fields and handle
func PreprocessorResourceLoadHandle(resource string, handle string) Preprocessor {
	out := preprocessorResourceLoad{
		resource: resource,
		handle:   handle,
	}

	return out
}

// PreprocessorResourceLoadId returns a new preprocessorResourceLoad from the required fields and id
func PreprocessorResourceLoadId(resource string, id uint64) Preprocessor {
	out := preprocessorResourceLoad{
		resource: resource,
		id:       id,
	}

	return out
}

// PreprocessorResourceLoadQuery returns a new preprocessorResourceLoad from the required fields and query
func PreprocessorResourceLoadQuery(resource string, query string) Preprocessor {
	out := preprocessorResourceLoad{
		resource: resource,
		query:    query,
	}

	return out
}
func PreprocessorResourceLoad(resource string) Preprocessor {
	out := preprocessorResourceLoad{
		resource: resource,
	}

	return out
}

func (t preprocessorResourceLoad) Ref() string {
	return PreprocessorHandleResourceLoad
}

func (t preprocessorResourceLoad) Params() map[string]interface{} {
	return map[string]interface{}{
		"handle":   t.handle,
		"id":       t.id,
		"query":    t.query,
		"resource": t.resource,
	}
}

// PreprocessorNamespaceLoadParams returns a new preprocessorNamespaceLoad from the params
func PreprocessorNamespaceLoadParams(params map[string]interface{}) Preprocessor {
	out := preprocessorNamespaceLoad{
		handle: cast.ToString(params["handle"]),

		id: cast.ToUint64(params["id"]),
	}
	return out
}

// PreprocessorNamespaceLoadHandle returns a new preprocessorNamespaceLoad from the required fields and handle
func PreprocessorNamespaceLoadHandle(handle string) Preprocessor {
	out := preprocessorNamespaceLoad{
		handle: handle,
	}

	return out
}

// PreprocessorNamespaceLoadId returns a new preprocessorNamespaceLoad from the required fields and id
func PreprocessorNamespaceLoadId(id uint64) Preprocessor {
	out := preprocessorNamespaceLoad{
		id: id,
	}

	return out
}

func (t preprocessorNamespaceLoad) Ref() string {
	return PreprocessorHandleNamespaceLoad
}

func (t preprocessorNamespaceLoad) Params() map[string]interface{} {
	return map[string]interface{}{
		"handle": t.handle,
		"id":     t.id,
	}
}

// ------------------------------------------------------------------------
// Task registry

func preprocessorDefinitions() TaskDefSet {
	return TaskDefSet{
		{
			Ref:         PreprocessorHandleNoop,
			Kind:        TaskPreprocessor,
			Description: "Noop does nothing.",
		},
		{
			Ref:         PreprocessorHandleAttachmentRemove,
			Kind:        TaskPreprocessor,
			Description: "Removes the attachment.",
			Params: []taskDefParam{
				{
					Name:     "mimeType",
					Kind:     "String",
					Required: true,
				},
			},
		},
		{
			Ref:         PreprocessorHandleAttachmentTransform,
			Kind:        TaskPreprocessor,
			Description: "Applies the specified transformations.",
			Params: []taskDefParam{
				{
					Name:     "height",
					Kind:     "Number",
					Required: false,
				},
				{
					Name:     "width",
					Kind:     "Number",
					Required: false,
				},
			},
		},
		{
			Ref:         PreprocessorHandleResourceRemove,
			Kind:        TaskPreprocessor,
			Description: "Removes the specified resource.",
			Params: []taskDefParam{
				{
					Name:     "identifier",
					Kind:     "String",
					Required: false,
				},
				{
					Name:     "resource",
					Kind:     "String",
					Required: true,
				},
			},
		},
		{
			Ref:         PreprocessorHandleResourceLoad,
			Kind:        TaskPreprocessor,
			Description: "Loads the specified resource from internal storage.",
			Params: []taskDefParam{
				{
					Name:     "handle",
					Kind:     "String",
					Required: false,
				},
				{
					Name:     "id",
					Kind:     "string",
					Required: false,
				},
				{
					Name:     "query",
					Kind:     "String",
					Required: false,
				},
				{
					Name:     "resource",
					Kind:     "String",
					Required: true,
				},
			},
		},
		{
			Ref:         PreprocessorHandleNamespaceLoad,
			Kind:        TaskPreprocessor,
			Description: "Loads the namespace with a predefined set of sub-resources.",
			Params: []taskDefParam{
				{
					Name:     "handle",
					Kind:     "String",
					Required: false,
				},
				{
					Name:     "id",
					Kind:     "string",
					Required: false,
				},
			},
		},
	}
}
