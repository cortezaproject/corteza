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
	constructorDecoderNoop    func(map[string]interface{}) (decoderNoop, error)
	constructorDecoderArchive func(map[string]interface{}) (decoderArchive, error)
)

func Test_decoder_tasks(t *testing.T) {

	t.Run("noop constructor unknown", func(t *testing.T) {
		_, err := DecoderNoopParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("noop constructor", func(t *testing.T) {
		// func test_decoder_tasks_noop_constructor(t *testing.T, c constructorDecoderNoop)
		test_decoder_tasks_noop_constructor(t, DecoderNoopParams)
	})

	t.Run("archive constructor unknown", func(t *testing.T) {
		_, err := DecoderArchiveParams(map[string]interface{}{
			"does not exist": "true",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown parameter")
	})
	t.Run("archive constructor", func(t *testing.T) {
		// func test_decoder_tasks_archive_constructor(t *testing.T, c constructorDecoderArchive)
		test_decoder_tasks_archive_constructor(t, DecoderArchiveParams)
	})
}
