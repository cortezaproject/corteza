package expr

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPathSplit(t *testing.T) {
	tcc := []struct {
		p   string
		r   []string
		err string
	}{
		{p: "a", r: []string{"a"}},
		{p: "foo.bar", r: []string{"foo", "bar"}},
		{p: "a.b[1]", r: []string{"a", "b", "1"}},
		{p: "a.b[1].baz[0]", r: []string{"a", "b", "1", "baz", "0"}},
		{p: "a.[]", err: invalidPathErr.Error()},
		{p: "a[1]bzz", err: invalidPathErr.Error()},
	}

	for _, tc := range tcc {
		t.Run(tc.p, func(t *testing.T) {
			req := require.New(t)
			pp, err := PathSplit(tc.p)
			if len(tc.err) == 0 {
				req.NoError(err)
			} else {
				req.EqualError(err, tc.err)
			}

			req.Equal(tc.r, pp)
		})
	}

}

func TestVars(t *testing.T) {
	var (
		req = require.New(t)

		vars = RVars{
			"int": Must(NewInteger(42)),
			"sub": RVars{
				"foo": Must(NewString("foo")),
			}.Vars(),
			"three": RVars{
				"two": RVars{
					"one": RVars{
						"go": Must(NewString("!")),
					}.Vars(),
				}.Vars(),
			}.Vars(),
		}.Vars()
	)

	req.NoError(Assign(vars, "int", 123))
	req.Equal(123, Must(Select(vars, "int")).Get().(int))

	req.NoError(Assign(vars, "sub.foo", "bar"))
	req.Equal("bar", Must(Select(vars, "sub.foo")).Get().(string))

	req.NoError(Assign(vars, "kv", &KV{}))
	req.NoError(Assign(vars, "kv.foo", "bar"))
	req.Equal("bar", Must(Select(vars, "kv.foo")).Get().(string))

	req.NoError(Assign(vars, "three.two.one.go", "!!!"))
	req.Equal("!!!", Must(Select(vars, "three.two.one.go")).Get().(string))
}
