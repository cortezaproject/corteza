package expr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	intArr    = []int{}
	stringArr = []string{"first"}
	boolArr   = []bool{true, true, false}
	floatArr  = []float64{69.420}

	vals = map[string]interface{}{
		"intArr":    intArr,
		"stringArr": stringArr,
		"boolArr":   boolArr,
		"floatArr":  floatArr,
		"intVal":    42,
		"stringVal": "foobar",
		"boolVal":   false,
	}
)

func Example_push_int() {
	eval(`push(intArr, intVal)`, vals)

	// output:
	// [42]
}

func Example_push_string() {
	eval(`push(stringArr, stringVal)`, vals)

	// output:
	// [first foobar]
}

func Example_push_bool() {
	eval(`push(boolArr, boolVal)`, vals)

	// output:
	// [true true false false]
}

func Example_push_float() {
	eval(`push(floatArr, 3.14)`, vals)

	// output:
	// [69.42 3.14]
}

func Example_pop_string() {
	eval(`pop(stringArr)`, vals)

	// output:
	// first
}

func Example_pop_int() {
	eval(`pop(intArr)`, vals)

	// output:
	// <nil>
}

func Example_pop_float() {
	eval(`pop(floatArr)`, vals)

	// output:
	// 69.42
}

func Test_push(t *testing.T) {
	tcc := []struct {
		base     interface{}
		new      []interface{}
		expect   interface{}
		expError bool
	}{
		{
			base:   []string{"1", "2", "3"},
			new:    []interface{}{"4"},
			expect: []string{"1", "2", "3", "4"},
		},
		{
			base:   Must(NewArray(nil)),
			new:    []interface{}{"foo"},
			expect: []TypedValue{Must(NewString("foo"))},
		},
		{
			base:   Must(NewArray([]interface{}{"foo"})),
			new:    []interface{}{"bar"},
			expect: []TypedValue{Must(NewString("foo")), Must(NewString("bar"))},
		},
		{
			base:   Must(NewArray([]interface{}{Must(NewString("foo"))})),
			new:    []interface{}{"bar"},
			expect: []TypedValue{Must(NewString("foo")), Must(NewString("bar"))},
		},
		{
			base:   Must(NewArray([]interface{}{Must(NewString("foo"))})),
			new:    []interface{}{Must(NewString("bar"))},
			expect: []TypedValue{Must(NewString("foo")), Must(NewString("bar"))},
		},
		{
			base:     []string{"1", "2", "3"},
			new:      []interface{}{4},
			expError: true,
		},
	}

	for p, tc := range tcc {
		t.Run(fmt.Sprintf("%d", p), func(t *testing.T) {
			var (
				req      = require.New(t)
				out, err = push(tc.base, tc.new...)
			)

			if tc.expError {
				req.Error(err)
				return
			}

			req.NoError(err)
			req.Equal(tc.expect, out)
		})
	}
}

func Test_shift(t *testing.T) {
	tcc := []struct {
		value    interface{}
		expect   interface{}
		expError bool
	}{
		{
			value:  []string{"1", "2", "3"},
			expect: "1",
		},
		{
			value:    map[string]string{"test": "123"},
			expect:   nil,
			expError: true,
		},
		{
			value:  []int{4, 5, 6, 7},
			expect: 4,
		},
		{
			value:  []int{},
			expect: nil,
		},
		{
			value:    int(1),
			expect:   nil,
			expError: true,
		},
		{
			value:  []float64{11.1},
			expect: 11.1,
		},
	}

	for p, tc := range tcc {
		t.Run(fmt.Sprintf("%d", p), func(t *testing.T) {
			var (
				req      = require.New(t)
				val, err = shift(tc.value)
			)

			if tc.expError {
				req.Error(err)
				return
			}

			req.NoError(err)
			req.Equal(tc.expect, val)
		})
	}
}

func Test_find(t *testing.T) {
	tcc := []struct {
		expect interface{}
		arr    interface{}
		val    interface{}
	}{
		{
			arr:    []string{"1", "2", "3"},
			val:    "3",
			expect: 2,
		},
		{
			arr:    []bool{true, false, true},
			val:    true,
			expect: 0,
		},
		{
			arr:    []int{4, 5, 6, 7},
			val:    7,
			expect: 3,
		},
		{
			arr:    []int{},
			val:    0,
			expect: -1,
		},
		{
			arr:    []float64{11.1, 12.4},
			val:    11.1,
			expect: 0,
		},
		{
			arr:    []float64{11.1, 12.4},
			val:    11.2,
			expect: -1,
		},
	}

	for p, tc := range tcc {
		t.Run(fmt.Sprintf("%d", p), func(t *testing.T) {
			var (
				req      = require.New(t)
				loc, err = find(tc.arr, tc.val)
			)

			req.NoError(err)
			req.Equal(tc.expect, loc)
		})
	}
}

