package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoleImport_CastSet(t *testing.T) {
	impFixTester(t, "roles", func(t *testing.T, ri *Role) {
		req := require.New(t)
		req.NotNil(ri.set)
		req.Len(ri.set, 2)

		req.NotNil(ri.set.FindByHandle("r1"))
		req.Equal("Role1", ri.set.FindByHandle("r1").Name)

		req.NotNil(ri.set.FindByHandle("r2"))
		req.Equal("Role2", ri.set.FindByHandle("r2").Name)
	})
}
