package expr

import (
	"reflect"
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

type (
	fac struct {
		expect interface{}
		arr    interface{}
		val    interface{}
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

func Test_shift(t *testing.T) {
	var (
		req = require.New(t)

		tcc = []tc{
			{
				value:  []string{"1", "2", "3"},
				expect: "1",
			},
			{
				value:  map[string]string{"test": "123"},
				expect: nil,
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
				value:  int(1),
				expect: nil,
			},
			{
				value:  []float64{11.1},
				expect: 11.1,
			},
		}
	)

	for _, tst := range tcc {
		val, _ := shift(tst.value)

		req.Equal(tst.expect, val)
	}
}

func Test_find(t *testing.T) {
	var (
		req = require.New(t)

		tcc = []fac{
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
	)

	for _, tst := range tcc {
		loc := find(tst.arr, tst.val)

		req.Equal(tst.expect, loc)
	}
}

func Test_count(t *testing.T) {
	var (
		req = require.New(t)

		tcc = []fac{
			{
				arr:    []string{"1", "2", "3"},
				val:    []string{"0", "3"},
				expect: 1,
			},
			{
				arr:    []bool{true, true},
				val:    []bool{false, false},
				expect: 0,
			},
			{
				arr:    []bool{true, true},
				val:    []bool{false, true},
				expect: 1,
			},
			{
				arr:    []int{4, 5, 6, 7},
				val:    []int{7, 4},
				expect: 2,
			},
			{
				arr:    []float64{11.1, 12.4},
				val:    []float64{0.1, 1.1},
				expect: 0,
			},
		}
	)

	for _, tst := range tcc {
		var loc int

		switch reflect.TypeOf(tst.val).Elem().Kind() {
		case reflect.String:
			loc = count(tst.arr, tst.val.([]string)[0], tst.val.([]string)[1])
			break
		case reflect.Bool:
			loc = count(tst.arr, tst.val.([]bool)[0], tst.val.([]bool)[1])
			break
		case reflect.Int:
			loc = count(tst.arr, tst.val.([]int)[0], tst.val.([]int)[1])
			break
		case reflect.Float64:
			loc = count(tst.arr, tst.val.([]float64)[0], tst.val.([]float64)[1])
			break
		}

		req.Equal(tst.expect, loc)
	}
}

func Test_has(t *testing.T) {
	var (
		req = require.New(t)

		tcc = []fac{
			{
				arr:    []string{"1", "2", "3"},
				val:    []string{"0", "3"},
				expect: true,
			},
			{
				arr:    []bool{true, true},
				val:    []bool{false, false},
				expect: false,
			},
			{
				arr:    []bool{true, true},
				val:    []bool{false, true},
				expect: true,
			},
			{
				arr:    []int{4, 5, 6, 7},
				val:    []int{7, 4},
				expect: true,
			},
			{
				arr:    []float64{11.1, 12.4},
				val:    []float64{0.1, 1.1},
				expect: false,
			},
		}
	)

	for _, tst := range tcc {
		var loc bool

		switch reflect.TypeOf(tst.val).Elem().Kind() {
		case reflect.String:
			loc = has(tst.arr, tst.val.([]string)[0], tst.val.([]string)[1])
			break
		case reflect.Bool:
			loc = has(tst.arr, tst.val.([]bool)[0], tst.val.([]bool)[1])
			break
		case reflect.Int:
			loc = has(tst.arr, tst.val.([]int)[0], tst.val.([]int)[1])
			break
		case reflect.Float64:
			loc = has(tst.arr, tst.val.([]float64)[0], tst.val.([]float64)[1])
			break
		}

		req.Equal(tst.expect, loc)
	}
}

func Test_hasAll(t *testing.T) {
	var (
		req = require.New(t)

		tcc = []fac{
			{
				arr:    []string{"1", "2", "3"},
				val:    []string{"0", "3"},
				expect: false,
			},
			{
				arr:    []bool{true, true},
				val:    []bool{false, false},
				expect: false,
			},
			{
				arr:    []bool{true, true},
				val:    []bool{false, true},
				expect: false,
			},
			{
				arr:    []int{4, 5, 6, 7},
				val:    []int{7, 4},
				expect: true,
			},
			{
				arr:    []float64{11.1, 12.4},
				val:    []float64{0.1, 1.1},
				expect: false,
			},
		}
	)

	for _, tst := range tcc {
		var loc bool

		switch reflect.TypeOf(tst.arr).Elem().Kind() {
		case reflect.String:
			loc = hasAll(tst.arr, tst.val.([]string)[0], tst.val.([]string)[1])
			break
		case reflect.Bool:
			loc = hasAll(tst.arr, tst.val.([]bool)[0], tst.val.([]bool)[1])
			break
		case reflect.Int:
			loc = hasAll(tst.arr, tst.val.([]int)[0], tst.val.([]int)[1])
			break
		case reflect.Float64:
			loc = hasAll(tst.arr, tst.val.([]float64)[0], tst.val.([]float64)[1])
			break
		}

		req.Equal(tst.expect, loc)
	}
}

func Test_slice(t *testing.T) {
	type (
		sct struct {
			vals   []int
			arr    interface{}
			expect interface{}
		}
	)

	var (
		req = require.New(t)

		tcc = []sct{
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
	)

	for _, tst := range tcc {
		new := slice(tst.arr, tst.vals[0], tst.vals[1])
		req.Equal(tst.expect, new)
	}
}
