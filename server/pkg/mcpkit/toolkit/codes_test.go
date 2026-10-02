package toolkit

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWithCode pins that a coded error keeps its cause reachable and its code
// readable through any number of fmt wraps, which is how handlers wrap errors.
func TestWithCode(t *testing.T) {
	cause := errors.New("boom")
	err := fmt.Errorf("record lookup failed: %w", WithCode(cause, CodeNotFound, "look it up"))

	var coded Coded
	require.ErrorAs(t, err, &coded)
	assert.Equal(t, CodeNotFound, coded.ToolErrorCode())
	assert.Equal(t, "look it up", coded.ToolErrorNext())
	assert.ErrorIs(t, err, cause)
	assert.Nil(t, WithCode(nil, CodeFailed, ""))
}

// TestRequiredArgumentsAreCoded: every missing or malformed declared argument
// is an invalid_argument, so the transport can point at the parameter docs.
func TestRequiredArgumentsAreCoded(t *testing.T) {
	var coded Coded
	_, err := ReqStr(map[string]any{}, "namespace")
	require.ErrorAs(t, err, &coded)
	assert.Equal(t, CodeInvalidArgument, coded.ToolErrorCode())

	_, err = ReqID(map[string]any{"moduleID": 12}, "moduleID")
	require.ErrorAs(t, err, &coded)
	assert.Equal(t, CodeInvalidArgument, coded.ToolErrorCode())

	_, err = ReqRef(map[string]any{"page": ""}, "page")
	require.ErrorAs(t, err, &coded)
	assert.Equal(t, CodeInvalidArgument, coded.ToolErrorCode())
}
