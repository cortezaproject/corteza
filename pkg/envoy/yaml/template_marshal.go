package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func templateFromResource(r *resource.Template, cfg *EncoderConfig) *template {
	return &template{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

func (t *template) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	us, ok := state.Res.(*resource.Template)
	if !ok {
		return encoderErrInvalidResource(resource.TEMPLATE_RESOURCE_TYPE, state.Res.ResourceType())
	}

	t.res = us.Res

	return nil
}

func (t *template) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if t.res.ID <= 0 {
		t.res.ID = nextID()
	}

	t.ts, err = resource.MakeTimestampsCUDAS(&t.res.CreatedAt, t.res.UpdatedAt, t.res.DeletedAt, nil, nil).
		Model(t.encoderConfig.TimeLayout, t.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	t.us, err = resolveUserstamps(state.ParentResources, t.us)
	if err != nil {
		return err
	}

	// @todo implement resource skipping?

	doc.addTemplate(t)
	return
}

func (c *template) MarshalYAML() (interface{}, error) {
	var err error

	meta, err := makeMap(
		"short", c.res.Meta.Short,
		"description", c.res.Meta.Description,
	)
	if err != nil {
		return nil, err
	}

	nsn, err := makeMap(
		"handle", c.res.Handle,
		"language", c.res.Language,
		"type", c.res.Type,
		"partial", c.res.Partial,
		"template", c.res.Template,

		"meta", meta,
		"Labels", c.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nsn, err = encodeTimestamps(nsn, c.ts)
	if err != nil {
		return nil, err
	}
	nsn, err = encodeUserstamps(nsn, c.us)
	if err != nil {
		return nil, err
	}

	return nsn, nil
}
