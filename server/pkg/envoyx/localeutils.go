package envoyx

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	localer interface {
		ResourceTranslation() string
	}
)

func ResourceTranslationsForNodes(tt types.ResourceTranslationSet, nn ...*Node) (translations NodeSet, err error) {
	translations = make(NodeSet, 0, len(tt)/2)
	dups := make(map[types.Lang]map[string]map[string]bool)

	for _, n := range nn {
		if n.Placeholder {
			continue
		}

		c, ok := n.Resource.(localer)
		if !ok {
			continue
		}

		// Split up the path of this resource
		//
		// @todo move over to those generated functions
		resPath := splitResourcePath(c.ResourceTranslation())

		// Find all of the translations that fall under this resource
		for _, r := range tt {

			// Split up the path of the rule
			//
			// @todo move over to that generated function
			rulePath := splitResourcePath(r.Resource)
			// @note resource translations don't support wildcards
			if !isPathSubset(rulePath, resPath, false) {
				// Mismatch; skip
				continue
			}

			// Check if this translation has already been seen
			if dups[r.Lang] != nil && dups[r.Lang][r.Resource][r.K] {
				continue
			}

			// Parse the path so we can process it further
			// @todo make a generic function for RBAC rules and res. tr.
			localeRt, res, path, err := ParseRule(r.Resource)
			if err != nil {
				return nil, err
			}

			if localeRt != n.ResourceType {
				// Resource type missmatch; skip
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

			translations = append(translations, &Node{
				Resource: r,

				ResourceType: types.ResourceTranslationResourceType,
				References:   rf,
				Scope:        n.Scope,
			})

			// Update the dup checking index
			if dups[r.Lang] == nil {
				dups[r.Lang] = make(map[string]map[string]bool)
			}
			if dups[r.Lang][r.Resource] == nil {
				dups[r.Lang][r.Resource] = make(map[string]bool)
			}
			dups[r.Lang][r.Resource][r.K] = true
		}
	}

	return
}
