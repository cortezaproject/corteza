package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	composeChart struct {
		res          *types.Chart
		refNamespace string

		// pointer to report and module reference
		refReportModules map[int]string

		rbac *rbacRules
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
	return eachMap(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composeChart{}
		)

		if v == nil {
			return nodeErr(n, "malformed chart definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if err = decodeRef(k, "chart", &wrap.res.Handle); err != nil {
			return
		}

		if wrap.res.Name == "" {
			// if name is not set, use handle
			wrap.res.Name = wrap.res.Handle
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

func (wset composeChartSet) MarshalEnvoy() ([]resource.Interface, error) {
	// namespace usually have bunch of sub-resources defined
	nn := make([]resource.Interface, 0, len(wset))

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
	if wrap.res == nil {
		wrap.rbac = &rbacRules{}
		wrap.res = &types.Chart{}
	}

	if wrap.rbac, err = decodeResourceAccessControl(types.ChartRBACResource, n); err != nil {
		return
	}

	return eachMap(n, func(k, v *yaml.Node) (err error) {
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

func (wrap composeChart) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.ComposeChart(wrap.res),
	)
}

func (wrap *composeChartConfig) UnmarshalYAML(n *yaml.Node) error {
	return eachMap(n, func(k, v *yaml.Node) (err error) {
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

	return eachMap(n, func(k, v *yaml.Node) (err error) {
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
