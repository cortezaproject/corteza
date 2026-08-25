package types

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/locale"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

// A modifier has to decide the match. Keying the candidate map by the raw value
// and requiring a hit before calling matchValue made every modifier equivalent
// to case-sensitive, so a rule declared ignore-case admitted the duplicate it
// was written to refuse.
func TestDeDupRule_modifierDecidesTheMatch(t *testing.T) {
	var (
		ctx = context.Background()
		ls  = locale.Global()

		ruleWith = func(m DeDupValueModifier) DeDupRule {
			return DeDupRule{
				Strict:        true,
				ConstraintSet: []*DeDupRuleConstraint{{Attribute: "email", Modifier: m}},
			}
		}

		recordWith = func(v string) Record {
			return Record{
				ID: 1,
				module: &Module{
					ID:     1,
					Fields: ModuleFieldSet{&ModuleField{Name: "email", Kind: "String"}},
				},
				Values: RecordValueSet{&RecordValue{RecordID: 1, Name: "email", Value: v}},
			}
		}

		existing = RecordValueSet{
			&RecordValue{RecordID: 2, Name: "email", Value: "priya.raghunathan@example.org"},
		}
	)

	for _, tc := range []struct {
		name     string
		modifier DeDupValueModifier
		incoming string
		wantDup  bool
	}{
		{"ignore-case catches a differing case", ignoreCase, "PRIYA.RAGHUNATHAN@example.org", true},
		{"ignore-case still catches an exact match", ignoreCase, "priya.raghunathan@example.org", true},
		{"case-sensitive ignores a differing case", caseSensitive, "PRIYA.RAGHUNATHAN@example.org", false},
		{"case-sensitive catches an exact match", caseSensitive, "priya.raghunathan@example.org", true},
		{"a genuinely different value is not a duplicate", ignoreCase, "someone.else@example.org", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rule := ruleWith(tc.modifier)
			got := rule.checkDuplication(ctx, ls, recordWith(tc.incoming), existing)

			if !tc.wantDup {
				require.Empty(t, got.Set, "expected no duplicate, got %v", got.Set)
				return
			}

			require.Equal(t, &RecordValueErrorSet{Set: []RecordValueError{{
				Kind:    deDupError.String(),
				Message: rule.IssueMessage(),
				Meta: map[string]interface{}{
					"field":         "email",
					"value":         "priya.raghunathan@example.org",
					"dupValueField": "email",
					"recordID":      cast.ToString(2),
					"rule":          rule.String(),
				},
			}}}, got)
		})
	}
}
