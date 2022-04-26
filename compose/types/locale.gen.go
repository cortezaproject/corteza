package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"strconv"
)

type (
	LocaleKey struct {
		Name          string
		Resource      string
		Path          string
		CustomHandler string
	}
)

// Types and stuff
const (
	ChartResourceTranslationType       = "compose:chart"
	ModuleResourceTranslationType      = "compose:module"
	ModuleFieldResourceTranslationType = "compose:module-field"
	NamespaceResourceTranslationType   = "compose:namespace"
	PageResourceTranslationType        = "compose:page"
)

var (
	// @todo can we remove LocaleKey struct for string constant?
	LocaleKeyChartYAxisLabel                                = LocaleKey{Path: "yAxis.label"}
	LocaleKeyChartMetricsMetricIDLabel                      = LocaleKey{Path: "metrics.{{metricID}}.label"}
	LocaleKeyModuleName                                     = LocaleKey{Path: "name"}
	LocaleKeyModuleFieldLabel                               = LocaleKey{Path: "label"}
	LocaleKeyModuleFieldMetaDescriptionView                 = LocaleKey{Path: "meta.description.view"}
	LocaleKeyModuleFieldMetaDescriptionEdit                 = LocaleKey{Path: "meta.description.edit"}
	LocaleKeyModuleFieldMetaHintView                        = LocaleKey{Path: "meta.hint.view"}
	LocaleKeyModuleFieldMetaHintEdit                        = LocaleKey{Path: "meta.hint.edit"}
	LocaleKeyModuleFieldExpressionValidatorValidatorIDError = LocaleKey{Path: "expression.validator.{{validatorID}}.error"}
	LocaleKeyModuleFieldMetaOptionsValueText                = LocaleKey{Path: "meta.options.{{value}}.text"}
	LocaleKeyModuleFieldMetaBoolValueLabel                  = LocaleKey{Path: "meta.bool.{{value}}.label"}
	LocaleKeyNamespaceName                                  = LocaleKey{Path: "name"}
	LocaleKeyNamespaceMetaSubtitle                          = LocaleKey{Path: "meta.subtitle"}
	LocaleKeyNamespaceMetaDescription                       = LocaleKey{Path: "meta.description"}
	LocaleKeyPageTitle                                      = LocaleKey{Path: "title"}
	LocaleKeyPageDescription                                = LocaleKey{Path: "description"}
	LocaleKeyPageRecordToolbarNewLabel                      = LocaleKey{Path: "recordToolbar.new.label"}
	LocaleKeyPageRecordToolbarEditLabel                     = LocaleKey{Path: "recordToolbar.edit.label"}
	LocaleKeyPageRecordToolbarSubmitLabel                   = LocaleKey{Path: "recordToolbar.submit.label"}
	LocaleKeyPageRecordToolbarDeleteLabel                   = LocaleKey{Path: "recordToolbar.delete.label"}
	LocaleKeyPageRecordToolbarCloneLabel                    = LocaleKey{Path: "recordToolbar.clone.label"}
	LocaleKeyPageRecordToolbarBackLabel                     = LocaleKey{Path: "recordToolbar.back.label"}
	LocaleKeyPagePageBlockBlockIDTitle                      = LocaleKey{Path: "pageBlock.{{blockID}}.title"}
	LocaleKeyPagePageBlockBlockIDDescription                = LocaleKey{Path: "pageBlock.{{blockID}}.description"}
	LocaleKeyPagePageBlockBlockIDButtonButtonIDLabel        = LocaleKey{Path: "pageBlock.{{blockID}}.button.{{buttonID}}.label"}
	LocaleKeyPagePageBlockBlockIDContentBody                = LocaleKey{Path: "pageBlock.{{blockID}}.content.body"}
)

// ResourceTranslation returns string representation of Locale resource for Chart by calling ChartResourceTranslation fn
//
// Locale resource is in "compose:chart/..." format
//
// This function is auto-generated
func (r Chart) ResourceTranslation() string {
	return ChartResourceTranslation(r.NamespaceID, r.ID)
}

