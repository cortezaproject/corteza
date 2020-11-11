package values

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/compose/types"
)

func makeInvalidExprErr(field *types.ModuleField, expr string, err error) types.RecordValueError {
	return types.RecordValueError{
		Kind:    "valueExpression",
		Message: fmt.Sprintf("invalid expression %q: %v", expr, err.Error()),
		Meta:    map[string]interface{}{"field": field.Name},
	}
}
func makeExprEvalErr(field *types.ModuleField, expr string, err error) types.RecordValueError {
	return types.RecordValueError{
		Kind:    "valueExpression",
		Message: fmt.Sprintf("failed to evaluate formula expression %q: %v", expr, err.Error()),
		Meta:    map[string]interface{}{"field": field.Name},
	}
}

func makeValueExprIncompErr(field *types.ModuleField) types.RecordValueError {
	return types.RecordValueError{
		Kind:    "evaluatedValueIncompatible",
		Message: "evaluated results incompatible",
		Meta:    map[string]interface{}{"field": field.Name},
	}
}

func parser() gval.Language {
	return gval.Full()
}

// Expression evaluates expression in ModuleField.Expressions.Value and
// assigns results to the record on that field
func Expression(ctx context.Context, m *types.Module, r *types.Record, old *types.Record, rve *types.RecordValueErrorSet) {
	var (
		exprParser = parser()

		scope = make(map[string]interface{})

		reserved = map[string]bool{
			"new": true,
			"old": true,
		}
	)

	// base scope with field=value(s) from new record
	scope = r.Values.Dict(m.Fields)

	// new record
	scope["new"] = r.Dict(m)

	if old != nil {
		// old values on record (before update)
		// this will not be set for new records
		scope["old"] = old.Dict(m)
	}

	for _, f := range m.Fields {
		if f.Expressions.Value == "" {
			continue
		}

		expr := f.Expressions.Value

		eval, err := exprParser.NewEvaluable(expr)
		if err != nil {
			rve.Push(makeInvalidExprErr(f, expr, err))
			return
		}

		tmp, err := eval(ctx, scope)
		if err != nil {
			rve.Push(makeExprEvalErr(f, expr, err))
			return
		}

		var strings []string
		if values, isSlice := tmp.([]interface{}); isSlice {
			if !f.Multi {
				rve.Push(makeValueExprIncompErr(f))
				continue
			}

			strings = make([]string, len(values))
			for i, value := range values {
				strings[i] = sanitize(f, value)
			}
		} else {
			if f.Multi {
				rve.Push(makeValueExprIncompErr(f))
				continue
			}

			strings = []string{sanitize(f, tmp)}

		}

		r.Values = r.Values.Replace(f.Name, strings...)

		if !reserved[f.Name] {
			// make sure we do not overrider reserved fields
			scope[f.Name] = tmp
		}

		// Reset $new with updated data
		scope["new"] = r.Dict(m)
	}
}
