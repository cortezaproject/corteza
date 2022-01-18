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
	preprocessorNoop             struct{}
	preprocessorAttachmentRemove struct {
		mimeType string
	}
	preprocessorAttachmentTransform struct {
		height int
		width  int
	}
	preprocessorExperimentalExport struct {
		exclLanguage     []string
		exclRoles        []string
		handle           string
		id               uint64
		inclLanguage     []string
		inclRBAC         bool
		inclRoles        []string
		inclTranslations bool
	}
)

const (
	PreprocessorHandleNoop                = "noop"
	PreprocessorHandleAttachmentRemove    = "attachmentRemove"
	PreprocessorHandleAttachmentTransform = "attachmentTransform"
	PreprocessorHandleExperimentalExport  = "experimentalExport"
)

// ------------------------------------------------------------------------
// Constructors and utils

// PreprocessorNoopParams returns a new preprocessorNoop from the params
func PreprocessorNoopParams(params map[string]interface{}) (Preprocessor, error) {
	var (
		out = preprocessorNoop{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{}
	for p := range params {
		if !index[p] {
			return nil, fmt.Errorf("unknown parameter provided to noop: %s", p)
		}
	}

	// Fill and check requirements
	return out, err
}

func (t preprocessorNoop) Ref() string {
	return PreprocessorHandleNoop
}

func (t preprocessorNoop) Params() map[string]interface{} {
	return nil
}

// PreprocessorAttachmentRemoveParams returns a new preprocessorAttachmentRemove from the params
func PreprocessorAttachmentRemoveParams(params map[string]interface{}) (Preprocessor, error) {
	var (
		out = preprocessorAttachmentRemove{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{
		"mimeType": true,
	}
	for p := range params {
		if !index[p] {
			return nil, fmt.Errorf("unknown parameter provided to attachmentRemove: %s", p)
		}
	}

	// Fill and check requirements
	if _, ok := params["mimeType"]; !ok {
		return nil, fmt.Errorf("required parameter not provided: mimeType")
	}
	out.mimeType = cast.ToString(params["mimeType"])
	return out, err
}

func PreprocessorAttachmentRemove(mimeType string) (Preprocessor, error) {
	var (
		err error
		out preprocessorAttachmentRemove
	)

	out = preprocessorAttachmentRemove{
		mimeType: mimeType,
	}

	return out, err
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
func PreprocessorAttachmentTransformParams(params map[string]interface{}) (Preprocessor, error) {
	var (
		out = preprocessorAttachmentTransform{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{
		"height": true,
		"width":  true,
	}
	for p := range params {
		if !index[p] {
			return nil, fmt.Errorf("unknown parameter provided to attachmentTransform: %s", p)
		}
	}

	// Fill and check requirements
	out.height = cast.ToInt(params["height"])
	out.width = cast.ToInt(params["width"])
	return out, err
}

// PreprocessorAttachmentTransformHeight returns a new preprocessorAttachmentTransform from the required fields and height
func PreprocessorAttachmentTransformHeight(height int) (Preprocessor, error) {
	var (
		err error
		out preprocessorAttachmentTransform
	)
	out = preprocessorAttachmentTransform{
		height: height,
	}

	return out, err
}

// PreprocessorAttachmentTransformWidth returns a new preprocessorAttachmentTransform from the required fields and width
func PreprocessorAttachmentTransformWidth(width int) (Preprocessor, error) {
	var (
		err error
		out preprocessorAttachmentTransform
	)
	out = preprocessorAttachmentTransform{
		width: width,
	}

	return out, err
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

// PreprocessorExperimentalExportParams returns a new preprocessorExperimentalExport from the params
func PreprocessorExperimentalExportParams(params map[string]interface{}) (Preprocessor, error) {
	var (
		out = preprocessorExperimentalExport{}
		err error
	)

	// Param validation
	// - supported params
	index := map[string]bool{
		"exclLanguage":     true,
		"exclRoles":        true,
		"handle":           true,
		"id":               true,
		"inclLanguage":     true,
		"inclRBAC":         true,
		"inclRoles":        true,
		"inclTranslations": true,
	}
	for p := range params {
		if !index[p] {
			return nil, fmt.Errorf("unknown parameter provided to experimentalExport: %s", p)
		}
	}

	// Fill and check requirements
	out.exclLanguage = cast.ToStringSlice(params["exclLanguage"])
	out.exclRoles = cast.ToStringSlice(params["exclRoles"])
	out.handle = cast.ToString(params["handle"])
	out.id = cast.ToUint64(params["id"])
	out.inclLanguage = cast.ToStringSlice(params["inclLanguage"])
	out.inclRBAC = cast.ToBool(params["inclRBAC"])
	out.inclRoles = cast.ToStringSlice(params["inclRoles"])
	out.inclTranslations = cast.ToBool(params["inclTranslations"])
	return out, err
}

// PreprocessorExperimentalExportExclLanguage returns a new preprocessorExperimentalExport from the required fields and exclLanguage
func PreprocessorExperimentalExportExclLanguage(exclLanguage []string) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		exclLanguage: exclLanguage,
	}

	return out, err
}

// PreprocessorExperimentalExportExclRoles returns a new preprocessorExperimentalExport from the required fields and exclRoles
func PreprocessorExperimentalExportExclRoles(exclRoles []string) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		exclRoles: exclRoles,
	}

	return out, err
}

// PreprocessorExperimentalExportHandle returns a new preprocessorExperimentalExport from the required fields and handle
func PreprocessorExperimentalExportHandle(handle string) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		handle: handle,
	}

	return out, err
}

// PreprocessorExperimentalExportId returns a new preprocessorExperimentalExport from the required fields and id
func PreprocessorExperimentalExportId(id uint64) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		id: id,
	}

	return out, err
}

