package registry

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/stretchr/testify/require"
)

func Test_registryAddGet(t *testing.T) {
	var (
		req = require.New(t)
		r   = NewRegistry(types.Config{})
	)

	r.Add("mockHandler", types.MockHandler{})

	h, err := r.Get("mockHandler")

	req.NoError(err)
	req.Len(r.h, 1)
	req.IsType(types.MockHandler{}, h)
}

func Test_registryAddGetErr(t *testing.T) {
	var (
		req = require.New(t)
		r   = NewRegistry(types.Config{})
	)

	r.Add("mockHandler", types.MockHandler{})

	h, err := r.Get("foo")

	req.EqualError(err, "could not get element from registry: foo")
	req.Len(r.h, 1)
	req.Nil(h)
}

func Test_registryMerge(t *testing.T) {
	type (
		tf struct {
			name   string
			err    string
			params string
			exp    string
		}
	)

	var (
		tcc = []tf{
			{
				name:   "set params",
				params: `{"foo":"bar"}`,
				exp:    "bar",
			},
			{
				name:   "set invalid params",
				params: `{"foo1":"bar"}`,
				exp:    "",
			},
			{
				name:   "set invalid params err",
				params: `{"foo1":"bar"`,
				exp:    "",
				err:    "unexpected EOF",
			},
		}
	)

	for _, tc := range tcc {
		var (
			req = require.New(t)
			r   = NewRegistry(types.Config{})
		)

		m, err := r.Merge(types.MockHandler{}, []byte(tc.params), types.Config{})

		if tc.err != "" {
			req.EqualError(err, tc.err)
		} else {
			req.Equal(m.(types.MockHandler).Foo, tc.exp)
			req.NoError(err)
		}
	}

}

func Test_registryAll(t *testing.T) {
	var (
		req = require.New(t)
		r   = NewRegistry(types.Config{})
	)

	r.Add("mockHandler", types.MockHandler{})

	list := r.All()

	req.Len(list, 1)
	req.NotEmpty(list[0].Name)
}
