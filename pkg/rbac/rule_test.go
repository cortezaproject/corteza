package rbac

import (
	"sort"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestRuleSetSort(t *testing.T) {
	var (
		req = require.New(t)
		rr  = RuleSet{
			{Resource: ":::/*/*/*"},
			{Resource: ":::/1/2/3"},
			{Resource: ":::/1/2/*"},
			{Resource: ":::/1/*/*"},
			{Resource: ":::/1/*/3"},
			{Resource: ":::/*/*/3"},
			{Resource: ":::/*/2/*"},
		}

		c int = 0
		i     = func() int {
			c++
			return c - 1
		}
	)

	req.NotNil(rr)
	sort.Sort(rr)
	spew.Dump(rr)
	c = i()
	req.Equal(":::/1/2/3", rr[i()].Resource)
	req.Equal(":::/1/*/3", rr[i()].Resource)
	req.Equal(":::/*/*/3", rr[i()].Resource)
	req.Equal(":::/1/2/*", rr[i()].Resource)
	req.Equal(":::/*/2/*", rr[i()].Resource)
	req.Equal(":::/1/*/*", rr[i()].Resource)
	req.Equal(":::/*/*/*", rr[i()].Resource)
}
