package filter

import (
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
)

const (
	PreFilterWeight = iota
	ProcesserWeight
	PostFilterWeight
)

func FilterWeight(w int, t types.FilterKind) int {
	mul := PreFilterWeight

	switch t {
	case types.PreFilter:
		mul = PreFilterWeight
	case types.Processer:
		mul = ProcesserWeight
	case types.PostFilter:
		mul = PostFilterWeight
	}

	return mul*100 + w
}
