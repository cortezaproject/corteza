package types

import (
	"database/sql/driver"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/spf13/cast"
)

type (
	// Modules - CRM module definitions
	ModuleField struct {
		ID          uint64 `json:"fieldID,string"`
        NamespaceID uint64 `json:"namespaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`
		Place       int    `json:"-"`

		Kind string `json:"kind"`
		Name string `json:"name"`

		// Options relevant to field type
		Options ModuleFieldOptions `json:"options"`

		// Configuration - how sub-services and sub-systems like DAL and record revisions
		// are configured to work with this field
		Config ModuleFieldConfig `json:"config"`

		Required     bool           `json:"isRequired"`
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

	ModuleFieldConfig struct {
		DAL     ModuleFieldConfigDAL         `json:"dal"`
		Privacy ModuleFieldConfigDataPrivacy `json:"privacy"`

		RecordRevisions ModuleFieldConfigRecordRevisions `json:"recordRevisions"`
	}

	// ModuleFieldConfigDAL holds DAL configuration for a specific field
	//
	// If strategy is not set for a specific field (nil)
	// then a default strategy is used
	ModuleFieldConfigDAL struct {
		EncodingStrategy *EncodingStrategy `json:"encodingStrategy"`
	}

	ModuleFieldConfigDataPrivacy struct {
		// Define the highest sensitivity level which
		// can be configured on the module fields
		SensitivityLevelID uint64 `json:"sensitivityLevelID,string,omitempty"`

		UsageDisclosure string `json:"usageDisclosure"`
	}

	ModuleFieldConfigRecordRevisions struct {
		// when true, skip record revisions for this field
		Skip bool `json:"enabled"`
	}

	// SystemFieldEncoding holds configuration for encoding record system fields
	//
	// If strategy is not set for a specific field (nil)
	// then a default strategy is used, assuming system field/column presence
	SystemFieldEncoding struct {
		ID *EncodingStrategy `json:"id"`

		ModuleID    *EncodingStrategy `json:"moduleID"`
		NamespaceID *EncodingStrategy `json:"namespaceID"`

		Revision *EncodingStrategy `json:"revision"`
		Meta     *EncodingStrategy `json:"meta"`

		OwnedBy *EncodingStrategy `json:"ownedBy"`

		CreatedAt *EncodingStrategy `json:"createdAt"`
		CreatedBy *EncodingStrategy `json:"createdBy"`

		UpdatedAt *EncodingStrategy `json:"updatedAt"`
		UpdatedBy *EncodingStrategy `json:"updatedBy"`

		DeletedAt *EncodingStrategy `json:"deletedAt"`
		DeletedBy *EncodingStrategy `json:"deletedBy"`
	}

	// EncodingStrategy is used by both: Module (for system fields) and ModuleField
	//
	EncodingStrategy struct {
		//Type       string         `json:"type"`
		//TypeParams map[string]any `json:"typeParams"`

		Omit bool `json:"omit,omitempty"`

		*EncodingStrategyAlias `json:"alias,omitempty"`
		*EncodingStrategyJSON  `json:"json,omitempty"`
		*EncodingStrategyPlain `json:"plain,omitempty"`
	}

	EncodingStrategyAlias struct {
		Ident string `json:"ident"`
	}

	EncodingStrategyJSON struct {
		Ident string `json:"ident"`
	}

	EncodingStrategyPlain struct{}

	ModuleFieldFilter struct {
		ModuleID []uint64
		Deleted  filter.State
		Limit    uint
	}
)

var (
	_ sort.Interface = &ModuleFieldSet{}
)

func (f *ModuleField) SelectOptions() (out []string) {
	if f.Kind != "Select" {
		return
	}

	var (
		options, has = f.Options["options"]
	)

	if !has {
		return
	}

	switch oo := options.(type) {
	case []string:
		out = oo
	case []interface{}:
		for _, o := range oo {
			switch c := o.(type) {
			case string:
				out = append(out, c)
			case map[string]string:
				if value, has := c["value"]; has {
					out = append(out, value)
				}
			case map[string]interface{}:
				if value, has := c["value"]; has {
					if value, ok := value.(string); ok {
						out = append(out, value)
					}
				}
			case ModuleFieldOptions:
				if value, has := c["value"]; has {
					if value, ok := value.(string); ok {
						out = append(out, value)
					}
				}
			}
		}
	}

	return
}

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
		f.setOptionKey(aux.Msg, "description", "view")
	}
}

