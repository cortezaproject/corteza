package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func composeChartFromResource(r *resource.ComposeChart, cfg *EncoderConfig) *composeChart {
	return &composeChart{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

func (n *composeChart) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	chr, ok := state.Res.(*resource.ComposeChart)
	if !ok {
		return encoderErrInvalidResource(resource.COMPOSE_CHART_RESOURCE_TYPE, state.Res.ResourceType())
	}

	// Get the related namespace
	n.relNs = resource.FindComposeNamespace(state.ParentResources, chr.RefNs.Identifiers)
	if n.relNs == nil {
		return resource.ComposeNamespaceErrUnresolved(chr.RefNs.Identifiers)
	}

	n.res = chr.Res
	n.refNamespace = relNsToRef(n.relNs)

	n.chartConfig = &composeChartConfig{
		config:  n.res.Config,
		reports: make([]*composeChartConfigReport, 0, 1),
	}

	for i, r := range chr.Res.Config.Reports {
		refMod := chr.RefMods[i]
		relMod := resource.FindComposeModule(state.ParentResources, refMod.Identifiers)
		if relMod == nil {
			return resource.ComposeModuleErrUnresolved(refMod.Identifiers)
		}

		ccr := &composeChartConfigReport{
			report:    r,
			relModule: relMod,
		}

		n.chartConfig.reports = append(n.chartConfig.reports, ccr)
	}

	return nil
}

func (n *composeChart) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	res := n.res
	if res.ID <= 0 {
		res.ID = nextID()
	}

	if state.Conflicting {
		return nil
	}

	// Handle timestamps
	n.ts, err = resource.MakeTimestampsCUDA(&n.res.CreatedAt, n.res.UpdatedAt, n.res.DeletedAt, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo skip eval?

	// if n.encoderConfig.CompactOutput {
	// 	err = doc.nestComposeChart(n.refNamespace, n)
	// } else {
	// 	doc.addComposeChart(n)
	// }
	doc.addComposeChart(n)

	return err
}

func (c *composeChart) MarshalYAML() (interface{}, error) {
	nn, err := makeMap(
		"handle", c.res.Handle,
		"name", c.res.Name,
		"config", c.chartConfig,
		"labels", c.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nn, err = encodeTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func (c *composeChartConfig) MarshalYAML() (interface{}, error) {
	nn, err := makeMap(
		"reports", c.reports,
		"colorScheme", c.config.ColorScheme,
	)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func (c *composeChartConfigReport) MarshalYAML() (interface{}, error) {
	nn, err := makeMap(
		"filter", c.report.Filter,
		"module", firstOkString(c.relModule.Handle, c.relModule.Name),
		"metrics", c.report.Metrics,
		"dimensions", c.report.Dimensions,
		"yAxis", c.report.YAxis,
	)
	if err != nil {
		return nil, err
	}

	return nn, nil
}
