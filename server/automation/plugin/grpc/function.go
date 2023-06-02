package grpc

import (
	"context"
	"encoding/json"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/plugin/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type (
	AutomationFunction interface {
		Meta() *types.Function
		Exec(context.Context, *expr.Vars) (*expr.Vars, error)
	}

	AutomationFunctionClient struct {
		client proto.AutomationFunctionServiceClient
	}

	AutomationFunctionPlugin struct {
		plugin.Plugin
		Impl AutomationFunction
	}

	AutomationFunctionServer struct {
		Impl AutomationFunction
		proto.UnimplementedAutomationFunctionServiceServer
	}
)

// GRPC plugin - centered implementations
func (p *AutomationFunctionPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterAutomationFunctionServiceServer(s, &AutomationFunctionServer{Impl: p.Impl})
	return nil
}

func (p *AutomationFunctionPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &AutomationFunctionClient{client: proto.NewAutomationFunctionServiceClient(c)}, nil
}

// GRPC server
func (m *AutomationFunctionServer) Meta(ctx context.Context, req *proto.MetaReq) (ff *proto.AutomationFunction, err error) {
	var f = m.Impl.Meta()

	ff = &proto.AutomationFunction{
		Ref: f.Ref,
		Meta: &proto.AutomationFunctionMeta{
			Short:       f.Meta.Short,
			Description: f.Meta.Description,
		},
		Labels: f.Labels,
	}

	for _, v := range f.Parameters {
		ff.Params = append(ff.Params, &proto.AutomationFunctionParams{Name: v.Name, Types: v.Types, Required: v.Required})
	}

	for _, v := range f.Results {
		ff.Results = append(ff.Results, &proto.AutomationFunctionParams{Name: v.Name, Types: v.Types, Required: v.Required})
	}

	return
}

func (m *AutomationFunctionServer) Exec(
	ctx context.Context,
	req *proto.ExecReq) (*proto.ExecResp, error) {

	oo := &expr.Vars{}

	if err := json.Unmarshal([]byte(req.Value), oo); err != nil {
		return nil, err
	}

	out, err := m.Impl.Exec(ctx, oo)

	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(out)

	if err != nil {
		return nil, err
	}

	return &proto.ExecResp{Value: string(b)}, err
}

// GRPC client
func (m *AutomationFunctionClient) Meta() (ff *types.Function) {
	af, err := m.client.Meta(context.Background(), &proto.MetaReq{})

	if err != nil {
		return nil
	}

	ff = &types.Function{
		Ref:  af.Ref,
		Kind: "function",
		Meta: &types.FunctionMeta{
			Short:       af.Meta.Short,
			Description: af.Meta.Description,
		},
		Labels: af.Labels,
	}

	for _, v := range af.Params {
		ff.Parameters = append(ff.Parameters, &types.Param{Name: v.Name, Types: v.Types, Required: v.Required})
	}

	for _, v := range af.Results {
		ff.Results = append(ff.Results, &types.Param{Name: v.Name, Types: v.Types, Required: v.Required})
	}

	return ff
}

func (m *AutomationFunctionClient) Exec(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {

	// we need to encode the map[string]any and decode in the server
	b, err := json.Marshal(in)

	if err != nil {
		return
	}

	er, err := m.client.Exec(ctx, &proto.ExecReq{
		Value: string(b),
	})

	if err != nil {
		return
	}

	oo := &expr.Vars{}

	if err = json.Unmarshal([]byte(er.Value), oo); err != nil {
		return
	}

	return oo, nil
}
