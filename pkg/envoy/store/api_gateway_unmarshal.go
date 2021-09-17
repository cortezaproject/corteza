package store

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

func newAPIGateway(gwr *types.ApigwRoute, ff types.ApigwFilterSet, ux *userIndex) *apiGateway {
	return &apiGateway{
		gwr: gwr,
		ff:  ff,

		ux: ux,
	}
}

func (awf *apiGateway) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewAPIGateway(awf.gwr)
	syncUserStamps(rs.Userstamps(), awf.ux)

	for _, f := range awf.ff {
		rt := rs.AddGatewayFilter(f)
		syncUserStamps(rt.Userstamps(), awf.ux)
	}

	return envoy.CollectNodes(
		rs,
	)
}
