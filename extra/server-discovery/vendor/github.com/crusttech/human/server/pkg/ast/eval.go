package ast

import (
	"fmt"
	"strings"

	"github.com/crusttech/human/server/pkg/expr"
	"github.com/spf13/cast"
)

// Eval evaluates the given AST tree with provided vars
func Eval(node *ASTNode, scope map[string]*expr.Vars) (expr.TypedValue, error) {
	if node == nil {
		return &expr.Boolean{}, nil
	}

	switch {
	case node.Symbol != "":
		// Leaf: variable lookup in the named scope.
		scopeName := "global"
		if s, ok := node.Meta["scope"].(string); ok && s != "" {
			scopeName = s
		}
		vars := scope[scopeName]
		if vars == nil {
			return nil, fmt.Errorf("eval: unknown scope %q", scopeName)
		}
		parts := strings.SplitN(node.Symbol, ".", 2)
		v, err := vars.Select(parts[0])
		if err != nil {
			return nil, fmt.Errorf("eval: unknown symbol %q in scope %q: %w", node.Symbol, scopeName, err)
		}
		if len(parts) == 1 {
			return v, nil
		}
		return selectNested(v, parts[1])

	case node.Value != nil:
		// Leaf: typed literal
		return node.Value.V, nil

	default:
		// Operator / ref node
		return evalRef(node, scope)
	}
}

// EvalBool evaluates node and casts the result to a bool.
func EvalBool(node *ASTNode, scope map[string]*expr.Vars) (bool, error) {
	v, err := Eval(node, scope)
	if err != nil {
		return false, err
	}

	b, err := expr.NewBoolean(v)
	if err != nil {
		return false, fmt.Errorf("evalBool: cannot cast %T to Boolean: %w", v, err)
	}

	return b.Get().(bool), nil
}

// selectNested traverses a dot-separated path into nested Vars or map[string]interface{} values.
func selectNested(v expr.TypedValue, path string) (expr.TypedValue, error) {
	parts := strings.SplitN(path, ".", 2)
	key := parts[0]

	switch c := v.(type) {
	case *expr.Vars:
		child, err := c.Select(key)
		if err != nil {
			return nil, fmt.Errorf("eval: no key %q: %w", key, err)
		}
		if len(parts) == 1 {
			return child, nil
		}
		return selectNested(child, parts[1])
	case expr.FieldSelector:
		child, err := c.Select(key)
		if err != nil {
			return nil, fmt.Errorf("eval: no key %q: %w", key, err)
		}
		if len(parts) == 1 {
			return child, nil
		}
		return selectNested(child, parts[1])
	default:
		raw := v.Get()
		if m, ok := raw.(map[string]interface{}); ok {
			val, exists := m[key]
			if !exists {
				return nil, fmt.Errorf("eval: no key %q in map", key)
			}
			tv, err := expr.Typify(val)
			if err != nil {
				return nil, err
			}
			if len(parts) == 1 {
				return tv, nil
			}
			return selectNested(tv, parts[1])
		}
		return nil, fmt.Errorf("eval: cannot traverse into %T with key %q", v, key)
	}
}

func evalRef(node *ASTNode, scope map[string]*expr.Vars) (expr.TypedValue, error) {
	boolResult := func(b bool) (expr.TypedValue, error) {
		v, err := expr.NewBoolean(b)
		return v, err
	}

	switch node.Ref {
	// ── logical ──────────────────────────────────────────────────
	case "and":
		for _, arg := range node.Args {
			b, err := EvalBool(arg, scope)
			if err != nil {
				return nil, err
			}
			if !b {
				return boolResult(false)
			}
		}
		return boolResult(true)

	case "or":
		for _, arg := range node.Args {
			b, err := EvalBool(arg, scope)
			if err != nil {
				return nil, err
			}
			if b {
				return boolResult(true)
			}
		}
		return boolResult(false)

	case "not":
		if len(node.Args) != 1 {
			return nil, fmt.Errorf("eval: not requires exactly 1 argument, got %d", len(node.Args))
		}
		b, err := EvalBool(node.Args[0], scope)
		if err != nil {
			return nil, err
		}
		return boolResult(!b)

	// ── null checks ───────────────────────────────────────────────
	case "isNull", "null":
		if len(node.Args) != 1 {
			return nil, fmt.Errorf("eval: %s requires exactly 1 argument", node.Ref)
		}
		v, err := Eval(node.Args[0], scope)
		if err != nil {
			return nil, err
		}
		return boolResult(v == nil || v.Get() == nil)

	case "isNotNull", "notNull":
		if len(node.Args) != 1 {
			return nil, fmt.Errorf("eval: %s requires exactly 1 argument", node.Ref)
		}
		v, err := Eval(node.Args[0], scope)
		if err != nil {
			return nil, err
		}
		return boolResult(v != nil && v.Get() != nil)

	// ── comparisons ───────────────────────────────────────────────
	case "eq", "ne", "lt", "gt", "lte", "gte":
		if len(node.Args) != 2 {
			return nil, fmt.Errorf("eval: %s requires exactly 2 arguments", node.Ref)
		}
		left, err := Eval(node.Args[0], scope)
		if err != nil {
			return nil, err
		}
		right, err := Eval(node.Args[1], scope)
		if err != nil {
			return nil, err
		}

		cmp, err := compareValues(left, right)
		if err != nil {
			return nil, fmt.Errorf("eval %s: %w", node.Ref, err)
		}

		var result bool
		switch node.Ref {
		case "eq":
			result = cmp == 0
		case "ne":
			result = cmp != 0
		case "lt":
			result = cmp < 0
		case "gt":
			result = cmp > 0
		case "lte":
			result = cmp <= 0
		case "gte":
			result = cmp >= 0
		}
		return boolResult(result)

	default:
		return nil, fmt.Errorf("eval: unsupported operator/ref %q", node.Ref)
	}
}

// compareValues compares two TypedValues numerically or as strings.
// Returns negative, zero, or positive.
func compareValues(a, b expr.TypedValue) (int, error) {
	var av, bv interface{}
	if a != nil {
		av = a.Get()
	}
	if b != nil {
		bv = b.Get()
	}
	if av == nil && bv == nil {
		return 0, nil
	}
	if av == nil || bv == nil {
		if av == nil {
			return -1, nil
		}
		return 1, nil
	}

	// Try numeric comparison first
	af, aerr := cast.ToFloat64E(av)
	bf, berr := cast.ToFloat64E(bv)
	if aerr == nil && berr == nil {
		switch {
		case af < bf:
			return -1, nil
		case af > bf:
			return 1, nil
		default:
			return 0, nil
		}
	}

	// Fall back to string comparison
	as, aerr := cast.ToStringE(av)
	bs, berr := cast.ToStringE(bv)
	if aerr != nil || berr != nil {
		return 0, fmt.Errorf("cannot compare values of type %T and %T", av, bv)
	}
	switch {
	case as < bs:
		return -1, nil
	case as > bs:
		return 1, nil
	default:
		return 0, nil
	}
}
