package rbac

import (
	"strings"
)

type (
	wrapperIndex struct {
		// indexed permits only max level identifiers
		indexed map[string]bool
		index   *ruleIndex
	}
)

// add indexes the given rules
func (svc *wrapperIndex) add(resource string, rules ...*Rule) (added bool) {
	if svc.indexed == nil {
		svc.indexed = make(map[string]bool, 24)
	}

	if svc.index == nil {
		svc.index = &ruleIndex{}
	}

	// Since we're only allowed to index under full resource identifiers
	// we'll only optionally update indexes if something new comes in.
	if strings.Contains(resource, "*") {
		return svc.addWild(resource, rules...)
	} else {
		return svc.addPlain(resource, rules...)
	}
}

// addWild handles scenario where we would grant permissions for a wildcard
//
// In case of a wild card we need to check if any matching resource falls under it
// if so, add it to the index, if not, ignore.
//
// In no case should we add this to the indexed map since it only permits max lvl identifiers.
func (svc *wrapperIndex) addWild(resource string, rules ...*Rule) (added bool) {
	give := false
	pp := strings.SplitN(resource, "*", 2)
	resource = strings.TrimRight(pp[0], "/")

	for r := range svc.indexed {
		if strings.HasPrefix(r, resource) {
			give = true
			break
		}
	}

	if !give {
		return false
	}

	svc.index.add(rules...)
	return true
}

// addPlain handles scenario where we specify rules with a max level resource identifier (no wildcards)
func (svc *wrapperIndex) addPlain(resource string, rules ...*Rule) (added bool) {
	svc.indexed[resource] = true
	svc.index.add(rules...)

	return true
}

// get returns the rules for the given role, operation and resource
func (svc *wrapperIndex) get(op string, res string) (out []*Rule) {
	// @note we'll expect the state is good and no nil checks are needed
	return svc.index.get(op, res)
}

// getSize returns the number of indexed resources (may not match to number of rules)
func (svc *wrapperIndex) getSize() int {
	return len(svc.indexed)
}

// isIndexed returns true if the resource is either indexed or potentially indexed
//
// If we're providing a max level resource identifier, it must occur in the index
// If we're providing a wildcard, we always assume it's in there
//
// # Underlying functions need to respect this
//
// @todo consider keeping track of prefixes so we can know for a fact.
// It doesn't really matter at this point since referencing functions don't care about this
func (svc *wrapperIndex) isIndexed(resource string) (ok bool) {
	// In case of wildcards, assume we have it indexed; further functions need
	// to handle this properly
	if strings.Contains(resource, "*") {
		return true
	}

	if svc.indexed == nil {
		return false
	}

	return svc.indexed[resource]
}
