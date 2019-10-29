package rh

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNow(t *testing.T) {
	var (
		r = require.New(t)

		val  time.Time
		ptr  *time.Time
		inv1 int
		inv2 string
	)

	SetCurrentTimeRounded(&val)
	r.NotEmpty(val)

	SetCurrentTimeRounded(&ptr)
	r.NotNil(ptr)

	SetCurrentTimeRounded(&inv1)
	r.Empty(inv1)

	SetCurrentTimeRounded(&inv2)
	r.Empty(inv2)

}
