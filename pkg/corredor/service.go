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
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	service struct {
		// stores corredor connection options
		// for when we're doing lazy setup
		opt options.CorredorOpt

		// list of all registered triggers
		//   map[<script-name>]
		registered map[string][]uintptr

		// list of all registered onManual triggers & scripts
		//   map[<script-name>][<resource>] = <run-as>
		manual map[string]map[string]string

		// Combined list of client and server scripts
		sScripts ScriptSet
		cScripts ScriptSet

		conn *grpc.ClientConn

		ssClient ServerScriptsClient
		csClient ClientScriptsClient

		log *zap.Logger

		eventRegistry  eventRegistry
		authTokenMaker authTokenMaker
		users          userFinder
	}

	Event interface {
		eventbus.Event

		// Encode (event) to arguments passed to
		// event handlers ([automation ]script runner)
		Encode() (map[string][]byte, error)

		// Decodes received data back to event
		Decode(map[string][]byte) error
	}

	eventRegistry interface {
		Register(h eventbus.Handler, ops ...eventbus.TriggerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	userFinder interface {
		FindByAny(interface{}) (*types.User, error)
	}

	authTokenMaker interface {
		Encode(auth.Identifiable) string
	}
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
func Setup(logger *zap.Logger, opt options.CorredorOpt) (err error) {
	if gCorredor != nil {
		// Prevent multiple initializations
		return
	}

	gCorredor = NewService(logger, opt)
	return
}

func NewService(logger *zap.Logger, opt options.CorredorOpt) *service {
	return &service{
		log: logger.Named("corredor"),
		opt: opt,

		registered: make(map[string][]uintptr),
		manual:     make(map[string]map[string]string),

		authTokenMaker: auth.DefaultJwtHandler,
		eventRegistry:  eventbus.Service(),
	}
}

func (svc *service) Connect(ctx context.Context) (err error) {
	if err = svc.connect(ctx); err != nil {
		return
	}

	svc.ssClient = NewServerScriptsClient(svc.conn)
	svc.csClient = NewClientScriptsClient(svc.conn)

	return
}

func (svc *service) connect(ctx context.Context) (err error) {
	if svc.conn, err = NewConnection(ctx, svc.opt, svc.log); err != nil {
		return
	}

	return
}

func (svc *service) SetEventRegistry(er eventRegistry) {
	svc.eventRegistry = er
}

func (svc *service) SetAuthTokenMaker(atm authTokenMaker) {
	svc.authTokenMaker = atm
}

func (svc *service) SetUserFinder(uf userFinder) {
	svc.users = uf
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

		ok    bool
		runAs string
	)

	if onManualEventType != evt {
		return errors.Errorf("triggered event type is not onManual (%q)", evt)
	}

	if _, ok = svc.manual[script]; !ok {
		return errors.Errorf("unregistered onManual script %q", script)
	}

	if runAs, ok = svc.manual[script][res]; !ok {
		return errors.Errorf("unregistered onManual script %q for resource %q", script, res)
	}

	return svc.exec(ctx, script, runAs, event)
}

func (svc *service) loadServerScripts(ctx context.Context) {
	var (
		err error
		rsp *ServerScriptListResponse
	)

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
			svc.eventRegistry.Unregister(ptrs...)
		}
	}

	// Reset indexes
	svc.registered = make(map[string][]uintptr)
	svc.manual = make(map[string]map[string]string)

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
		trigger := script.Triggers[i]

		if len(trigger.Events) == 0 {
			// We've removed the last event
			//
			// break now to prevent code below to
			// complain about missing event types
			continue
		}

		if ops, err = makeTriggerOpts(trigger); err != nil {
			log.Warn(
				"could not make trigger options",
				zap.Error(err),
			)

			continue
		}

		ptr := svc.eventRegistry.Register(func(ctx context.Context, ev eventbus.Event) (err error) {
			// Is this compatible event?
			if cce, ok := ev.(Event); ok {
				// Can only work with corteza compatible events
				return svc.exec(ctx, script.Name, trigger.RunAs, cce)
			}

			return nil
		}, ops...)
		ptrs = append(ptrs, ptr)
	}

	return ptrs
}

// Exec finds and runs specific script with given event
//
// It does not do any constraints checking - this is the responsibility of the
// individual event implemntation
func (svc service) exec(ctx context.Context, script string, runAs string, event Event) (err error) {
	var (
		rsp *ExecResponse

		encodedEvent   map[string][]byte
		encodedResults = make(map[string][]byte)

		log = svc.log.With(
			zap.String("script", script),
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

	// Resolve/expand invoker user details from the context
	invoker, err := svc.users.FindByAny(ctx)
	if err != nil {
		return err
	}

	log = log.With(zap.Stringer("invoker", invoker))
	if err = encodeArguments(req.Args, "invoker", invoker); err != nil {
		return
	}

	if len(runAs) > 0 {
		if !svc.opt.RunAsEnabled {
			return errors.New("could not make runner context, run-as disabled")
		}

		var definer auth.Identifiable

		// Run this script as defined user (definer)
		//
		// We search for the defined (run-as) user,
		// assign it to authUser argument and make an
		// authentication token for it
		definer, err = svc.users.FindByAny(runAs)
		if err != nil {
			return err
		}

		log = log.With(zap.Stringer("run-as", definer))

		// current (authenticated) user
		if err = encodeArguments(req.Args, "authUser", definer); err != nil {
			return
		}

		if err = encodeArguments(req.Args, "authToken", svc.authTokenMaker.Encode(definer)); err != nil {
			return
		}

	} else {
		// Run script with the same user that invoked it

		// current (authenticated) user
		if err = encodeArguments(req.Args, "authUser", invoker); err != nil {
			return
		}

		if err = encodeArguments(req.Args, "authToken", svc.authTokenMaker.Encode(invoker)); err != nil {
			return
		}
	}

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////
	// Additional (string) arguments

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
		// See if this was a "soft abort"
		//
		// This means, we do not make any logs of this just
		// tell the caller that the call was aborted
		s := status.Convert(err)
		if s != nil && s.Code() == codes.Aborted {
			// Special care for errors with Aborted code
			msg := s.Message()

			if len(msg) == 0 {
				// No extra message, fallback to "aborted"
				msg = "Aborted"
			}

			return errors.New(msg)
		}

		log.Warn("corredor responded with error", zap.Error(err))
		return errors.New("failed to execute corredor script")
	}

	log.Info("executed", zap.Any("result", rsp.Result))

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
