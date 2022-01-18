package gig

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"

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
func PostprocessorNoopParams(params map[string]interface{}) (postprocessorNoop, error) {
	var (
		out = postprocessorNoop{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{}
	for p := range params {
		if !index[p] {
			return out, fmt.Errorf("unknown parameter provided to noop: %s", p)
		}
	}

	// Fill and check requirements
	return out, err
}

func (t postprocessorNoop) Ref() string {
	return PostprocessorHandleNoop
}

func (t postprocessorNoop) Params() map[string]interface{} {
	return nil
}

// PostprocessorDiscardParams returns a new postprocessorDiscard from the params
func PostprocessorDiscardParams(params map[string]interface{}) (postprocessorDiscard, error) {
	var (
		out = postprocessorDiscard{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{}
	for p := range params {
		if !index[p] {
			return out, fmt.Errorf("unknown parameter provided to discard: %s", p)
		}
	}

	// Fill and check requirements
	return out, err
}

func (t postprocessorDiscard) Ref() string {
	return PostprocessorHandleDiscard
}

func (t postprocessorDiscard) Params() map[string]interface{} {
	return nil
}

// PostprocessorSaveParams returns a new postprocessorSave from the params
func PostprocessorSaveParams(params map[string]interface{}) (postprocessorSave, error) {
	var (
		out = postprocessorSave{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{}
	for p := range params {
		if !index[p] {
			return out, fmt.Errorf("unknown parameter provided to save: %s", p)
		}
	}

	// Fill and check requirements
	return out, err
}

func (t postprocessorSave) Ref() string {
	return PostprocessorHandleSave
}

func (t postprocessorSave) Params() map[string]interface{} {
	return nil
}

// PostprocessorArchiveParams returns a new postprocessorArchive from the params
func PostprocessorArchiveParams(params map[string]interface{}) (postprocessorArchive, error) {
	var (
		out = postprocessorArchive{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{
		"encoding": true,
		"name":     true,
	}
	for p := range params {
		if !index[p] {
			return out, fmt.Errorf("unknown parameter provided to archive: %s", p)
		}
	}

	// Fill and check requirements
	if _, ok := params["encoding"]; !ok {
		return out, fmt.Errorf("required parameter not provided: encoding")
	}
	out.encoding = archiveFromParams(params["encoding"])
	out.name = cast.ToString(params["name"])
	out, err = postprocessorArchiveTransformer(out)
	return out, err
}

// PostprocessorArchiveName returns a new postprocessorArchive from the required fields and name
func PostprocessorArchiveName(encoding archive, name string) (postprocessorArchive, error) {
	var (
		err error
		out postprocessorArchive
	)
	out = postprocessorArchive{
		encoding: encoding,
		name:     name,
	}
	out, err = postprocessorArchiveTransformer(out)

	return out, err
}
func PostprocessorArchive(encoding archive) (postprocessorArchive, error) {
	var (
		err error
		out postprocessorArchive
	)

	out = postprocessorArchive{
		encoding: encoding,
	}
	out, err = postprocessorArchiveTransformer(out)

	return out, err
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