// ChartResourceTranslation returns string representation of Locale resource for Chart
//
// Locale resource is in the compose:chart/... format
//
// This function is auto-generated
func ChartResourceTranslation(NamespaceID uint64, ID uint64) string {
	cpts := []interface{}{
		ChartResourceTranslationType,
		strconv.FormatUint(NamespaceID, 10),
		strconv.FormatUint(ID, 10),
	}

	return fmt.Sprintf(ChartResourceTranslationTpl(), cpts...)
}

func ChartResourceTranslationTpl() string {
	return "%s/%s/%s"
}

func (r *Chart) DecodeTranslations(tt locale.ResourceTranslationIndex) {

	r.decodeTranslations(tt)
}

func (r *Chart) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}

	out = append(out, r.encodeTranslations()...)

	return out
}

// ResourceTranslation returns string representation of Locale resource for Module by calling ModuleResourceTranslation fn
//
// Locale resource is in "compose:module/..." format
//
// This function is auto-generated
func (r Module) ResourceTranslation() string {
	return ModuleResourceTranslation(r.NamespaceID, r.ID)
}

// ModuleResourceTranslation returns string representation of Locale resource for Module
//
// Locale resource is in the compose:module/... format
//
// This function is auto-generated
func ModuleResourceTranslation(NamespaceID uint64, ID uint64) string {
	cpts := []interface{}{
		ModuleResourceTranslationType,
		strconv.FormatUint(NamespaceID, 10),
		strconv.FormatUint(ID, 10),
	}

	return fmt.Sprintf(ModuleResourceTranslationTpl(), cpts...)
}

func ModuleResourceTranslationTpl() string {
	return "%s/%s/%s"
}

func (r *Module) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleName.Path); aux != nil {
		r.Name = aux.Msg
	}
	r.decodeTranslations(tt)
}

func (r *Module) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}

	out = append(out, &locale.ResourceTranslation{
		Resource: r.ResourceTranslation(),
		Key:      LocaleKeyModuleName.Path,
		Msg:      locale.SanitizeMessage(r.Name),
	})
	out = append(out, r.encodeTranslations()...)

	return out
}

// ResourceTranslation returns string representation of Locale resource for ModuleField by calling ModuleFieldResourceTranslation fn
//
// Locale resource is in "compose:module-field/..." format
//
// This function is auto-generated
func (r ModuleField) ResourceTranslation() string {
	return ModuleFieldResourceTranslation(r.NamespaceID, r.ModuleID, r.ID)
}

// ModuleFieldResourceTranslation returns string representation of Locale resource for ModuleField
//
// Locale resource is in the compose:module-field/... format
//
// This function is auto-generated
func ModuleFieldResourceTranslation(NamespaceID uint64, ModuleID uint64, ID uint64) string {
	cpts := []interface{}{
		ModuleFieldResourceTranslationType,
		strconv.FormatUint(NamespaceID, 10),
		strconv.FormatUint(ModuleID, 10),
		strconv.FormatUint(ID, 10),
	}

	return fmt.Sprintf(ModuleFieldResourceTranslationTpl(), cpts...)
}

func ModuleFieldResourceTranslationTpl() string {
	return "%s/%s/%s/%s"
}

func (r *ModuleField) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldLabel.Path); aux != nil {
		r.Label = aux.Msg
	}

	r.decodeTranslationsMetaDescriptionView(tt)

	r.decodeTranslationsMetaDescriptionEdit(tt)

	r.decodeTranslationsMetaHintView(tt)

	r.decodeTranslationsMetaHintEdit(tt)

	r.decodeTranslationsExpressionValidatorValidatorIDError(tt)

	r.decodeTranslationsMetaOptionsValueText(tt)

	r.decodeTranslationsMetaBoolValueLabel(tt)

}

