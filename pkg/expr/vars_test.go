package expr

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

// extract typed-value

func Example_set_vars() {
	var (
		p = map[string]interface{}{
			"vars":  &Vars{},
			"key":   "foo",
			"value": &String{value: "foo"},
		}
	)

	eval(`toJSON(set(vars, key, value))`, p)

	// output:
	// {"foo":{"@value":"foo","@type":"String"}}
}

func Example_merge_vars() {
	var (
		p = map[string]interface{}{
			"vars": &Vars{},
			"foo": &Vars{value: map[string]TypedValue{
				"k1": &String{value: "v1"},
			}},
			"bar": &Vars{value: map[string]TypedValue{
				"k2": &String{value: "v2"},
			}},
		}
	)

	eval(`toJSON(merge(foo, bar))`, p)

	// output:
	// {"k1":{"@value":"v1","@type":"String"},"k2":{"@value":"v2","@type":"String"}}
}

func Example_filter_vars() {
	var (
		p = map[string]interface{}{
			"vars": &Vars{value: map[string]TypedValue{
				"k1": &String{value: "v1"},
				"k5": &String{value: "v5"},
			}},
			"key1": "k1",
			"key2": "k3",
			"key3": "k5",
		}
	)

	eval(`toJSON(filter(vars, key1, key2, key3))`, p)

	// output:
	// {"k1":{"@value":"v1","@type":"String"},"k5":{"@value":"v5","@type":"String"}}
}

func Example_omit_vars() {
	var (
		p = map[string]interface{}{
			"vars": &Vars{value: map[string]TypedValue{
				"k1": &String{value: "v1"},
				"k2": &String{value: "v2"},
				"k3": &String{value: "v3"},
			}},
			"key1": "k1",
			"key2": "k3",
		}
	)

	eval(`toJSON(omit(vars, key1, key2))`, p)

	// output:
	// {"k2":{"@value":"v2","@type":"String"}}
}

func TestVars_Decode(t *testing.T) {
	t.Run("mix", func(t *testing.T) {
		var (
			req = require.New(t)

			dst = &struct {
				Int        int64
				Uint64     uint64
				String     string `var:"STRING"`
				RawString  string `var:"rawString"`
				Bool       bool
				Unexisting byte
			}{}

			vars, _ = NewVars(map[string]interface{}{
				"int":     Must(NewInteger(42)),
				"STRING":  Must(NewString("foo")),
				"bool":    Must(NewBoolean(true)),
				"missing": Must(NewBoolean(true)),
			})
		)

		req.NoError(vars.Decode(dst))
		req.Equal(int64(42), dst.Int)
		req.Equal("foo", dst.String)
		req.Equal(true, dst.Bool)
		req.Empty(dst.Unexisting)
	})

	t.Run("interfaces", func(t *testing.T) {
		var (
			req = require.New(t)

			dst = &struct {
				IString  interface{} `var:"iString"`
				IInteger interface{} `var:"iInteger"`
				IBool    interface{} `var:"iBool"`
			}{}

			vars, _ = NewVars(map[string]interface{}{
				"iString":  Must(NewString("foo")),
				"iInteger": Must(NewInteger(42)),
				"iBool":    Must(NewBoolean(true)),
			})
		)

		req.NoError(vars.Decode(dst))
	})

	t.Run("vars-in-vars", func(t *testing.T) {
		var (
			req = require.New(t)

			dst = &struct {
				Vars *Vars `var:"vars"`
			}{}

			vars, _ = NewVars(map[string]interface{}{
				"vars": map[string]interface{}{"foo": Must(NewString("bar"))},
			})
		)

		req.NoError(vars.Decode(dst))
	})

	t.Run("int-uint", func(t *testing.T) {
		var (
			req = require.New(t)

			dst = &struct {
				Int    int64
				Uint64 uint64
			}{}

			vars, _ = NewVars(map[string]TypedValue{
				"uint64": Must(NewAny("42")),
				"int":    Must(NewAny("42")),
			})
		)

		dst.Uint64 = 0
		dst.Int = 0

		req.NoError(vars.Decode(dst))
		req.Equal(uint64(42), dst.Uint64)
		req.Equal(int64(42), dst.Int)
	})

}

