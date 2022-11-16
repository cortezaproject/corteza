package yaml

import (
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	report struct {
		res     *types.Report
		sources reportSourceSet
		blocks  reportBlockSet

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

	reportBlock struct {
		res *types.ReportBlock

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		locale resourceTranslationSet
	}
	reportBlockSet []*reportBlock
)

func (nn reportSet) configureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}
