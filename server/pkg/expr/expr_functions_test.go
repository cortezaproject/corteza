package expr

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type (
	exprFunctionsDoc struct {
		Groups []struct {
			Name      string
			Functions []exprFunctionDoc
		}
	}

	exprFunctionDoc struct {
		Name    string
		Returns string
		Example struct {
			Expr   string
			Scope  string
			Result string
			Note   string
		}
	}
)

func loadExprFunctionsDoc(t *testing.T) []exprFunctionDoc {
	raw, err := os.ReadFile("expr_functions.yaml")
	require.NoError(t, err)

	var doc exprFunctionsDoc
	require.NoError(t, yaml.Unmarshal(raw, &doc))

	var out []exprFunctionDoc
	for _, g := range doc.Groups {
		out = append(out, g.Functions...)
	}

	return out
}

// renderExampleResult is the JSON form an example's result is documented in:
// typed values unwrapped, times as RFC 3339, durations as Go duration strings.
func renderExampleResult(v interface{}) interface{} {
	v = UntypedValue(v)

	switch c := v.(type) {
	case time.Time:
		return c.Format(time.RFC3339Nano)
	case *time.Time:
		return c.Format(time.RFC3339Nano)
	case time.Duration:
		return c.String()
	case map[string]TypedValue:
		out := make(map[string]interface{}, len(c))
		for k, tv := range c {
			out[k] = renderExampleResult(tv)
		}
		return out
	case []TypedValue:
		out := make([]interface{}, len(c))
		for i, tv := range c {
			out[i] = renderExampleResult(tv)
		}
		return out
	}

	return v
}

// exampleReturnKinds holds each documented return type to the Go values an
// evaluation of that type yields.
var exampleReturnKinds = map[string]func(interface{}) bool{
	"String":  func(v interface{}) bool { _, ok := v.(string); return ok },
	"Boolean": func(v interface{}) bool { _, ok := v.(bool); return ok },
	"Float":   func(v interface{}) bool { _, ok := v.(float64); return ok },
	"Integer": func(v interface{}) bool {
		switch v.(type) {
		case int, int64:
			return true
		}
		return false
	},
	"Array": func(v interface{}) bool {
		return v != nil && reflect.TypeOf(v).Kind() == reflect.Slice
	},
	"DateTime": func(v interface{}) bool {
		switch v.(type) {
		case time.Time, *time.Time:
			return true
		}
		return false
	},
	"Duration": func(v interface{}) bool { _, ok := v.(time.Duration); return ok },
	"Vars": func(v interface{}) bool {
		switch v.(type) {
		case *Vars, *KV, *KVV:
			return true
		}
		return false
	},
	"Any": func(interface{}) bool { return true },
}

// TestExprFunctionExamples evaluates every documented example and holds it to
// its documented result and return type.
func TestExprFunctionExamples(t *testing.T) {
	p := NewParser()

	for _, fn := range loadExprFunctionsDoc(t) {
		fn := fn
		t.Run(fn.Name, func(t *testing.T) {
			ex := fn.Example
			require.NotEmptyf(t, ex.Expr, "%s has no example", fn.Name)
			require.Truef(t, ex.Result != "" || ex.Note != "",
				"%s example needs a result, or a note when the result varies", fn.Name)

			isKind, known := exampleReturnKinds[fn.Returns]
			require.Truef(t, known, "%s returns unknown type %q", fn.Name, fn.Returns)

			scope := map[string]interface{}{}
			if ex.Scope != "" {
				require.NoError(t, json.Unmarshal([]byte(ex.Scope), &scope))
			}

			vars, err := NewVars(scope)
			require.NoError(t, err)

			ev, err := p.Parse(ex.Expr)
			require.NoErrorf(t, err, "parse %q", ex.Expr)

			v, err := ev.Eval(context.Background(), vars)
			require.NoErrorf(t, err, "eval %q", ex.Expr)

			require.Truef(t, isKind(UntypedValue(v)) || isKind(v),
				"%q returned %T, documented as %s", ex.Expr, v, fn.Returns)

			if ex.Result == "" {
				return
			}

			got, err := json.Marshal(renderExampleResult(v))
			require.NoError(t, err)
			require.JSONEqf(t, ex.Result, string(got), "%q", ex.Expr)
		})
	}
}

// TestExprFunctionsDocumented holds expr_functions.yaml to the functions this
// package registers, and func_names.gen.go to expr_functions.yaml.
func TestExprFunctionsDocumented(t *testing.T) {
	reg := regexp.MustCompile(`(?:gval\.Function|gvalFunc)\(\s*"([^"]+)"`)

	sources, err := filepath.Glob("*.go")
	require.NoError(t, err)

	registered := map[string]bool{}
	for _, src := range sources {
		if strings.HasSuffix(src, "_test.go") {
			continue
		}

		raw, err := os.ReadFile(src)
		require.NoError(t, err)

		for _, m := range reg.FindAllStringSubmatch(string(raw), -1) {
			registered[m[1]] = true
		}
	}

	require.NotEmpty(t, registered)

	documented := map[string]bool{}
	for _, fn := range loadExprFunctionsDoc(t) {
		require.Falsef(t, documented[fn.Name], "%s is documented twice", fn.Name)
		documented[fn.Name] = true
	}

	var missing, stale []string
	for n := range registered {
		if !documented[n] {
			missing = append(missing, n)
		}
	}
	for n := range documented {
		if !registered[n] {
			stale = append(stale, n)
		}
	}
	sort.Strings(missing)
	sort.Strings(stale)

	require.Emptyf(t, missing, "registered but not in expr_functions.yaml")
	require.Emptyf(t, stale, "in expr_functions.yaml but not registered")

	names := make([]string, 0, len(documented))
	for n := range documented {
		names = append(names, n)
	}
	sort.Strings(names)
	require.Equal(t, names, BuiltInFunctionNames(),
		"func_names.gen.go is out of date — run make codegen-legacy")
}
