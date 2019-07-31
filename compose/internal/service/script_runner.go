package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/cortezaproject/corteza-server/compose/proto"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

// Script runner provides an interface to corteza-corredor (Spanish for runner) service
// that helps us with execution of JavaScript code -- compose's triggers & automation code
//
// corteza-server communicates with corteza-corredor via gRPC protocol.
//
// This service accepts ns/trigger/module/record (combinations), makes a call via gRPC protocol and
// returns record/module/ns or just tests trigger's script

type (
	scriptRunner struct {
		c          options.ScriptRunnerOpt
		logger     *zap.Logger
		conn       *grpc.ClientConn
		client     proto.ScriptRunnerClient
		jwtEncoder auth.TokenEncoder
	}

	Runnable interface {
		proto.Runnable

		IsCritical() bool
		GetRunnerID() uint64
	}

	ScriptRunnerService interface {
		Close() error
		Namespace(context.Context, Runnable, *types.Namespace) (*types.Namespace, error)
		Module(context.Context, Runnable, *types.Namespace, *types.Module) (*types.Module, error)
		Record(context.Context, Runnable, *types.Namespace, *types.Module, *types.Record) (*types.Record, error)
	}
)

// @todo move to opt so all services can use it
func ScriptRunner(c options.ScriptRunnerOpt) (svc *scriptRunner, err error) {
	svc = &scriptRunner{
		c:          c,
		logger:     DefaultLogger.Named("script-runner"),
		jwtEncoder: auth.DefaultJwtHandler,
	}

	if !c.Enabled {
		// Do not connect when script runner is not enabled
		return
	}

	return svc, svc.connect()
}

func (svc *scriptRunner) connect() (err error) {
	if svc.c.Log {
		// Send logs to zap
		//
		// waiting for https://github.com/uber-go/zap/pull/538
		grpclog.SetLogger(zapgrpc.NewLogger(svc.logger.Named("grpc")))
	}

	var dopts = []grpc.DialOption{
		// @todo insecure?
		grpc.WithInsecure(),
	}

	if svc.c.MaxBackoffDelay > 0 {
		dopts = append(dopts, grpc.WithBackoffMaxDelay(svc.c.MaxBackoffDelay))
	}

	svc.conn, err = grpc.Dial(svc.c.Addr, dopts...)

	if err != nil {
		return
	}

	svc.client = proto.NewScriptRunnerClient(svc.conn)
	return
}

func (svc scriptRunner) Close() error {
	return svc.conn.Close()
}

func (svc scriptRunner) callOptions() []grpc.CallOption {
	return []grpc.CallOption{
		grpc.WaitForReady(true),
	}
}

// Creates a new JWT for
func (svc scriptRunner) getJWT(ctx context.Context, r Runnable) string {
	if r.GetRunnerID() > 0 {
		// @todo implement this
		//       at the moment we do not he the ability fetch user info from non-system service
		//       extend/implement this feature when our services will know how to communicate with each-other
	}

	return svc.jwtEncoder.Encode(auth.GetIdentityFromContext(ctx))
}

func (svc scriptRunner) Namespace(ctx context.Context, s Runnable, ns *types.Namespace) (*types.Namespace, error) {
	panic("scriptRunner.Namespace() not implemented")
}

func (svc scriptRunner) Module(ctx context.Context, s Runnable, ns *types.Namespace, m *types.Module) (*types.Module, error) {
	panic("scriptRunner.Module() not implemented")
}

func (svc scriptRunner) Record(ctx context.Context, s Runnable, ns *types.Namespace, m *types.Module, r *types.Record) (*types.Record, error) {
	if s == nil {
		return nil, errors.New("script not provided")
	}

	if ns == nil {
		return nil, errors.New("namespace not provided")
	}

	if m == nil {
		return nil, errors.New("module not provided")
	}

	if !svc.c.Enabled {
		if s.IsCritical() {
			// Oh dear, we are in quite a pickle:
			// Script runner is disabled but we have critical script to run
			return nil, errors.New("script runner disabled")
		}

		// Log this
		svc.logger.Debug("executing script", zap.Any("record", r))

		// and pretend like nothing happened
		return r, nil
	}

	svc.logger.Debug("executing script", zap.Any("record", r))

	ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
	defer cancelFn()

	rsp, err := svc.client.Record(
		ctx,
		&proto.RunRecordRequest{
			JWT:       svc.getJWT(ctx, s),
			Script:    proto.ScriptFromRunnable(s),
			Namespace: proto.FromNamespace(ns),
			Module:    proto.FromModule(m),
			Record:    proto.FromRecord(r),
		},
		svc.callOptions()...,
	)

	svc.logger.Debug("call sent")

	if err != nil {
		svc.logger.Debug("script executed, did not return record", zap.Error(err))
		if !s.IsCritical() {
			// This was not a critical call and we do not care about
			// errors from script running service.
			return r, nil
		}

		return nil, err
	}

	if s.IsAsync() {
		svc.logger.Debug("script executed / async")
		// Async call, we do not care about what we get back
		return r, nil
	}

	svc.logger.Debug("script executed", zap.Any("record", rsp.Record))

	// Result from the automation script
	return proto.ToRecord(rsp.Record), nil
}
