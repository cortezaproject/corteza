package yaml

import (
	"context"
	"encoding/json"

	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

func settingFromResource(r *resource.Setting, cfg *EncoderConfig) *setting {
	return &setting{
		res:           &r.Res,
		encoderConfig: cfg,
	}
}

func (n *setting) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	st, ok := state.Res.(*resource.Setting)
	if !ok {
		return encoderErrInvalidResource(resource.SettingsResourceType, state.Res.ResourceType())
	}

	n.res = &st.Res
	n.us = st.Userstamps()

	return nil
}

func (n *setting) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	n.ts, err = resource.MakeTimestampsCUDA(nil, &n.res.UpdatedAt, nil, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}

	n.us, err = resolveUserstamps(state.ParentResources, n.us)
	if err != nil {
		return err
	}

	// @todo skip eval?

	doc.addSetting(n)
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

	nn, err = encodeTimestamps(nn, c.ts)
	if err != nil {
		return nil, err
	}

	nn, err = encodeUserstamps(nn, c.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}