func (f *ModuleField) decodeTranslationsMetaDescriptionEdit(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldMetaDescriptionEdit.Path); aux != nil {
		f.setOptionKey(aux.Msg, "description", "edit")
	}
}

func (f *ModuleField) decodeTranslationsMetaHintView(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldMetaHintView.Path); aux != nil {
		f.setOptionKey(aux.Msg, "hint", "view")
	}
}

func (f *ModuleField) decodeTranslationsMetaHintEdit(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	if aux = tt.FindByKey(LocaleKeyModuleFieldMetaHintEdit.Path); aux != nil {
		f.setOptionKey(aux.Msg, "hint", "edit")
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
		outOpt := map[string]any{
			"value": "",
			"text":  "",
			"style": map[string]any{
				"textColor":       "",
				"backgroundColor": "",
			},
		}

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

			if style, ok := opt["style"].(map[string]any); ok {
				outOpt["style"] = style
			}
		}

		// find the translation for that value
		// and update the option (effectively overwriting
		// the original text value (in case of map option)
		trKey := strings.NewReplacer("{{value}}", outOpt["value"].(string)).Replace(LocaleKeyModuleFieldMetaOptionsValueText.Path)
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
	t := &locale.ResourceTranslation{
		Resource: f.ResourceTranslation(),
		Key:      LocaleKeyModuleFieldMetaDescriptionView.Path,
	}
	if v := f.getOptionKey("description", "view"); v != nil {
		t.Msg = cast.ToString(v)
	}
	out = append(out, t)
	return out
}

func (f *ModuleField) encodeTranslationsMetaDescriptionEdit() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	t := &locale.ResourceTranslation{
		Resource: f.ResourceTranslation(),
		Key:      LocaleKeyModuleFieldMetaDescriptionEdit.Path,
	}
	if v := f.getOptionKey("description", "edit"); v != nil {
		t.Msg = cast.ToString(v)
	}
	out = append(out, t)
	return out
}

func (f *ModuleField) encodeTranslationsMetaHintView() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	t := &locale.ResourceTranslation{
		Resource: f.ResourceTranslation(),
		Key:      LocaleKeyModuleFieldMetaHintView.Path,
	}
	if v := f.getOptionKey("hint", "view"); v != nil {
		t.Msg = cast.ToString(v)
	}
	out = append(out, t)
	return out
}

func (f *ModuleField) encodeTranslationsMetaHintEdit() (out locale.ResourceTranslationSet) {
	out = locale.ResourceTranslationSet{}
	t := &locale.ResourceTranslation{
		Resource: f.ResourceTranslation(),
		Key:      LocaleKeyModuleFieldMetaHintEdit.Path,
	}
	if v := f.getOptionKey("hint", "edit"); v != nil {
		t.Msg = cast.ToString(v)
	}
	out = append(out, t)
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
	if len(f.Options) == 0 {
		return nil
	}

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

func (set *ModuleFieldSet) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set ModuleFieldSet) Value() (driver.Value, error) { return json.Marshal(set) }

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

func (f ModuleField) IsSensitive() bool {
	return f.Config.Privacy.SensitivityLevelID > 0
}

func (f *ModuleField) setValue(name string, pos uint, value any) (err error) {
	switch name {
	// @todo consider moving this to the .cue definition; figure out why it wasn't yet
	case "NamespaceID", "namespaceID":
		f.NamespaceID = cast.ToUint64(value)
	case "Options.ModuleID":
		f.Options["moduleID"] = cast.ToString(value)
	}

	return
}

func (p *ModuleFieldConfig) Scan(src any) error          { return sql.ParseJSON(src, p) }
func (p ModuleFieldConfig) Value() (driver.Value, error) { return json.Marshal(p) }
