package expr

import (
	"fmt"
	"reflect"
	"strings"

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
func push(arr interface{}, nn ...interface{}) (out interface{}, err error) {
	if arr == nil {
		// If base is empty, return pushed items directly
		return nn, nil
	} else if i, is := arr.([]interface{}); is {
		// Simple append if we're dealing with []interface{} base
		return append(i, nn...), nil
	} else if arr, err = toSlice(arr); err != nil {
		return
	}

	if stv, is := arr.([]TypedValue); is {
		// slice of typed values, this will make things easier
		for _, n := range nn {
			if tv, is := n.(TypedValue); is {
				stv = append(stv, tv)
			} else {
				// wrap unknown types...
				stv = append(stv, Must(Typify(n)))
			}
		}

		return stv, nil
	}

	var (
		c    = reflect.ValueOf(arr)
		nval = reflect.MakeSlice(c.Type(), c.Len()+len(nn), c.Cap()+len(nn))
	)

	reflect.Copy(nval, c)

	for i, n := range nn {
		nt := reflect.ValueOf(n).Type()
		it := nval.Index(c.Len() + i).Type()
		if nt != it {
			return nil, fmt.Errorf("can not push %v to %v slice", nt, it)
		}

		nval.Index(c.Len() + i).Set(reflect.ValueOf(n))
	}

	return nval.Interface(), nil
}

// pop takes the last value in slice, does not modify original
func pop(arr interface{}) (out interface{}, err error) {
	if arr, err = toSlice(arr); err != nil {
		return
	}

	c := reflect.ValueOf(arr)

	if c.Len() == 0 {
		return nil, nil
	}

	return c.Index(c.Len() - 1).Interface(), nil
}

// shifts takes the first value in slice, does not modify original
func shift(arr interface{}) (out interface{}, err error) {
	if arr, err = toSlice(arr); err != nil {
		return
	}

	c := reflect.ValueOf(arr)

	if c.Len() == 0 {
		return nil, nil
	}

	return c.Index(0).Interface(), nil
}

// count gets the count of occurrences in the first params(can be string, slice)
// If given only one parameter then it returns length as count
func count(arr interface{}, v ...interface{}) (count int, err error) {
	typeErr := fmt.Errorf("unexpected type: %T, expecting slice", arr)
	arr = UntypedValue(arr)

	var (
		occ int
		c   = reflect.ValueOf(arr)
	)

	switch c.Kind() {
	case reflect.String:
		if len(v) == 0 {
			return len(c.String()), nil
		}

		for _, ww := range v {
			count += strings.Count(c.String(), reflect.ValueOf(ww).String())
		}
		return

	case reflect.Slice:
		if len(v) == 0 {
			return c.Len(), nil
		}

		for _, vv := range v {
			if occ, err = find(arr, vv); err != nil {
				return 0, err
			} else if occ != -1 {
				count++
			}
		}

	default:
		return 0, typeErr
	}

	return
}

// has finds any occurrence of the values in slice or key in a map
func has(arr interface{}, vv ...interface{}) (b bool, err error) {
	arr = UntypedValue(arr)

	if isMap(arr) {
		for _, v := range vv {
			if reflect.ValueOf(arr).MapIndex(reflect.ValueOf(v)).IsValid() {
				return true, nil
			}
		}

		return
	}

	var c int
	if c, err = count(arr, vv...); err != nil {
		return
	}

	return c > 0, nil
}

// hasAll finds all the occurrences in the slice
func hasAll(arr interface{}, v ...interface{}) (b bool, err error) {
	if arr, err = toSlice(arr); err != nil {
		return
	}

	var c int
	if c, err = count(arr, v...); err != nil {
		return
	}

	return c == len(v), nil
}

// find takes a value and gets the position in slice
// if no results, returns -1
func find(arr interface{}, v interface{}) (p int, err error) {
	if arr, err = toSlice(arr); err != nil {
		return
	}

	for p = 0; p < reflect.ValueOf(arr).Len(); p++ {
		c := reflect.ValueOf(arr)

		if c.Index(p).Interface() == v {
			return
		}
	}

	return -1, nil
}

// slice slices slices
func slice(arr interface{}, start, end int) interface{} {
	arr = UntypedValue(arr)

	v := reflect.ValueOf(arr)

	if start >= v.Len() {
		return v.Interface()
	}

	if end == -1 || end > v.Len() {
		end = v.Len()
	}

	return v.Slice(start, end).Interface()
}
