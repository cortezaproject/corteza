package types

import (
	"database/sql/driver"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/spf13/cast"
)

type (
	// Modules - CRM module definitions
	ModuleField struct {
		ID          uint64 `json:"fieldID,string"`
		NamespaceID uint64 `json:"namspaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`
		Place       int    `json:"-"`

		Kind string `json:"kind"`
		Name string `json:"name"`

		Options ModuleFieldOptions `json:"options"`

		Private      bool           `json:"isPrivate"`
		Required     bool           `json:"isRequired"`
		Visible      bool           `json:"isVisible"`
		Multi        bool           `json:"isMulti"`
		DefaultValue RecordValueSet `json:"defaultValue"`

		Expressions ModuleFieldExpr `json:"expressions"`

		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Label string `json:"label"`
	}

	ModuleFieldFilter struct {
		ModuleID []uint64
		Deleted  filter.State
	}
)

var (
	_ sort.Interface = &ModuleFieldSet{}
)

func (f *ModuleField) decodeTranslationsExpressionValidatorValidatorIDError(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	for i, e := range f.Expressions.Validators {
		validatorID := locale.ContentID(e.ValidatorID, i)
		rpl := strings.NewReplacer(
			"{{validatorID}}", strconv.FormatUint(validatorID, 10),
		)

		if aux = tt.FindByKey(rpl.Replace(LocaleKeyModuleFieldExpressionValidatorValidatorIDError.Path)); aux != nil {
			f.Expressions.Validators[i].Error = aux.Msg
		}
	}
}

func (f *ModuleField) decodeTranslationsMetaDescriptionView(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldMetaDescriptionView.Path); aux != nil {
		f.setOptionKey(aux.Msg, "description", "edit")
	}
}

func (f *ModuleField) decodeTranslationsMetaDescriptionEdit(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldMetaDescriptionEdit.Path); aux != nil {
		f.setOptionKey(aux.Msg, "description", "view")
	}
}

func (f *ModuleField) decodeTranslationsMetaHintView(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldMetaHintView.Path); aux != nil {
		f.setOptionKey(aux.Msg, "hint", "edit")
	}
}

func (f *ModuleField) decodeTranslationsMetaHintEdit(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldMetaHintEdit.Path); aux != nil {
		f.setOptionKey(aux.Msg, "hint", "view")
	}
}

// Decodes translations and modifies options
//
// Why "options-option-texts"? Because we're translating option txts under options key-value
func (f *ModuleField) decodeTranslationsMetaOptionsValueText(tt locale.ResourceTranslationIndex) {
	var (
		tr *locale.ResourceTranslation
	)

	optsUnknown, has := f.Options["options"]
	if !has {
		return
	}

	optsSlice, is := optsUnknown.([]interface{})
	if !is {
		return
	}

	for i, optUnknown := range optsSlice {
		outOpt := map[string]string{}

		// what is this we're dealing with? slice of strings (values) or map (value+text)
		switch opt := optUnknown.(type) {
		case string:
			// cast string (value) into map (value+text)
			// and use value as text (as a fallback in case
			// of missing translation)
			outOpt["value"] = opt
			outOpt["text"] = opt

		case map[string]interface{}:
			outOpt["value"], is = opt["value"].(string)
			if !is {
				continue
			}

			outOpt["text"], _ = opt["text"].(string)
		}

		// find the translation for that value
		// and update the option (effectively overwriting
		// the original text value (in case of map option)
		trKey := strings.NewReplacer("{{value}}", outOpt["value"]).Replace(LocaleKeyModuleFieldMetaOptionsValueText.Path)
		if tr = tt.FindByKey(trKey); tr != nil {
			outOpt["text"] = tr.Msg
		}

		// Update slice item with translated option
		optsSlice[i] = outOpt
	}
}

