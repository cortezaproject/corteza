package yaml

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
)

func apiGatewayFromResource(r *resource.APIGateway, cfg *EncoderConfig) *apiGateway {
	ff := make(apiGwFilterSet, len(r.Filters))
	for i, t := range r.Filters {
		ff[i] = &apiGwFilter{
			res:           t.Res,
			encoderConfig: cfg,
		}
	}

	return &apiGateway{
		res:     r.Res,
		filters: ff,

		encoderConfig: cfg,
	}
}

func (n *apiGateway) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	wf, ok := state.Res.(*resource.APIGateway)
	if !ok {
		return encoderErrInvalidResource(systemTypes.ApigwRouteResourceType, state.Res.ResourceType())
	}

	n.res = wf.Res
	n.us = wf.Userstamps()

	return nil
}

func (n *apiGateway) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
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

	doc.addApiGateway(n)

	return err
}

func (g *apiGateway) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"endpoint", g.res.Endpoint,
		"method", g.res.Method,
		"enabled", g.res.Enabled,
		"group", g.res.Group,
		"meta", g.res.Meta,

		"filters", g.filters,
	)
	if err != nil {
		return nil, err
	}

	nn, err = encodeTimestamps(nn, g.ts)
	if err != nil {
		return nil, err
	}

	nn, err = encodeUserstamps(nn, g.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func (f *apiGwFilter) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"weight", f.res.Weight,
		"ref", f.res.Ref,
		"kind", f.res.Kind,
		"enabled", f.res.Enabled,
		"params", f.res.Params,
	)
	if err != nil {
		return nil, err
	}

	nn, err = encodeTimestamps(nn, f.ts)
	if err != nil {
		return nil, err
	}

	nn, err = encodeUserstamps(nn, f.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}
