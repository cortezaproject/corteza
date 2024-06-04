package expr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_isEmpty(t *testing.T) {
	var (
		req              = require.New(t)
		unsetSliceString []string
		unsetSliceInt    []int8
		unsetSliceBool   []bool
		unsetSliceFloat  []float32
		unsetString      string
		unsetInt64       int64

		unsetDateTime *time.Time
		unsetDuration time.Duration

		tcc = []struct {
			value  interface{}
			expect interface{}
			sc     string
		}{
			{
				value:  []string{},
				expect: true,
				sc:     "empty string slice",
			},
			{
				value:  map[string]string{},
				expect: true,
				sc:     "empty map of strings",
			},
			{
				value:  unsetSliceString,
				expect: true,
				sc:     "undefined string slice",
			},
			{
				value:  []int{},
				expect: true,
				sc:     "empty int slice",
			},
			{
				value:  []int{1},
				expect: false,
				sc:     "1-elem int slice",
			},
			{
				value:  unsetSliceInt,
				expect: true,
				sc:     "undefined int slice",
			},
			{
				value:  unsetSliceBool,
				expect: true,
				sc:     "undefined slice bool",
			},
			{
				value:  int(1),
				expect: false,
				sc:     "defined int",
			},
			{
				value:  int(0),
				expect: true,
				sc:     "defined int 0",
			},
			{
				value:  "",
				expect: true,
				sc:     "emty string",
			},
			{
				value:  unsetString,
				expect: true,
				sc:     "undefined string",
			},
			{
				value:  unsetSliceFloat,
				expect: true,
				sc:     "undefined slice float",
			},
			{
				value:  unsetInt64,
				expect: true,
				sc:     "undefined slice int64",
			},
			{
				value:  []float32{11.1},
				expect: false,
				sc:     "non-empty slice float32",
			},
			{
				value:  []float32{},
				expect: true,
				sc:     "empty slice float32",
			},
			{
				value:  unsetDateTime,
				expect: true,
				sc:     "undefined datetime",
			},
			{
				value:  Must(NewInteger(nil)),
				expect: false,
				sc:     "undefined Integer expr",
			},
			{
				value:  Must(NewDateTime(nil)),
				expect: true,
				sc:     "undefined DateTime expr",
			},
			{
				value:  Must(NewDateTime(unsetDateTime)),
				expect: true,
				sc:     "undefined DateTime expr",
			},
			{
				value:  Must(NewArray([]string{"A"})),
				expect: false,
				sc:     "non empty Array expr",
			},
			{
				value:  Must(NewBoolean(nil)),
				expect: false,
				sc:     "undefined Boolean expr",
			},
			{
				value:  Must(NewBytes(nil)),
				expect: true,
				sc:     "undefined Bytes expr",
			},
			{
				value:  Must(NewDuration(unsetDuration)),
				expect: false,
				sc:     "undefined Duration expr",
			},
		}
	)

	for _, tst := range tcc {
		req.Equal(tst.expect, isEmpty(tst.value), "Failed isEmpty test: %s", tst.sc)
	}
}
func Test_isNil(t *testing.T) {
	var (
		req              = require.New(t)
		unsetSliceString []string
		unsetSliceInt    []int8
		unsetSliceBool   []bool
		unsetSliceFloat  []float32
		unsetString      string
		unsetInt64       int64

		unsetDateTime *time.Time
		unsetDuration time.Duration

		tcc = []struct {
			value  interface{}
			expect interface{}
			sc     string
		}{
			{
				value:  []string{},
				expect: false,
				sc:     "empty string slice",
			},
			{
				value:  map[string]string{},
				expect: false,
				sc:     "empty map of strings",
			},
			{
				// @todo
				value:  unsetSliceString,
				expect: false,
				sc:     "undefined string slice",
			},
			{
				value:  []int{},
				expect: false,
				sc:     "empty int slice",
			},
			{
				value:  []int{1},
				expect: false,
				sc:     "1-elem int slice",
			},
			{
				// @todo
				value:  unsetSliceInt,
				expect: false,
				sc:     "undefined int slice",
			},
			{
				value:  unsetSliceBool,
				expect: false,
				sc:     "undefined slice bool",
			},
			{
				value:  int(1),
				expect: false,
				sc:     "defined int",
			},
			{
				value:  int(0),
				expect: false,
				sc:     "defined int 0",
			},
			{
				value:  "",
				expect: false,
				sc:     "emty string",
			},
			{
				value:  unsetString,
				expect: false,
				sc:     "undefined string",
			},
			{
				value:  unsetSliceFloat,
				expect: false,
				sc:     "undefined slice float",
			},
			{
				value:  unsetInt64,
				expect: false,
				sc:     "undefined slice int64",
			},
			{
				value:  []float32{11.1},
				expect: false,
				sc:     "non-empty slice float32",
			},
			{
				value:  []float32{},
				expect: false,
				sc:     "empty slice float32",
			},
			{
				value:  unsetDateTime,
				expect: true,
				sc:     "undefined datetime",
			},
			{
				value:  Must(NewInteger(nil)),
				expect: false,
				sc:     "undefined Integer expr",
			},
			{
				value:  Must(NewDateTime(nil)),
				expect: true,
				sc:     "nil DateTime expr",
			},
			{
				value:  Must(NewDateTime(unsetDateTime)),
				expect: true,
				sc:     "undefined DateTime expr",
			},
			{
				value:  Must(NewArray([]string{"A"})),
				expect: false,
				sc:     "non empty Array expr",
			},
			{
				value:  Must(NewBoolean(nil)),
				expect: false,
				sc:     "undefined Boolean expr",
			},
			{
				value:  Must(NewBytes(nil)),
				expect: false,
				sc:     "undefined Bytes expr",
			},
			{
				value:  Must(NewDuration(unsetDuration)),
				expect: false,
				sc:     "undefined Duration expr",
			},
		}
	)

	for _, tst := range tcc {
		req.Equal(tst.expect, isNil(tst.value), "Failed isNil test: %s", tst.sc)
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
