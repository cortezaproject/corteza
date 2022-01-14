package compose

import (
	"testing"
)

func Test_namespace_import(t *testing.T) {
	t.Run("nested", func(_ *testing.T) {
		_, h, _ := setup(t)
		grantImportExport(h)

		sessionID := namespaceImportInitPathSafe(t, h, "nested.zip")
		h.a.NotEqual(0, sessionID)
	})

	t.Run("no namespace", func(_ *testing.T) {
		_, h, _ := setup(t)
		grantImportExport(h)

		_, err := namespaceImportInitPath(t, h, "no-ns.zip")
		h.a.Error(err)
		h.a.Contains(err.Error(), "namespace.errors.importMissingNamespace")
	})
}
