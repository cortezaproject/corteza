package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResource(t *testing.T) {
	var (
		req = require.New(t)

		sCases = []struct {
			r Resource
			s string
		}{
			{
				Resource("a:b:c"),
				"a:b:c"},
			{
				Resource("a:b:c").RBACResource(),
				"a:b:c"},
			{
				Resource("a:b:").AppendID(1),
				"a:b:1"},
			{
				Resource("a:b:").AppendWildcard(),
				"a:b:*"},
			{
				Resource("a:b:1").TrimID(),
				"a:b:"},
			{
				Resource("a:b:1").GetService(),
				"a"},
		}
	)

	for _, sc := range sCases {
		req.Equal(sc.s, sc.r.String())
	}

	var r string
	r = "a:"
	req.True(Resource(r).IsAppendable(), "Expecting resource %q to be appendable", r)
	r = "a:1"
	req.True(Resource(r).IsAppendable(), "Expecting resource %q to be appendable", r)
	r = "a:*"
	req.True(Resource(r).IsAppendable(), "Expecting resource %q to be appendable", r)

	r = "a"
	req.True(Resource(r).IsValid(), "Expecting resource %q to be valid", r)
	r = "a:"
	req.False(Resource(r).IsValid(), "Expecting resource %q not to be valid", r)
	r = "a:1"
	req.True(Resource(r).IsValid(), "Expecting resource %q to be valid", r)
	r = "a:*"
	req.True(Resource(r).IsValid(), "Expecting resource %q to be valid", r)

	r = "a:1"
	req.False(Resource(r).HasWildcard(), "Expecting resource %q to not have wildcard", r)
	r = "a:*"
	req.True(Resource(r).HasWildcard(), "Expecting resource %q to have wildcard", r)
}
