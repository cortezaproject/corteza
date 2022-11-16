package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	APIGateway struct {
		*base
		Res *types.ApigwRoute

		Filters []*APIGatewayFilter
	}

	APIGatewayFilter struct {
		*base
		Res *types.ApigwFilter
	}
)

func NewAPIGateway(res *types.ApigwRoute) *APIGateway {
	// @todo we could use method + path combo to uniquely identify them also
	r := &APIGateway{
		base: &base{},
	}
	r.SetResourceType(types.ApigwRouteResourceType)
	r.Res = res

	r.AddIdentifier(identifiers(res.ID)...)

	// Initial stamps
	r.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))
	us := MakeUserstampsCUDO(res.CreatedBy, res.UpdatedBy, res.DeletedBy, 0)
	r.SetUserstamps(us)

	return r
}

func (r *APIGateway) Resource() interface{} {
	return r.Res
}

func (r *APIGateway) RBACParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.ApigwRouteRbacResourceTpl(), types.ApigwRouteResourceType, firstOkString(strconv.FormatUint(r.Res.ID, 10)))

	return
}

func (r *APIGateway) AddGatewayFilter(res *types.ApigwFilter) *APIGatewayFilter {
	f := &APIGatewayFilter{
		base: &base{},
	}

	f.Res = res

	// Initial stamps
	f.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))
	f.SetUserstamps(MakeUserstampsCUDO(res.CreatedBy, res.UpdatedBy, res.DeletedBy, 0))

	r.Filters = append(r.Filters, f)

	return f
}

func (r *APIGateway) SysID() uint64 {
	return r.Res.ID
}

// FindAPIGateway looks for the ApigwRoute in the resource set
func FindAPIGateway(rr InterfaceSet, ii Identifiers) (ns *types.ApigwRoute) {
	var wfRes *APIGateway

	rr.Walk(func(r Interface) error {
		wr, ok := r.(*APIGateway)
		if !ok {
			return nil
		}

		if wr.Identifiers().HasAny(ii) {
			wfRes = wr
		}
		return nil
	})

	// Found it
	if wfRes != nil {
		return wfRes.Res
	}
	return nil
}

func APIGatewayErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("automation apu gateway unresolved %v", ii.StringSlice())
}
