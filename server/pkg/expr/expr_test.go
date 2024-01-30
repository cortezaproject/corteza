package expr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVars(t *testing.T) {
	var (
		req = require.New(t)

		vars = Must(Typify(map[string]interface{}{
			"int": Must(NewInteger(42)),
			"sub": map[string]interface{}{
				"foo": Must(NewString("foo")),
			},
			"three": map[string]interface{}{
				"two": map[string]interface{}{
					"one": map[string]interface{}{
						"go": Must(NewString("!")),
					},
				},
			},
		}))

		tv = func(in interface{}) TypedValue {
			switch cnv := in.(type) {
			case int:
				return Must(NewInteger(cnv))
			case string:
				return Must(NewString(cnv))
			}

			return Must(NewAny(in))
		}
	)

	req.NoError(Assign(vars, "int", tv(123)))
	req.Equal(int64(123), Must(Select(vars, "int")).Get().(int64))

	req.NoError(Assign(vars, "sub.foo", tv("bar")))
	req.Equal("bar", Must(Select(vars, "sub.foo")).Get().(string))

	req.NoError(Assign(vars, "kv", &KV{}))
	req.NoError(Assign(vars, "kv.foo", tv("bar")))
	req.Equal("bar", Must(Select(vars, "kv.foo")).Get().(string))

	req.NoError(Assign(vars, "three.two.one.go", tv("!!!")))
	req.Equal("!!!", Must(Select(vars, "three.two.one.go")).Get().(string))
}

func TestAssign(t *testing.T) {
	base := &Vars{
		value: map[string]TypedValue{
			"a": &Vars{value: map[string]TypedValue{
				"b": &Vars{value: map[string]TypedValue{
					"c": &Vars{value: map[string]TypedValue{
						"d": &Vars{value: map[string]TypedValue{
							"e": &Vars{value: map[string]TypedValue{}},
						}},
					}},
				}},
			}},
		},
	}
	val := Must(NewInteger(10))

	err := Assign(base, "a.b.c.d.e.f", val)
	require.NoError(t, err)
}

func BenchmarkAssign(b *testing.B) {
	base := &Vars{
		value: map[string]TypedValue{
			"a": &Vars{value: map[string]TypedValue{
				"b": &Vars{value: map[string]TypedValue{
					"c": &Vars{value: map[string]TypedValue{
						"d": &Vars{value: map[string]TypedValue{
							"e": &Vars{value: map[string]TypedValue{}},
						}},
					}},
				}},
			}},
		},
	}
	val := Must(NewInteger(10))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Assign(base, "a.b.c.d.e.f", val)
	}
}
