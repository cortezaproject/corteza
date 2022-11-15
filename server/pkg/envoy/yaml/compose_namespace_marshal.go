package yaml

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

func (n *composeNamespace) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	ns, ok := state.Res.(*resource.ComposeNamespace)
	if !ok {
		return encoderErrInvalidResource(types.NamespaceResourceType, state.Res.ResourceType())
	}

	n.meta = composeNamespaceMeta(ns.Res.Meta)

	return nil
}

func (n *composeNamespace) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if n.res.ID <= 0 {
		n.res.ID = nextID()
	}

	n.ts, err = resource.MakeTimestampsCUDA(&n.res.CreatedAt, n.res.UpdatedAt, n.res.DeletedAt, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo skip eval?

	doc.addComposeNamespace(n)
	return nil
}

func (c *composeNamespace) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"namespaceID", c.res.ID,
		"name", c.res.Name,
		"slug", c.res.Slug,
		"enabled", c.res.Enabled,
		"labels", c.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	if !c.meta.empty() {
		nn, err = addMap(nn,
			"meta", c.res.Meta,
		)
	}

	nn, err = encodeTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	if len(c.modules) > 0 {
		c.modules.configureEncoder(c.encoderConfig)

		nn, err = encodeResource(nn, "modules", c.modules, c.encoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if len(c.records) > 0 {
		c.records.configureEncoder(c.encoderConfig)

		// Records don't have this
		nn, err = encodeResource(nn, "records", c.records, false, "")
		if err != nil {
			return nil, err
		}
	}

	if len(c.pages) > 0 {
		c.pages.configureEncoder(c.encoderConfig)

		nn, err = encodeResource(nn, "pages", c.pages, c.encoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if len(c.charts) > 0 {
		c.charts.configureEncoder(c.encoderConfig)

		nn, err = encodeResource(nn, "charts", c.charts, c.encoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	return nn, nil
}
