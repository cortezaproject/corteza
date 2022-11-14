package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocument_UnmarshalYAML(t *testing.T) {
	t.Run("doc 1", func(t *testing.T) {
		req := require.New(t)
		doc, err := parseDocument("global_rbac_1")
		req.NoError(err)
		req.NotNil(doc)
		req.NotNil(doc.rbac)
		req.NotEmpty(doc.rbac)
	})
}
