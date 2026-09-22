package expr

import (
	"fmt"
	"strings"
)

var (
	invalidPathErr = fmt.Errorf("invalid path format")
)

func PathBase(path string) string {
	return strings.Split(path, ".")[0]
}

func Assign(base TypedValue, path string, val TypedValue) (err error) {
	if _, ok := base.(FlatStore); ok {
		return assignFlat(base, path, val)
	}
	return assignPath(base, path, val)
}

func assignFlat(base TypedValue, key string, val TypedValue) error {
	switch fa := base.(type) {
	case FieldAssigner:
		return fa.AssignFieldValue(key, val)
	case DeepFieldAssigner:
		return fa.AssignFieldValue(&flatPath{key: key}, val)
	default:
		return fmt.Errorf("cannot set on unexpected type: %T", base)
	}
}

// flatPath is a single-segment Pather for flat key assignment — no path splitting.
type flatPath struct {
	key  string
	done bool
}

func (p *flatPath) String() string { return p.key }
func (p *flatPath) More() bool     { return !p.done }
func (p *flatPath) IsLast() bool   { return true }
func (p *flatPath) Get() string    { return p.key }
func (p *flatPath) Rest() string   { return "" }
func (p *flatPath) Next() error    { p.done = true; return nil }

func assignPath(base TypedValue, path string, val TypedValue) (err error) {
	pp := Path(path)
	err = pp.Next()
	if err != nil {
		return
	}

	if !pp.More() {
		panic("setting value with empty path")
	}

	var (
		key = ""
	)

	// descend lower by the path but
	// stop before the last part of the path

	for !pp.IsLast() {
		switch s := base.(type) {
		case DeepFieldAssigner:
			return s.AssignFieldValue(pp, val)

		case FieldSelector:
			key = pp.Get()
			err = pp.Next()
			if err != nil {
				return
			}

			if base, err = s.Select(key); err != nil {
				return err
			}

		default:
			return fmt.Errorf("cannot set value on %s with path '%s'", base.Type(), path)

		}
	}

	key = pp.Get()

	// try with field setter first
	// if not a FieldSetter it has to be a Selector
	// that returns TypedValue that we can set
	switch setter := base.(type) {
	case DeepFieldAssigner:
		return setter.AssignFieldValue(pp, val)

	case FieldAssigner:
		return setter.AssignFieldValue(key, val)

	case FieldSelector:
		if base, err = setter.Select(key); err != nil {
			return err
		}

		return base.Assign(val)

	default:
		return fmt.Errorf("%T does not support value assigning with '%s'", base, path)
	}
}

func Select(base TypedValue, path string) (out TypedValue, err error) {
	pp := Path(path)
	err = pp.Next()
	if err != nil {
		return
	}

	if !pp.More() {
		panic("setting value with empty path")
	}

	var (
		failure = fmt.Errorf("cannot get value from %s with path '%s'", base.Type(), path)
		key     string
	)

	// descend lower by the path but
	// stop before the last part of the path
	for pp.More() {
		s, is := base.(FieldSelector)
		if !is {
			return nil, failure
		}

		key = pp.Get()
		err = pp.Next()
		if err != nil {
			return
		}

		if base, err = s.Select(key); err != nil {
			return nil, err
		}
	}

	return base, nil
}
