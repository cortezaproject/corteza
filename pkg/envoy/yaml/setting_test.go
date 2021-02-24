package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSettings_UnmarshalYAML(t *testing.T) {
	t.Run("settings 1", func(t *testing.T) {
		req := require.New(t)

		doc, err := parseDocument("settings_1")
		req.NoError(err)
		req.NotNil(doc)
		req.Len(doc.settings, 20)
		// req.Contains(doc.settings, "privacy.mask.email")
		// req.Equal(true, doc.settings["privacy.mask.email"])
	})
}
