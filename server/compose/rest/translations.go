package rest

import (
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
)

// checkTranslatedResources makes sure translations are only set on the given resources
//
// Resource of the translation comes from the payload, access is checked on the resource
// from the request path. Allowed resources that end with a slash are used as a prefix.
func checkTranslatedResources(tt locale.ResourceTranslationSet, allowed ...string) error {
	for _, t := range tt {
		if t == nil || !isTranslatedResourceAllowed(t.Resource, allowed...) {
			return errors.InvalidData("translated resource does not match the resource from the request")
		}
	}

	return nil
}

func isTranslatedResourceAllowed(resource string, allowed ...string) bool {
	for _, a := range allowed {
		if strings.HasSuffix(a, "/") {
			if strings.HasPrefix(resource, a) && len(resource) > len(a) {
				return true
			}
		} else if resource == a {
			return true
		}
	}

	return false
}
