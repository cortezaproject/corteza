package corredor

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

type (
	service struct {
		// stores corredor connection options
		// for when we're doing lazy setup
		opt options.CorredorOpt

		// list of all registered triggers
		//
		registered map[string][]uintptr

		// list of all registered onManual triggers & scripts
		//   map[script-name][resource]bool
		manual map[string]map[string]bool

		// Combined list of client and server scripts
		sScripts ScriptSet
		cScripts ScriptSet

		ssClient ServerScriptsClient
		csClient ClientScriptsClient

		log *zap.Logger

		eventbus eventRegistrator
		jwtMaker AuthTokenMaker
	}

	Event interface {
		eventbus.Event

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

const onManualEventType = "onManual"

var (
	// Global corredor service
	gCorredor *service
)

func Service() *service {
	return gCorredor
}

// Start connects to Corredor & initialize service
func Start(ctx context.Context, logger *zap.Logger, opt options.CorredorOpt) (err error) {
	if gCorredor != nil {
		// Prevent multiple initializations
		return
	}

	var (
		conn *grpc.ClientConn
	)

	if conn, err = NewConnection(ctx, opt, logger); err != nil {
		return
	}

	gCorredor = NewService(conn, eventbus.Service(), logger, opt)
	return
}

func NewService(conn *grpc.ClientConn, er eventRegistrator, logger *zap.Logger, opt options.CorredorOpt) *service {
	return &service{
		ssClient:   NewServerScriptsClient(conn),
		csClient:   NewClientScriptsClient(conn),
		log:        logger.Named("corredor"),
		registered: make(map[string][]uintptr),
		manual:     make(map[string]map[string]bool),
		eventbus:   er,
		opt:        opt,
	}
}

func (svc *service) SetJwtMaker(fn AuthTokenMaker) {
	svc.jwtMaker = fn
}

func (svc *service) Load(ctx context.Context) {
	go svc.loadServerScripts(ctx)
	go svc.loadClientScripts(ctx)
}

// FindManual returns filtered list of scripts that can be manually triggered
func (svc service) FindOnManual(filter ManualScriptFilter) (out ScriptSet, f ManualScriptFilter, err error) {
	f = filter

	var (
		tmp ScriptSet

		scriptFilter = makeScriptFilter(f)
	)

	if !f.ExcludeServerScripts {
		tmp, err = svc.sScripts.Filter(scriptFilter)
		out = append(out, tmp...)
	}

	if !f.ExcludeClientScripts {
		tmp, err = svc.cScripts.Filter(scriptFilter)
		out = append(out, tmp...)
	}

	return
}

// ExecOnManual verifies request & executes
func (svc service) ExecOnManual(ctx context.Context, script string, event Event) (err error) {
	var (
		res = event.ResourceType()
		evt = event.EventType()
	)

	if onManualEventType != evt {
		return errors.Errorf("triggered event type is not onManual (%q)", evt)
	}

	if _, ok := svc.manual[script]; !ok {
		return errors.Errorf("unregistered onManual script %q", script)
	}

	if _, ok := svc.manual[script][res]; !ok {
		return errors.Errorf("unregistered onManual script %q for resource %q", script, res)
	}

	return svc.exec(ctx, script, event)
}

func (svc *service) loadServerScripts(ctx context.Context) {
	var (
		err error
		rsp *ServerScriptListResponse
	)

	if svc.jwtMaker == nil {
		// @todo
		//   return errors.New("can not load corredor scripts without jwt maker")
	}

	svc.log.Debug("reloading server scripts")
	rsp, err = svc.ssClient.List(ctx, &ServerScriptListRequest{}, grpc.WaitForReady(true))
	if err != nil {
		svc.log.Error("could not load corredor server scripts", zap.Error(err))
		return
	}

	svc.registerServerScripts(rsp.Scripts...)
}

// Registers Corredor scripts to eventbus and list of manual scripts
func (svc *service) registerServerScripts(ss ...*ServerScript) {
	svc.sScripts = make([]*Script, len(ss))

	// Remove all previously registered triggers
	for _, ptrs := range svc.registered {
		if len(ptrs) > 0 {
			svc.eventbus.Unregister(ptrs...)
		}
	}

	// Reset indexes
	svc.registered = make(map[string][]uintptr)
	svc.manual = make(map[string]map[string]bool)

	for i, script := range ss {
		svc.sScripts[i] = &Script{
			Name:        script.Name,
			Label:       script.Label,
			Description: script.Description,
			Errors:      script.Errors,
			Triggers:    script.Triggers,
		}

		if len(script.Errors) == 0 {
			svc.manual[script.Name] = pluckManualTriggers(script)
			svc.registered[script.Name] = svc.registerTriggers(script)
		}

		svc.log.Debug(
			"registered",
			zap.String("script", script.Name),
			zap.Int("manual", len(svc.manual[script.Name])),
			zap.Int("triggers", len(svc.registered[script.Name])),
		)
	}
}

// Creates handler function for eventbus subsystem
//
// If trigger has "onManual" event type, it removes it and
// registers that script to the list of manual triggers
func (svc *service) registerTriggers(script *ServerScript) []uintptr {
	var (
		ops  []eventbus.TriggerRegOp
		err  error
		ptrs = make([]uintptr, 0, len(script.Triggers))

		log = svc.log.With(zap.String("script", script.Name))
	)

	for i := range script.Triggers {
		// We're modifying trigger in the loop,
		// so let's make a copy we can play with
		trigger := *script.Triggers[i]

		if len(trigger.Events) == 0 {
			// We've removed the last event
			//
			// break now to prevent code below to
			// complain about missing event types
			continue
		}

		if ops, err = makeTriggerOpts(&trigger); err != nil {
			log.Warn(
				"could not make trigger options",
				zap.Error(err),
			)

			continue
		}

		ptr := svc.eventbus.Register(makeEventHandler(svc, script.Name, trigger.RunAs), ops...)
		ptrs = append(ptrs, ptr)
	}

	return ptrs
}

// Exec finds and runs specific script with given event
//
// It does not do any constraints checking - this is the responsibility of the
// individual event implemntation
func (svc service) exec(ctx context.Context, script string, event Event) (err error) {
	var (
		rsp *ExecResponse

		encodedEvent   map[string][]byte
		encodedResults = make(map[string][]byte)

		log = svc.log.With(
			zap.String("script", script),
			zap.Stringer("runAs", auth.GetIdentityFromContext(ctx)),
			zap.String("event", event.EventType()),
			zap.String("resource", event.ResourceType()),
		)
	)

	log.Debug("triggered")

	if encodedEvent, err = event.Encode(); err != nil {
		return
	}

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////
	// Additional ([]byte) arguments

	req := &ExecRequest{
		Name: script,
		Args: make(map[string]string),
	}

	// Cast arguments from map[string]json.RawMessage to map[string]string
	for key := range encodedEvent {
		req.Args[key] = string(encodedEvent[key])
	}

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////
	// Additional (string) arguments

	// pass security credentials
	if err = encodeArguments(req.Args, "authUser", auth.GetIdentityFromContext(ctx)); err != nil {
		return
	}
	if err = encodeArguments(req.Args, "jwt", auth.GetJwtFromContext(ctx)); err != nil {
		return
	}

	// basic event/event info
	if err = encodeArguments(req.Args, "event", event.EventType()); err != nil {
		return
	}
	if err = encodeArguments(req.Args, "resource", event.ResourceType()); err != nil {
		return
	}

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////

	ctx, cancel := context.WithTimeout(ctx, svc.opt.DefaultExecTimeout)
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
		s := status.Convert(err)
		if s != nil && s.Code() == codes.Aborted {
			// Special care for errors with Aborted code
			msg := s.Message()

			if len(msg) == 0 {
				msg = "Aborted"
			}

			return errors.New(msg)
		}

		log.Debug("corredor responded with error", zap.Error(err))
		return errors.New("failed to execute corredor script")
	}

	log.Debug("corredor responded", zap.Any("result", rsp.Result))

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////

	// @todo process metadata (log, errors, stacktrace)
	// spew.Dump("grpc exec header", header)
	// spew.Dump("grpc exec trailer", trailer)

	if rsp.Result == nil {
		// No results
		return
	}

	// Cast map[string]json.RawMessage to map[string]string
	for key := range rsp.Result {
		encodedResults[key] = []byte(rsp.Result[key])
	}

	// Send results back to the event for decoding
	err = event.Decode(encodedResults)
	if err != nil {
		log.Debug("could not decode results", zap.Error(err))
		return
	}

	// Everything ok
	return
}

func (svc *service) loadClientScripts(ctx context.Context) {
	var (
		err error
		rsp *ClientScriptListResponse
	)

	svc.log.Debug("reloading client scripts")
	rsp, err = svc.csClient.List(ctx, &ClientScriptListRequest{}, grpc.WaitForReady(true))
	if err != nil {
		svc.log.Error("could not load corredor client scripts", zap.Error(err))
		return
	}

	svc.registerClientScripts(rsp.Scripts...)
}

func (svc *service) registerClientScripts(ss ...*ClientScript) {
	svc.cScripts = make([]*Script, len(ss))

	for i, script := range ss {
		svc.cScripts[i] = &Script{
			Name:        script.Name,
			Label:       script.Label,
			Description: script.Description,
			Errors:      script.Errors,
			Triggers:    script.Triggers,
			Bundle:      script.Bundle,
			Type:        script.Type,
		}
	}
}
