package gig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func test_postprocessor_tasks_noop_constructor(t *testing.T, c constructorPostprocessorNoop) {
	_, err := c(map[string]interface{}{})

	require.NoError(t, err)
}

func test_postprocessor_tasks_discard_constructor(t *testing.T, c constructorPostprocessorDiscard) {
	_, err := c(map[string]interface{}{})

	require.NoError(t, err)
}

// @todo
func test_postprocessor_tasks_save_constructor(t *testing.T, c constructorPostprocessorSave) {
	_, err := c(map[string]interface{}{})

	require.NoError(t, err)
}

func test_postprocessor_tasks_archive_constructor(t *testing.T, c constructorPostprocessorArchive) {
	t.Run("tar", func(t *testing.T) {
		out, err := c(map[string]interface{}{
			"encoding": "tar",
		})

		require.NoError(t, err)
		require.NotEmpty(t, out.name)
		require.Equal(t, ArchiveTar, out.encoding)
	})

	t.Run("zip", func(t *testing.T) {
		out, err := c(map[string]interface{}{
			"encoding": "zip",
		})

		require.NoError(t, err)
		require.NotEmpty(t, out.name)
		require.Equal(t, ArchiveZIP, out.encoding)
	})
}
