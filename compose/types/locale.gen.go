package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - compose.module-field.yaml
// - compose.module.yaml
// - compose.namespace.yaml
// - compose.page.yaml

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
	ModuleFieldResourceTranslationType = "compose:module-field"
	ModuleResourceTranslationType      = "compose:module"
	NamespaceResourceTranslationType   = "compose:namespace"
	PageResourceTranslationType        = "compose:page"
)

var (
	LocaleKeyModuleFieldLabel = LocaleKey{
		Name:     "label",
		Resource: ModuleFieldResourceTranslationType,
		Path:     "label",
	}
	LocaleKeyModuleFieldValidatorError = LocaleKey{
		Name:          "validatorError",
		Resource:      ModuleFieldResourceTranslationType,
		Path:          "expression.validator.{{validatorID}}.error",
		CustomHandler: "validatorError",
	}
	LocaleKeyModuleName = LocaleKey{
		Name:     "name",
		Resource: ModuleResourceTranslationType,
		Path:     "name",
	}
	LocaleKeyNamespaceName = LocaleKey{
		Name:     "name",
		Resource: NamespaceResourceTranslationType,
		Path:     "name",
	}
	LocaleKeyNamespaceSubtitle = LocaleKey{
		Name:     "subtitle",
		Resource: NamespaceResourceTranslationType,
		Path:     "subtitle",
	}
	LocaleKeyNamespaceDescription = LocaleKey{
		Name:     "description",
		Resource: NamespaceResourceTranslationType,
		Path:     "description",
	}
	LocaleKeyPageTitle = LocaleKey{
		Name:     "title",
		Resource: PageResourceTranslationType,
		Path:     "title",
	}
	LocaleKeyPageDescription = LocaleKey{
		Name:     "description",
		Resource: PageResourceTranslationType,
		Path:     "description",
	}
	LocaleKeyPageBlockTitle = LocaleKey{
		Name:     "blockTitle",
		Resource: PageResourceTranslationType,
		Path:     "pageBlock.{{blockID}}.title",
	}
	LocaleKeyPageBlockDescription = LocaleKey{
		Name:     "blockDescription",
		Resource: PageResourceTranslationType,
		Path:     "pageBlock.{{blockID}}.description",
	}
	LocaleKeyPageBlockAutomationButtonlabel = LocaleKey{
		Name:     "blockAutomationButtonlabel",
		Resource: PageResourceTranslationType,
		Path:     "pageBlock.{{blockID}}.button.{{buttonID}}.label",
	}
)

// ResourceTranslation returns string representation of Locale resource for ModuleField by calling ModuleFieldResourceTranslation fn
//
// Locale resource is in the compose:module-field/... format
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
func ModuleFieldResourceTranslation(namespaceID uint64, moduleID uint64, id uint64) string {
	cpts := []interface{}{ModuleFieldResourceTranslationType}
	cpts = append(cpts, strconv.FormatUint(namespaceID, 10), strconv.FormatUint(moduleID, 10), strconv.FormatUint(id, 10))

	return fmt.Sprintf(ModuleFieldResourceTranslationTpl(), cpts...)
}

// @todo template
func ModuleFieldResourceTranslationTpl() string {
	return "%s/%s/%s/%s"
}

func (r *ModuleField) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation
	if aux = tt.FindByKey(LocaleKeyModuleFieldLabel.Path); aux != nil {
		r.Label = aux.Msg
	} else {
		r.Label = LocaleKeyModuleFieldLabel.Path
	}
	r.decodeTranslationsValidatorError(tt)
}

func (r *ModuleField) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyModuleFieldLabel.Path,
			Msg:      r.Label,
		},
	}

	out = append(out, r.encodeTranslationsValidatorError()...)

	return out
}

// ResourceTranslation returns string representation of Locale resource for Module by calling ModuleResourceTranslation fn
//
// Locale resource is in the compose:module/... format
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
func ModuleResourceTranslation(namespaceID uint64, id uint64) string {
	cpts := []interface{}{ModuleResourceTranslationType}
	cpts = append(cpts, strconv.FormatUint(namespaceID, 10), strconv.FormatUint(id, 10))

	return fmt.Sprintf(ModuleResourceTranslationTpl(), cpts...)
}

// @todo template
func ModuleResourceTranslationTpl() string {
	return "%s/%s/%s"
}

func (r *Module) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation
	if aux = tt.FindByKey(LocaleKeyModuleName.Path); aux != nil {
		r.Name = aux.Msg
	} else {
		r.Name = LocaleKeyModuleName.Path
	}

	r.decodeTranslations(tt)
}

func (r *Module) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyModuleName.Path,
			Msg:      r.Name,
		},
	}

	out = append(out, r.encodeTranslations()...)

	return out
}

// ResourceTranslation returns string representation of Locale resource for Namespace by calling NamespaceResourceTranslation fn
//
// Locale resource is in the compose:namespace/... format
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
func NamespaceResourceTranslation(id uint64) string {
	cpts := []interface{}{NamespaceResourceTranslationType}
	cpts = append(cpts, strconv.FormatUint(id, 10))

	return fmt.Sprintf(NamespaceResourceTranslationTpl(), cpts...)
}

// @todo template
func NamespaceResourceTranslationTpl() string {
	return "%s/%s"
}

func (r *Namespace) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation
	if aux = tt.FindByKey(LocaleKeyNamespaceName.Path); aux != nil {
		r.Name = aux.Msg
	} else {
		r.Name = LocaleKeyNamespaceName.Path
	}
	if aux = tt.FindByKey(LocaleKeyNamespaceSubtitle.Path); aux != nil {
		r.Meta.Subtitle = aux.Msg
	} else {
		r.Meta.Subtitle = LocaleKeyNamespaceSubtitle.Path
	}
	if aux = tt.FindByKey(LocaleKeyNamespaceDescription.Path); aux != nil {
		r.Meta.Description = aux.Msg
	} else {
		r.Meta.Description = LocaleKeyNamespaceDescription.Path
	}
}

func (r *Namespace) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyNamespaceName.Path,
			Msg:      r.Name,
		},
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyNamespaceSubtitle.Path,
			Msg:      r.Meta.Subtitle,
		},
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyNamespaceDescription.Path,
			Msg:      r.Meta.Description,
		},
	}

	return out
}

// ResourceTranslation returns string representation of Locale resource for Page by calling PageResourceTranslation fn
//
// Locale resource is in the compose:page/... format
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
func PageResourceTranslation(namespaceID uint64, id uint64) string {
	cpts := []interface{}{PageResourceTranslationType}
	cpts = append(cpts, strconv.FormatUint(namespaceID, 10), strconv.FormatUint(id, 10))

	return fmt.Sprintf(PageResourceTranslationTpl(), cpts...)
}

// @todo template
func PageResourceTranslationTpl() string {
	return "%s/%s/%s"
}

func (r *Page) DecodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation
	if aux = tt.FindByKey(LocaleKeyPageTitle.Path); aux != nil {
		r.Title = aux.Msg
	} else {
		r.Title = LocaleKeyPageTitle.Path
	}
	if aux = tt.FindByKey(LocaleKeyPageDescription.Path); aux != nil {
		r.Description = aux.Msg
	} else {
		r.Description = LocaleKeyPageDescription.Path
	}

	r.decodeTranslations(tt)
}

func (r *Page) EncodeTranslations() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyPageTitle.Path,
			Msg:      r.Title,
		},
		{
			Resource: r.ResourceTranslation(),
			Key:      LocaleKeyPageDescription.Path,
			Msg:      r.Description,
		},
	}

	out = append(out, r.encodeTranslations()...)

	return out
}
