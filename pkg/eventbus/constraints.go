package eventbus

import (
	"errors"
	"path"
	"regexp"
	"strings"
)

type (
	mustBeEqual struct {
		not    bool
		name   string
		values []string
	}

	mustBeLike struct {
		not    bool
		name   string
		values []string
	}

	mustMatch struct {
		not    bool
		name   string
		values []*regexp.Regexp
	}

	ConstraintMatcher interface {
		Name() string
		Values() []string
		Match(value string) bool
	}

	constraintSet []ConstraintMatcher
)

var ErrUnsupportedOp = errors.New("operator not supported")
var ErrUnsupportedName = errors.New("constraint name not supported")

func (c mustBeEqual) Name() string     { return c.name }
func (c mustBeLike) Name() string      { return c.name }
func (c mustMatch) Name() string       { return c.name }
func (c mustBeEqual) Values() []string { return c.values }
func (c mustBeLike) Values() []string  { return c.values }
func (c mustMatch) Values() []string   { return nil }

// Converts raw ConstraintMatcher into one of the
// ConstraintMatcher handlers
func ConstraintMaker(name, op string, vv ...string) (ConstraintMatcher, error) {
	switch strings.ToLower(op) {
	case "", "eq", "=", "==", "===":
		return MustBeEqual(name, vv...)
	case "not eq", "ne", "!=", "!==":
		return MustNotBeEqual(name, vv...)
	case "like":
		return MustBeLike(name, vv...)
	case "not like":
		return MustNotBeLike(name, vv...)
	case "~":
		return MustMatch(name, vv...)
	case "!~":
		return MustNotMatch(name, vv...)
	default:
		return nil, ErrUnsupportedOp
	}
}

func MustMakeConstraint(name, op string, vv ...string) ConstraintMatcher {
	c, err := ConstraintMaker(name, op, vv...)
	if err != nil {
		panic(err)
	}

	return c
}

func MustBeEqual(name string, vv ...string) (ConstraintMatcher, error) {
	return &mustBeEqual{name: name, values: vv}, nil
}

func MustNotBeEqual(name string, vv ...string) (ConstraintMatcher, error) {
	return &mustBeEqual{name: name, values: vv, not: true}, nil
}

func (c mustBeEqual) Match(value string) bool {
	for _, v := range c.values {
		if value == v {
			return !c.not
		}
	}

	return c.not
}

func mustLikeMaker(name string, not bool, vv ...string) (*mustBeLike, error) {
	var (
		m = &mustBeLike{
			name:   name,
			values: make([]string, len(vv)),
			not:    not,
		}
	)

	for i, v := range vv {
		v = strings.ReplaceAll(v, "%", "*")
		v = strings.ReplaceAll(v, "_", "?")
		m.values[i] = v
	}

	return m, nil
}

func MustBeLike(name string, vv ...string) (ConstraintMatcher, error) {
	return mustLikeMaker(name, false, vv...)
}

func MustNotBeLike(name string, vv ...string) (ConstraintMatcher, error) {
	return mustLikeMaker(name, true, vv...)
}

func (c mustBeLike) Match(value string) bool {
	for _, v := range c.values {
		if m, _ := path.Match(v, value); m {
			return !c.not
		}
	}

	return c.not
}

func mustMatchMaker(name string, not bool, vv ...string) (*mustMatch, error) {
	var (
		m = &mustMatch{
			name:   name,
			values: make([]*regexp.Regexp, len(vv)),
			not:    not,
		}
		err error
	)

	for i, v := range vv {
		if m.values[i], err = regexp.Compile(v); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func MustMatch(name string, vv ...string) (ConstraintMatcher, error) {
	return mustMatchMaker(name, false, vv...)
}

func MustNotMatch(name string, vv ...string) (ConstraintMatcher, error) {
	return mustMatchMaker(name, true, vv...)

}

func (c mustMatch) Match(value string) bool {
	for _, v := range c.values {
		if v.MatchString(value) {
			return !c.not
		}
	}

	return c.not
}

func MatchFirst(checks ...func() bool) bool {
	for _, check := range checks {
		if check() {
			return true
		}
	}

	return false
}
