package rbac

import (
	"path"
	"strings"
)

type (
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

func matchResource(matcher, resource string) (m bool) {
	if matcher == resource {
		// if resources match make sure no wildcards are resent
		return strings.Index(resource, wildcard) == -1
	}

	m, _ = path.Match(matcher, resource)
	return
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
