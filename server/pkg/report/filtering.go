package report

import "github.com/cortezaproject/corteza/server/pkg/qlng"

func (a *Filter) mergeAnd(b *Filter) *Filter {
	return merger(a, b, "and")
}

func (a *Filter) mergeOr(b *Filter) *Filter {
	return merger(a, b, "or")
}

func merger(a, b *Filter, ref string) *Filter {
	// 1. merge the two
	aa := make(qlng.ASTNodeSet, 0, 2)

	// It needs to be under a group, so we get an `(a) and/or (b)`
	if a != nil {
		aa = append(aa, &qlng.ASTNode{Ref: "group", Args: qlng.ASTNodeSet{a.ASTNode}})
	}
	if b != nil {
		aa = append(aa, &qlng.ASTNode{Ref: "group", Args: qlng.ASTNodeSet{b.ASTNode}})
	}

	// 2. flatten
	// @todo do some more in-depth processing?
	if len(aa) == 1 {
		return &Filter{
			// this [0][0] will always hold if we get to this point
			ASTNode: aa[0].Args[0],
		}
	}
	if len(aa) == 0 {
		return nil
	}

	return &Filter{
		ASTNode: &qlng.ASTNode{
			Ref:  ref,
			Args: aa,
		},
	}
}