// Decodes translations and modifies options
func (f *ModuleField) decodeTranslationsMetaBoolValueLabel(tt locale.ResourceTranslationIndex) {
	if f.Kind != "Bool" {
		return
	}

	var aux *locale.ResourceTranslation
	if aux = tt.FindByKey(strings.NewReplacer("{{value}}", "true").Replace(LocaleKeyModuleFieldMetaBoolValueLabel.Path)); aux != nil {
		f.setOptionKey(aux.Msg, "trueLabel")
	}
	if aux = tt.FindByKey(strings.NewReplacer("{{value}}", "false").Replace(LocaleKeyModuleFieldMetaBoolValueLabel.Path)); aux != nil {
		f.setOptionKey(aux.Msg, "falseLabel")
	}
}

func (m *ModuleField) encodeTranslationsExpressionValidatorValidatorIDError() (out locale.ResourceTranslationSet) {
	out = make(locale.ResourceTranslationSet, 0, 3)

	// Module field expressions
	for i, e := range m.Expressions.Validators {
		validatorID := locale.ContentID(e.ValidatorID, i)
		rpl := strings.NewReplacer(
			"{{validatorID}}", strconv.FormatUint(validatorID, 10),
		)

		out = append(out, &locale.ResourceTranslation{
			Resource: m.ResourceTranslation(),
			Key:      rpl.Replace(LocaleKeyModuleFieldExpressionValidatorValidatorIDError.Path),
			Msg:      e.Error,
		})
	}

	return
}

func (f *ModuleField) encodeTranslationsMetaDescriptionView() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	if v := f.getOptionKey("description", "edit"); v != nil {
		aux := cast.ToString(v)

		out = append(out, &locale.ResourceTranslation{
			Resource: f.ResourceTranslation(),
			Key:      LocaleKeyModuleFieldMetaDescriptionView.Path,
			Msg:      aux,
		})
	}
	return out
}

func (f *ModuleField) encodeTranslationsMetaDescriptionEdit() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	if v := f.getOptionKey("description", "view"); v != nil {
		aux := cast.ToString(v)
		out = append(out, &locale.ResourceTranslation{
			Resource: f.ResourceTranslation(),
			Key:      LocaleKeyModuleFieldMetaDescriptionEdit.Path,
			Msg:      aux,
		})
	}
	return out
}

func (f *ModuleField) encodeTranslationsMetaHintView() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	if v := f.getOptionKey("hint", "edit"); v != nil {
		aux := cast.ToString(v)
		out = append(out, &locale.ResourceTranslation{
			Resource: f.ResourceTranslation(),
			Key:      LocaleKeyModuleFieldMetaHintView.Path,
			Msg:      aux,
		})
	}
	return out
}

func (f *ModuleField) encodeTranslationsMetaHintEdit() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	if v := f.getOptionKey("hint", "view"); v != nil {
		aux := cast.ToString(v)
		out = append(out, &locale.ResourceTranslation{
			Resource: f.ResourceTranslation(),
			Key:      LocaleKeyModuleFieldMetaHintEdit.Path,
			Msg:      aux,
		})
	}
	return out
}

// extracts option texts and converts (encodes) them to translations
func (f *ModuleField) encodeTranslationsMetaOptionsValueText() (out locale.ResourceTranslationSet) {
	out = make(locale.ResourceTranslationSet, 0, 3)

	optsUnknown, has := f.Options["options"]
	if !has {
		return
	}

	optsSlice, is := optsUnknown.([]interface{})
	if !is {
		return
	}

	for _, optUnknown := range optsSlice {
		// we only care about maps (items with value & text)
		switch opt := optUnknown.(type) {
		case map[string]interface{}:
			value, _ := opt["value"].(string)
			text, _ := opt["text"].(string)

			out = append(out, &locale.ResourceTranslation{
				Resource: f.ResourceTranslation(),
				Key: strings.NewReplacer("{{value}}", value).
					Replace(LocaleKeyModuleFieldMetaOptionsValueText.Path),
				Msg: text,
			})
		}
	}

	return
}

