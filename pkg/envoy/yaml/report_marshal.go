package yaml

import (
	"context"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func reportFromResource(r *resource.Report, cfg *EncoderConfig) *report {
	ss := make(reportSourceSet, len(r.Sources))
	for i, s := range r.Sources {
		ss[i] = &reportSource{
			res:           s.Res,
			encoderConfig: cfg,
		}
	}

	pp := make(reportBlockSet, len(r.Blocks))
	for i, p := range r.Blocks {
		pp[i] = &reportBlock{
			res:           p.Res,
			encoderConfig: cfg,
		}
	}

	return &report{
		res:     r.Res,
		sources: ss,
		blocks:  pp,

		encoderConfig: cfg,
	}
}

func (n *report) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	wf, ok := state.Res.(*resource.Report)
	if !ok {
		return encoderErrInvalidResource(automationTypes.WorkflowResourceType, state.Res.ResourceType())
	}

	n.res = wf.Res
	n.us = wf.Userstamps()

	return nil
}

func (n *report) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
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

	// @todo skip eval?

	doc.addReport(n)

	return err
}

func (wf *report) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"handle", wf.res.Handle,
		"meta", wf.res.Meta,

		"sources", wf.sources,
		"blocks", wf.blocks,
	)
	if err != nil {
		return nil, err
	}

	nn, err = encodeTimestamps(nn, wf.ts)
	if err != nil {
		return nil, err
	}

	nn, err = encodeUserstamps(nn, wf.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func (t *reportSource) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"meta", t.res.Meta,
		"step", t.res.Step,
	)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func (t *reportBlock) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"title", t.res.Title,
		"description", t.res.Description,
		"key", t.res.Key,
		"kind", t.res.Kind,
		"options", t.res.Options,
		"elements", t.res.Elements,
		"sources", t.res.Sources,
		"xywh", t.res.XYWH,
		"layout", t.res.Layout,
	)
	if err != nil {
		return nil, err
	}

	return nn, nil
}
