package rbac

import (
	"path"
	"strings"
)

type (
	resource string
	Resource interface {
		RbacResource() string
	}

	resourceDicter interface {
		Dict() map[string]interface{}
	}
)

const (
	nsSep    = "::"
	cmpSep   = ":"
	pathSep  = "/"
	wildcard = "*"
)

// NewResource constructs untyped resource from the given string
//
// This is a utility method that should not be used for standard permission checking and granting
// it's intended to be used for testing end permission evaluation where we do not have access to the resource struct
func NewResource(s string) Resource {
	return resource(s)
}

// RbacResource returns string of an untyped resource
func (t resource) RbacResource() string {
	return string(t)
}

// ResourceType extracts 1st part of the resource
//
// ns::cmp:res/c returns ns::cmp:res
// ns::cmp:res/  returns ns::cmp:res
// ns::cmp:res   returns ns::cmp:res
func ResourceType(r string) string {
	if p := strings.Index(r, pathSep); p > 0 {
		return r[:p]
	} else {
		return r
	}
}

func ResourceComponent(r string) string {
	var (
		t  = ResourceType(r)
		ns = strings.Index(t, nsSep)
		c  = strings.LastIndex(t, cmpSep)
	)

	// make sure that we have both namespace + component separators
	if c > ns+1 && ns > -1 {
		return t[:c]
	} else {
		return t
	}
}

// match returns true if the given resource matches the given pattern
func matchResource(matcher, resource string) (m bool) {
	if matcher == resource {
		return true
	}

	m, _ = path.Match(matcher, resource)
	return
}

func hasWildcards(resource string) bool {
	return strings.Index(resource, wildcard) != -1
}

// returns level for the given resource match
// In a nutshell, level indicates number of wildcard characters
//
// More defined resources use less wildcards and are on a lower level
func level(r string) (score int) {
	var nl bool
	for l := len(r) - 1; l > strings.Index(r, pathSep); l-- {
		switch r[l] {
		case wildcard[0]:
			// nop
		case pathSep[0]:
			// found next resource reference level
			score *= 10
			nl = false
		default:
			if !nl {
				score += 1
				nl = true
			}
		}
	}

	return
}
