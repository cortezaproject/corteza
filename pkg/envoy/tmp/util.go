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

func GenericFilter(ii resource.Identifiers) genericFilter {
	f := genericFilter{}
	for i := range ii {
		if i == "" {
			continue
		}

		if refy.MatchString(i) {
			id, err := strconv.ParseUint(i, 10, 64)
			if err != nil {
				continue
			}
			f.ID = id
		} else if handy.MatchString(i) {
			f.Ref = i
		} else {
			f.Name = i
		}
	}

	return f
}

func walkResources(rr []resource.Interface, f func(r resource.Interface) error) (err error) {
	for _, r := range rr {
		err = f(r)
		if err != nil {
			return
		}
	}
	return nil
}
