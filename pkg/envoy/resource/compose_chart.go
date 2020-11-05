package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	// ComposeChart represents a ComposeChart
	ComposeChart struct {
		*base
		Res *types.Chart

		// Might keep track of related namespace
	}
)

func NewComposeChart(res *types.Chart, nsRef string) *ComposeChart {
	r := &ComposeChart{base: &base{}}
	r.SetResourceType(COMPOSE_CHART_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Name, res.ID)...)

	r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)

	return r
}

func (m *ComposeChart) SearchQuery() types.ChartFilter {
	f := types.ChartFilter{
		Handle: m.Res.Handle,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("chartID=%d", m.Res.ID)
	}

	return f
}
