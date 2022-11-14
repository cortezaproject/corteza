package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

func roleFromResource(r *resource.Role, cfg *EncoderConfig) *role {
	return &role{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

func (r *role) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	rl, ok := state.Res.(*resource.Role)
	if !ok {
		return encoderErrInvalidResource(types.RoleResourceType, state.Res.ResourceType())
	}

	r.res = rl.Res

	return nil
}

func (r *role) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if r.res.ID <= 0 {
		r.res.ID = nextID()
	}

	r.ts, err = resource.MakeTimestampsCUDA(&r.res.CreatedAt, r.res.UpdatedAt, r.res.DeletedAt, r.res.ArchivedAt).
		Model(r.encoderConfig.TimeLayout, r.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo skip eval?

	doc.addRole(r)
	return err
}

func (c *role) MarshalYAML() (interface{}, error) {
	var err error

	nsn, err := makeMap(
		"name", c.res.Name,
		"handle", c.res.Handle,
		"labels", c.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nsn, err = encodeTimestamps(nsn, c.ts)
	if err != nil {
		return nil, err
	}

	return nsn, nil
}
