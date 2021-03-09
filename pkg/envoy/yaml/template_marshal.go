package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
)

func templateFromResource(r *resource.Template, cfg *EncoderConfig) *template {
	return &template{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

// Prepare prepares the template to be encoded
//
// Any validation, additional constraining should be performed here.
func (u *template) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	us, ok := state.Res.(*resource.Template)
	if !ok {
		return encoderErrInvalidResource(resource.TEMPLATE_RESOURCE_TYPE, state.Res.ResourceType())
	}

	u.res = us.Res

	return nil
}

// Encode encodes the template to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (u *template) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if u.res.ID <= 0 {
		u.res.ID = util.NextID()
	}

	// Encode timestamps
	u.ts, err = resource.MakeCUDASTimestamps(&u.res.CreatedAt, u.res.UpdatedAt, u.res.DeletedAt, nil, nil).
		Model(u.encoderConfig.TimeLayout, u.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo implement resource skipping?

	doc.AddTemplate(u)
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

	nsn, err = mapTimestamps(nsn, c.ts)
	if err != nil {
		return nil, err
	}

	return nsn, nil
}
