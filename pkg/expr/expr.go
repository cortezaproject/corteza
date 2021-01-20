package expr

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	pathDelimiter = "."
)

var (
	invalidPathErr = fmt.Errorf("invalid path format")
)

func PathSplit(path string) ([]string, error) {
	out := make([]string, 0)
	s := bufio.NewScanner(strings.NewReader(path))
	s.Split(pathSplitter)

	for s.Scan() {
		if len(s.Text()) == 0 {
			return nil, invalidPathErr
		}
		out = append(out, s.Text())
	}

	if s.Err() != nil {
		return nil, s.Err()
	}

	return out, nil
}

func pathSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for i := 0; i < len(data); i += 1 {
		switch data[i] {
		case '.', '[':
			return i + 1, data[start:i], nil
		case ']':
			// When at closing bracket but not at the end, make sure we properly split the token
			if i == len(data)-1 {
				return i + 1, data[start:i], nil
			}

			if data[i+1] != '.' {
				return 0, nil, invalidPathErr
			}

			return i + 2, data[start:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func PathBase(path string) string {
	return strings.Split(path, ".")[0]
}

func Assign(base TypedValue, path string, val interface{}) error {
	pp, err := PathSplit(path)
	if err != nil {
		return err
	}

	if len(pp) == 0 {
		panic("setting value with empty path")
	}

	var (
		key = pp[0]
	)

	// descend lower by the path but
	// stop before the last part of the path
	for len(pp) > 1 {
		switch s := base.(type) {
		case DeepFieldAssigner:
			return s.AssignFieldValue(pp, val)

		case FieldSelector:
			key, pp = pp[0], pp[1:]
			if base, err = s.Select(key); err != nil {
				return err
			}

		default:
			return fmt.Errorf("can not set value on %s with path '%s'", base.Type(), path)

		}
	}

	key = pp[0]

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

func Select(base TypedValue, path string) (TypedValue, error) {
	pp, err := PathSplit(path)
	if err != nil {
		return nil, err
	}

	if len(pp) == 0 {
		panic("selecting value with empty path")
	}

	var (
		failure = fmt.Errorf("can not get value from %s with path '%s'", base.Type(), path)
		key     string
	)

	// descend lower by the path but
	// stop before the last part of the path
	for len(pp) > 0 {
		s, is := base.(FieldSelector)
		if !is {
			return nil, failure
		}

		key, pp = pp[0], pp[1:]
		if base, err = s.Select(key); err != nil {
			return nil, err
		}

	}

	return base, nil
}
