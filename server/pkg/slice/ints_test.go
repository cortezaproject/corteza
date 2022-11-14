package slice

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestUInt64s_MarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		ii   UInt64s
		json string
	}{
		{
			"empty",
			UInt64s{},
			`[]`,
		},
		{
			"one",
			UInt64s{285372959844073660},
			`["285372959844073660"]`,
		},
		{
			"two",
			UInt64s{285372959844073660, 285372959844073661},
			`["285372959844073660","285372959844073661"]`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				req         = require.New(t)
				result, err = json.Marshal(c.ii)
			)

			req.NoError(err)
			req.Equal(string(result), c.json)
		})
	}
}
