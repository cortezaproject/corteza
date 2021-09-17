package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	report struct {
		res         *types.Report
		sources     reportSourceSet
		projections reportProjectionSet

		ts *resource.Timestamps
		us *resource.Userstamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		rbac   rbacRuleSet
		locale resourceTranslationSet
	}
	reportSet []*report

	reportSource struct {
		res *types.ReportDataSource

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig
	}
	reportSourceSet []*reportSource

	reportProjection struct {
		res *types.ReportProjection

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		locale resourceTranslationSet
	}
	reportProjectionSet []*reportProjection
)

func (nn reportSet) configureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}
