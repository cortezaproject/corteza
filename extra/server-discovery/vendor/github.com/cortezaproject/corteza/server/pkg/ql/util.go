package ql

import (
	"regexp"
	"strings"
)

var (
	truthy = regexp.MustCompile(`^(t(rue)?|y(es)?|1)$`)
)

// Check what boolean value the given string conforms to
func evalBool(v string) bool {
	return truthy.MatchString(strings.ToLower(v))
}

// Check if the given string is a float
func isFloaty(v string) bool {
	return strings.Contains(v, ".")
}

func MergeAnd(a, b *ASTNode) *ASTNode {
	return merger(a, b, "and")
}

func MergeOr(a, b *ASTNode) *ASTNode {
	return merger(a, b, "or")
}

func merger(a, b *ASTNode, ref string) *ASTNode {
	// 1. merge the two
	aa := make(ASTNodeSet, 0, 2)

	// It needs to be under a group, so we get an `(a) and/or (b)`
	if a != nil {
		aa = append(aa, &ASTNode{Ref: "group", Args: ASTNodeSet{a}})
	}
	if b != nil {
		aa = append(aa, &ASTNode{Ref: "group", Args: ASTNodeSet{b}})
	}

	// 2. flatten
	// @todo do some more in-depth processing?
	if len(aa) == 1 {
		// this [0][0] will always hold if we get to this point
		return aa[0].Args[0]
	}
	if len(aa) == 0 {
		return nil
	}

	return &ASTNode{
		Ref:  ref,
		Args: aa,
	}
}
