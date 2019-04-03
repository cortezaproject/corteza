// +build unit

package rules

import (
	"testing"

	"github.com/crusttech/crust/internal/test"
)

func TestResource(t *testing.T) {
	var (
		assert = test.Assert

		sCases = []struct {
			r Resource
			s string
		}{
			{
				Resource("a:b:c"),
				"a:b:c"},
			{
				Resource("a:b:c").PermissionResource(),
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
		assert(t, sc.r.String() == sc.s, "Resource check failed (%s != %s)", sc.r, sc.s)
	}

	var r string
	r = "a:"
	assert(t, Resource(r).IsAppendable(), "Expecting resource %q to be appendable", r)
	r = "a:1"
	assert(t, Resource(r).IsAppendable(), "Expecting resource %q to be appendable", r)
	r = "a:*"
	assert(t, Resource(r).IsAppendable(), "Expecting resource %q to be appendable", r)

	r = "a"
	assert(t, Resource(r).IsValid(), "Expecting resource %q to be valid", r)
	r = "a:"
	assert(t, !Resource(r).IsValid(), "Expecting resource %q not to be valid", r)
	r = "a:1"
	assert(t, Resource(r).IsValid(), "Expecting resource %q to be valid", r)
	r = "a:*"
	assert(t, Resource(r).IsValid(), "Expecting resource %q to be valid", r)

	r = "a:1"
	assert(t, !Resource(r).HasWildcard(), "Expecting resource %q to not have wildcard", r)
	r = "a:*"
	assert(t, Resource(r).HasWildcard(), "Expecting resource %q to have wildcard", r)
}
