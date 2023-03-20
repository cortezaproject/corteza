package envoyx

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	rbacer interface {
		RbacResource() string
	}
)

func RBACRulesForNodes(rr rbac.RuleSet, nn ...*Node) (rules NodeSet, err error) {
	rules = make(NodeSet, 0, len(rr)/2)
	dups := make(map[uint64]map[string]map[string]bool)

	for _, n := range nn {
		if n.Placeholder {
			continue
		}

		c, ok := n.Resource.(rbacer)
		if !ok {
			continue
		}

		// Split up the path of this resource
		//
		// @todo move over to those generated functions
		resPath := splitResourcePath(c.RbacResource())

		// Find all of the rules that fall under this resource
		for _, r := range rr {
			if r.RoleID == 0 {
				// Can't exist; skip to avoid edge cases
				continue
			}

			// Split up the path of the rule
			//
			// @todo move over to that generated function
			rulePath := splitResourcePath(r.Resource)

			// Don't handle rules that are not a subset of the resource
			// @todo this should be handled by the parent probably
			if len(rulePath) > 0 && rulePath[0] == "*" || rulePath[0] == "" {
				continue
			}

			if !isPathSubset(rulePath, resPath, true) {
				// Mismatch; skip
				continue
			}

			// Check if this rule has already been seen
			if dups[r.RoleID] != nil && dups[r.RoleID][r.Resource] != nil && dups[r.RoleID][r.Resource][r.Operation] {
				continue
			}

			// Parse the path so we can process it further
			ruleRt, res, path, err := ParseRule(r.Resource)
			if err != nil {
				return nil, err
			}

			if ruleRt != n.ResourceType {
				// Type missmatch; skip
				continue
			}

			// Get the refs
			rf := make(map[string]Ref, 2)

			for i, ref := range append(path, res) {
				// Whenever you'd use a wildcard, it will produce a nil so it
				// needs to be skipped
				if ref == nil {
					continue
				}

				ref.Scope = n.Scope

				// @todo make the thing not a pointer
				rf[fmt.Sprintf("Path.%d", i)] = *ref
			}

			// Ref to the rule
			rf["RoleID"] = Ref{
				ResourceType: types.RoleResourceType,
				Identifiers:  MakeIdentifiers(r.RoleID),
			}

			rules = append(rules, &Node{
				Resource: r,

				ResourceType: rbac.RuleResourceType,
				References:   rf,
				Scope:        n.Scope,
			})

			// Update the dup checking index
			if dups[r.RoleID] == nil {
				dups[r.RoleID] = make(map[string]map[string]bool)
			}
			if dups[r.RoleID][r.Resource] == nil {
				dups[r.RoleID][r.Resource] = make(map[string]bool)
			}
			dups[r.RoleID][r.Resource][r.Operation] = true
		}
	}

	return
}

func splitResourcePath(p string) []string {
	return strings.Split(p, "/")[1:]
}

func isPathSubset(rulePath, resPath []string, wildcards bool) bool {
	if len(rulePath) == 0 && len(resPath) == 0 {
		return true
	}

	// The lengths must match since missing bits are replaced with wildcards
	if len(rulePath) != len(resPath) {
		return false
	}

	for i := 0; i < len(resPath); i++ {
		if wildcards && rulePath[i] == "*" {
			// Rule matches everything from now on; if we got this far, we're good
			return true
		}
		if rulePath[i] != resPath[i] {
			return false
		}
	}

	return true
}
