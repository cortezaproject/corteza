package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// Application represents a Application
	Application struct {
		*base
		Res *types.Application
	}
)

func NewApplication(res *types.Application) *Application {
	r := &Application{base: &base{}}
	r.SetResourceType(types.ApplicationResourceType)
	r.Res = res

	r.AddIdentifier(identifiers(res.Name, res.ID)...)

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))
	// Initial userstamps
	if res.OwnerID > 0 {
		r.SetUserstamps(&Userstamps{
			OwnedBy: &Userstamp{UserID: res.OwnerID},
		})
	}

	return r
}

func (r *Application) Resource() interface{} {
	return r.Res
}

func (r *Application) SysID() uint64 {
	return r.Res.ID
}

func (r *Application) RBACParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.ApplicationRbacResourceTpl(), types.ApplicationResourceType, firstOkString(strconv.FormatUint(r.Res.ID, 10)))

	return
}

// FindApplication looks for the app in the resource set
func FindApplication(rr InterfaceSet, ii Identifiers) (ap *types.Application) {
	var apRes *Application

	rr.Walk(func(r Interface) error {
		ar, ok := r.(*Application)
		if !ok {
			return nil
		}

		if ar.Identifiers().HasAny(ii) {
			apRes = ar
		}

		return nil
	})

	// Found it
	if apRes != nil {
		return apRes.Res
	}

	return nil
}

func ApplicationErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("application unresolved %v", ii.StringSlice())
}
