package wfexec

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVariables_Merge(t *testing.T) {
	var (
		req = require.New(t)
		vv  Variables
	)

	req.Empty(vv)

	vv = vv.Merge(Variables{"a": 1})
	req.Contains(vv, "a")

	vv = Variables{"a": 1}.Merge()
	req.Contains(vv, "a")

	vv = Variables{"a": 1}.Merge(Variables{"b": 2}, Variables{"c": 3})
	req.Contains(vv, "a")
	req.Contains(vv, "b")
	req.Contains(vv, "b")
}
