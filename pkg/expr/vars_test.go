package expr

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// extract typed-value

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

			vars = RVars{
				"int":     Must(NewInteger(42)),
				"STRING":  Must(NewString("foo")),
				"bool":    Must(NewBoolean(true)),
				"missing": Must(NewBoolean(true)),
			}.Vars()
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

			vars = RVars{
				"iString":  Must(NewString("foo")),
				"iInteger": Must(NewInteger(42)),
				"iBool":    Must(NewBoolean(true)),
			}.Vars()
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

			vars = RVars{
				"uint64": Must(NewAny("42")),
				"int":    Must(NewAny("42")),
			}.Vars()
		)

		dst.Uint64 = 0
		dst.Int = 0

		req.NoError(vars.Decode(dst))
		req.Equal(uint64(42), dst.Uint64)
		req.Equal(int64(42), dst.Int)
	})
}
