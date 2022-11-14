package cast2

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUint64(t *testing.T) {
	var (
		req           = require.New(t)
		target uint64 = 42
	)

	req.NoError(Uint64(nil, &target))
	req.Equal(uint64(0), target)

	req.NoError(Uint64("42", &target))
	req.Equal(uint64(42), target)

	target = 5
	req.Error(Uint64("-42", &target))
	req.Equal(uint64(5), target)
}
