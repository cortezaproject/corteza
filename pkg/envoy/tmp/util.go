package tmp

import (
	"regexp"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/id"
)

type (
	genericFilter struct {
		ID   uint64
		Ref  string
		Name string
	}
)

var (
	refy  = regexp.MustCompile(`^[1-9](\d*)$`)
	handy = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$`)

	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func GetGenericFilter(ii resource.Identifiers) genericFilter {
	f := genericFilter{}
	for id := range ii {
		if id == "" {
			continue
		}
		if refy.MatchString(id) {
			id, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				continue
			}
			f.ID = id
		} else if handy.MatchString(id) {
			f.Ref = id
		} else {
			f.Name = id
		}
	}

	return f
}