func Test_count(t *testing.T) {
	tcc := []struct {
		expect interface{}
		arr    interface{}
		val    []interface{}
	}{
		{
			arr:    []string{"1", "2", "3"},
			val:    []interface{}{"0", "3"},
			expect: 1,
		},
		{
			arr:    []bool{true, true},
			val:    []interface{}{false, false},
			expect: 0,
		},
		{
			arr:    []bool{true, true},
			val:    []interface{}{false, true},
			expect: 1,
		},
		{
			arr:    []int{4, 5, 6, 7},
			val:    []interface{}{7, 4},
			expect: 2,
		},
		{
			arr:    []float64{11.1, 12.4},
			val:    []interface{}{0.1, 1.1},
			expect: 0,
		},
	}

	for p, tc := range tcc {
		t.Run(fmt.Sprintf("%d", p), func(t *testing.T) {
			var (
				req = require.New(t)
				err error
				loc int
			)

			loc, err = count(tc.arr, tc.val...)
			req.NoError(err)
			req.Equal(tc.expect, loc)
		})
	}
}

func Test_has(t *testing.T) {
	tcc := []struct {
		expect interface{}
		arr    interface{}
		val    []interface{}
	}{
		{
			arr:    []string{"1", "2", "3"},
			val:    []interface{}{"0", "3"},
			expect: true,
		},
		{
			arr:    []bool{true, true},
			val:    []interface{}{false, false},
			expect: false,
		},
		{
			arr:    []bool{true, true},
			val:    []interface{}{false, true},
			expect: true,
		},
		{
			arr:    []int{4, 5, 6, 7},
			val:    []interface{}{7, 4},
			expect: true,
		},
		{
			arr:    []float64{11.1, 12.4},
			val:    []interface{}{0.1, 1.1},
			expect: false,
		},
	}

	for p, tc := range tcc {
		t.Run(fmt.Sprintf("%d", p), func(t *testing.T) {
			var (
				req = require.New(t)
				loc bool
				err error
			)

			loc, err = has(tc.arr, tc.val...)
			req.NoError(err)
			req.Equal(tc.expect, loc)
		})
	}
}

func Test_hasAll(t *testing.T) {
	tcc := []struct {
		arr      interface{}
		val      []interface{}
		hasAll   bool
		expError bool
	}{
		{
			arr:    []string{"1", "2", "3"},
			val:    []interface{}{"0", "3"},
			hasAll: false,
		},
		{
			arr:    []bool{true, true},
			val:    []interface{}{false, false},
			hasAll: false,
		},
		{
			arr:    []bool{true, true},
			val:    []interface{}{false, true},
			hasAll: false,
		},
		{
			arr:    []int{4, 5, 6, 7},
			val:    []interface{}{7, 4},
			hasAll: true,
		},
		{
			arr:    []float64{11.1, 12.4},
			val:    []interface{}{0.1, 1.1},
			hasAll: false,
		},
	}

	for p, tc := range tcc {
		t.Run(fmt.Sprintf("%d", p), func(t *testing.T) {
			var (
				req       = require.New(t)
				rval, err = hasAll(tc.arr, tc.val...)
			)

			if tc.expError {
				req.Error(err)
				return
			}

			req.NoError(err)
			req.Equal(tc.hasAll, rval)
		})
	}
}

func Test_slice(t *testing.T) {
	tcc := []struct {
		vals   []int
		arr    interface{}
		expect interface{}
	}{
		{
			vals:   []int{0, 3},
			arr:    []string{"1", "2", "3"},
			expect: []string{"1", "2", "3"},
		},
		{
			vals:   []int{0, 1},
			arr:    []string{"1", "2", "3"},
			expect: []string{"1"},
		},
		{
			vals:   []int{2, 3},
			arr:    []bool{true, true},
			expect: []bool{true, true},
		},
		{
			vals:   []int{1, -1},
			arr:    []int{4, 5, 6, 7},
			expect: []int{5, 6, 7},
		},
		{
			vals:   []int{3, -1},
			arr:    []float64{11.1, 12.4},
			expect: []float64{11.1, 12.4},
		},
	}

	for p, tc := range tcc {
		t.Run(fmt.Sprintf("%d", p), func(t *testing.T) {
			var (
				req = require.New(t)
				ss  = slice(tc.arr, tc.vals[0], tc.vals[1])
			)

			req.Equal(tc.expect, ss)
		})
	}
}
