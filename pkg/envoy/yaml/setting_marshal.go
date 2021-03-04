package yaml

import (
	"context"
	"encoding/json"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func settingFromResource(r *resource.Setting, cfg *EncoderConfig) *setting {
	return &setting{
		res:           &r.Res,
		encoderConfig: cfg,
	}
}

// Prepare prepares the setting to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *setting) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	st, ok := state.Res.(*resource.Setting)
	if !ok {
		return encoderErrInvalidResource(resource.SETTINGS_RESOURCE_TYPE, state.Res.ResourceType())
	}

	n.res = &st.Res
	n.us = st.Userstamps()

	return nil
}

// Encode encodes the setting to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *setting) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	// Timestamps
	n.ts, err = resource.MakeCUDATimestamps(nil, &n.res.UpdatedAt, nil, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	n.us, err = resolveUserstamps(state.ParentResources, n.us)
	if err != nil {
		return err
	}

	// @todo skip eval?

	doc.AddSetting(n)
	return nil
}

func (c *setting) MarshalYAML() (interface{}, error) {
	var aux interface{}
	err := json.Unmarshal(c.res.Value, &aux)
	if err != nil {
		return nil, err
	}
	nn, err := makeMap(
		"name", c.res.Name,
		"value", aux,
	)
	if err != nil {
		return nil, err
	}

	nn, err = mapTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	nn, err = mapUserstamps(nn, c.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}
