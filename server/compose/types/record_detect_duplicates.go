package types

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/pkg/str"
	"github.com/spf13/cast"
	"strings"
)

type (
	deDup struct {
		ls localeService
	}

	localeService interface {
		T(ctx context.Context, ns, key string, rr ...string) string
	}

	DeDupRule struct {
		Name          DeDupRuleName          `json:"name"`
		Strict        bool                   `json:"strict"`
		ErrorMessage  string                 `json:"errorMessage"`
		ConstraintSet DeDupRuleConstraintSet `json:"constraints"`
	}

	DeDupRuleConstraint struct {
		Attribute  string                    `json:"attribute"`
		Modifier   DeDupValueModifier        `json:"modifier"`
		MultiValue DeDupMultiValueConstraint `json:"multiValue"`
	}

	DeDupRuleConstraintSet []*DeDupRuleConstraint

	// DeDupRuleName represent the identifier for duplicate detection rule
	DeDupRuleName string

	// DeDupValueModifier represent the algorithm used to check value string
	DeDupValueModifier string

	// DeDupMultiValueConstraint for matching multi values accordingly
	DeDupMultiValueConstraint string

	// DeDupIssueKind based on strict mode rule or duplication config
	DeDupIssueKind string
)

const (
	ignoreCase    DeDupValueModifier = "ignore-case"
	caseSensitive DeDupValueModifier = "case-sensitive"
	fuzzyMatch    DeDupValueModifier = "fuzzy-match"
	soundsLike    DeDupValueModifier = "sounds-like"

	oneOf DeDupMultiValueConstraint = "one-of"
	equal DeDupMultiValueConstraint = "equal"

	deDupWarning DeDupIssueKind = "duplication_warning"
	deDupError   DeDupIssueKind = "duplication_error"
)

func DeDup() *deDup {
	return &deDup{
		ls: locale.Global(),
	}
}

func (d deDup) CheckDuplication(ctx context.Context, rules DeDupRuleSet, rec Record, rr RecordSet) (out *RecordValueErrorSet, err error) {
	out = &RecordValueErrorSet{}
	err = rules.Walk(func(rule *DeDupRule) error {
		if rule.HasAttributes() {
			values := rr.GetValuesByName(distinct(rule.Attributes())...)

			set := rule.validateValue(ctx, d.ls, rec, values)

			if !set.IsValid() {
				out.Push(set.Set...)
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	if out.IsValid() {
		out = nil
	}
	return
}

func (rule DeDupIssueKind) String() string {
	return string(rule)
}

func (rule DeDupRule) HasAttributes() bool {
	return len(rule.ConstraintSet) > 0 && len(rule.Attributes()) > 0
}

func (rule DeDupRule) Attributes() (out []string) {
	for _, c := range rule.ConstraintSet {
		out = append(out, c.Attribute)
	}
	return
}

func (rule DeDupRule) IsStrict() bool {
	return rule.Strict
}

func (rule DeDupRule) IssueKind() string {
	out := deDupWarning
	if rule.Strict {
		out = deDupError
	}

	return out.String()
}

func (rule DeDupRule) IssueMessage() (out string) {
	return "record-field.errors.duplicateValue"
}

func (rule DeDupRule) String() string {
	return fmt.Sprintf("%s duplicate detection on `%s` field", rule.Name, strings.Join(rule.Attributes(), ", "))
}

// validateValue will check duplicate detection based on rules name
func (rule DeDupRule) validateValue(ctx context.Context, ls localeService, rec Record, vv RecordValueSet) (out *RecordValueErrorSet) {
	return rule.checkCaseSensitiveDuplication(ctx, ls, rec, vv)
}

func (rule DeDupRule) checkCaseSensitiveDuplication(ctx context.Context, ls localeService, rec Record, vv RecordValueSet) (out *RecordValueErrorSet) {
	var (
		recVal = rec.Values
	)

	for _, c := range rule.ConstraintSet {
		rvv := recVal.FilterByName(c.Attribute)
		if rvv.Len() == 0 {
			continue
		}

		var (
			valErr = &RecordValueErrorSet{}
		)

		_ = vv.Walk(func(v *RecordValue) error {
			if v.RecordID != rec.ID {
				_ = rvv.Walk(func(rv *RecordValue) error {
					if len(rv.Value) > 0 && matchValue(c.Modifier, rv.Value, v.Value) {
						valErr.Push(RecordValueError{
							Kind:    rule.IssueKind(),
							Message: ls.T(ctx, "compose", rule.IssueMessage()),
							Meta: map[string]interface{}{
								"field":         v.Name,
								"value":         v.Value,
								"dupValueField": rv.Name,
								"recordID":      cast.ToString(v.RecordID),
								"rule":          rule.String(),
							},
						})
					}
					return nil
				})

				// 1. multiValue is empty, then all value needs to be a match then return error/warning
				// 2. multiValue is oneOf, then one or more value needs to be a match then return error/warning
				// 3. multiValue is equal, then all value needs to be a match then return error/warning
				if (!valErr.IsValid() && (!c.HasMultiValue() || c.IsAllEqual()) && valErr.Len() == rvv.Len()) || (c.IsOneOf() && valErr.Len() > 0) {
					if out == nil {
						out = &RecordValueErrorSet{}
					}
					out.Push(valErr.Set...)
				}
			}
			return nil
		})
	}

	return
}

func (c DeDupRuleConstraint) HasMultiValue() bool {
	switch c.MultiValue {
	case oneOf, equal:
		return true
	default:
		return false
	}
}

func (c DeDupRuleConstraint) IsAllEqual() bool {
	return c.MultiValue == equal
}

func (c DeDupRuleConstraint) IsOneOf() bool {
	return c.MultiValue == oneOf
}

func (v *RecordValueErrorSet) SetMetaID(id uint64) {
	if v.IsValid() {
		return
	}

	for _, val := range v.Set {
		if val.Meta != nil {
			if _, ok := val.Meta["id"]; !ok {
				val.Meta["id"] = cast.ToString(id)
			}
		}
	}
}

func (v *RecordValueErrorSet) HasStrictErrors() bool {
	return v.HasKind(deDupError.String())
}

// distinct only list the different (distinct) values
func distinct(input []string) (out []string) {
	keys := make(map[string]bool)
	for _, val := range input {
		if _, ok := keys[val]; !ok {
			keys[val] = true
			out = append(out, val)
		}
	}
	return
}

// matchValue will check if the input matches with target string as per the modifier
func matchValue(modifier DeDupValueModifier, input string, target string) bool {
	switch modifier {
	case ignoreCase:
		return str.Match(input, target, str.CaseInSensitiveMatch)
	case caseSensitive:
		return str.Match(input, target, str.CaseSensitiveMatch)
	case fuzzyMatch:
		return str.Match(input, target, str.LevenshteinDistance)
	case soundsLike:
		return str.Match(input, target, str.Soundex)
	default:
		// ignoreCase as default, if not specified
		return str.Match(input, target, str.CaseInSensitiveMatch)
	}
}
