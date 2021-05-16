package rbac

import (
	"path"
	"strings"
)

type (
	Resource interface {
		RbacResource() string
	}
)

func ResourceSchema(r string) string {
	i := strings.Index(r, ":")
	if i < 0 {
		return ""
	}

	return r[:i]
}

func matchResource(matcher, resource string) (m bool) {
	if level(matcher) == 0 {
		return matcher == resource
	}

	m, _ = path.Match(matcher, resource)
	return
}

// returns level for the given resource match
// In a nutshell, level indicates number of wildcard characters
//
// More defined resources use less wildcards and are on a lower level
func level(r string) int {
	return strings.Count(r, string("*"))
}
