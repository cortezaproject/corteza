package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
)

// Prepare prepares the composeNamespace to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *composeNamespace) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	ns, ok := state.Res.(*resource.ComposeNamespace)
	if !ok {
		return encoderErrInvalidResource(resource.COMPOSE_NAMESPACE_RESOURCE_TYPE, state.Res.ResourceType())
	}

	n.meta = composeNamespaceMeta(ns.Res.Meta)

	return nil
}

// Encode encodes the composeNamespace to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *composeNamespace) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if n.res.ID <= 0 {
		n.res.ID = util.NextID()
	}

	// Timestamps
	n.ts, err = resource.MakeCUDATimestamps(&n.res.CreatedAt, n.res.UpdatedAt, n.res.DeletedAt, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo skip eval?

	doc.AddComposeNamespace(n)
	return nil
}

func (c *composeNamespace) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"name", c.res.Name,
		"slug", c.res.Slug,
		"enabled", c.res.Enabled,
		"labels", c.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	if !c.meta.Empty() {
		nn, err = addMap(nn,
			"meta", c.res.Meta,
		)
	}

	nn, err = mapTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	if len(c.modules) > 0 {
		c.modules.ConfigureEncoder(c.encoderConfig)

		nn, err = encodeResource(nn, "modules", c.modules, c.encoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if len(c.records) > 0 {
		c.records.ConfigureEncoder(c.encoderConfig)

		// Records don't have this
		nn, err = encodeResource(nn, "records", c.records, false, "")
		if err != nil {
			return nil, err
		}
	}

	if len(c.pages) > 0 {
		c.pages.ConfigureEncoder(c.encoderConfig)

		nn, err = encodeResource(nn, "pages", c.pages, c.encoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if len(c.charts) > 0 {
		c.charts.ConfigureEncoder(c.encoderConfig)

		nn, err = encodeResource(nn, "charts", c.charts, c.encoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	return nn, nil
}

func (c composeNamespaceMeta) MarshalYAML() (interface{}, error) {
	return c, nil
}
