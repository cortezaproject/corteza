package expr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBuiltInFunctionNames holds builtInFunctionNames to what the language
// actually registers, from both sides: a bare name reads a scope variable of
// that name, and `name()` still reaches the function.
func TestBuiltInFunctionNames(t *testing.T) {
	const sentinel = 987654.0
	p := NewParser()
	ctx := context.Background()

	unwrap := func(v interface{}) interface{} {
		if tv, is := v.(TypedValue); is {
			return tv.Get()
		}
		return v
	}

	for _, name := range BuiltInFunctionNames() {
		vars, err := NewVars(map[string]interface{}{name: Must(NewFloat(sentinel))})
		require.NoError(t, err)

		bare, err := p.Parse(name)
		require.NoErrorf(t, err, "parse %q", name)

		v, err := bare.Eval(ctx, vars)
		require.NoErrorf(t, err, "eval %q", name)
		require.Equalf(t, sentinel, unwrap(v),
			"a bare %q must read the scope variable", name)

		call, err := p.Parse(name + "()")
		require.NoErrorf(t, err, "parse %q()", name)

		// Whatever the function makes of no arguments, it must be the one
		// reached. Resolving the variable instead fails with gval's
		// "could not call" — the scope value is a number, not a func.
		v, err = call.Eval(ctx, vars)
		if err != nil {
			require.NotContainsf(t, err.Error(), "could not call",
				"%q() resolved to the scope variable, so the name is not registered — remove it from builtInFunctionNames", name)
			continue
		}

		require.NotEqualf(t, sentinel, unwrap(v),
			"%q() must call the function, not read the variable", name)
	}
}

// TestFunctionNameShadowing is the module-field case: a field named after a
// function is readable, and the function is still callable beside it.
func TestFunctionNameShadowing(t *testing.T) {
	p := NewParser()
	ctx := context.Background()

	vars, err := NewVars(map[string]interface{}{
		"split": Must(NewString("a b")),
		"sum":   Must(NewFloat(3)),
	})
	require.NoError(t, err)

	eval := func(e string) interface{} {
		ev, err := p.Parse(e)
		require.NoErrorf(t, err, "parse %q", e)

		v, err := ev.Eval(ctx, vars)
		require.NoErrorf(t, err, "eval %q", e)

		if tv, is := v.(TypedValue); is {
			return tv.Get()
		}
		return v
	}

	require.Equal(t, "a b", eval(`split`))
	require.Equal(t, []string{"a", "b"}, eval(`split(split, " ")`))
	require.Equal(t, 4.0, eval(`sum + 1`))
	require.Equal(t, 5.0, eval(`sum(2, 3)`))
	require.Equal(t, "a b!", eval(`split + "!"`))
}