func (m *ModuleField) encodeTranslationsMetaBoolValueLabel() (out locale.ResourceTranslationSet) {
	if m.Kind != "Bool" {
		return
	}

	out = make(locale.ResourceTranslationSet, 0, 3)
	out = append(out, &locale.ResourceTranslation{
		Resource: m.ResourceTranslation(),
		Key:      strings.NewReplacer("{{value}}", "true").Replace(LocaleKeyModuleFieldMetaBoolValueLabel.Path),
		Msg:      m.Options.String("trueLabel"),
	})

	out = append(out, &locale.ResourceTranslation{
		Resource: m.ResourceTranslation(),
		Key:      strings.NewReplacer("{{value}}", "false").Replace(LocaleKeyModuleFieldMetaBoolValueLabel.Path),
		Msg:      m.Options.String("falseLabel"),
	})

	return
}

func (f ModuleField) Clone() *ModuleField {
	return &f
}

func (f ModuleField) setOptionKey(v interface{}, kk ...string) {
	opt := f.Options

	for _, k := range kk[0 : len(kk)-1] {
		_, ok := opt[k]
		if !ok {
			opt = map[string]interface{}{k: make(map[string]interface{})}
		}
		var aux ModuleFieldOptions
		switch c := opt[k].(type) {
		case map[string]interface{}:
			aux = ModuleFieldOptions(c)
		case ModuleFieldOptions:
			aux = c
		}

		opt = aux
	}

	k := kk[len(kk)-1]
	opt[k] = v
}

func (f ModuleField) getOptionKey(kk ...string) interface{} {
	opt := f.Options

	for _, k := range kk[0 : len(kk)-1] {
		_, ok := opt[k]
		if !ok {
			opt = map[string]interface{}{k: make(map[string]interface{})}
		}

		var aux ModuleFieldOptions
		switch c := opt[k].(type) {
		case map[string]interface{}:
			aux = ModuleFieldOptions(c)
		case ModuleFieldOptions:
			aux = c
		}

		opt = aux
	}

	return opt[kk[len(kk)-1]]
}

func (set ModuleFieldSet) Clone() (out ModuleFieldSet) {
	out = make([]*ModuleField, len(set))
	for i := range set {
		out[i] = set[i].Clone()
	}

	return out
}

func (set *ModuleFieldSet) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, set)
	}
	return nil
}

func (set ModuleFieldSet) Value() (driver.Value, error) {
	return json.Marshal(set)
}

func (set ModuleFieldSet) Names() (names []string) {
	names = make([]string, len(set))

	for i := range set {
		names[i] = set[i].Name
	}

	return
}

func (set ModuleFieldSet) HasName(name string) bool {
	for i := range set {
		if name == set[i].Name {
			return true
		}
	}

	return false
}

func (set ModuleFieldSet) FindByName(name string) *ModuleField {
	for i := range set {
		if name == set[i].Name {
			return set[i]
		}
	}

	return nil
}

func (set ModuleFieldSet) FilterByModule(moduleID uint64) (ff ModuleFieldSet) {
	for i := range set {
		if set[i].ModuleID == moduleID {
			ff = append(ff, set[i])
		}
	}

	return
}

func (set ModuleFieldSet) Len() int {
	return len(set)
}

func (set ModuleFieldSet) Less(i, j int) bool {
	return set[i].Place < set[j].Place
}

func (set ModuleFieldSet) Swap(i, j int) {
	set[i], set[j] = set[j], set[i]
}

func (f ModuleField) IsBoolean() bool {
	return f.Kind == "Bool"
}

func (f ModuleField) IsNumeric() bool {
	return f.Kind == "Number"
}

func (f ModuleField) IsDateTime() bool {
	return f.Kind == "DateTime" && !f.IsDateOnly() && !f.IsTimeOnly()
}

func (f ModuleField) IsDateOnly() bool {
	return f.Kind == "DateTime" && f.Options.Bool("onlyDate")
}

func (f ModuleField) IsTimeOnly() bool {
	return f.Kind == "DateTime" && f.Options.Bool("onlyTime")
}

// IsRef tells us if value of this field be a reference to something
// (another record, file , user)?
func (f ModuleField) IsRef() bool {
	return f.Kind == "Record" || f.Kind == "User" || f.Kind == "File"
}
