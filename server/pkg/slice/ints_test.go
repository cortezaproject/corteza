package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasUint64(t *testing.T) {
	cases := []struct {
		name string
		ss   []uint64
		s    uint64
		o    bool
	}{
		{
			"empty",
			[]uint64{},
			42,
			false,
		},
		{
			"has not",
			[]uint64{42},
			12345,
			false,
		},
		{
			"has",
			[]uint64{42},
			42,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.EqualValues(t, HasUint64(c.ss, c.s), c.o)
		})
	}
}
