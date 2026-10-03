package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModuleFieldOptions_SanitizeXSS(t *testing.T) {
	req := require.New(t)

	// sanitized unless explicitly turned off
	req.True(ModuleFieldOptions{}.SanitizeXSS())
	req.True(ModuleFieldOptions{"sanitizeXSS": true}.SanitizeXSS())
	req.False(ModuleFieldOptions{"sanitizeXSS": false}.SanitizeXSS())

	// rich text is always sanitized
	req.True(ModuleFieldOptions{"sanitizeXSS": false, "useRichTextEditor": true}.SanitizeXSS())

	opt := ModuleFieldOptions{}
	opt.SetSanitizeXSS(false)
	req.False(opt.SanitizeXSS())
}
