package ast

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/expr"
)

// globalScope wraps a single *expr.Vars into the "global" scope map.
func globalScope(v *expr.Vars) map[string]*expr.Vars {
	return map[string]*expr.Vars{"global": v}
}

func TestEval(t *testing.T) {
	t.Run("NilNode", func(t *testing.T) {
		v, err := Eval(nil, globalScope(&expr.Vars{}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := v.(*expr.Boolean); !ok {
			t.Fatalf("expected *expr.Boolean, got %T", v)
		}
	})

	t.Run("Symbol/Found", func(t *testing.T) {
		vars := makeVars("x", expr.Must(expr.NewInteger(42)))
		v, err := Eval(sym("x"), globalScope(vars))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Get().(int64) != 42 {
			t.Fatalf("expected 42, got %v", v.Get())
		}
	})

	t.Run("Symbol/Missing", func(t *testing.T) {
		_, err := Eval(sym("missing"), globalScope(&expr.Vars{}))
		if err == nil {
			t.Fatal("expected error for missing symbol")
		}
	})

	t.Run("Literal/String", func(t *testing.T) {
		v, err := Eval(strVal("hello"), globalScope(&expr.Vars{}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Get().(string) != "hello" {
			t.Fatalf("expected 'hello', got %v", v.Get())
		}
	})

	t.Run("Literal/Integer", func(t *testing.T) {
		v, err := Eval(intVal(7), globalScope(&expr.Vars{}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Get().(int64) != 7 {
			t.Fatalf("expected 7, got %v", v.Get())
		}
	})

	t.Run("And/AllTrue", func(t *testing.T) {
		n := ref("and", boolVal(true), boolVal(true))
		if !mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("And/OneFalse", func(t *testing.T) {
		n := ref("and", boolVal(true), boolVal(false), boolVal(true))
		if mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected false")
		}
	})

	t.Run("And/Empty", func(t *testing.T) {
		// zero args → vacuously true
		n := ref("and")
		if !mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected true for empty and")
		}
	})

	t.Run("And/ShortCircuit", func(t *testing.T) {
		// second arg is a bad symbol; should not be reached
		n := ref("and", boolVal(false), sym("nope"))
		b, err := EvalBool(n, globalScope(&expr.Vars{}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b {
			t.Fatal("expected false")
		}
	})

	t.Run("Or/OneTrue", func(t *testing.T) {
		n := ref("or", boolVal(false), boolVal(true))
		if !mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("Or/AllFalse", func(t *testing.T) {
		n := ref("or", boolVal(false), boolVal(false))
		if mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected false")
		}
	})

	t.Run("Or/Empty", func(t *testing.T) {
		// zero args → vacuously false
		n := ref("or")
		if mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected false for empty or")
		}
	})

	t.Run("Or/ShortCircuit", func(t *testing.T) {
		// second arg is a bad symbol; should not be reached
		n := ref("or", boolVal(true), sym("nope"))
		b, err := EvalBool(n, globalScope(&expr.Vars{}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !b {
			t.Fatal("expected true")
		}
	})

	t.Run("Not/True", func(t *testing.T) {
		if mustBool(t, ref("not", boolVal(true)), globalScope(&expr.Vars{})) {
			t.Fatal("expected false")
		}
	})

	t.Run("Not/False", func(t *testing.T) {
		if !mustBool(t, ref("not", boolVal(false)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("Not/WrongArity", func(t *testing.T) {
		_, err := EvalBool(ref("not", boolVal(true), boolVal(false)), globalScope(&expr.Vars{}))
		if err == nil {
			t.Fatal("expected error for wrong arity")
		}
	})

	t.Run("Not/NoArgs", func(t *testing.T) {
		_, err := EvalBool(ref("not"), globalScope(&expr.Vars{}))
		if err == nil {
			t.Fatal("expected error for zero args")
		}
	})

	t.Run("IsNull/NonNilValue", func(t *testing.T) {
		if mustBool(t, ref("isNull", strVal("x")), globalScope(&expr.Vars{})) {
			t.Fatal("expected false")
		}
	})

	t.Run("IsNull/Alias", func(t *testing.T) {
		// "null" is an alias for "isNull"
		if mustBool(t, ref("null", strVal("x")), globalScope(&expr.Vars{})) {
			t.Fatal("expected false")
		}
	})

	t.Run("IsNull/WrongArity", func(t *testing.T) {
		_, err := EvalBool(ref("isNull", strVal("a"), strVal("b")), globalScope(&expr.Vars{}))
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("IsNotNull/NonNilValue", func(t *testing.T) {
		if !mustBool(t, ref("isNotNull", strVal("x")), globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("IsNotNull/Alias", func(t *testing.T) {
		if !mustBool(t, ref("notNull", strVal("x")), globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("IsNotNull/WrongArity", func(t *testing.T) {
		_, err := EvalBool(ref("isNotNull", strVal("a"), strVal("b")), globalScope(&expr.Vars{}))
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Eq/True", func(t *testing.T) {
		if !mustBool(t, ref("eq", intVal(5), intVal(5)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("Eq/False", func(t *testing.T) {
		if mustBool(t, ref("eq", intVal(5), intVal(6)), globalScope(&expr.Vars{})) {
			t.Fatal("expected false")
		}
	})

	t.Run("Eq/Strings", func(t *testing.T) {
		if !mustBool(t, ref("eq", strVal("abc"), strVal("abc")), globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("Ne/True", func(t *testing.T) {
		if !mustBool(t, ref("ne", intVal(1), intVal(2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("Ne/False", func(t *testing.T) {
		if mustBool(t, ref("ne", intVal(3), intVal(3)), globalScope(&expr.Vars{})) {
			t.Fatal("expected false")
		}
	})

	t.Run("Comparison/WrongArity", func(t *testing.T) {
		for _, op := range []string{"eq", "ne", "lt", "gt", "lte", "gte"} {
			_, err := EvalBool(ref(op, intVal(1)), globalScope(&expr.Vars{}))
			if err == nil {
				t.Fatalf("expected error for %s with 1 arg", op)
			}
		}
	})

	t.Run("Lt", func(t *testing.T) {
		cases := []struct {
			a, b int
			want bool
		}{
			{1, 2, true},
			{2, 2, false},
			{3, 2, false},
		}
		for _, c := range cases {
			got := mustBool(t, ref("lt", intVal(c.a), intVal(c.b)), globalScope(&expr.Vars{}))
			if got != c.want {
				t.Errorf("lt(%d,%d): got %v, want %v", c.a, c.b, got, c.want)
			}
		}
	})

	t.Run("Gt", func(t *testing.T) {
		cases := []struct {
			a, b int
			want bool
		}{
			{3, 2, true},
			{2, 2, false},
			{1, 2, false},
		}
		for _, c := range cases {
			got := mustBool(t, ref("gt", intVal(c.a), intVal(c.b)), globalScope(&expr.Vars{}))
			if got != c.want {
				t.Errorf("gt(%d,%d): got %v, want %v", c.a, c.b, got, c.want)
			}
		}
	})

	t.Run("Lte", func(t *testing.T) {
		if !mustBool(t, ref("lte", intVal(2), intVal(2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true for lte equal")
		}
		if !mustBool(t, ref("lte", intVal(1), intVal(2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true for lte less")
		}
		if mustBool(t, ref("lte", intVal(3), intVal(2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected false for lte greater")
		}
	})

	t.Run("Gte", func(t *testing.T) {
		if !mustBool(t, ref("gte", intVal(2), intVal(2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true for gte equal")
		}
		if !mustBool(t, ref("gte", intVal(3), intVal(2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true for gte greater")
		}
		if mustBool(t, ref("gte", intVal(1), intVal(2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected false for gte less")
		}
	})

	t.Run("Float/Comparison", func(t *testing.T) {
		if !mustBool(t, ref("lt", floatVal(1.1), floatVal(2.2)), globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("StringComparison/Lexical", func(t *testing.T) {
		if !mustBool(t, ref("lt", strVal("apple"), strVal("banana")), globalScope(&expr.Vars{})) {
			t.Fatal("expected 'apple' < 'banana'")
		}
	})

	t.Run("Eq/WithVar", func(t *testing.T) {
		vars := makeVars("age", expr.Must(expr.NewInteger(30)))
		if !mustBool(t, ref("eq", sym("age"), intVal(30)), globalScope(vars)) {
			t.Fatal("expected true")
		}
	})

	t.Run("Nested/AndOr", func(t *testing.T) {
		// (true OR false) AND (false OR true) → true
		n := ref("and",
			ref("or", boolVal(true), boolVal(false)),
			ref("or", boolVal(false), boolVal(true)),
		)
		if !mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("Nested/NotAnd", func(t *testing.T) {
		// NOT (true AND false) → true
		n := ref("not", ref("and", boolVal(true), boolVal(false)))
		if !mustBool(t, n, globalScope(&expr.Vars{})) {
			t.Fatal("expected true")
		}
	})

	t.Run("UnsupportedRef", func(t *testing.T) {
		_, err := Eval(ref("bogus"), globalScope(&expr.Vars{}))
		if err == nil {
			t.Fatal("expected error for unsupported ref")
		}
	})

	t.Run("CompareValues/Numeric", func(t *testing.T) {
		a := expr.Must(expr.NewFloat(1.5))
		b := expr.Must(expr.NewFloat(2.5))
		cmp, err := compareValues(a, b)
		if err != nil {
			t.Fatal(err)
		}
		if cmp >= 0 {
			t.Fatalf("expected negative, got %d", cmp)
		}
	})

	t.Run("CompareValues/StringFallback", func(t *testing.T) {
		a := expr.Must(expr.NewString("a"))
		b := expr.Must(expr.NewString("z"))
		cmp, err := compareValues(a, b)
		if err != nil {
			t.Fatal(err)
		}
		if cmp >= 0 {
			t.Fatalf("expected negative, got %d", cmp)
		}
	})

	t.Run("CompareValues/Equal", func(t *testing.T) {
		a := expr.Must(expr.NewInteger(10))
		b := expr.Must(expr.NewInteger(10))
		cmp, err := compareValues(a, b)
		if err != nil {
			t.Fatal(err)
		}
		if cmp != 0 {
			t.Fatalf("expected 0, got %d", cmp)
		}
	})
}

func TestEvalScope(t *testing.T) {
	t.Run("GlobalScopeDefault", func(t *testing.T) {
		// A symbol node with no Meta["scope"] resolves from "global"
		vars := makeVars("x", expr.Must(expr.NewInteger(7)))
		v, err := Eval(sym("x"), globalScope(vars))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Get().(int64) != 7 {
			t.Fatalf("expected 7, got %v", v.Get())
		}
	})

	t.Run("NamedScope", func(t *testing.T) {
		// A symbol node with Meta["scope"]="step1" resolves from that scope.
		step1Vars := makeVars("result", expr.Must(expr.NewBoolean(true)))
		scope := map[string]*expr.Vars{
			"global": makeVars(),
			"step1":  step1Vars,
		}
		node := &ASTNode{
			Symbol: "result",
			Meta:   map[string]any{"scope": "step1"},
		}
		v, err := Eval(node, scope)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Get().(bool) != true {
			t.Fatalf("expected true, got %v", v.Get())
		}
	})

	t.Run("MissingScope", func(t *testing.T) {
		// Referencing a scope not present in the map returns an error.
		node := &ASTNode{
			Symbol: "x",
			Meta:   map[string]any{"scope": "nonexistent"},
		}
		_, err := Eval(node, globalScope(&expr.Vars{}))
		if err == nil {
			t.Fatal("expected error for missing scope")
		}
	})

	t.Run("TwoSiblingsDifferentScopes", func(t *testing.T) {
		// Two Symbol leaves in the same expression tree, each with a different
		// scope, both resolve correctly in a single EvalBool call.
		scopeA := makeVars("val", expr.Must(expr.NewInteger(10)))
		scopeB := makeVars("val", expr.Must(expr.NewInteger(10)))

		nodeA := &ASTNode{Symbol: "val", Meta: map[string]any{"scope": "a"}}
		nodeB := &ASTNode{Symbol: "val", Meta: map[string]any{"scope": "b"}}
		n := ref("eq", nodeA, nodeB)

		scope := map[string]*expr.Vars{
			"global": makeVars(),
			"a":      scopeA,
			"b":      scopeB,
		}

		b, err := EvalBool(n, scope)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !b {
			t.Fatal("expected true (both scopes have val=10)")
		}
	})
}

func strVal(s string) *ASTNode {
	return &ASTNode{Value: MakeValueOf("String", s)}
}

func intVal(i int) *ASTNode {
	return &ASTNode{Value: MakeValueOf("Integer", i)}
}

func floatVal(f float64) *ASTNode {
	return &ASTNode{Value: MakeValueOf("Float", f)}
}

func boolVal(b bool) *ASTNode {
	return &ASTNode{Value: MakeValueOf("Boolean", b)}
}

func sym(name string) *ASTNode {
	return &ASTNode{Symbol: name}
}

func ref(op string, args ...*ASTNode) *ASTNode {
	return &ASTNode{Ref: op, Args: args}
}

func makeVars(kvs ...interface{}) *expr.Vars {
	v := &expr.Vars{}
	for i := 0; i+1 < len(kvs); i += 2 {
		key := kvs[i].(string)
		val := kvs[i+1].(expr.TypedValue)
		_ = v.Set(key, val)
	}
	return v
}

func mustBool(t *testing.T, node *ASTNode, scope map[string]*expr.Vars) bool {
	t.Helper()
	b, err := EvalBool(node, scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return b
}
