package expr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty(t *testing.T) {
	var (
		req              = require.New(t)
		unsetSliceString []string
		unsetSliceInt    []int8
		unsetSliceBool   []bool
		unsetSliceFloat  []float32
		unsetString      string
		unsetInt64       int64

		tcc = []struct {
			value  interface{}
			expect interface{}
		}{
			{
				value:  []string{},
				expect: true,
			},
			{
				value:  map[string]string{},
				expect: true,
			},
			{
				value:  unsetSliceString,
				expect: true,
			},
			{
				value:  []int{},
				expect: true,
			},
			{
				value:  []int{1},
				expect: false,
			},
			{
				value:  unsetSliceInt,
				expect: true,
			},
			{
				value:  unsetSliceBool,
				expect: true,
			},
			{
				value:  int(1),
				expect: false,
			},
			{
				value:  int(0),
				expect: true,
			},
			{
				value:  "",
				expect: true,
			},
			{
				value:  unsetString,
				expect: true,
			},
			{
				value:  unsetSliceFloat,
				expect: true,
			},
			{
				value:  unsetInt64,
				expect: true,
			},
			{
				value:  []float32{11.1},
				expect: false,
			},
			{
				value:  []float32{},
				expect: true,
			},
		}
	)

	for _, tst := range tcc {
		req.Equal(tst.expect, isEmpty(tst.value))
	}
}

func Test_length(t *testing.T) {
	var (
		req = require.New(t)

		tcc = []struct {
			len   int
			value interface{}
		}{
			{0, []string{}},
			{0, map[string]string{}},
			{3, "foo"},
			{0, make(chan string)},
			{0, 34234},
		}
	)

	for _, tst := range tcc {
		req.Equal(tst.len, length(tst.value))
	}
}
