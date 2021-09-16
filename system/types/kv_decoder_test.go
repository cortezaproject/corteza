package types

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	type (
		subdst struct {
			S string `kv:"s"`
			B bool   `kv:"b"`

			Bar struct {
				Foo string `kv:"foo"`
			} `kv:"bar"`
		}

		withHandler struct{}

		dst struct {
			S string `kv:"s"`
			B bool   `kv:"b"`
			N int    `kv:"n"`

			NoKV string

			WH withHandler

			Ptr *string

			Sub subdst `kv:"sub"`

			Map map[string]string `kv:"sub.map"`
			S2I map[string]int    `kv:"sub.s2i"`

			PrefilledString1 string
			PrefilledString2 string
		}
	)

	var (
		ptr = "point-me"

		aux = dst{
			PrefilledString1: "values",
			PrefilledString2: "values",
		}

		kv = SettingsKV{
			"s":           types.JSONText(`"string"`),
			"b":           types.JSONText("true"),
			"n":           types.JSONText("42"),
			"sub.s":       types.JSONText(`"string"`),
			"sub.b":       types.JSONText("true"),
			"sub.bar":     nil,
			"sub.bar.foo": types.JSONText(`"foobar"`),

			"noKV": types.JSONText(`"NO-SettingsKV-!"`),
			"ptr":  types.JSONText(`"point-me"`),

			"sub.map.foo": types.JSONText(`"foo"`),
			"sub.map.bar": types.JSONText(`"bar"`),
			"sub.map.baz": types.JSONText(`"baz"`),

			"sub.s2i.one": types.JSONText(`1`),
			"sub.s2i.two": types.JSONText(`2`),

			"prefilledString1": nil,
			"prefilledString2": types.JSONText(`""`),
		}

		eq = dst{
			S: "string",
			B: true,
			N: 42,

			NoKV: "NO-SettingsKV-!",
			Ptr:  &ptr,

			Sub: subdst{
				S: "string",
				B: true,
			},

			Map: map[string]string{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
			},

			S2I: map[string]int{
				"one": 1,
				"two": 2,
			},

			PrefilledString1: "",
			PrefilledString2: "",
		}
	)

	// setting this externally (embedded structs)
	eq.Sub.Bar.Foo = "foobar"

	require.NoError(t, DecodeKV(kv, &aux))
	require.Equal(t, eq, aux)
}

type (
	decodeHandlerBase struct {
		Foo decodeHandlerSub
		Bar *decodeHandlerSub
	}

	decodeHandlerSub struct {
		set int
	}
)

func (b *decodeHandlerSub) DecodeKV(kv SettingsKV, prefix string) error {
	b.set++
	return nil
}

var _ KVDecoder = &decodeHandlerSub{}

func TestDecodeHandler(t *testing.T) {
	var (
		kv = SettingsKV{
			// should panic if DecodeKV is not called:
			//   cannot unmarshal number into Go value of type settings.decodeHandlerSub
			"foo": types.JSONText(`1`),
		}

		aux = decodeHandlerBase{
			Foo: decodeHandlerSub{},
			Bar: &decodeHandlerSub{},
		}
	)
	require.NoError(t, DecodeKV(kv, &aux))
	require.NotNil(t, aux.Foo)
	require.NotNil(t, aux.Bar)
	require.Equal(t, 1, aux.Foo.set)
	require.Equal(t, 1, aux.Bar.set)
}

func TestDecodeKV_WithFinalTag(t *testing.T) {
	type (
		dst struct {
			NotFinal struct {
				Foo string
				Sub struct {
					SubFoo int
				}
			}

			IsFinal struct {
				Foo string
				Sub struct {
					SubFoo int
				}
			} `kv:",final"`
		}
	)

	var (
		r   = require.New(t)
		aux = dst{}
		kv  = SettingsKV{
			"notFinal.foo":        types.JSONText(`"42"`),
			"notFinal.sub.subFoo": types.JSONText(`42`),
			"isFinal":             types.JSONText(`{"Foo":"final42","Sub":{"SubFoo":42}}`),
		}
	)

	r.NoError(DecodeKV(kv, &aux))
	r.Equal(aux.NotFinal.Foo, "42")
	r.Equal(aux.NotFinal.Sub.SubFoo, 42)
	r.Equal(aux.IsFinal.Foo, "final42")
	r.Equal(aux.IsFinal.Sub.SubFoo, 42)
}
