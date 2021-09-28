package yaml

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v3"
)

func (wset *reportSet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &report{}
		)

		if v == nil {
			return y7s.NodeErr(n, "malformed report definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if err = decodeRef(k, "report handle", &wrap.res.Handle); err != nil {
			return y7s.NodeErr(n, "Report reference must be a valid handle")
		}

		if wrap.res.Meta == nil {
			wrap.res.Meta = &types.ReportMeta{}
		}
		if wrap.res.Meta.Name == "" {
			// if name is not set, use handle
			wrap.res.Meta.Name = wrap.res.Handle
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wrap *report) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.rbac = make(rbacRuleSet, 0, 10)
		wrap.res = &types.Report{}
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	if wrap.locale, err = decodeLocale(n); err != nil {
		return
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}

	if wrap.us, err = decodeUserstamps(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "handle":
			return y7s.DecodeScalar(v, "report handle", &wrap.res.Handle)

		case "meta":
			return v.Decode(&wrap.res.Meta)

		case "sources":
			wrap.sources = make(reportSourceSet, 0, 10)

			err = v.Decode(&wrap.sources)
			if err != nil {
				return err
			}

		case "blocks":
			wrap.blocks = make(reportBlockSet, 0, 10)

			err = v.Decode(&wrap.blocks)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (wrap *reportSource) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.ReportDataSource{}
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "meta":
			return v.Decode(&wrap.res.Meta)
		case "step":
			return v.Decode(&wrap.res.Step)
		}

		return nil
	})
}

func (wrap *reportBlock) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.res = &types.ReportBlock{}
	}

	if wrap.envoyConfig, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.locale, err = decodeLocale(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "title":
			return y7s.DecodeScalar(v, "title", &wrap.res.Title)
		case "description":
			return y7s.DecodeScalar(v, "description", &wrap.res.Description)
		case "key":
			return y7s.DecodeScalar(v, "key", &wrap.res.Key)
		case "kind":
			return y7s.DecodeScalar(v, "kind", &wrap.res.Kind)
		case "options":
			return v.Decode(&wrap.res.Options)
		case "elements":
			return v.Decode(&wrap.res.Elements)
		case "sources":
			return v.Decode(&wrap.res.Sources)
		case "xywh":
			return v.Decode(&wrap.res.XYWH)
		case "layout":
			return y7s.DecodeScalar(v, "layout", &wrap.res.Layout)
		}

		return nil
	})
}

func (wset reportSet) MarshalEnvoy() ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, len(wset)*2)

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (wrap report) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewReport(wrap.res)
	rs.SetTimestamps(wrap.ts)
	rs.SetUserstamps(wrap.us)
	rs.SetConfig(wrap.envoyConfig)

	// default report translations
	// includes translations for nested resources also
	// var defaultReportTranslations []resource.Interface
	// dft, err := rs.EncodeTranslations()
	// if err != nil {
	// 	return nil, err
	// }
	// for _, d := range dft {
	// 	d.MarkDefault()
	// 	defaultReportTranslations = append(defaultReportTranslations, d)
	// }

	for _, s := range wrap.sources {
		rs.AddReportSource(s.res)
	}

	for _, p := range wrap.blocks {
		rs.AddReportBlock(p.res)
	}

	return envoy.CollectNodes(
		rs,
		// defaultReportTranslations,
		wrap.rbac.bindResource(rs),
	)
}
