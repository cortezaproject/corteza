package expr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBuiltInFunctionNames verifies that every name in builtInFunctionNames is
// actually reserved: evaluating the bare identifier against a scope that also
// binds that name must NOT return the scope value (it resolves to the function
// instead). If a name is no longer registered, the bare reference would return
// the sentinel and this test fails, flagging a stale entry.
func TestBuiltInFunctionNames(t *testing.T) {
	const sentinel = 987654.0
	p := NewParser()

	for _, name := range BuiltInFunctionNames() {
		vars, err := NewVars(map[string]interface{}{name: Must(NewFloat(sentinel))})
		require.NoError(t, err)

		e, err := p.Parse(name)
		require.NoErrorf(t, err, "parse %q", name)

		v, err := e.Eval(context.Background(), vars)
		if err != nil {
			// function requires args / errors bare — still proves it's reserved
			continue
		}

		got, ok := v.(TypedValue)
		if ok {
			v = got.Get()
		}
		require.NotEqualf(t, sentinel, v,
			"name %q resolves to the scope variable, so it is not reserved — remove it from builtInFunctionNames", name)
	}
}

func TestIsBuiltInFunction(t *testing.T) {
	require.True(t, IsBuiltInFunction("sum"))
	require.True(t, IsBuiltInFunction("count"))
	require.False(t, IsBuiltInFunction("total"))
	require.False(t, IsBuiltInFunction("mySum"))
	require.False(t, IsBuiltInFunction(""))
}
