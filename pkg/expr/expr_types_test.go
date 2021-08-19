package expr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Example_set_kv() {
	var (
		kv = &KV{}
		p  = map[string]interface{}{
			"kv":    kv,
			"key":   "k1",
			"value": Must(NewString("v11")),
		}
	)

	eval(`set(kv, key, value)`, p)
	fmt.Printf("\nOriginal KV should be %v", kv.value)

	// output:
	// map[k1:v11]
	// Original KV should be map[]
}

func Example_merge_kv() {
	var (
		kv = &KV{}
		p  = map[string]interface{}{
			"kv": kv,
			"foo": &KV{value: map[string]string{
				"k1": "v1",
			}},
			"bar": &KV{value: map[string]string{
				"k2": "v2",
			}},
		}
	)

	eval(`merge(kv, foo, bar)`, p)
	fmt.Printf("\nOriginal KV should be %v", kv.value)

	// output:
	// map[k1:v1 k2:v2]
	// Original KV should be map[]
}

func Example_filter_kv() {
	var (
		kv = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
		p = map[string]interface{}{
			"kv":   kv,
			"key1": "k1",
			"key2": "k3",
		}
	)

	eval(`filter(kv, key1, key2)`, p)
	fmt.Printf("\nOriginal KV should be %v", kv.value)

	// output:
	// map[k1:v1]
	// Original KV should be map[k1:v1 k2:v2]
}

func Example_omit_kv() {
	var (
		kv = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		}}
		p = map[string]interface{}{
			"kv":   kv,
			"key1": "k1",
			"key2": "k3",
		}
	)

	eval(`omit(kv, key1, key2)`, p)
	fmt.Printf("\nOriginal KV should be %v", kv.value)

	// output:
	// map[k2:v2]
	// Original KV should be map[k1:v1 k2:v2 k3:v3]
}

func Example_set_kvv() {
	var (
		kvv = &KVV{}
		p   = map[string]interface{}{
			"kvv":   kvv,
			"key":   "foo",
			"value": Must(NewString("bar")),
		}
	)

	eval(`set(kvv, key, value)`, p)
	fmt.Printf("\nOriginal KVV should be %v", kvv.value)

	// output:
	// map[foo:[bar]]
	// Original KVV should be map[]
}

func Example_merge_kvv() {
	var (
		kvv = &KVV{}
		p   = map[string]interface{}{
			"kvv": kvv,
			"foo": &KVV{value: map[string][]string{
				"k1": {"v1"},
			}},
			"bar": &KVV{value: map[string][]string{
				"k1": {"v11"},
				"k2": {"v2"},
			}},
		}
	)

	eval(`merge(kvv, foo, bar)`, p)
	fmt.Printf("\nOriginal KVV should be %v", kvv.value)

	// output:
	// map[k1:[v1 v11] k2:[v2]]
	// Original KVV should be map[]
}

func Example_filter_kvv() {
	var (
		kvv = &KVV{value: map[string][]string{
			"k1": {"v1"},
			"k2": {"v2"},
		}}
		p = map[string]interface{}{
			"kv":   kvv,
			"key1": "k1",
			"key2": "k3",
		}
	)

	eval(`filter(kv, key1, key2)`, p)
	fmt.Printf("\nOriginal KVV should be %v", kvv.value)

	// output:
	// map[k1:[v1]]
	// Original KVV should be map[k1:[v1] k2:[v2]]
}

func Example_omit_kvv() {
	var (
		kvv = &KVV{value: map[string][]string{
			"k1": {"v1"},
			"k2": {"v2"},
			"k3": {"v3"},
		}}
		p = map[string]interface{}{
			"kvv":  kvv,
			"key1": "k1",
			"key2": "k3",
		}
	)

	eval(`omit(kvv, key1, key2)`, p)
	fmt.Printf("\nOriginal KVV should be %v", kvv.value)

	// output:
	// map[k2:[v2]]
	// Original KVV should be map[k1:[v1] k2:[v2] k3:[v3]]
}

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

func TestKV_Assign(t *testing.T) {
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

func TestKV_Set(t *testing.T) {
	var (
		req = require.New(t)

		kv = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
	)

	out, err := set(kv, "k1", Must(NewString("v11")))
	req.NoError(err)
	req.Equal("v11", out.(*KV).GetValue()["k1"])

	// Making sure empty KV updates without error
	kv = &KV{}
	out, err = set(kv, "k1", Must(NewString("v11")))
	req.NoError(err)
	req.Equal("v11", out.(*KV).GetValue()["k1"])
}

func TestKV_Merge(t *testing.T) {
	var (
		req = require.New(t)

		kv  = &KV{}
		foo = &KV{value: map[string]string{
			"k1": "v1",
		}}
		bar = &KV{value: map[string]string{
			"k2": "v2",
		}}
		expected = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
	)

	out, err := merge(kv, foo, bar)
	req.NoError(err)
	req.Equal(expected, out)
}

func TestKV_Clone(t *testing.T) {
	var (
		req = require.New(t)

		kv = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
		expected = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
	)

	out, err := kv.Merge()
	req.NoError(err)
	req.Equal(expected, out)
}

func TestKV_Filter(t *testing.T) {
	var (
		req = require.New(t)

		kv = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
		}}
		expected = &KV{value: map[string]string{
			"k1": "v1",
		}}
	)

	out, err := filter(kv, "k1", "k3")
	req.NoError(err)
	req.Equal(expected, out)
	req.NotEqual(kv, out)
}