func TestVars_Assign(t *testing.T) {
	var (
		req  = require.New(t)
		vars = &Vars{}
	)

	req.NoError(Assign(vars, "foo", &String{value: "foo"}))
	req.NoError(Assign(vars, "vars", &Vars{}))
	req.NoError(Assign(vars, "vars.foo", &String{value: "foo"}))

	fmt.Println("Vars: ", vars)
}

func TestVars_Set(t *testing.T) {
	var (
		req = require.New(t)

		vars = &Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
			"k2": &String{value: "v2"},
		}}
	)

	out, err := set(vars, "k1", &String{value: "v11"})
	fmt.Println("Out: ", out)
	req.NoError(err)
	req.Equal(&String{value: "v11"}, out.(*Vars).GetValue()["k1"])

	// Making sure empty vars updates without error
	vars = &Vars{}
	out, err = set(vars, "k2", &Array{value: []TypedValue{&String{value: "foo"}, &String{value: "bar"}}})
	req.NoError(err)
	req.Equal(&Array{value: []TypedValue{&String{value: "foo"}, &String{value: "bar"}}}, out.(*Vars).GetValue()["k2"])
}

func TestVars_MergeVars(t *testing.T) {
	var (
		req = require.New(t)

		vars Vars
		foo  = Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
		}}
		bar = Vars{value: map[string]TypedValue{
			"k2": &String{value: "v2"},
		}}
		expected = &Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
			"k2": &String{value: "v2"},
		}}
	)

	fmt.Println("vars: ", vars)
	out := vars.MustMerge(&foo, &bar)
	fmt.Println("vars: ", out)
	req.Equal(expected, out)
}

func TestVars_Merge(t *testing.T) {
	var (
		req = require.New(t)

		vars = &Vars{}
		foo  = &Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
		}}
		bar = &Vars{value: map[string]TypedValue{
			"k2": &String{value: "v2"},
		}}
		expected = &Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
			"k2": &String{value: "v2"},
		}}
	)

	out, err := merge(vars, foo, bar)
	req.NoError(err)
	req.Equal(expected, out)
}

func TestVars_Filter(t *testing.T) {
	var (
		req = require.New(t)

		vars = &Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
			"k2": &String{value: "v2"},
		}}
		expected = &Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
		}}
	)

	out, err := filter(vars, "k1", "k3")
	req.NoError(err)
	req.Equal(expected, out)
}

func TestVars_Omit(t *testing.T) {
	var (
		req = require.New(t)

		vars = Vars{value: map[string]TypedValue{
			"k1": &String{value: "v1"},
			"k2": &String{value: "v2"},
			"k3": &String{value: "v3"},
		}}
		expected = &Vars{value: map[string]TypedValue{
			"k2": &String{value: "v2"},
		}}
	)

	out, err := omit(&vars, "k1", "k3")
	req.NoError(err)
	req.Equal(expected, out)
}

func TestVars_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		json string
		vars map[string]interface{}
	}{
		{"empty", "", make(map[string]interface{})},
		{"object", "{}", make(map[string]interface{})},
		{"string", `{"a":{"@value":"b"}}`, map[string]interface{}{"a": &Unresolved{value: "b"}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				r           = require.New(t)
				unmarshaled = &Vars{}
				aux, _      = NewVars(c.vars)
			)

			r.NoError(unmarshaled.UnmarshalJSON([]byte(c.json)))
			r.Equal(aux, unmarshaled)
		})
	}
}

func TestVars_MarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		json string
		vars map[string]interface{}
	}{
		{"empty", "{}", nil},
		{"string", `{"a":{"@value":"b","@type":"String"}}`, map[string]interface{}{"a": &String{value: "b"}}},
		{"array",
			`{"arr":{"@value":[{"@value":"foo","@type":"String"},{"@value":"bar","@type":"String"}],"@type":"Array"}}`,
			map[string]interface{}{"arr": &Array{value: []TypedValue{&String{value: "foo"}, &String{value: "bar"}}}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				r      = require.New(t)
				aux, _ = NewVars(c.vars)
			)

			j, err := json.Marshal(aux)
			r.NoError(err)
			r.Equal(c.json, string(j))
		})
	}
}
