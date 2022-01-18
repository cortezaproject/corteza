package gig

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	constructorPreprocessorNoop                func(map[string]interface{}) (preprocessorNoop, error)
	constructorPreprocessorAttachmentRemove    func(map[string]interface{}) (preprocessorAttachmentRemove, error)
	constructorPreprocessorAttachmentTransform func(map[string]interface{}) (preprocessorAttachmentTransform, error)
	constructorPreprocessorExperimentalExport  func(map[string]interface{}) (preprocessorExperimentalExport, error)
)

func Test_preprocessor_tasks(t *testing.T) {

	t.Run("noop constructor unknown", func(t *testing.T) {
		_, err := PreprocessorNoopParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("noop constructor", func(t *testing.T) {
		// func test_preprocessor_tasks_noop_constructor(t *testing.T, c constructorPreprocessorNoop)
		test_preprocessor_tasks_noop_constructor(t, PreprocessorNoopParams)
	})

	t.Run("attachmentRemove constructor required missing", func(t *testing.T) {
		_, err := PreprocessorAttachmentRemoveParams(map[string]interface{}{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "required")
	})
	t.Run("attachmentRemove constructor unknown", func(t *testing.T) {
		_, err := PreprocessorAttachmentRemoveParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("attachmentRemove constructor", func(t *testing.T) {
		// func test_preprocessor_tasks_attachmentRemove_constructor(t *testing.T, c constructorPreprocessorAttachmentRemove)
		test_preprocessor_tasks_attachmentRemove_constructor(t, PreprocessorAttachmentRemoveParams)
	})

	t.Run("attachmentTransform constructor unknown", func(t *testing.T) {
		_, err := PreprocessorAttachmentTransformParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("attachmentTransform constructor", func(t *testing.T) {
		// func test_preprocessor_tasks_attachmentTransform_constructor(t *testing.T, c constructorPreprocessorAttachmentTransform)
		test_preprocessor_tasks_attachmentTransform_constructor(t, PreprocessorAttachmentTransformParams)
	})

	t.Run("experimentalExport constructor unknown", func(t *testing.T) {
		_, err := PreprocessorExperimentalExportParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("experimentalExport constructor", func(t *testing.T) {
		// func test_preprocessor_tasks_experimentalExport_constructor(t *testing.T, c constructorPreprocessorExperimentalExport)
		test_preprocessor_tasks_experimentalExport_constructor(t, PreprocessorExperimentalExportParams)
	})
}
