package gig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func test_preprocessor_tasks_noop_constructor(t *testing.T, c constructorPreprocessorNoop) {
	_, err := c(map[string]interface{}{})

	require.NoError(t, err)
}

func test_preprocessor_tasks_attachmentRemove_constructor(t *testing.T, c constructorPreprocessorAttachmentRemove) {
	// @todo not implemented
}

func test_preprocessor_tasks_attachmentTransform_constructor(t *testing.T, c constructorPreprocessorAttachmentTransform) {
	// @todo not implemented
}

func test_preprocessor_tasks_experimentalExport_constructor(t *testing.T, c constructorPreprocessorExperimentalExport) {
	_, err := c(map[string]interface{}{
		"handle":           "test",
		"inclRBAC":         true,
		"inclTranslations": true,
	})

	require.NoError(t, err)
}
