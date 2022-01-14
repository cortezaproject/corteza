package yaml

import (
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	composeChart struct {
		res         *types.Chart
		chartConfig *composeChartConfig
		ts          *resource.Timestamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		relNs        *types.Namespace
		refNamespace string

		// pointer to report and module reference
		refReportModules map[int]string

		rbac rbacRuleSet
	}
	composeChartSet []*composeChart

	composeChartConfig struct {
		config types.ChartConfig

		reports          []*composeChartConfigReport
		refReportModules map[int]string
	}

	composeChartConfigReport struct {
		report *types.ChartConfigReport

		refModule string
		relModule *types.Module
	}
)

func (nn composeChartSet) configureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}

func relChartToRef(chr *types.Chart) string {
	return firstOkString(chr.Handle, strconv.FormatUint(chr.ID, 10))
}
