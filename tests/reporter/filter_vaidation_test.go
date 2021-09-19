package reporter

import (
	"testing"
)

func Test_filter_validation(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
	)

	t.Run("empty conjunction", func(t *testing.T) {
		loadErr(ctx, h, m, dd[0], "could not build query: expecting 1 or more arguments, got 0")
	})
}
