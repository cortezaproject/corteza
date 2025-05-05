package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexing(t *testing.T) {
	req := require.New(t)

	svc := wrapperIndex{}
	req.True(svc.add("corteza::compose:module-field/1/2/3"))
	req.True(svc.add("corteza::compose:module-field/1/4/6"))

	req.True(svc.add("corteza::compose:module-field/1/*/*"))
	req.True(svc.add("corteza::compose:module-field/1/4/*"))

	// False since no resource matches this wildcard
	req.False(svc.add("corteza::compose:module-field/1/5/*"))
	req.False(svc.add("corteza::compose:module-field/2/*/*"))

	// False since it's a completely different resource
	req.False(svc.add("corteza::compose:record/1/2/*"))
}
