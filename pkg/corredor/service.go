package corredor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

type (
	service struct {
		// stores corredor connection options
		// for when we're doing lazy setup
		corredorOpt options.CorredorOpt

		// list of all registered triggers
		//
		registered map[string][]uintptr

		ssClient ServerScriptsClient
		log      *zap.Logger

		eventbus eventRegistrator
		jwtMaker AuthTokenMaker
	}

	Event interface {
		// ResourceType from resource that fired the event
		ResourceType() string

		// Event that was fired
		EventType() string

		// Match tests if given constraints match
		// event's internal values
		Match(name string, op string, values ...string) bool

		// Encode (event) to arguments passed to
		// event handlers ([automation ]script runner)
		Encode() (map[string][]byte, error)

		// Decodes received data back to event
		Decode(map[string][]byte) error
	}

	eventRegistrator interface {
		Register(h eventbus.Handler, ops ...eventbus.TriggerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	AuthTokenMaker func(user string) (string, error)
)

var (
	// Global corredor service
	gService *service
)

func Service() *service {
	return gService
}

// Start connects to Corredor & initialize service
func Start(ctx context.Context, logger *zap.Logger, opt options.CorredorOpt) (err error) {
	if gService != nil {
		// Prevent multiple initializations
		return
	}

	var (
		conn *grpc.ClientConn
	)

	if conn, err = NewConnection(ctx, opt, logger); err != nil {
		return
	}

	gService = NewService(conn, logger)
	return
}

func NewService(conn *grpc.ClientConn, logger *zap.Logger) *service {
	return &service{
		ssClient:   NewServerScriptsClient(conn),
		log:        logger.Named("corredor"),
		registered: make(map[string][]uintptr),
		eventbus:   eventbus.Default(),
	}
}

func (svc *service) SetJwtMaker(fn AuthTokenMaker) {
	svc.jwtMaker = fn
}

func (svc *service) Load(ctx context.Context) (err error) {
	var (
		rsp *ServerScriptListResponse
		ss  ScriptSet
	)

	if svc.jwtMaker == nil {
		// @todo
		//   return errors.New("can not load corredor scripts without jwt maker")
	}

	svc.log.Debug("reloading server scripts")
	rsp, err = svc.ssClient.List(ctx, &ServerScriptListRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return errors.Wrap(err, "could not load corredor scripts")
	}

	for _, script := range rsp.Scripts {
		if len(script.Errors) > 0 {
			continue
		}

		s := &Script{
			Name:        script.Name,
			Label:       script.Label,
			Description: script.Description,
			Errors:      script.Errors,
		}

		svc.registerTriggers(script.Name, script.Triggers...)

		svc.log.Debug(
			"loaded server script",
			zap.String("name", s.Name),
			zap.Int("triggers", len(script.Triggers)),
		)

		ss = append(ss, s)
	}

	svc.log.Info("server scripts reloaded", zap.Int("count", len(ss)))

	return
}

func (svc service) registerTriggers(script string, tt ...*Trigger) {
	var (
		ops     []eventbus.TriggerRegOp
		handler eventbus.Handler
		err     error

		log = svc.log.With(zap.String("script", script))
	)

	if ptrs, has := svc.registered[script]; has && len(ptrs) > 0 {
		// Unregister previously registered triggers
		svc.eventbus.Unregister(ptrs...)
		log.Debug(
			"triggers unregistered",
			zap.Uintptrs("triggers", ptrs),
		)
	}

	// Make room for new
	svc.registered[script] = make([]uintptr, 0)

	for i := range tt {
		if ops, err = svc.makeTriggerRegOpts(tt[i]); err != nil {
			log.Warn(
				"trigger could not be registered",
				zap.Error(err),
			)

			continue
		}

		var runAs = tt[i].RunAs

		handler = func(ctx context.Context, ev eventbus.Event) error {
			// Is this compatible event?

			if ce, ok := ev.(Event); ok {
				if len(runAs) > 0 {
					jwt, err := svc.jwtMaker(runAs)
					if err != nil {
						return err
					}

					ctx = auth.SetJwtToContext(ctx, jwt)
				}

				return svc.Exec(ctx, script, ce)
			}

			return nil
		}

		ptr := svc.eventbus.Register(handler, ops...)

		log.Debug(
			"trigger registered",
			zap.Uintptr("triggers", ptr),
		)
	}
}

// Exec finds and runs specific script with given event
//
// It does not do any (trigger, constraints) checking
//
// For consistency,
func (svc service) Exec(ctx context.Context, script string, event Event) (err error) {
	var (
		rsp *ExecResponse

		encArgs    map[string][]byte
		encResults = make(map[string][]byte)

		log = svc.log.With(
			zap.String("script", script),
			zap.Stringer("runAs", auth.GetIdentityFromContext(ctx)),
			zap.String("event", event.EventType()),
			zap.String("resource", event.ResourceType()),
		)
	)

	log.Debug("triggered")

	if encArgs, err = event.Encode(); err != nil {
		return
	}

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////
	// Additional ([]byte) arguments

	req := &ExecRequest{
		Name: script,
		Args: make(map[string]string),
	}

	if encArgs["authUser"], err = json.Marshal(auth.GetIdentityFromContext(ctx)); err != nil {
		return
	}

	// Cast arguments from map[string]json.RawMessage to map[string]string
	if encArgs != nil {
		for key := range encArgs {
			req.Args[key] = string(encArgs[key])
		}
	}

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////
	// Additional (string) arguments

	// pass security credentials
	req.Args["jwt"] = auth.GetJwtFromContext(ctx)

	// basic event/event info
	req.Args["event"] = event.EventType()
	req.Args["resource"] = event.ResourceType()

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////

	ctx, cancel := context.WithTimeout(ctx, svc.corredorOpt.DefaultExecTimeout)
	defer cancel()

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////

	var header, trailer metadata.MD
	rsp, err = svc.ssClient.Exec(
		ctx,
		req,
		grpc.WaitForReady(true),
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	if err != nil {
		log.Debug("corredor responded with error", zap.Error(err))
		return
	}

	log.Debug("corredor responded", zap.Any("result", rsp.Result))

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////

	// @todo process metadata (log, errors, stacktrace)
	spew.Dump("grpc exec header", header)
	spew.Dump("grpc exec trailer", trailer)

	if rsp.Result == nil {
		// No results
		return
	}

	// Cast map[string]json.RawMessage to map[string]string
	for key := range rsp.Result {
		encResults[key] = []byte(rsp.Result[key])
	}

	// Send results back to the event for decoding
	err = event.Decode(encResults)
	if err != nil {
		log.Debug("could not decode results", zap.Error(err))
		return
	}

	// Everything ok
	return
}

func (svc service) makeTriggerRegOpts(t *Trigger) (oo []eventbus.TriggerRegOp, err error) {
	if len(t.Events) == 0 {
		return nil, fmt.Errorf("can not generate trigger without at least one events")
	}
	if len(t.Resources) == 0 {
		return nil, fmt.Errorf("can not generate trigger without at least one resource")
	}

	oo = append(oo, eventbus.On(t.Events...))
	oo = append(oo, eventbus.For(t.Resources...))

	for i := range t.Constraints {
		oo = append(oo, eventbus.Constraint(
			t.Constraints[i].Name,
			t.Constraints[i].Op,
			t.Constraints[i].Value...,
		))
	}

	return
}
