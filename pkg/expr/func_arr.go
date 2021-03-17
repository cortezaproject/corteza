package expr

import (
	"reflect"

	"github.com/PaesslerAG/gval"
)

func ArrayFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("push", push),
		gval.Function("pop", pop),
		gval.Function("shift", shift),
		gval.Function("count", count),
		gval.Function("has", has),
		gval.Function("hasAll", hasAll),
	}
}

// push adds a value to the end of slice, returns copy
func push(arr interface{}, p interface{}) (interface{}, error) {
	if !isSlice(arr) {
		return nil, &reflect.ValueError{Method: "Index", Kind: reflect.ValueOf(arr).Kind()}
	}

	c := reflect.ValueOf(arr)

	nval := reflect.MakeSlice(c.Type(), c.Len()+1, c.Cap()+1)
	reflect.Copy(nval, c)
	nval.Index(c.Len()).Set(reflect.ValueOf(p))

	return nval.Interface(), nil
}

// pop takes the last value in slice, does not modify original
func pop(arr interface{}) (interface{}, error) {
	if !isSlice(arr) {
		return nil, &reflect.ValueError{Method: "Index", Kind: reflect.ValueOf(arr).Kind()}
	}

	c := reflect.ValueOf(arr)

	if c.Len() == 0 {
		return nil, nil
	}

	return c.Index(c.Len() - 1).Interface(), nil
}

// shifts takes the first value in slice, does not modify original
func shift(arr interface{}) (interface{}, error) {
	if !isSlice(arr) {
		return nil, &reflect.ValueError{Method: "Index", Kind: reflect.ValueOf(arr).Kind()}
	}

	c := reflect.ValueOf(arr)

	if c.Len() == 0 {
		return nil, nil
	}

	return c.Index(0).Interface(), nil
}

// cound gets the count of occurences in the slice
func count(arr interface{}, v ...interface{}) int {
	if !isSlice(arr) {
		return 0
	}

	count := 0

	for _, vv := range v {
		if find(arr, vv) != -1 {
			count++
		}
	}

	return count
}

// has finds any occurence of the values in slice
func has(arr interface{}, v ...interface{}) bool {
	return count(arr, v...) > 0
}

// hasAll finds all the occurences in the slice
func hasAll(arr interface{}, v ...interface{}) bool { return count(arr, v...) == len(v) }

// find takes a value and gets the position in slice
// if no results, returns -1
func find(arr interface{}, v interface{}) int {
	if !isSlice(arr) {
		return 0
	}

	for i := 0; i < reflect.ValueOf(arr).Len(); i++ {
		c := reflect.ValueOf(arr)

		if c.Index(i).Interface() == v {
			return i
		}
	}

	return -1
}

// slice slices slices
func slice(arr interface{}, start, end int) interface{} {
	v := reflect.ValueOf(arr)

	if start >= v.Len() {
		return v.Interface()
	}

	if end == -1 || end > v.Len() {
		end = v.Len()
	}

	return v.Slice(start, end).Interface()
}
