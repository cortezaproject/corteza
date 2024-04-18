package types

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestDeDupRule_checkCaseSensitiveDuplication(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()
		ls  = locale.Global()

		rule1 = DeDupRule{
			Name:   "",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute: "name",
					Modifier:  ignoreCase,
				},
			},
		}

		tests = []struct {
			name    string
			rule    DeDupRule
			rec     Record
			vv      RecordValueSet
			wantOut *RecordValueErrorSet
		}{
			{
				name: "no duplication",
				rule: rule1,
				rec: Record{
					ID: 1,
					module: &Module{
						ID: 1,
						Fields: ModuleFieldSet{
							&ModuleField{
								Name:  "name",
								Kind:  "String",
								Multi: false,
							},
						},
					},
					Values: RecordValueSet{
						&RecordValue{
							RecordID: 1,
							Name:     "name",
							Value:    "test",
						},
					},
				},
				vv: RecordValueSet{
					&RecordValue{
						RecordID: 2,
						Name:     "name",
						Value:    "test",
					},
				},
				wantOut: &RecordValueErrorSet{
					Set: []RecordValueError{
						{
							Kind:    deDupError.String(),
							Message: rule1.IssueMessage("test"),
							Meta: map[string]interface{}{
								"field":         "name",
								"value":         "test",
								"dupValueField": "name",
								"recordID":      cast.ToString(2),
								"rule":          rule1.String(),
							},
						},
					},
				},
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := tt.rule.checkDuplication(ctx, ls, tt.rec, tt.vv)
			req.Equal(tt.wantOut, gotOut, "checkDuplication() = %v, want %v", gotOut, tt.wantOut)
		})
	}
}

