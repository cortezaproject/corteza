package rbac

import (
	"sort"
	"strings"
)

type (
	// ruleIndex indexes all given RBAC rules to optimize lookup times
	//
	// The algorithm is based on the standard trie structure.
	// The max depth for a check operation is M+2 where M is the number of
	// RBAC resource path elements + component + some meta.
	ruleIndex struct {
		children map[uint64]*ruleIndexNode
	}

	ruleIndexNode struct {
		children map[string]*ruleIndexNode
		isLeaf   bool
		access   Access
		rule     *Rule

		count int
	}
)

// buildRuleIndex indexes the given rules for optimal lookups
//
// The build isn't that cleanned up but the lookup is good, I promise <3
func buildRuleIndex(rules []*Rule) (index *ruleIndex) {
	index = &ruleIndex{}
	index.add(rules...)
	return index
}

func (index *ruleIndex) add(rules ...*Rule) {
	if index.children == nil {
		index.children = make(map[uint64]*ruleIndexNode, len(rules)/2)
	}

	for _, r := range rules {
		// skip duplicates
		if index.has(r) {
			continue
		}

		if _, ok := index.children[r.RoleID]; !ok {
			index.children[r.RoleID] = &ruleIndexNode{
				children: make(map[string]*ruleIndexNode, 4),
			}
		}
		index.children[r.RoleID].count++
		n := index.children[r.RoleID]

		bits := append([]string{r.Operation}, strings.Split(r.Resource, "/")...)
		for _, b := range bits {
			if _, ok := n.children[b]; !ok {
				n.children[b] = &ruleIndexNode{
					children: make(map[string]*ruleIndexNode, 4),
				}
			}
			n.children[b].count++

			n = n.children[b]
		}

		n.isLeaf = true
		n.access = r.Access
		n.rule = r
	}
}

func (index *ruleIndex) remove(rules ...*Rule) {
	if len(rules) == 0 {
		return
	}

	for _, r := range rules {
		if _, ok := index.children[r.RoleID]; !ok {
			continue
		}

		if !index.has(r) {
			continue
		}

		bits := append([]string{r.Operation}, strings.Split(r.Resource, "/")...)
		index.removeRec(index.children[r.RoleID], bits)

		// Finishing touch cleanup
		if len(index.children[r.RoleID].children[r.Operation].children) == 0 {
			delete(index.children[r.RoleID].children, r.Operation)
		}
		if len(index.children[r.RoleID].children) == 0 {
			delete(index.children, r.RoleID)
		}
	}
}

func (index *ruleIndex) removeRec(n *ruleIndexNode, bits []string) {
	// Recursive in; decrement counters
	n.count--

	if len(bits) == 0 {
		return
	}

	n = n.children[bits[0]]
	index.removeRec(n, bits[1:])

	// Recursive out; yoink out obsolete štuff
	if len(bits) == 1 {
		if len(n.children) > 0 {
			n.children[bits[0]].isLeaf = false
			n.children[bits[0]].rule = nil
		}
		return
	}

	if n.children[bits[1]].count == 0 {
		delete(n.children, bits[1])
	}
}

func (t *ruleIndex) has(r *Rule) bool {
	return len(t.collect(true, r.RoleID, r.Operation, r.Resource)) > 0
}

func (t *ruleIndex) get(role uint64, op, res string) (out []*Rule) {
	return t.collect(false, role, op, res)
}

// get returns all RBAC rules matching these constraints
//
// The get operation's lookup complexity is the longest RBAC key + 1 for
// the operation + 1 for the role.
//
// Our longest bit will be 6 so this is essentially constant time.
func (t *ruleIndex) collect(exact bool, role uint64, op, res string) (out []*Rule) {
	if t.children == nil {
		return
	}

	if _, ok := t.children[role]; !ok {
		return
	}

	// An edge case implied by the test suite
	if op == "" && res == "" {
		out = append(out, t.children[role].children[""].children[""].rule)
		return
	}

	// Pull out the nodes for the role
	aux, ok := t.children[role]
	if !ok {
		return
	}

	aux, ok = aux.children[op]
	if !ok {
		return
	}

	return aux.get(exact, res, 0)
}

// get returns all of the rules matching these constraints
//
// Under the hood...
// We're avoiding string processing (concatenation, splitting, ...) as that can
// be a memory hog in scenarios where we're pounding this function.
//
// The from denotes the substring we've not yet processed.
func (n *ruleIndexNode) get(exact bool, res string, from int) (out []*Rule) {
	if n == nil || n.children == nil {
		return
	}

	// If we've reached the leaf node but haven't yet processed the entire resource,
	// we've reached an invalid scenario since we can't go any deeper
	to := len(res)
	if n.isLeaf && from < to {
		return
	}

	// Once from passes to, we've processed the entire resource
	if from >= to {
		if n.isLeaf {
			out = append(out, n.rule)
			return
		}
	}

	// Get the next / delimiter.
	// Clamp the index to the length of the resource.
	// Adjust the index to account the from (the start index of the remaining resource)
	nextDelim := strings.Index(res[from:to], "/")
	if nextDelim < 0 {
		nextDelim = len(res)
	} else {
		nextDelim += from
	}

	// Get RBAC rules down the actual path
	pathBit := res[from:nextDelim]
	if n.children[pathBit] != nil {
		out = append(out, n.children[pathBit].get(exact, res, nextDelim+1)...)
	}

	// Get RBAC rules down the wildcard path
	if !exact && n.children[wildcard] != nil {
		out = append(out, n.children[wildcard].get(exact, res, nextDelim+1)...)
	}

	return
}

// empty returns true if the index is empty
func (t *ruleIndex) empty() bool {
	return t == nil || t.children == nil || len(t.children) == 0
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
