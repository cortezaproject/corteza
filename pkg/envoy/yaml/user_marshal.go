package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
)

func userFromResource(r *resource.User, cfg *EncoderConfig) *user {
	return &user{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

// Prepare prepares the user to be encoded
//
// Any validation, additional constraining should be performed here.
func (u *user) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	us, ok := state.Res.(*resource.User)
	if !ok {
		return encoderErrInvalidResource(resource.USER_RESOURCE_TYPE, state.Res.ResourceType())
	}

	u.res = us.Res

	return nil
}

// Encode encodes the user to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (u *user) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if u.res.ID <= 0 {
		u.res.ID = util.NextID()
	}

	// Encode timestamps
	u.ts, err = resource.MakeCUDASTimestamps(&u.res.CreatedAt, u.res.UpdatedAt, u.res.DeletedAt, nil, u.res.SuspendedAt).
		Model(u.encoderConfig.TimeLayout, u.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo implement resource skipping?

	doc.AddUser(u)
	return
}

func (c *user) MarshalYAML() (interface{}, error) {
	var err error

	nsn, err := makeMap(
		"username", c.res.Username,
		"email", c.res.Email,
		"name", c.res.Name,
		"handle", c.res.Handle,
		"kind", c.res.Kind,

		"meta", c.res.Meta,

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
