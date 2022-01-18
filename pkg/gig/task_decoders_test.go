package gig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func test_decoder_tasks_noop_constructor(t *testing.T, c constructorDecoderNoop) {
	_, err := c(map[string]interface{}{
		"source": "0",
	})

	require.NoError(t, err)
}

func test_decoder_tasks_archive_constructor(t *testing.T, c constructorDecoderArchive) {
	_, err := c(map[string]interface{}{
		"source": "0",
	})

	require.NoError(t, err)
}
