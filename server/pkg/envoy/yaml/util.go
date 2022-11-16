package yaml

import (
	"strconv"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

var (
	// wrapper around nextID that will aid testing
	nextID = func() uint64 {
		return id.Next()
	}

	// wrapper around time.Now() that will aid testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}
)

func firstOkString(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func resolveUserstamps(rr []resource.Interface, us *resource.Userstamps) (*resource.Userstamps, error) {
	if us == nil {
		return nil, nil
	}

	fetch := func(us *resource.Userstamp) (*resource.Userstamp, error) {
		if us == nil || (us.UserID == 0 && us.Ref == "") {
			return nil, nil
		}

		// This one can be considered as valid
		if us.Ref != "" && us.UserID > 0 && us.U != nil {
			return us, nil
		}

		ii := resource.MakeIdentifiers()

		if us.UserID > 0 {
			ii = ii.Add(strconv.FormatUint(us.UserID, 10))
		}
		if us.Ref != "" {
			ii = ii.Add(us.Ref)
		}

		u := resource.FindUser(rr, ii)
		if u == nil {
			return nil, resource.UserErrUnresolved(ii)
		}

		return resource.MakeUserstamp(u), nil
	}
	var err error
	us.CreatedBy, err = fetch(us.CreatedBy)
	us.UpdatedBy, err = fetch(us.UpdatedBy)
	us.DeletedBy, err = fetch(us.DeletedBy)
	us.OwnedBy, err = fetch(us.OwnedBy)
	us.RunAs, err = fetch(us.RunAs)

	if err != nil {
		return nil, err
	}

	return us, nil
}
