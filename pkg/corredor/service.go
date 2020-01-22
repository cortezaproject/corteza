package corredor

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
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
		//   map[<script-name>][<resource>] = true
		manual map[string]map[string]bool

		// Combined list of client and server scripts
		sScripts ScriptSet
		cScripts ScriptSet

		conn *grpc.ClientConn

		ssClient ServerScriptsClient
		csClient ClientScriptsClient

		log *zap.Logger

		eventRegistry  eventRegistry
		authTokenMaker authTokenMaker

		// Services to help with script security
		// we'll find users (runAs) and roles (allow, deny) for
		users userFinder
		roles roleFinder

		// set of permission rules, generated from security info of each script
		permissions permissions.RuleSet
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
		Register(h eventbus.HandlerFn, ops ...eventbus.HandlerRegOp) uintptr
		Unregister(ptrs ...uintptr)
	}

	userFinder interface {
		FindByAny(interface{}) (*types.User, error)
	}

	roleFinder interface {
		FindByAny(interface{}) (*types.Role, error)
	}

	authTokenMaker interface {
		Encode(auth.Identifiable) string
	}

	permissionRuleChecker interface {
		Check(res permissions.Resource, op permissions.Operation, roles ...uint64) permissions.Access
	}
)

const onManualEventType = "onManual"

var (
	// Global corredor service
	gCorredor *service
)