func TestDedupRule_checkMultiValueEqualDuplication(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()
		ls  = locale.Global()

		rule1 = DeDupRule{
			Name:   "",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute:  "name",
					Modifier:   caseSensitive,
					MultiValue: equal,
				},
			},
		}

		rule2 = DeDupRule{
			Name:   "ignore case rule",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute:  "name",
					Modifier:   ignoreCase,
					MultiValue: equal,
				},
			},
		}

		numberRule = DeDupRule{
			Name:   "number rule",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute:  "count",
					Modifier:   ignoreCase,
					MultiValue: equal,
				},
			},
		}

		locationRule = DeDupRule{
			Name:   "location rule",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute:  "location",
					Modifier:   ignoreCase,
					MultiValue: equal,
				},
			},
		}

		tests = []struct {
			name    string
			rule    DeDupRule
			rec     Record
			vv      RecordValueSet
			wantOut *RecordValueErrorSet
		}{
			{
				name: "no duplication",
				rule: rule1,
				rec: Record{
					ID: 1,
					module: &Module{
						ID: 1,
						Fields: ModuleFieldSet{
							&ModuleField{
								Name:  "name",
								Kind:  "String",
								Multi: true,
							},
						},
					},
					Values: RecordValueSet{
						&RecordValue{
							RecordID: 1,
							Name:     "name",
							Value:    "test",
						},
						&RecordValue{
							RecordID: 1,
							Name:     "name",
							Value:    "test test",
						},
					},
				},
				vv: RecordValueSet{
					&RecordValue{
						RecordID: 0,
						Name:     "name",
						Value:    "test",
					},
					&RecordValue{
						RecordID: 0,
						Name:     "name",
						Value:    "test test",
					},
				},
				wantOut: &RecordValueErrorSet{
					Set: []RecordValueError{
						{
							Kind:    deDupError.String(),
							Message: rule1.IssueMultivalueMessage([]string{"test", "test test"}),
							Meta: map[string]interface{}{
								"field":         "name",
								"dupValueField": "name",
								"rule":          rule1.String(),
							},
						},
					},
				},
			},
			{
				name: "no duplication",
				rule: rule2,
				rec: Record{
					ID: 1,
					module: &Module{
						ID: 1,
						Fields: ModuleFieldSet{
							&ModuleField{
								Name:  "name",
								Kind:  "String",
								Multi: true,
							},
						},
					},
					Values: RecordValueSet{
						&RecordValue{
							RecordID: 1,
							Name:     "name",
							Value:    "test",
						},
						&RecordValue{
							RecordID: 1,
							Name:     "name",
							Value:    "test tEst",
						},
					},
				},
				vv: RecordValueSet{
					&RecordValue{
						RecordID: 0,
						Name:     "name",
						Value:    "test",
					},
					&RecordValue{
						RecordID: 0,
						Name:     "name",
						Value:    "Test Test",
					},
				},
				wantOut: &RecordValueErrorSet{
					Set: []RecordValueError{
						{
							Kind:    deDupError.String(),
							Message: rule2.IssueMultivalueMessage([]string{"test", "test tEst"}),
							Meta: map[string]interface{}{
								"field":         "name",
								"dupValueField": "name",
								"rule":          rule2.String(),
							},
						},
					},
				},
			},
			{
				name: "no duplication",
				rule: numberRule,
				rec: Record{
					ID: 1,
					module: &Module{
						ID: 1,
						Fields: ModuleFieldSet{
							&ModuleField{
								Name:  "count",
								Multi: true,
							},
						},
					},
					Values: RecordValueSet{
						&RecordValue{
							RecordID: 1,
							Name:     "count",
							Value:    "234",
						},
						&RecordValue{
							RecordID: 1,
							Name:     "count",
							Value:    "897",
						},
					},
				},
				vv: RecordValueSet{
					&RecordValue{
						RecordID: 0,
						Name:     "count",
						Value:    "897",
					},
					&RecordValue{
						RecordID: 0,
						Name:     "count",
						Value:    "234",
					},
				},
				wantOut: &RecordValueErrorSet{
					Set: []RecordValueError{
						{
							Kind:    deDupError.String(),
							Message: numberRule.IssueMultivalueMessage([]string{"234", "897"}),
							Meta: map[string]interface{}{
								"field":         "count",
								"dupValueField": "count",
								"rule":          numberRule.String(),
							},
						},
					},
				},
			},
			{
				name: "no duplication",
				rule: locationRule,
				rec: Record{
					ID: 1,
					module: &Module{
						ID: 1,
						Fields: ModuleFieldSet{
							&ModuleField{
								Name:  "location",
								Multi: true,
							},
						},
					},
					Values: RecordValueSet{
						&RecordValue{
							RecordID: 1,
							Name:     "location",
							Value:    "{\"coordinates\":[-6.7833479,20.3768206]}",
						},
						&RecordValue{
							RecordID: 1,
							Name:     "location",
							Value:    "{\"coordinates\":[0.7833479,10.3768206]}",
						},
					},
				},
				vv: RecordValueSet{
					&RecordValue{
						RecordID: 0,
						Name:     "location",
						Value:    "{\"coordinates\":[0.7833479,10.3768206]}",
					},
					&RecordValue{
						RecordID: 0,
						Name:     "location",
						Value:    "{\"coordinates\":[-6.7833479,20.3768206]}",
					},
				},
				wantOut: &RecordValueErrorSet{
					Set: []RecordValueError{
						{
							Kind:    deDupError.String(),
							Message: locationRule.IssueMultivalueMessage([]string{"{\"coordinates\":[-6.7833479,20.3768206]}", "{\"coordinates\":[0.7833479,10.3768206]}"}),
							Meta: map[string]interface{}{
								"field":         "location",
								"dupValueField": "location",
								"rule":          locationRule.String(),
							},
						},
					},
				},
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := tt.rule.checkDuplication(ctx, ls, tt.rec, tt.vv)
			req.Equal(tt.wantOut, gotOut, "checkDuplication() = %v, want %v", gotOut, tt.wantOut)
		})
	}
}

func Test_matchValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		target   string
		modifier DeDupValueModifier
		want     bool
	}{
		{
			name:     "ignoreCase match value",
			input:    "test",
			target:   "tEst",
			modifier: ignoreCase,
			want:     true,
		},
		{
			name:     "caseSensitive match value",
			input:    "tEst",
			target:   "tEst",
			modifier: caseSensitive,
			want:     true,
		},
		{
			name:     "fuzzyMatch match value",
			input:    "kitten",
			target:   "sitting",
			modifier: fuzzyMatch,
			want:     true,
		},
		{
			name:     "soundsLike match value",
			input:    "Robert",
			target:   "Rupert",
			modifier: soundsLike,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchValue(tt.modifier, tt.input, tt.target); got != tt.want {
				t.Errorf("matchValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rulesetValidation(t *testing.T) {
	var (
		req = require.New(t)

		tests = []struct {
			name    string
			ruleset DeDupRuleSet
		}{
			{
				name: "no constraint",
				ruleset: DeDupRuleSet{&DeDupRule{
					Name:          "",
					Strict:        true,
					ConstraintSet: []*DeDupRuleConstraint{},
				}},
			},
			{
				name: "invalid constraint",
				ruleset: DeDupRuleSet{&DeDupRule{
					Name:   "",
					Strict: true,
					ConstraintSet: []*DeDupRuleConstraint{
						{
							Attribute: "",
						},
					},
				}},
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req.Error(tt.ruleset.Validate())
		})
	}
}
