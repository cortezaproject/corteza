package yaml

import (
	"context"
	"encoding/json"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
)

func messagingChannelFromResource(r *resource.MessagingChannel, cfg *EncoderConfig) *messagingChannel {
	return &messagingChannel{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

// Prepare prepares the messagingChannel to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *messagingChannel) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	ch, ok := state.Res.(*resource.MessagingChannel)
	if !ok {
		return encoderErrInvalidResource(resource.MESSAGING_CHANNEL_RESOURCE_TYPE, state.Res.ResourceType())
	}

	n.res = ch.Res
	n.us = ch.Userstamps()

	return nil
}

// Encode encodes the messagingChannel to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *messagingChannel) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if n.res.ID <= 0 {
		n.res.ID = util.NextID()
	}

	n.ts, err = resource.MakeCUDATimestamps(&n.res.CreatedAt, n.res.UpdatedAt, n.res.DeletedAt, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	n.us, err = resolveUserstamps(state.ParentResources, n.us)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	// @todo 1skip eval?

	doc.AddMessagingChannel(n)

	return err
}

func (c *messagingChannel) MarshalYAML() (interface{}, error) {
	// Get a struct from the raw JSON string
	auxMeta := make(map[string]interface{})
	json.Unmarshal(c.res.Meta, auxMeta)

	nsn, err := makeMap(
		"name", c.res.Name,
		"topic", c.res.Topic,
		"type", c.res.Type,
		"meta", auxMeta,
		"membershipPolicy", c.res.MembershipPolicy,
	)
	if err != nil {
		return nil, err
	}

	nsn, err = mapTimestamps(nsn, c.ts)
	if err != nil {
		return nil, err
	}

	nsn, err = addMap(nsn, "creator", c.us.CreatedBy)
	if err != nil {
		return nil, err
	}

	return nsn, nil
}
