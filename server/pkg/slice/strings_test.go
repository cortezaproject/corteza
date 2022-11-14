package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntersectStrings(t *testing.T) {
	cases := []struct {
		name string
		a    []string
		b    []string
		c    []string
	}{
		{
			"empty",
			[]string{},
			[]string{},
			[]string{},
		},
		{
			"none",
			[]string{"a"},
			[]string{"b"},
			[]string{},
		},
		{
			"some",
			[]string{"a", "b"},
			[]string{"a", "c"},
			[]string{"a"},
		},
		{
			"all",
			[]string{"a", "b"},
			[]string{"a", "b"},
			[]string{"a", "b"},
		},
		{
			"dups",
			[]string{"a", "b", "b", "b"},
			[]string{"a", "b"},
			[]string{"a", "b"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.EqualValues(t, IntersectStrings(c.a, c.b), c.c)
		})
	}
}

func TestToStringBoolMap(t *testing.T) {
	cases := []struct {
		name string
		i    []string
		o    map[string]bool
	}{
		{
			"empty",
			[]string{},
			map[string]bool{},
		},
		{
			"some",
			[]string{"a"},
			map[string]bool{"a": true},
		},
		{
			"many",
			[]string{"a", "b", "c", "c", "d"},
			map[string]bool{"a": true, "b": true, "c": true, "d": true},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.EqualValues(t, ToStringBoolMap(c.i), c.o)
		})
	}
}

func TestHasString(t *testing.T) {
	cases := []struct {
		name string
		ss   []string
		s    string
		o    bool
	}{
		{
			"empty",
			[]string{},
			"a",
			false,
		},
		{
			"has not",
			[]string{"a"},
			"b",
			false,
		},
		{
			"has",
			[]string{"a"},
			"a",
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.EqualValues(t, HasString(c.ss, c.s), c.o)
		})
	}
}

func TestPluckStrung(t *testing.T) {
	cases := []struct {
		name string
		ss   []string
		ff   []string
		o    []string
	}{
		{
			"empty",
			[]string{},
			[]string{},
			[]string{},
		},
		{
			"some",
			[]string{"a", "b"},
			[]string{"a"},
			[]string{"b"},
		},
		{
			"all",
			[]string{"a", "b"},
			[]string{"a", "b"},
			[]string{},
		},
		{
			"not there",
			[]string{"a"},
			[]string{"b"},
			[]string{"a"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.EqualValues(t, PluckString(c.ss, c.ff...), c.o)
		})
	}
}