func TestKV_Omit(t *testing.T) {
	var (
		req = require.New(t)

		kv = &KV{value: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		}}
		expected = &KV{value: map[string]string{
			"k2": "v2",
		}}
	)

	out, err := omit(kv, "k1", "k3")
	req.NoError(err)
	req.Equal(expected, out)
	req.NotEqual(kv, out)
}

func TestKVV_Assign(t *testing.T) {
	var (
		req = require.New(t)
		kvv KVV
	)

	req.NoError(Assign(&kvv, "foo", Must(NewString("bar"))))
	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])

	// Making sure http.Header is properly converted
	kvv = KVV{}
	req.NoError(kvv.Assign(http.Header{"foo-bar": []string{"bar"}}))
	req.Contains(kvv.value, "foo-bar")
	req.Equal([]string{"bar"}, kvv.value["foo-bar"])

	// Making sure url.Values are properly converted
	kvv = KVV{}
	req.NoError(kvv.Assign(url.Values{"foo": []string{"bar"}}))

	req.Contains(kvv.value, "foo")
	req.Equal([]string{"bar"}, kvv.value["foo"])

	kvv = KVV{}
	req.NoError(Assign(&kvv, "deep", Must(NewString("bar"))))
	req.NoError(Assign(&kvv, "deep[0]", Must(NewString("bar"))))
	req.NoError(Assign(&kvv, "deep[]", Must(NewString("baz"))))
	req.NoError(Assign(&kvv, "deep[]", Must(NewString("bar"))))
	req.NoError(Assign(&kvv, "deep[3]", Must(NewString("baz"))))
	req.NoError(Assign(&kvv, "deep[3]", Must(NewString("b4z"))))
	req.Contains(kvv.value, "deep")
	req.Equal([]string{"bar", "baz", "bar", "b4z"}, kvv.value["deep"])
}

func TestKVV_Set(t *testing.T) {
	var (
		req = require.New(t)
		kvv = &KVV{value: map[string][]string{
			"k1": {"v1"},
			"k2": {"v2"},
		}}
	)

	out, err := set(kvv, "k1", Must(NewString("v11")))
	req.NoError(err)
	req.Equal([]string{"v11"}, out.(*KVV).GetValue()["k1"])

	// Making sure empty KV updates without error
	kvv = &KVV{}
	out, err = set(kvv, "foo", Must(NewString("bar")))
	req.NoError(err)
	req.Equal([]string{"bar"}, out.(*KVV).GetValue()["foo"])
}

func TestKVV_Merge(t *testing.T) {
	var (
		req = require.New(t)

		kvv *KVV
		foo = KVV{value: map[string][]string{
			"k1": {"v1"},
		}}
		bar = KVV{value: map[string][]string{
			"k1": {"v11"},
			"k2": {"v2"},
		}}
		expected = &KVV{value: map[string][]string{
			"k1": {"v1", "v11"},
			"k2": {"v2"},
		}}
	)

	out, err := merge(kvv, &foo, &bar)
	req.NoError(err)
	req.Equal(expected, out)
}

func TestKVV_Clone(t *testing.T) {
	var (
		req = require.New(t)

		kvv = KVV{value: map[string][]string{
			"k1": {"v1", "v11"},
			"k2": {"v2"},
		}}
		expected = &KVV{value: map[string][]string{
			"k1": {"v1", "v11"},
			"k2": {"v2"},
		}}
	)

	out, err := kvv.Merge()
	req.NoError(err)
	req.Equal(expected, out)
}

func TestKVV_Filter(t *testing.T) {
	var (
		req = require.New(t)

		kvv = &KVV{value: map[string][]string{
			"k1": {"v1"},
			"k2": {"v2"},
		}}
		expected = &KVV{value: map[string][]string{
			"k1": {"v1"},
		}}
	)

	out, err := filter(kvv, "k1", "k3")
	req.NoError(err)
	req.Equal(expected, out)
}

func TestKVV_Omit(t *testing.T) {
	var (
		req = require.New(t)

		kvv = &KVV{value: map[string][]string{
			"k1": {"v1"},
			"k2": {"v2"},
			"k3": {"v3"},
		}}
		expected = &KVV{value: map[string][]string{
			"k2": {"v2"},
		}}
	)

	out, err := omit(kvv, "k1", "k3")
	req.NoError(err)
	req.Equal(expected, out)
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
		"strings": &Array{value: arr},
		"iface":   Must(NewString("typed")),
		"typed":   Must(NewString("typed")),
		"values":  &Array{value: arr},
	})

	req.NoError(err)
	req.NoError(vars.Decode(&foo))
	req.Len(foo.Strings, 2)
	req.Len(foo.Values, 2)
}
