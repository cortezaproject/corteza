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

// PreprocessorExperimentalExportParams returns a new preprocessorExperimentalExport from the params
func PreprocessorExperimentalExportParams(params map[string]interface{}) Preprocessor {
	out := preprocessorExperimentalExport{
		exclLanguage: cast.ToStringSlice(params["exclLanguage"]),

		exclRoles: cast.ToStringSlice(params["exclRoles"]),

		handle: cast.ToString(params["handle"]),

		id: cast.ToUint64(params["id"]),

		inclLanguage: cast.ToStringSlice(params["inclLanguage"]),

		inclRBAC: cast.ToBool(params["inclRBAC"]),

		inclRoles: cast.ToStringSlice(params["inclRoles"]),

		inclTranslations: cast.ToBool(params["inclTranslations"]),
	}
	return out
}

// PreprocessorExperimentalExportExclLanguage returns a new preprocessorExperimentalExport from the required fields and exclLanguage
func PreprocessorExperimentalExportExclLanguage(exclLanguage []string) Preprocessor {
	out := preprocessorExperimentalExport{
		exclLanguage: exclLanguage,
	}

	return out
}

// PreprocessorExperimentalExportExclRoles returns a new preprocessorExperimentalExport from the required fields and exclRoles
func PreprocessorExperimentalExportExclRoles(exclRoles []string) Preprocessor {
	out := preprocessorExperimentalExport{
		exclRoles: exclRoles,
	}

	return out
}

// PreprocessorExperimentalExportHandle returns a new preprocessorExperimentalExport from the required fields and handle
func PreprocessorExperimentalExportHandle(handle string) Preprocessor {
	out := preprocessorExperimentalExport{
		handle: handle,
	}

	return out
}

// PreprocessorExperimentalExportId returns a new preprocessorExperimentalExport from the required fields and id
func PreprocessorExperimentalExportId(id uint64) Preprocessor {
	out := preprocessorExperimentalExport{
		id: id,
	}

	return out
}

// PreprocessorExperimentalExportInclLanguage returns a new preprocessorExperimentalExport from the required fields and inclLanguage
func PreprocessorExperimentalExportInclLanguage(inclLanguage []string) Preprocessor {
	out := preprocessorExperimentalExport{
		inclLanguage: inclLanguage,
	}

	return out
}

// PreprocessorExperimentalExportInclRBAC returns a new preprocessorExperimentalExport from the required fields and inclRBAC
func PreprocessorExperimentalExportInclRBAC(inclRBAC bool) Preprocessor {
	out := preprocessorExperimentalExport{
		inclRBAC: inclRBAC,
	}

	return out
}

// PreprocessorExperimentalExportInclRoles returns a new preprocessorExperimentalExport from the required fields and inclRoles
func PreprocessorExperimentalExportInclRoles(inclRoles []string) Preprocessor {
	out := preprocessorExperimentalExport{
		inclRoles: inclRoles,
	}

	return out
}

// PreprocessorExperimentalExportInclTranslations returns a new preprocessorExperimentalExport from the required fields and inclTranslations
func PreprocessorExperimentalExportInclTranslations(inclTranslations bool) Preprocessor {
	out := preprocessorExperimentalExport{
		inclTranslations: inclTranslations,
	}

	return out
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
