package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
)

const (
	COMPOSE_CHART_RESOURCE_TYPE = "composeChart"
)

type (
	// ComposeChart represents a ComposeChart
	composeChart struct {
		*base
		Res *types.Chart

		// Might keep track of related namespace
	}
)

func ComposeChart(res *types.Chart) *composeChart {
	r := &composeChart{base: &base{}}
	r.SetResourceType(COMPOSE_CHART_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Name, res.ID)...)

	return r
}

func (m *composeChart) SearchQuery() types.ChartFilter {
	f := types.ChartFilter{
		Handle: m.Res.Handle,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("chartID=%d", m.Res.ID)
	}

	return f
}
