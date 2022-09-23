package types

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/locale"
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
		Name       DeDupRuleName `json:"name"`
		Strict     bool          `json:"strict"`
		Attributes []string      `json:"attributes"`
	}

	// DeDupRuleName represent the identifier for duplicate detection rule
	DeDupRuleName string

	// DeDupIssueKind based on strict mode rule or duplication config
	DeDupIssueKind string
)

const (
	caseSensitive DeDupRuleName = "case-sensitive"

	dupWarning DeDupIssueKind = "duplication_warning"
	dupError   DeDupIssueKind = "duplication_error"
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
			values := rr.GetValuesByName(distinct(rule.Attributes)...)

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
	return len(rule.Attributes) > 0
}

func (rule DeDupRule) IsStrict() bool {
	return rule.Strict
}

func (rule DeDupRule) IssueKind() string {
	out := dupWarning
	if rule.Strict {
		out = dupError
	}

	return out.String()
}

func (rule DeDupRule) IssueMessage() (out string) {
	return "record-field.errors.duplicateValue"
}

func (rule DeDupRule) String() string {
	return fmt.Sprintf("%s duplicate detection on `%s` field", rule.Name, strings.Join(rule.Attributes, ", "))
}

// validateValue will check duplicate detection based on rules name
func (rule DeDupRule) validateValue(ctx context.Context, ls localeService, rec Record, vv RecordValueSet) (out *RecordValueErrorSet) {
	switch rule.Name {
	case caseSensitive:
		return rule.checkCaseSensitiveDuplication(ctx, ls, rec, vv)
	default:
		return rule.checkCaseSensitiveDuplication(ctx, ls, rec, vv)
	}
}

func (rule DeDupRule) checkCaseSensitiveDuplication(ctx context.Context, ls localeService, rec Record, vv RecordValueSet) (out *RecordValueErrorSet) {
	out = &RecordValueErrorSet{}
	recVal := rec.Values

	_ = recVal.Walk(func(newV *RecordValue) error {
		_ = vv.Walk(func(v *RecordValue) error {
			if v.RecordID != rec.ID {
				if toLower(v.Value) == toLower(newV.Value) {
					out.Push(RecordValueError{
						Kind:    rule.IssueKind(),
						Message: ls.T(ctx, "compose", rule.IssueMessage()),
						Meta: map[string]interface{}{
							"dupValueField": newV.Name,
							"recordID":      cast.ToString(v.RecordID),
							"field":         v.Name,
							"value":         v.Value,
							"rule":          rule.String(),
						},
					})
				}
			}
			return nil
		})
		return nil
	})
	return
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
	return v.HasKind(dupError.String())
}

// CaseSensitiveDuplicationRule prepares the case-sensitive duplicate detection rule
func CaseSensitiveDuplicationRule(strict bool, identifiers ...string) DeDupRule {
	return makeDuplicationRule(caseSensitive, strict, identifiers...)
}

// makeDuplicationRule prepares duplication detection rules
func makeDuplicationRule(name DeDupRuleName, strict bool, attributes ...string) DeDupRule {
	return DeDupRule{
		Name:       name,
		Strict:     strict,
		Attributes: attributes,
	}
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

func toLower(s string) string {
	return strings.ToLower(s)
}
