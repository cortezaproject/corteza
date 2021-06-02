package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

func userFromResource(r *resource.User, cfg *EncoderConfig) *user {
	return &user{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

func (u *user) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	us, ok := state.Res.(*resource.User)
	if !ok {
		return encoderErrInvalidResource(types.UserResourceType, state.Res.ResourceType())
	}

	u.res = us.Res

	return nil
}

func (u *user) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if u.res.ID <= 0 {
		u.res.ID = nextID()
	}

	// Encode timestamps
	u.ts, err = resource.MakeTimestampsCUDAS(&u.res.CreatedAt, u.res.UpdatedAt, u.res.DeletedAt, nil, u.res.SuspendedAt).
		Model(u.encoderConfig.TimeLayout, u.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	// @todo implement resource skipping?

	doc.addUser(u)
	return
}

func (u *user) MarshalYAML() (interface{}, error) {
	var err error

	nsn, err := makeMap(
		"username", u.res.Username,
		"email", u.res.Email,
		"name", u.res.Name,
		"handle", u.res.Handle,
		"kind", u.res.Kind,

		"meta", u.res.Meta,

		"labels", u.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nsn, err = encodeTimestamps(nsn, u.ts)
	if err != nil {
		return nil, err
	}

	return nsn, nil
}
