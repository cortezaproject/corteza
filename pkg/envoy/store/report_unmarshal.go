package store

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

func newReport(wf *types.Report, ux *userIndex) *report {
	return &report{
		rp: wf,
		ss: wf.Sources,
		pp: wf.Projections,

		ux: ux,
	}
}

func (awf *report) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewReport(awf.rp)
	syncUserStamps(rs.Userstamps(), awf.ux)

	for _, s := range awf.ss {
		rs.AddReportSource(s)
	}

	for _, p := range awf.pp {
		rs.AddReportProjection(p)
	}

	return envoy.CollectNodes(
		rs,
	)
}
