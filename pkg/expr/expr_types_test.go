package expr

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestTypedValueOperations(t *testing.T) {
	scope, _ := NewVars(map[string]interface{}{
		"xUint":   uint(1),
		"xInt":    1,
		"xBoolT":  true,
		"xBoolF":  false,
		"xString": "foo",
	})

	tcc := []struct {
		expects interface{}
		expr    string
	}{
		// sanity check
		{true, "xBoolT"},
		{false, "xBoolF"},

		// uint ops
		{true, "xUint == 1"},
		{uint(1), "xUint"},

		// uint ops
		{true, "xInt == 1"},
		{int(1), "xInt"},

		// string ops
		{true, `xString == "foo"`},
		{"foo", "xString"},
	}

	for _, tc := range tcc {
		t.Run(tc.expr, func(t *testing.T) {
			var (
				req       = require.New(t)
				parser    = NewParser()
				eval, err = parser.Parse(tc.expr)
			)

			req.NoError(err)
			result, err := eval.Eval(context.Background(), scope)
			req.NoError(err)
			req.Equal(tc.expects, UntypedValue(result))
		})
	}

}

func TestKV_Set(t *testing.T) {
	var (
		req = require.New(t)

		vars = KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
	)

	req.NoError(Assign(&vars, "k1", Must(NewString("v11"))))
	req.Equal("v11", vars.value["k1"])
	req.Equal("v2", vars.value["k2"])

}

func TestKVV_Set(t *testing.T) {
	var (
		req = require.New(t)
		kvv KVV
	)

	req.NoError(Assign(&kvv, "foo", Must(NewString("bar"))))
	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])

	// Making sure http.Header is properly converted
	kvv = KVV{}
	req.NoError(kvv.Assign(http.Header{"foo": []string{"bar"}}))
	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])

	// Making sure url.Values are properly converted
	kvv = KVV{}
	req.NoError(kvv.Assign(url.Values{"foo": []string{"bar"}}))
	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])
}

func TestCastEmptyString(t *testing.T) {
	var (
		req = require.New(t)
	)
	{
		f, err := CastToFloat("")
		req.NoError(err)
		req.Zero(f)
	}
	{
		i, err := CastToInteger("")
		req.NoError(err)
		req.Zero(i)
	}
	{
		u, err := CastToUnsignedInteger("")
		req.NoError(err)
		req.Zero(u)
	}
	{
		u, err := CastToDuration("")
		req.NoError(err)
		req.Zero(u)
	}
	{
		u, err := CastToBoolean("")
		req.NoError(err)
		req.False(u)
	}
}

func TestCastToArray(t *testing.T) {
	var (
		req = require.New(t)
	)

	arr, err := CastToArray([]interface{}{
		Must(NewString("abc")),
		Must(NewString("123")),
	})

	req.NoError(err)
	req.Len(arr, 2)
}

func TestArrayDecode(t *testing.T) {
	var (
		req = require.New(t)

		foo = struct {
			Typed   TypedValue
			Iface   interface{}
			Strings []string
			Values  []TypedValue
		}{}
	)

	arr, err := CastToArray([]interface{}{
		Must(NewString("abc")),
		Must(NewString("123")),
	})

	req.NoError(err)

	vars, err := NewVars(map[string]interface{}{
		"strings": &Array{arr},
		"iface":   Must(NewString("typed")),
		"typed":   Must(NewString("typed")),
		"values":  &Array{arr},
	})

	req.NoError(err)
	req.NoError(vars.Decode(&foo))
	req.Len(foo.Strings, 2)
	req.Len(foo.Values, 2)
}