// PreprocessorExperimentalExportInclLanguage returns a new preprocessorExperimentalExport from the required fields and inclLanguage
func PreprocessorExperimentalExportInclLanguage(inclLanguage []string) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		inclLanguage: inclLanguage,
	}

	return out, err
}

// PreprocessorExperimentalExportInclRBAC returns a new preprocessorExperimentalExport from the required fields and inclRBAC
func PreprocessorExperimentalExportInclRBAC(inclRBAC bool) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		inclRBAC: inclRBAC,
	}

	return out, err
}

// PreprocessorExperimentalExportInclRoles returns a new preprocessorExperimentalExport from the required fields and inclRoles
func PreprocessorExperimentalExportInclRoles(inclRoles []string) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		inclRoles: inclRoles,
	}

	return out, err
}

// PreprocessorExperimentalExportInclTranslations returns a new preprocessorExperimentalExport from the required fields and inclTranslations
func PreprocessorExperimentalExportInclTranslations(inclTranslations bool) (Preprocessor, error) {
	var (
		err error
		out preprocessorExperimentalExport
	)
	out = preprocessorExperimentalExport{
		inclTranslations: inclTranslations,
	}

	return out, err
}

func (t preprocessorExperimentalExport) Ref() string {
	return PreprocessorHandleExperimentalExport
}

func (t preprocessorExperimentalExport) Params() map[string]interface{} {
	return map[string]interface{}{
		"exclLanguage":     t.exclLanguage,
		"exclRoles":        t.exclRoles,
		"handle":           t.handle,
		"id":               t.id,
		"inclLanguage":     t.inclLanguage,
		"inclRBAC":         t.inclRBAC,
		"inclRoles":        t.inclRoles,
		"inclTranslations": t.inclTranslations,
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
			Ref:         PreprocessorHandleExperimentalExport,
			Kind:        TaskPreprocessor,
			Description: "Loads the namespace along with some sub-resources (modules, pages, charts, ...)",
			Params: []taskDefParam{
				{
					Name:     "exclLanguage",
					Kind:     "[]string",
					Required: false,
				},
				{
					Name:     "exclRoles",
					Kind:     "[]string",
					Required: false,
				},
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
					Name:     "inclLanguage",
					Kind:     "[]string",
					Required: false,
				},
				{
					Name:     "inclRBAC",
					Kind:     "Bool",
					Required: false,
				},
				{
					Name:     "inclRoles",
					Kind:     "[]string",
					Required: false,
				},
				{
					Name:     "inclTranslations",
					Kind:     "Bool",
					Required: false,
				},
			},
		},
	}
}
