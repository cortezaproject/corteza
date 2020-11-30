package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	// ComposeChart represents a ComposeChart
	ComposeChart struct {
		*base
		Res *types.Chart

		// Might keep track of related namespace
		NsRef  *Ref
		ModRef RefSet
	}
)

func NewComposeChart(res *types.Chart, nsRef string, mmRef []string) *ComposeChart {
	r := &ComposeChart{
		base:   &base{},
		ModRef: make(RefSet, len(mmRef)),
	}
	r.SetResourceType(COMPOSE_CHART_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Name, res.ID)...)

	r.NsRef = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	for i, mRef := range mmRef {
		r.ModRef[i] = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, mRef).Constraint(r.NsRef)
	}

	return r
}

func (r *ComposeChart) SysID() uint64 {
	return r.Res.ID
}
