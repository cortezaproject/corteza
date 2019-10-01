package compose

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigure(t *testing.T) {
	var config = Configure()
	require.True(t, config != nil, "Configure valid")
	require.True(t, func() bool { config.Init(); return true }(), "Initialization ok")
	require.True(t, config.MakeCLI(context.Background()) != nil, "CLI created")
}
