package types

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	// Used for expression steps, arguments/results mapping and for input validation
	Expr struct {
		// Variable name to set results of the expression to
		Target string `json:"target"`

		// Source of the value / name of the variable from scope; if set value is copied target
		//
		// Takes precedence to value
		//
		Source string `json:"source,omitempty"`

		// Expression to evaluate over the input variables; results will be set to target
		//
		Expr string `json:"expr,omitempty"`

		// Raw value to be set to target
		//
		// If expression is set and fails, evaluation defaults to value and suppresses the error
		Value interface{} `json:"value,omitempty"`

		eval expr.Evaluable

		// Expected type of the input value
		Type string `json:"type,omitempty"`

		typ expr.Type

		// Set of tests that can be run before input is evaluated and result copied to scope
		Tests TestSet `json:"tests,omitempty"`
	}

	ExprSet []*Expr

	// WorkflowStepExpression is created from WorkflowStep with kind=expressions
	expressionsStep struct {
		wfexec.StepIdentifier
		Set ExprSet
	}
)

func NewExpr(target, typ, expr string) (e *Expr, err error) {
	return &Expr{Expr: expr, Target: target, Type: typ}, nil
}

func (e *Expr) SetType(fn func(string) (expr.Type, error)) error {
	if typ, err := fn(e.Type); err != nil {
		return err
	} else {
		e.typ = typ
		return nil
	}
}

func (e Expr) GetExpr() string              { return e.Expr }
func (e *Expr) SetEval(eval expr.Evaluable) { e.eval = eval }
func (e Expr) Eval(ctx context.Context, scope *expr.Vars) (interface{}, error) {
	return e.eval.Eval(ctx, scope)
}
func (e Expr) Test(ctx context.Context, scope *expr.Vars) (bool, error) {
	return e.eval.Test(ctx, scope)
}

func (set ExprSet) GetByTarget(t string) *Expr {
	for _, e := range set {
		if e.Target == t {
			return e
		}
	}
	return nil
}

func (set ExprSet) Validate(ctx context.Context, in *expr.Vars) (TestSet, error) {
	var (
		out TestSet
		vv  TestSet
		err error

		// Copy/create scope
		scope = (&expr.Vars{}).Merge(in)
	)

	for _, e := range set {
		vv, err = e.Tests.Validate(ctx, scope)
		if err != nil {
			return nil, err
		}

		out = append(out, vv...)
	}

	return out, nil
}

// Eval on expression set (ExprSet) evaluates all expressions in the set and returns new scope with all set targets
func (set ExprSet) Eval(ctx context.Context, in *expr.Vars) (*expr.Vars, error) {
	var (
		err error

		// Copy input to scope
		scope = (&expr.Vars{}).Merge(in)

		// Prepare output scope
		out = &expr.Vars{}

		// Untyped evaluation result
		value interface{}

		knownType = func(p expr.Type) bool {
			return p != nil && p.Type() != expr.Any{}.Type() && p.Type() != expr.Unresolved{}.Type()
		}
	)

	for _, e := range set {
		value = e.Value

		if e.typ == nil {
			return nil, errors.Internal("type for target %q not initialized", e.Target)
		}

		err = func() (err error) {
			if len(e.Source) > 0 {
				// can copy from existing variable
				if !scope.Has(e.Source) {
					return errors.NotFound("variable %q does not exist", e.Source)
				}

				value, _ = expr.Select(scope, e.Source)
				return
			}

			if len(e.Expr) > 0 {
				if e.eval == nil {
					// no expression set, fallback to default value
					return errors.Internal("expression language for target %q not initialized", e.Target)
				} else if value, err = e.eval.Eval(ctx, scope); err != nil {
					return errors.Internal("expression %q failed: %s", e.Expr, err.Error()).Wrap(err)
				}
			}

			return
		}()

		if err != nil && e.Value == nil {
			return nil, err
		}

		typedValue, is := value.(expr.TypedValue)
		if !is {
			if e.typ == nil {
				typedValue, _ = expr.NewAny(value)
			} else if typedValue, err = e.typ.Cast(value); err != nil {
				return nil, fmt.Errorf("cannot cast value on %s to type %s", e.Target, typedValue.Type())
			}
		}

		if !knownType(e.typ) && !knownType(typedValue) && typedValue.Type() != e.typ.Type() {
			// Both, expression & value have type set;
			// check if it's the same type or return an error
			return nil, fmt.Errorf("cannot set to %q (type %s) value of type %s", e.Target, e.typ.Type(), typedValue.Type())
		}

		if e.typ != nil {
			if !knownType(typedValue) {
				// Expression has fixed type but value does not
				// cast the value of evaluation to type of the expressicason
				if typedValue, err = e.typ.Cast(value); err != nil {
					return nil, err
				}
			} else if e.typ.Type() != typedValue.Type() && e.typ.Type() != (expr.Any{}).Type() {
				//
				if typedValue, err = e.typ.Cast(value); err != nil {
					return nil, err
				}
			}
		}

		// Set result of the expression to scope
		//
		// Set() fn handles multi-level path (eg "base.level1.level2")
		// that can set result of the expression deep into scope's value
		if err = expr.Assign(scope, e.Target, typedValue); err != nil {
			return nil, err
		}

		// Take base of the path (1st part) and
		// copy value of it to output scope
		//
		// This ensures us that the entire variable
		// from the original scope will be present in the output
		scope.Copy(out, expr.PathBase(e.Target))
	}

	return out, nil
}

func ExpressionsStep(ee ...*Expr) *expressionsStep {
	return &expressionsStep{Set: ee}
}

func (s *expressionsStep) Exec(ctx context.Context, r *wfexec.ExecRequest) (wfexec.ExecResponse, error) {
	result, err := s.Set.Eval(ctx, r.Scope.Merge(r.Input))
	if err != nil {
		return nil, err
	}

	return result, nil
}
