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
		Exec(in *expr.Vars) (out *expr.Vars, err error)
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
func (m *AutomationFunctionServer) Meta(ctx context.Context, req *proto.MetaReq) (*proto.AutomationFunction, error) {
	f := m.Impl.Meta()

	return &proto.AutomationFunction{
		Ref: f.Ref,
		Meta: &proto.AutomationFunctionMeta{
			Short:       f.Meta.Short,
			Description: f.Meta.Description,
		},
		Params: []*proto.AutomationFunctionParams{
			{
				Name:     f.Parameters[0].Name,
				Types:    f.Parameters[0].Types,
				Required: f.Parameters[0].Required,
			},
		},
		Results: []*proto.AutomationFunctionParams{
			{
				Name:     f.Results[0].Name,
				Types:    f.Results[0].Types,
				Required: f.Results[0].Required,
			},
		},
	}, nil
}

func (m *AutomationFunctionServer) Exec(
	ctx context.Context,
	req *proto.ExecReq) (*proto.ExecResp, error) {

	oo := &expr.Vars{}

	if err := json.Unmarshal([]byte(req.Value), oo); err != nil {
		return nil, err
	}

	out, err := m.Impl.Exec(oo)

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
func (m *AutomationFunctionClient) Meta() *types.Function {
	af, err := m.client.Meta(context.Background(), &proto.MetaReq{})

	if err != nil {
		return nil
	}

	return &types.Function{
		Ref:  af.Ref,
		Kind: "function",
		Meta: &types.FunctionMeta{
			Short:       af.Meta.Short,
			Description: af.Meta.Description,
		},
		Parameters: types.ParamSet{
			{
				Name:     af.Params[0].Name,
				Types:    af.Params[0].Types,
				Required: af.Params[0].Required,
			},
		},
		Results: types.ParamSet{
			{
				Name:     af.Results[0].Name,
				Types:    af.Results[0].Types,
				Required: af.Results[0].Required,
			},
		},
	}
}

func (m *AutomationFunctionClient) Exec(in *expr.Vars) (out *expr.Vars, err error) {

	// we need to encode the map[string]any and decode in the server
	b, err := json.Marshal(in)

	if err != nil {
		return
	}

	er, err := m.client.Exec(context.Background(), &proto.ExecReq{
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
