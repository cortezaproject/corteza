package values

import (
	"context"
	"strings"
	"testing"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/stretchr/testify/require"
)

func issuesModule(f *types.ModuleField) *types.Module {
	return &types.Module{
		ID: 1,
		Fields: types.ModuleFieldSet{
			&types.ModuleField{Name: "stage", Kind: "Select"},
			&types.ModuleField{Name: "qty", Kind: "Number"},
			&types.ModuleField{Name: "price", Kind: "Number", Options: types.ModuleFieldOptions{"precision": 2}},
			f,
		},
	}
}

// The scope a value expression gets is not the record shape the REST API
// returns: `record` and a bare `values` resolve to nothing, and the failure
// used to surface only on the next record create.
func Test_ExpressionIssues_valueExprScope(t *testing.T) {
	for _, tc := range []struct {
		name     string
		expr     string
		severity string
		contains string
	}{
		{"record.values is not in scope", `record.values.stage != "hired"`, IssueSeverityError, "unknown parameter record.values"},
		{"bare values is not in scope", `values.stage != "hired"`, IssueSeverityError, "unknown parameter values.stage"},
		{"single quotes do not parse", `stage != 'hired'`, IssueSeverityError, "does not parse"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)
			m := issuesModule(&types.ModuleField{
				Name:        "is_active",
				Kind:        "Bool",
				Expressions: types.ModuleFieldExpr{ValueExpr: tc.expr},
			})

			ii := ExpressionIssues(context.Background(), m)
			req.Len(ii, 1)
			req.Equal("is_active", ii[0].Field)
			req.Equal(IssueSlotValue, ii[0].Slot)
			req.Equal(tc.severity, ii[0].Severity)
			req.Equal(tc.expr, ii[0].Expression)
			req.Contains(ii[0].Message, tc.contains)
		})
	}
}

// A dry run that flagged working expressions would be worse than no check: the
// caller would learn to ignore it.
func Test_ExpressionIssues_workingExpressionsAreClean(t *testing.T) {
	for _, e := range []string{
		`stage != "hired" && stage != "rejected"`,
		`new.values.stage != "hired"`,
		`old.values.stage != "hired"`,
		`qty * price`,
		`qty > 0 ? "many" : "none"`,
		`length(stage) > 3`,
		`new.recordID`,
	} {
		t.Run(e, func(t *testing.T) {
			m := issuesModule(&types.ModuleField{
				Name:        "computed",
				Kind:        "String",
				Expressions: types.ModuleFieldExpr{ValueExpr: e},
			})
			require.Empty(t, ExpressionIssues(context.Background(), m))
		})
	}
}

// Each slot has its own scope, so an issue has to say which slot it is about —
// `value` is a name only a validator/sanitizer/formatter carries, and the
// record's fields are names only a value expression carries.
func Test_ExpressionIssues_perSlotScope(t *testing.T) {
	req := require.New(t)

	m := issuesModule(&types.ModuleField{
		Name: "score",
		Kind: "Number",
		Expressions: types.ModuleFieldExpr{
			// legal in its slot
			Validators: []types.ModuleFieldValidator{{Test: `value < 0 || value > 5`, Error: "out of range"}},
			Sanitizers: []string{`trim(value)`},
			// `values` is a validator name; a formatter's scope carries only `value`
			Formatters: []string{`values.stage`},
		},
	})

	ii := ExpressionIssues(context.Background(), m)
	req.Len(ii, 1)
	req.Equal(IssueSlotFormatter, ii[0].Slot)
	req.Equal("score", ii[0].Field)

	// ...and a value expression may not name `value`
	m2 := issuesModule(&types.ModuleField{
		Name:        "computed",
		Kind:        "String",
		Expressions: types.ModuleFieldExpr{ValueExpr: `value + "x"`},
	})
	ii2 := ExpressionIssues(context.Background(), m2)
	req.Len(ii2, 1)
	req.Equal(IssueSlotValue, ii2[0].Slot)
}

// An expression that only fails because the dry run stands every field empty is
// reported, but not as an error: it may well work on a record that has values.
func Test_ExpressionIssues_dataDependentIsAWarning(t *testing.T) {
	req := require.New(t)

	m := issuesModule(&types.ModuleField{
		Name: "when",
		Kind: "DateTime",
		// a DateTime field with no value casts to nil, so format() has nothing to take
		Expressions: types.ModuleFieldExpr{ValueExpr: `format(when, "2006-01-02")`},
	})

	ii := ExpressionIssues(context.Background(), m)
	req.Len(ii, 1)
	req.Equal(IssueSeverityWarning, ii[0].Severity)
	req.Contains(ii[0].Message, "dry run")
}

func Test_ExpressionIssues_noExpressionsNoIssues(t *testing.T) {
	require.Empty(t, ExpressionIssues(context.Background(), issuesModule(&types.ModuleField{Name: "plain", Kind: "String"})))
	require.Empty(t, ExpressionIssues(context.Background(), nil))
}

// The severity split keys on this prefix. If gval rewords it, every scope
// mistake quietly becomes a warning, which is the failure this whole check
// exists to prevent.
func Test_unknownParameterPrefix(t *testing.T) {
	eval, err := expr.Parser().NewEvaluable(`record.values.stage`)
	require.NoError(t, err)

	_, err = eval(context.Background(), map[string]any{})
	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), unknownParamPrefix),
		"gval no longer reports an unresolved name as %q: %v", unknownParamPrefix, err)
}
