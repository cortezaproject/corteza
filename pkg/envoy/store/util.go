package store

import (
	"regexp"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/id"
)

type (
	genericFilter struct {
		id     uint64
		handle string
		name   string
	}
)

var (
	// simple uint check.
	// we'll use the pkg/handle to check for handles.
	refy = regexp.MustCompile(`^[1-9](\d*)$`)

	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

// makeGenericFilter is a helper to determine the base resource filter.
//
// It attempts to determine an identifier, handle, and name.
func makeGenericFilter(ii resource.Identifiers) (f genericFilter) {
	for i := range ii {
		if i == "" {
			continue
		}

		if refy.MatchString(i) && f.id <= 0 {
			id, err := strconv.ParseUint(i, 10, 64)
			if err != nil {
				continue
			}
			f.id = id
		} else if handle.IsValid(i) && f.handle == "" {
			f.handle = i
		} else if f.name == "" {
			f.name = i
		}
	}

	return f
}
