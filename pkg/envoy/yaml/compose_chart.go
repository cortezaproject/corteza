package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"gopkg.in/yaml.v3"
)

type (
	composeChart struct {
		res          *types.Chart
		refNamespace string

		// pointer to report and module reference
		refReportModules map[int]string

		*rbacRules
	}
	composeChartSet []*composeChart

	composeChartConfig struct {
		config           types.ChartConfig
		refReportModules map[int]string
	}

	composeChartConfigReport struct {
		report    *types.ChartConfigReport
		refModule string
	}
)

func (wset *composeChartSet) UnmarshalYAML(n *yaml.Node) error {
	return iterator(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composeChart{}
		)

		if v == nil {
			return nodeErr(n, "malformed chart definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if k != nil {
			if wrap.res.Handle != "" {
				return nodeErr(k, "cannot define handle in mapped chart definition")
			}

			if !handle.IsValid(k.Value) {
				return nodeErr(n, "Chart reference must be a valid handle")
			}

			wrap.res.Handle = k.Value
			wrap.res.Name = k.Value
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wset composeChartSet) setNamespaceRef(ref string) error {
	for _, res := range wset {
		if res.refNamespace != "" && ref != res.refNamespace {
			return fmt.Errorf("cannot override namespace reference %s with %s", res.refNamespace, ref)
		}

		res.refNamespace = ref
	}

	return nil
}

func (wset composeChartSet) MarshalEnvoy() ([]envoy.Node, error) {
	// namespace usually have bunch of sub-resources defined
	nn := make([]envoy.Node, 0, len(wset)*10)

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (wrap *composeChart) UnmarshalYAML(n *yaml.Node) (err error) {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "chart definition must be a map")
	}

	if wrap.res == nil {
		wrap.rbacRules = &rbacRules{}
		wrap.res = &types.Chart{}
	}

	if wrap.rbacRules, err = decodeResourceAccessControl(types.ChartRBACResource, n); err != nil {
		return
	}

	return iterator(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "handle":
			return decodeScalar(v, "chart handle", &wrap.res.Handle)

		case "name":
			return decodeScalar(v, "chart name", &wrap.res.Name)

		case "config":
			cfg := composeChartConfig{
				refReportModules: make(map[int]string),
			}

			if err = v.Decode(&cfg); err != nil {
				return err
			}

			// copy decoded values from aux type
			wrap.res.Config = cfg.config
			wrap.refReportModules = cfg.refReportModules

		}

		return nil
	})
}

func (wrap composeChart) MarshalEnvoy() ([]envoy.Node, error) {
	nn := make([]envoy.Node, 0, 16)
	//nn = append(nn, &envoy.ComposeChartNode{Chart: wrap.res})

	return nn, nil
}

func (wrap *composeChartConfig) UnmarshalYAML(n *yaml.Node) error {
	if !isKind(n, yaml.MappingNode) {
		return nodeErr(n, "chart configuration must be a map")
	}

	return iterator(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "reports":
			reports := make([]composeChartConfigReport, 0)
			if err = v.Decode(&reports); err != nil {
				return nodeErr(v, "could not decode reports: %w", err)
			}

			// collect reports and referenced modules from wrapped type
			wrap.config.Reports = make([]*types.ChartConfigReport, len(reports))
			for r := range reports {
				wrap.config.Reports[r] = reports[r].report

				if reports[r].refModule != "" {
					wrap.refReportModules[r] = reports[r].refModule
				}
			}

		case "colorScheme":
			if err = v.Decode(&wrap.config.ColorScheme); err != nil {
				return nodeErr(v, "could not decode color scheme: %w", err)
			}

		}

		return nil
	})
}

func (wrap *composeChartConfigReport) UnmarshalYAML(n *yaml.Node) error {
	if err := n.Decode(&wrap.report); err != nil {
		return err
	}

	return iterator(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "module":
			// custom decoder for referenced module
			// we'll copy this to the dedicated prop of the wrapping structure
			// so that the parent decoder can collect it
			return decodeRef(v, "chart report module", &wrap.refModule)

		}

		return nil
	})
}
