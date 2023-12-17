package rbac

import (
	"sort"
	"strconv"
	"strings"
)

type (
	// ruleIndex indexes all given RBAC rules to optimize lookup times
	//
	// The algorithm is based on the standard trie structure.
	// The max depth for a check operation is M+2 where M is the number of
	// RBAC resource path elements + component + some meta.
	ruleIndex struct {
		root *ruleIndexNode
	}

	ruleIndexNode struct {
		children map[string]*ruleIndexNode

		isLeaf bool
		access Access
		rule   *Rule
	}
)

func mkInitial(op string, roleID uint64) (out string) {
	return op + "-" + strconv.FormatUint(roleID, 10)
}

func explodeThing(r string) (out []string) {
	return strings.Split(r, "/")
}

// buildRuleIndex indexes the given rules for optimal lookups
func buildRuleIndex(rules []*Rule) (ix *ruleIndex) {
	ix = &ruleIndex{}

	for _, r := range rules {
		n := ix.root
		if n == nil {
			n = &ruleIndexNode{children: make(map[string]*ruleIndexNode, 2)}
			ix.root = n
		}

		bits := append([]string{mkInitial(r.Operation, r.RoleID)}, explodeThing(r.Resource)...)
		for _, b := range bits {
			if _, ok := n.children[b]; !ok {
				n.children[b] = &ruleIndexNode{
					children: make(map[string]*ruleIndexNode, 4),
				}
			}

			n = n.children[b]
		}

		n.isLeaf = true
		n.access = r.Access
		n.rule = r
	}

	return ix
}

func (t *ruleIndex) matchingRule(role uint64, op, res string) (out *Rule) {
	set := RuleSet(t.get(role, op, res))
	sort.Sort(set)

	for _, s := range set {
		if s.Access == Inherit {
			continue
		}

		return s
	}

	return nil
}

// get returns all rules matching the given params
func (t *ruleIndex) get(role uint64, op, res string) (out []*Rule) {
	if t.root == nil {
		return
	}

	bits := append([]string{mkInitial(op, role)}, explodeThing(res)...)
	return t.root.get(bits)
}

func (n *ruleIndexNode) get(bits []string) (out []*Rule) {
	if n == nil || n.children == nil {
		return
	}

	if n.isLeaf && len(bits) > 0 {
		return
	}

	if len(bits) == 0 {
		if n.isLeaf {
			out = append(out, n.rule)
			return
		}
	}

	b := bits[0]
	bits = bits[1:]

	if n.children[b] != nil {
		out = append(out, n.children[b].get(bits)...)
	}

	if n.children[wildcard] != nil {
		out = append(out, n.children[wildcard].get(bits)...)
	}

	return
}

// empty returns true if the index is empty
func (t *ruleIndex) empty() bool {
	return t == nil || t.root == nil
}
