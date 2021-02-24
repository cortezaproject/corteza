package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
)

func roleFromResource(r *resource.Role, cfg *EncoderConfig) *role {
	return &role{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

// Prepare prepares the role to be encoded
//
// Any validation, additional constraining should be performed here.
func (r *role) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	rl, ok := state.Res.(*resource.Role)
	if !ok {
		return encoderErrInvalidResource(resource.ROLE_RESOURCE_TYPE, state.Res.ResourceType())
	}

	r.res = rl.Res

	return nil
}

// Encode encodes the role to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (r *role) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if r.res.ID <= 0 {
		r.res.ID = util.NextID()
	}

	r.ts, err = resource.MakeCUDATimestamps(&r.res.CreatedAt, r.res.UpdatedAt, r.res.DeletedAt, r.res.ArchivedAt).
		Model(r.encoderConfig.TimeLayout, r.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo 1skip eval?

	doc.AddRole(r)
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

	nsn, err = mapTimestamps(nsn, c.ts)
	if err != nil {
		return nil, err
	}

	return nsn, nil
}
