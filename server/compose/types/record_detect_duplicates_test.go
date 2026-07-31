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
                            Message: rule1.IssueMessage(),
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
                            Message: rule1.IssueMultivalueMessage(),
							Meta: map[string]interface{}{
								"field":         "name",
								"dupValueField": "name",
								"rule":          rule1.String(),
                                "value": "test, test test",
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
                            Message: rule2.IssueMultivalueMessage(),
							Meta: map[string]interface{}{
								"field":         "name",
								"dupValueField": "name",
								"rule":          rule2.String(),
                                "value": "test, test tEst",
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
                            Message: numberRule.IssueMultivalueMessage(),
							Meta: map[string]interface{}{
								"field":         "count",
								"dupValueField": "count",
								"rule":          numberRule.String(),
                                "value": "234, 897",
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
                            Message: locationRule.IssueMultivalueMessage(),
							Meta: map[string]interface{}{
								"field":         "location",
								"dupValueField": "location",
								"rule":          locationRule.String(),
                                "value": "{\"coordinates\":[-6.7833479,20.3768206]}, {\"coordinates\":[0.7833479,10.3768206]}",
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

func TestDeDupRule_checkCompositeConstraintDuplication(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()
		ls  = locale.Global()

		compositeRule = DeDupRule{
			Name:   "composite rule",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute: "name",
					Modifier:  ignoreCase,
				},
				{
					Attribute: "email",
					Modifier:  caseSensitive,
				},
			},
		}

		makeRecord = func(id uint64, name, email string) Record {
			return Record{
				ID: id,
				module: &Module{
					ID: 1,
					Fields: ModuleFieldSet{
						&ModuleField{
							Name:  "name",
							Kind:  "String",
							Multi: false,
						},
						&ModuleField{
							Name:  "email",
							Kind:  "Email",
							Multi: false,
						},
					},
				},
				Values: RecordValueSet{
					&RecordValue{
						RecordID: id,
						Name:     "name",
						Value:    name,
					},
					&RecordValue{
						RecordID: id,
						Name:     "email",
						Value:    email,
					},
				},
			}
		}

		tests = []struct {
			name      string
			rule      DeDupRule
			rec       Record
			vv        RecordValueSet
			expectErr bool
		}{
			{
				name: "composite match - both fields match same record",
				rule: compositeRule,
				rec:  makeRecord(1, "John Doe", "john@test.com"),
				vv: RecordValueSet{
					&RecordValue{RecordID: 2, Name: "name", Value: "John Doe"},
					&RecordValue{RecordID: 2, Name: "email", Value: "john@test.com"},
				},
				expectErr: true,
			},
			{
				name: "composite partial match - only name matches (bug fix)",
				rule: compositeRule,
				rec:  makeRecord(1, "John Doe", "john@test.com"),
				vv: RecordValueSet{
					&RecordValue{RecordID: 2, Name: "name", Value: "John Doe"},
					&RecordValue{RecordID: 2, Name: "email", Value: "other@test.com"},
				},
				expectErr: false,
			},
			{
				name: "composite partial match - only email matches (bug fix)",
				rule: compositeRule,
				rec:  makeRecord(1, "John Doe", "john@test.com"),
				vv: RecordValueSet{
					&RecordValue{RecordID: 2, Name: "name", Value: "Jane Smith"},
					&RecordValue{RecordID: 2, Name: "email", Value: "john@test.com"},
				},
				expectErr: false,
			},
			{
				name: "composite no match - neither field matches",
				rule: compositeRule,
				rec:  makeRecord(1, "John Doe", "john@test.com"),
				vv: RecordValueSet{
					&RecordValue{RecordID: 2, Name: "name", Value: "Jane Smith"},
					&RecordValue{RecordID: 2, Name: "email", Value: "jane@test.com"},
				},
				expectErr: false,
			},
			{
				name: "composite cross-field value swap - new(2,1) vs existing(1,2)",
				rule: compositeRule,
				rec:  makeRecord(1, "2", "1"),
				vv: RecordValueSet{
					&RecordValue{RecordID: 2, Name: "name", Value: "1"},
					&RecordValue{RecordID: 2, Name: "email", Value: "2"},
				},
				expectErr: false,
			},
			{
				name: "composite with missing field - new record has no email",
				rule: compositeRule,
				rec: Record{
					ID: 1,
					module: &Module{
						ID: 1,
						Fields: ModuleFieldSet{
							&ModuleField{Name: "name", Kind: "String", Multi: false},
							&ModuleField{Name: "email", Kind: "Email", Multi: false},
						},
					},
					Values: RecordValueSet{
						&RecordValue{RecordID: 1, Name: "name", Value: "John Doe"},
					},
				},
				vv: RecordValueSet{
					&RecordValue{RecordID: 2, Name: "name", Value: "John Doe"},
					&RecordValue{RecordID: 2, Name: "email", Value: "john@test.com"},
				},
				expectErr: false,
			},
			{
				name: "composite with missing field - existing record has no email",
				rule: compositeRule,
				rec:  makeRecord(1, "John Doe", "john@test.com"),
				vv: RecordValueSet{
					&RecordValue{RecordID: 2, Name: "name", Value: "John Doe"},
				},
				expectErr: false,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := tt.rule.checkDuplication(ctx, ls, tt.rec, tt.vv)
			if tt.expectErr {
				req.False(gotOut.IsValid(), "expected duplication errors but got none")
				req.Greater(gotOut.Len(), 0, "expected at least one error")
			} else {
				req.True(gotOut.IsValid(), "expected no duplication errors but got: %v", gotOut)
			}
		})
	}
}

func TestDeDupRule_checkDuplicationOnRemovedField(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()
		ls  = locale.Global()

		rule = DeDupRule{
			Name:   "removed field rule",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute: "gone",
					Modifier:  caseSensitive,
				},
			},
		}

		// module no longer has the "gone" field, but the record and the
		// existing values still carry it
		rec = Record{
			ID: 1,
			module: &Module{
				ID:     1,
				Fields: ModuleFieldSet{&ModuleField{Name: "name"}},
			},
			Values: RecordValueSet{
				&RecordValue{RecordID: 1, Name: "gone", Value: "same"},
			},
		}

		vv = RecordValueSet{
			&RecordValue{RecordID: 2, Name: "gone", Value: "same"},
		}
	)

	req.NotPanics(func() {
		out := rule.checkDuplication(ctx, ls, rec, vv)
		req.True(out.IsValid(), "expected no duplication errors but got: %v", out)
	})
}

func TestDeDupRule_checkMultiValueEqualNoDuplication(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()
		ls  = locale.Global()

		tagRule = DeDupRule{
			Name:   "tag rule",
			Strict: true,
			ConstraintSet: []*DeDupRuleConstraint{
				{
					Attribute:  "tags",
					Modifier:   caseSensitive,
					MultiValue: equal,
				},
			},
		}

		rec = Record{
			ID: 1,
			module: &Module{
				ID: 1,
				Fields: ModuleFieldSet{
					&ModuleField{Name: "tags", Multi: true},
				},
			},
			Values: RecordValueSet{
				&RecordValue{RecordID: 1, Name: "tags", Value: "one"},
				&RecordValue{RecordID: 1, Name: "tags", Value: "two"},
			},
		}

		vv = RecordValueSet{
			&RecordValue{RecordID: 2, Name: "tags", Value: "three"},
			&RecordValue{RecordID: 2, Name: "tags", Value: "four"},
		}
	)

	// none of the values match, so nothing may be reported — an empty
	// RecordValueError must not end up in the set
	out := tagRule.checkDuplication(ctx, ls, rec, vv)
	req.True(out.IsValid(), "expected no duplication errors but got: %v", out)
	req.Zero(out.Len())
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
