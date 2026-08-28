package values

import (
	"context"
	"fmt"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/expr"
)

type (
	// ExpressionIssue is one thing wrong with one expression on one field.
	//
	// Slot names what the expression is for, because the four slots do not share
	// a scope: a value expression reads the record's fields, a validator reads
	// the value under test, a sanitizer and a formatter read only `value`.
	ExpressionIssue struct {
		Field      string `json:"field"`
		Slot       string `json:"slot"`
		Severity   string `json:"severity"`
		Expression string `json:"expression"`
		Message    string `json:"message"`
	}
)

const (
	// A parse failure and an unresolved name are fatal whatever the record
	// holds. Anything else the dry run turns up depends on the values, so it is
	// reported without claiming the expression is broken.
	IssueSeverityError   = "error"
	IssueSeverityWarning = "warning"

	IssueSlotValue     = "value"
	IssueSlotValidator = "validator"
	IssueSlotSanitizer = "sanitizer"
	IssueSlotFormatter = "formatter"

	// gval reports a name its scope does not carry with this prefix. Pinned by
	// TestUnknownParameterPrefix so an upgrade that rewords it is caught here
	// rather than silently demoting every scope mistake to a warning.
	unknownParamPrefix = "unknown parameter"
)

// ExpressionIssues parses and dry-runs every expression a module's fields carry
// and reports what cannot work.
//
// The dry run stands each field at its empty value and evaluates against the
// same scope the record save builds, so a name that does not resolve is found
// at module save instead of surfacing later as a failed record create blamed on
// the record.
//
// It reports rather than refuses: a module already stored with a broken
// expression must stay editable, and an expression can fail the dry run for
// reasons that depend on the values a real record carries.
//
// What it cannot catch: a bare identifier the scope does not carry evaluates to
// nil instead of failing, so a misspelled field name is only found where the
// expression goes on to do something with it. Only a dotted path reports itself
// as unresolved.
func ExpressionIssues(ctx context.Context, m *types.Module) (issues []ExpressionIssue) {
	if m == nil {
		return
	}

	var (
		parser = expr.Parser()
		scope  = ValueExprScope(m, emptyRecord(m), nil)
	)

	for _, f := range m.Fields {
		if f == nil {
			continue
		}

		if f.Expressions.ValueExpr != "" {
			issues = append(issues, checkExpr(ctx, parser, f.Name, IssueSlotValue, f.Expressions.ValueExpr, scope)...)
		}

		for _, v := range f.Expressions.Validators {
			if v.Test == "" {
				continue
			}

			issues = append(issues, checkExpr(ctx, parser, f.Name, IssueSlotValidator, v.Test, map[string]interface{}{
				"value":    "",
				"oldValue": "",
				"values":   scope["new"].(map[string]interface{})["values"],
			})...)
		}

		for _, e := range f.Expressions.Sanitizers {
			if e == "" {
				continue
			}
			issues = append(issues, checkExpr(ctx, parser, f.Name, IssueSlotSanitizer, e, map[string]interface{}{"value": ""})...)
		}

		for _, e := range f.Expressions.Formatters {
			if e == "" {
				continue
			}
			issues = append(issues, checkExpr(ctx, parser, f.Name, IssueSlotFormatter, e, map[string]interface{}{"value": ""})...)
		}
	}

	return
}

// emptyRecord is a record of this module with every field present and empty —
// the most permissive realistic scope. Casting is what turns an empty value
// into the type the field's kind implies, so a Number stands at 0 and a String
// at "", and an expression doing arithmetic on its own fields dry-runs clean.
func emptyRecord(m *types.Module) *types.Record {
	r := &types.Record{ModuleID: m.ID, NamespaceID: m.NamespaceID}

	for _, f := range m.Fields {
		if f == nil {
			continue
		}
		r.Values = append(r.Values, &types.RecordValue{Name: f.Name})
	}

	return r
}

func checkExpr(ctx context.Context, parser gval.Language, field, slot, e string, scope map[string]interface{}) []ExpressionIssue {
	eval, err := parser.NewEvaluable(e)
	if err != nil {
		return []ExpressionIssue{{
			Field:      field,
			Slot:       slot,
			Severity:   IssueSeverityError,
			Expression: e,
			Message: fmt.Sprintf("expression does not parse: %v. String literals need double quotes — "+
				"a single-quoted string is a syntax error", err),
		}}
	}

	if _, err = eval(ctx, scope); err != nil {
		severity, msg := IssueSeverityWarning, fmt.Sprintf(
			"expression failed a dry run with every field empty: %v. It may still work on a record that carries values", err)

		if strings.HasPrefix(err.Error(), unknownParamPrefix) {
			severity = IssueSeverityError
			msg = fmt.Sprintf("%v — that name is not in scope for a %s expression, so this fails on every save", err, slot)
		}

		return []ExpressionIssue{{
			Field:      field,
			Slot:       slot,
			Severity:   severity,
			Expression: e,
			Message:    msg,
		}}
	}

	return nil
}
