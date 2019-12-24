package corredor

import (
	"context"
	"encoding/json"
	"fmt"

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
		opt options.CorredorOpt

		// list of all registered triggers
		//
		registered map[string][]uintptr

		// list of all registered onManual triggers & scripts
		//   map[resource][script-name]bool
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

const onManualEventType = "onManual"

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

	gService = NewService(conn, logger, opt)
	return
}

func NewService(conn *grpc.ClientConn, logger *zap.Logger, opt options.CorredorOpt) *service {
	return &service{
		ssClient:   NewServerScriptsClient(conn),
		csClient:   NewClientScriptsClient(conn),
		log:        logger.Named("corredor"),
		registered: make(map[string][]uintptr),
		manual:     make(map[string]map[string]bool),
		eventbus:   eventbus.Default(),
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

	if _, ok := svc.manual[res]; !ok {
		return errors.Errorf("unregistered onManual resource %q", res)
	}

	if _, ok := svc.manual[res][script]; !ok {
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

	svc.manual = make(map[string]map[string]bool)
	svc.sScripts = make([]*Script, len(rsp.Scripts))

	for i, script := range rsp.Scripts {
		svc.sScripts[i] = &Script{
			Name:        script.Name,
			Label:       script.Label,
			Description: script.Description,
			Errors:      script.Errors,
			Triggers:    script.Triggers,
		}

		if len(script.Errors) == 0 {
			svc.registerTriggers(script)
		}
	}
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

	svc.cScripts = make([]*Script, len(rsp.Scripts))

	for i, script := range rsp.Scripts {
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

// Creates handler function for eventbus subsystem
//
// If trigger has "onManual"
func (svc service) registerTriggers(script *ServerScript) {
	var (
		ops     []eventbus.TriggerRegOp
		handler eventbus.Handler
		err     error

		log = svc.log.With(zap.String("scriptName", script.Name))
	)

	if ptrs, has := svc.registered[script.Name]; has && len(ptrs) > 0 {
		// Unregister previously registered triggers
		svc.eventbus.Unregister(ptrs...)
		log.Debug(
			"triggers unregistered",
			zap.Uintptrs("triggers", ptrs),
		)
	}

	// Make room for new
	svc.registered[script.Name] = make([]uintptr, 0)

	for i := range script.Triggers {
		// We're modifying trigger in the loop,
		// so let's make a copy we can play with
		trigger := *script.Triggers[i]

		if popOnManualEventType(&trigger) {
			for _, res := range trigger.Resources {
				if svc.manual[res] == nil {
					svc.manual[res] = make(map[string]bool)
				}

				svc.manual[res][script.Name] = true
			}

			log.Debug("manual trigger registered", zap.Strings("resources", trigger.Resources))

			if len(trigger.Events) == 0 {
				// We've removed the last event
				//
				// break now to prevent code below to
				// complain about missing event types
				continue
			}
		}

		if ops, err = svc.makeTriggerRegOpts(&trigger); err != nil {
			log.Warn(
				"trigger could not be registered",
				zap.Error(err),
			)

			continue
		}

		var runAs = trigger.RunAs

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

				return svc.exec(ctx, script.Name, ce)
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
		log.Debug("corredor responded with error", zap.Error(err))
		return
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

func encodeArguments(args map[string]string, key string, val interface{}) (err error) {
	var tmp []byte

	if tmp, err = json.Marshal(val); err != nil {
		return
	}

	args[key] = string(tmp)
	return
}
