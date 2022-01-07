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
	postprocessorNoop    struct{}
	postprocessorDiscard struct{}
	postprocessorSave    struct{}
	postprocessorArchive struct {
		encoding archive
		name     string
	}
)

const (
	PostprocessorHandleNoop    = "noop"
	PostprocessorHandleDiscard = "discard"
	PostprocessorHandleSave    = "save"
	PostprocessorHandleArchive = "archive"
)

// ------------------------------------------------------------------------
// Constructors and utils

// PostprocessorNoopParams returns a new postprocessorNoop from the params
func PostprocessorNoopParams(params map[string]interface{}) Postprocessor {
	out := postprocessorNoop{}
	return out
}

func (t postprocessorNoop) Ref() string {
	return PostprocessorHandleNoop
}

func (t postprocessorNoop) Params() map[string]interface{} {
	return nil
}

// PostprocessorDiscardParams returns a new postprocessorDiscard from the params
func PostprocessorDiscardParams(params map[string]interface{}) Postprocessor {
	out := postprocessorDiscard{}
	return out
}

func (t postprocessorDiscard) Ref() string {
	return PostprocessorHandleDiscard
}

func (t postprocessorDiscard) Params() map[string]interface{} {
	return nil
}

// PostprocessorSaveParams returns a new postprocessorSave from the params
func PostprocessorSaveParams(params map[string]interface{}) Postprocessor {
	out := postprocessorSave{}
	return out
}

func (t postprocessorSave) Ref() string {
	return PostprocessorHandleSave
}

func (t postprocessorSave) Params() map[string]interface{} {
	return nil
}

// PostprocessorArchiveParams returns a new postprocessorArchive from the params
func PostprocessorArchiveParams(params map[string]interface{}) Postprocessor {
	out := postprocessorArchive{
		encoding: archiveFromParams(params["encoding"]),

		name: cast.ToString(params["name"]),
	}
	out = postprocessorArchiveTransformer(out)
	return out
}

// PostprocessorArchiveName returns a new postprocessorArchive from the required fields and name
func PostprocessorArchiveName(encoding archive, name string) Postprocessor {
	out := postprocessorArchive{
		encoding: encoding,
		name:     name,
	}
	out = postprocessorArchiveTransformer(out)

	return out
}
func PostprocessorArchive(encoding archive) Postprocessor {
	out := postprocessorArchive{
		encoding: encoding,
	}
	out = postprocessorArchiveTransformer(out)

	return out
}

func (t postprocessorArchive) Ref() string {
	return PostprocessorHandleArchive
}

func (t postprocessorArchive) Params() map[string]interface{} {
	return map[string]interface{}{
		"encoding": t.encoding,
		"name":     t.name,
	}
}

// ------------------------------------------------------------------------
// Task registry

func postprocessorDefinitions() TaskDefSet {
	return TaskDefSet{
		{
			Ref:         PostprocessorHandleNoop,
			Kind:        TaskPostprocessor,
			Description: "Noop does nothing.",
		},
		{
			Ref:         PostprocessorHandleDiscard,
			Kind:        TaskPostprocessor,
			Description: "Discards the resulting sources.",
		},
		{
			Ref:         PostprocessorHandleSave,
			Kind:        TaskPostprocessor,
			Description: "Saves the resulting sources to perminant storage; not implemented.",
		},
		{
			Ref:         PostprocessorHandleArchive,
			Kind:        TaskPostprocessor,
			Description: "Compresses the resulting sources into an archive.",
			Params: []taskDefParam{
				{
					Name:     "encoding",
					Kind:     "String",
					Required: true,
				},
				{
					Name:     "name",
					Kind:     "String",
					Required: false,
				},
			},
		},
	}
}
