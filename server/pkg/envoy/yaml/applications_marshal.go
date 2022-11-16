package yaml

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/system/types"
)

func applicationFromResource(r *resource.Application, cfg *EncoderConfig) *application {
	return &application{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

func (n *application) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	app, ok := state.Res.(*resource.Application)
	if !ok {
		return encoderErrInvalidResource(types.ApplicationResourceType, state.Res.ResourceType())
	}

	n.res = app.Res
	n.us = app.Userstamps()

	return nil
}

func (n *application) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if n.res.ID <= 0 {
		n.res.ID = nextID()
	}

	n.ts, err = resource.MakeTimestampsCUDA(&n.res.CreatedAt, n.res.UpdatedAt, n.res.DeletedAt, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}
	n.us, err = resolveUserstamps(state.ParentResources, n.us)
	if err != nil {
		return err
	}

	// @todo 1skip eval?

	doc.addApplication(n)

	return err
}

func (c *application) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"name", c.res.Name,
		"enabled", c.res.Enabled,
		"weight", c.res.Weight,

		"unify", c.res.Unify,

		"labels", c.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nn, err = encodeTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	if c.us != nil {
		nn, err = addMap(nn, "owner", c.us.OwnedBy)
		if err != nil {
			return nil, err
		}
	}

	return nn, nil
}