func (r *ModuleField) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}

	out = append(out, &locale.ResourceTranslation{
		Resource: r.ResourceTranslation(),
		Key:      LocaleKeyModuleFieldLabel.Path,
		Msg:      locale.SanitizeMessage(r.Label),
	})

	out = append(out, r.encodeTranslationsMetaDescriptionView()...)

	out = append(out, r.encodeTranslationsMetaDescriptionEdit()...)

	out = append(out, r.encodeTranslationsMetaHintView()...)

	out = append(out, r.encodeTranslationsMetaHintEdit()...)

	out = append(out, r.encodeTranslationsExpressionValidatorValidatorIDError()...)

	out = append(out, r.encodeTranslationsMetaOptionsValueText()...)

	out = append(out, r.encodeTranslationsMetaBoolValueLabel()...)

	return out
}

// ResourceTranslation returns string representation of Locale resource for Namespace by calling NamespaceResourceTranslation fn
//
// Locale resource is in "compose:namespace/..." format
//
// This function is auto-generated
func (r Namespace) ResourceTranslation() string {
	return NamespaceResourceTranslation(r.ID)
}

// NamespaceResourceTranslation returns string representation of Locale resource for Namespace
//
// Locale resource is in the compose:namespace/... format
//
// This function is auto-generated
func NamespaceResourceTranslation(ID uint64) string {
	cpts := []interface{}{
		NamespaceResourceTranslationType,
		strconv.FormatUint(ID, 10),
	}

	return fmt.Sprintf(NamespaceResourceTranslationTpl(), cpts...)
}

func NamespaceResourceTranslationTpl() string {
	return "%s/%s"
}

func (r *Namespace) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyNamespaceName.Path); aux != nil {
		r.Name = aux.Msg
	}

	if aux = tt.FindByKey(LocaleKeyNamespaceMetaSubtitle.Path); aux != nil {
		r.Meta.Subtitle = aux.Msg
	}

	if aux = tt.FindByKey(LocaleKeyNamespaceMetaDescription.Path); aux != nil {
		r.Meta.Description = aux.Msg
	}
}

func (r *Namespace) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}

	out = append(out, &locale.ResourceTranslation{
		Resource: r.ResourceTranslation(),
		Key:      LocaleKeyNamespaceName.Path,
		Msg:      locale.SanitizeMessage(r.Name),
	})

	out = append(out, &locale.ResourceTranslation{
		Resource: r.ResourceTranslation(),
		Key:      LocaleKeyNamespaceMetaSubtitle.Path,
		Msg:      locale.SanitizeMessage(r.Meta.Subtitle),
	})

	out = append(out, &locale.ResourceTranslation{
		Resource: r.ResourceTranslation(),
		Key:      LocaleKeyNamespaceMetaDescription.Path,
		Msg:      locale.SanitizeMessage(r.Meta.Description),
	})

	return out
}

// ResourceTranslation returns string representation of Locale resource for Page by calling PageResourceTranslation fn
//
// Locale resource is in "compose:page/..." format
//
// This function is auto-generated
func (r Page) ResourceTranslation() string {
	return PageResourceTranslation(r.NamespaceID, r.ID)
}

// PageResourceTranslation returns string representation of Locale resource for Page
//
// Locale resource is in the compose:page/... format
//
// This function is auto-generated
func PageResourceTranslation(NamespaceID uint64, ID uint64) string {
	cpts := []interface{}{
		PageResourceTranslationType,
		strconv.FormatUint(NamespaceID, 10),
		strconv.FormatUint(ID, 10),
	}

	return fmt.Sprintf(PageResourceTranslationTpl(), cpts...)
}

func PageResourceTranslationTpl() string {
	return "%s/%s/%s"
}

func (r *Page) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyPageTitle.Path); aux != nil {
		r.Title = aux.Msg
	}

	if aux = tt.FindByKey(LocaleKeyPageDescription.Path); aux != nil {
		r.Description = aux.Msg
	}

	r.decodeTranslations(tt)
}

func (r *Page) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}

	out = append(out, &locale.ResourceTranslation{
		Resource: r.ResourceTranslation(),
		Key:      LocaleKeyPageTitle.Path,
		Msg:      locale.SanitizeMessage(r.Title),
	})

	out = append(out, &locale.ResourceTranslation{
		Resource: r.ResourceTranslation(),
		Key:      LocaleKeyPageDescription.Path,
		Msg:      locale.SanitizeMessage(r.Description),
	})

	out = append(out, r.encodeTranslations()...)

	return out
}
