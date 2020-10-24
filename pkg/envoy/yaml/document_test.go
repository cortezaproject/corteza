package yaml

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDocument_UnmarshalYAML(t *testing.T) {
	t.Run("doc 1", func(t *testing.T) {
		req := require.New(t)
		doc, err := parseDocument("global_rbac_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.rbac)
		req.NotEmpty(doc.rbac.rules)
	})
}