const (
	permOpExec permissions.Operation = "exec"
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
		manual:     make(map[string]map[string]bool),

		authTokenMaker: auth.DefaultJwtHandler,
		eventRegistry:  eventbus.Service(),
		permissions:    permissions.RuleSet{},
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

func (svc *service) SetRoleFinder(rf roleFinder) {
	svc.roles = rf
}

func (svc *service) Load(ctx context.Context) {
	go svc.loadServerScripts(ctx)
	go svc.loadClientScripts(ctx)
}

// Find returns filtered list of scripts that can be manually triggered
func (svc service) Find(ctx context.Context, filter Filter) (out ScriptSet, f Filter, err error) {
	f = filter

	var (
		tmp ScriptSet

		scriptFilter = svc.makeScriptFilter(ctx, f)
	)

	if !f.ExcludeServerScripts {
		tmp, err = svc.sScripts.Filter(scriptFilter)
		out = append(out, tmp...)
	}

	if !f.ExcludeClientScripts {
		tmp, err = svc.cScripts.Filter(scriptFilter)
		out = append(out, tmp...)
	}

	f.Count = uint(len(out))

	return
}

// An enhanced version of basic script filter maker (from util.go)
// that (after basic filtering) also does RBAC check for each script
func (svc service) makeScriptFilter(ctx context.Context, f Filter) func(s *Script) (b bool, err error) {
	var (
		base = f.makeFilterFn()
	)

	return func(s *Script) (b bool, err error) {
		if b, err = base(s); !b {
			return
		}

		return svc.canExec(ctx, s.Name), nil
	}
}

// ExecOnManual verifies permissions, event and script and sends exec request to corredor
func (svc service) ExecOnManual(ctx context.Context, scriptName string, event Event) (err error) {
	var (
		res = event.ResourceType()
		evt = event.EventType()

		script *Script

		ok    bool
		runAs string
	)

	if onManualEventType != evt {
		return errors.Errorf("triggered event type is not onManual (%q)", evt)
	}

	if len(scriptName) == 0 {
		return errors.Errorf("script name not provided (%q)", scriptName)
	}

	if _, ok = svc.manual[scriptName]; !ok {
		return errors.Errorf("unregistered onManual script %q", scriptName)
	}

	if _, ok = svc.manual[scriptName][res]; !ok {
		return errors.Errorf("unregistered onManual script %q for resource %q", scriptName, res)
	}

	if script = svc.sScripts.FindByName(scriptName); script == nil {
		return errors.Errorf("nonexistent script (%q)", scriptName)
	}

	if !svc.canExec(ctx, scriptName) {
		return errors.Errorf("permission to execute %s denied", scriptName)
	}

	return svc.exec(ctx, scriptName, runAs, event)
}

// Can current user execute this script
func (svc service) canExec(ctx context.Context, script string) bool {
	u := auth.GetIdentityFromContext(ctx)
	if auth.IsSuperUser(u) {
		return true
	}

	return svc.permissions.Check(permissions.Resource(script), permOpExec, u.Roles()...) != permissions.Deny
}

func (svc *service) loadServerScripts(ctx context.Context) {
	var (
		err error
		rsp *ServerScriptListResponse
	)

	svc.log.Debug("reloading server scripts")

	ctx, cancel := context.WithTimeout(ctx, svc.opt.ListTimeout)
	defer cancel()

	rsp, err = svc.ssClient.List(ctx, &ServerScriptListRequest{}, grpc.WaitForReady(true))
	if err != nil {
		svc.log.Error("could not load corredor server scripts", zap.Error(err))
		return
	}

	svc.registerServerScripts(rsp.Scripts...)
}

// Registers Corredor scripts to eventbus and list of manual scripts
func (svc *service) registerServerScripts(ss ...*ServerScript) {
	var (
		permRuleGenerator = func(script string, access permissions.Access, roles ...string) (permissions.RuleSet, error) {
			out := make([]*permissions.Rule, len(roles))
			for i, role := range roles {
				if r, err := svc.roles.FindByAny(role); err != nil {
					return nil, err
				} else {
					out[i] = &permissions.Rule{
						RoleID:    r.ID,
						Resource:  permissions.Resource(script),
						Operation: permOpExec,
						Access:    access,
					}
				}

			}
			return out, nil
		}

		u     *types.User
		err   error
		runAs = ""
	)

	svc.sScripts = make([]*Script, 0, len(ss))

	// Remove all previously registered triggers
	for _, ptrs := range svc.registered {
		if len(ptrs) > 0 {
			svc.eventRegistry.Unregister(ptrs...)
		}
	}

	// Reset indexes
	svc.registered = make(map[string][]uintptr)
	svc.manual = make(map[string]map[string]bool)

	// Reset security
	svc.permissions = permissions.RuleSet{}

	for _, script := range ss {
		var (
			// collectors for allow&deny rules
			// we'll merge
			allow = permissions.RuleSet{}
			deny  = permissions.RuleSet{}
		)

		if nil != svc.sScripts.FindByName(script.Name) {
			// Do not allow duplicated scripts
			continue
		}

		s := &Script{
			Name:        script.Name,
			Label:       script.Label,
			Description: script.Description,
			Errors:      script.Errors,
			Triggers:    script.Triggers,
			Security:    &ScriptSecurity{Security: script.Security},
		}

		scriptErrPush := func(err error, msg string) {
			s.Errors = append(s.Errors, errors.Wrap(err, msg).Error())
		}

		if len(s.Errors) == 0 {

			if manual := pluckManualTriggers(script); len(manual) > 0 {

				if script.Security != nil {
					runAs = script.Security.RunAs

					if runAs != "" {
						// Prefetch run-as user
						if u, err = svc.users.FindByAny(runAs); err != nil {
							scriptErrPush(err, "could not load run-as user security info")
						} else {
							s.Security.runAs = u.ID
						}
					}

					if allow, err = permRuleGenerator(script.Name, permissions.Allow, script.Security.Allow...); err != nil {
						scriptErrPush(err, "could not load allow role security info")
					}

					if deny, err = permRuleGenerator(script.Name, permissions.Deny, script.Security.Deny...); err != nil {
						scriptErrPush(err, "could not load deny role security info")
					}

					svc.permissions = append(svc.permissions, allow...)
					svc.permissions = append(svc.permissions, deny...)
				}

				if len(s.Errors) == 0 {
					svc.manual[script.Name] = manual
				}
			}

			if len(s.Errors) == 0 {
				svc.registered[script.Name] = svc.registerTriggers(script)
			}
		}

		// Even if there are errors, we'll append the scripts
		// because we need to serve them as a list for script management
		svc.sScripts = append(svc.sScripts, s)

		if len(s.Errors) == 0 {
			svc.log.Debug(
				"script registered",
				zap.String("script", s.Name),
				zap.Stringer("security", s.Security),
				zap.Int("manual", len(svc.manual[script.Name])),
				zap.Int("triggers", len(svc.registered[script.Name])),
			)
		} else {
			svc.log.Warn(
				"script loaded with errors",
				zap.String("script", s.Name),
				zap.Strings("errors", s.Errors),
			)
		}
	}
}

// Creates handler function for eventbus subsystem
//
// If trigger has "onManual" event type, it removes it and
// registers that script to the list of manual triggers
func (svc *service) registerTriggers(script *ServerScript) []uintptr {
	var (
		ops  []eventbus.HandlerRegOp
		err  error
		ptrs = make([]uintptr, 0, len(script.Triggers))

		log = svc.log.With(zap.String("script", script.Name))

		runAs string
	)

	if script.Security != nil {
		runAs = script.Security.RunAs
	}

	for i := range script.Triggers {
		// We're modifying trigger in the loop,
		// so let's make a copy we can play with
		trigger := script.Triggers[i]

		if len(trigger.EventTypes) == 0 {
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
			if ev, ok := ev.(Event); ok {
				// Can only work with corteza compatible events
				return svc.exec(ctx, script.Name, runAs, ev)
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

	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"x-request-id": []string{middleware.GetReqID(ctx)},
	})

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
		log.Debug(
			"could not decode results",
			zap.Error(err),
			zap.Any("results", encodedResults),
		)
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

	ctx, cancel := context.WithTimeout(ctx, svc.opt.ListTimeout)
	defer cancel()

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
