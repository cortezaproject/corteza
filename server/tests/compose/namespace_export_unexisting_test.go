package compose

import (
	"testing"
)

func Test_namespace_export_unexisting(t *testing.T) {
	ctx, h, s := setup(t)
	grantImportExport(h)
	_ = ctx
	_ = s

	_, err := namespaceExport(t, h, 42)
	h.a.Error(err)
	h.a.Contains(err.Error(), "namespace does not exist")

	cleanup(t)
}
