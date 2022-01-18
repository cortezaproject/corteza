package gig

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	constructorPostprocessorNoop    func(map[string]interface{}) (postprocessorNoop, error)
	constructorPostprocessorDiscard func(map[string]interface{}) (postprocessorDiscard, error)
	constructorPostprocessorSave    func(map[string]interface{}) (postprocessorSave, error)
	constructorPostprocessorArchive func(map[string]interface{}) (postprocessorArchive, error)
)

func Test_postprocessor_tasks(t *testing.T) {

	t.Run("noop constructor unknown", func(t *testing.T) {
		_, err := PostprocessorNoopParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("noop constructor", func(t *testing.T) {
		// func test_postprocessor_tasks_noop_constructor(t *testing.T, c constructorPostprocessorNoop)
		test_postprocessor_tasks_noop_constructor(t, PostprocessorNoopParams)
	})

	t.Run("discard constructor unknown", func(t *testing.T) {
		_, err := PostprocessorDiscardParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("discard constructor", func(t *testing.T) {
		// func test_postprocessor_tasks_discard_constructor(t *testing.T, c constructorPostprocessorDiscard)
		test_postprocessor_tasks_discard_constructor(t, PostprocessorDiscardParams)
	})

	t.Run("save constructor unknown", func(t *testing.T) {
		_, err := PostprocessorSaveParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("save constructor", func(t *testing.T) {
		// func test_postprocessor_tasks_save_constructor(t *testing.T, c constructorPostprocessorSave)
		test_postprocessor_tasks_save_constructor(t, PostprocessorSaveParams)
	})

	t.Run("archive constructor required missing", func(t *testing.T) {
		_, err := PostprocessorArchiveParams(map[string]interface{}{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "required")
	})
	t.Run("archive constructor unknown", func(t *testing.T) {
		_, err := PostprocessorArchiveParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("archive constructor", func(t *testing.T) {
		// func test_postprocessor_tasks_archive_constructor(t *testing.T, c constructorPostprocessorArchive)
		test_postprocessor_tasks_archive_constructor(t, PostprocessorArchiveParams)
	})
}
