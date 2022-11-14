package yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

func (wset *composeNamespaceSet) UnmarshalYAML(n *yaml.Node) error {
	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composeNamespace{}
		)

		if v == nil {
			return y7s.NodeErr(n, "malformed namespace definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if k != nil {
			if wrap.res.Slug != "" {
				return y7s.NodeErr(k, "cannot define slug in mapped namespace definition")
			}

			if !handle.IsValid(k.Value) {
				return y7s.NodeErr(n, "namespace reference must be a valid handle")
			}

			// set namespace slug from map key value
			wrap.res.Slug = k.Value
		}

		*wset = append(*wset, wrap)
		return
	})
}

func (wrap *composeNamespace) UnmarshalYAML(n *yaml.Node) (err error) {
	if !y7s.IsKind(n, yaml.MappingNode) {
		return y7s.NodeErr(n, "namespace definition must be a map or scalar")
	}

	if wrap.res == nil {
		wrap.res = &types.Namespace{
			// namespaces are enabled by default
			Enabled: true,
		}
	}

	if err = n.Decode(&wrap.res); err != nil {
		return
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

	return y7s.Each(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "modules":
			return v.Decode(&wrap.modules)

		case "records":
			return v.Decode(&wrap.records)

		case "charts":
			return v.Decode(&wrap.charts)

		case "pages":
			return v.Decode(&wrap.pages)

		}

		return nil
	})
}

func (wset composeNamespaceSet) MarshalEnvoy() ([]resource.Interface, error) {
	if wset == nil {
		return []resource.Interface{}, nil
	}

	// namespace usually have bunch of sub-resources defined
	nn := make([]resource.Interface, 0, len(wset)*10)

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (wrap composeNamespace) MarshalEnvoy() ([]resource.Interface, error) {
	nsr := resource.NewComposeNamespace(wrap.res)
	nsr.SetTimestamps(wrap.ts)
	nsr.SetConfig(wrap.envoyConfig)

	var defaultNsTranslations []resource.Interface
	dft, err := nsr.EncodeTranslations()
	if err != nil {
		return nil, err
	}
	for _, t := range dft {
		t.MarkDefault()
		defaultNsTranslations = append(defaultNsTranslations, t)
	}

	return envoy.CollectNodes(
		nsr,
		defaultNsTranslations,
		wrap.modules,
		wrap.pages,
		wrap.records,
		wrap.charts,
		wrap.rbac.bindResource(nsr),
		wrap.locale.bindResource(nsr),
	)
}
