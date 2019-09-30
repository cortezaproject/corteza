package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImporter_Cast(t *testing.T) {
	impFixTester(t, "importer_with_namespace", func(t *testing.T, i *Importer) {
		req := require.New(t)

		req.NotNil(i)
		req.NotNil(i.namespaces)
		req.NotNil(i.namespaces.set)
		req.NotNil(i.namespaces.modules)
		req.NotNil(i.namespaces.modules["foo"])
		req.NotNil(i.namespaces.modules["foo"].set)
		req.Len(i.namespaces.modules["foo"].set, 1)
		req.NotNil(i.namespaces.modules["foo"].set.FindByHandle("m1"))
		req.Equal("Module1", i.namespaces.modules["foo"].set.FindByHandle("m1").Name)

	})
}
