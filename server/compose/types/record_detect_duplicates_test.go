package types

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"testing"
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
			gotOut := tt.rule.checkCaseSensitiveDuplication(ctx, ls, tt.rec, tt.vv)
			req.Equal(tt.wantOut, gotOut, "checkCaseSensitiveDuplication() = %v, want %v", gotOut, tt.wantOut)
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
