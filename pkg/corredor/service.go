package corredor

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"

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

		// list of all registered event handlers
		//   map[<script-name>]
		registered map[string][]uintptr

		// list of all registered explicitly executable script
		//   map[<script-name>][<resource>] = true
		explicit map[string]map[string]bool

		// Combined list of client and server scripts
		sScripts   ScriptSet
		sScriptsTS time.Time
		cScripts   ScriptSet
		cScriptsTS time.Time

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

	ScriptArgs interface {
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
		explicit:   make(map[string]map[string]bool),

		authTokenMaker: auth.DefaultJwtHandler,
		eventRegistry:  eventbus.Service(),
		permissions:    permissions.RuleSet{},
	}
}

func (svc *service) Connect(ctx context.Context) (err error) {
	if !svc.opt.Enabled {
		return
	}

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

// Watch watches for changes
func (svc *service) Watch(ctx context.Context) {
	go func() {
		defer sentry.Recover()
		var ticker = time.NewTicker(svc.opt.ListRefresh)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				svc.Load(ctx)
			}
		}
	}()

	svc.log.Debug("watcher initialized")
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
	if !svc.opt.Enabled {
		return
	}

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

// Exec verifies permissions, event and script and sends exec request to corredor
func (svc service) Exec(ctx context.Context, scriptName string, args ScriptArgs) (err error) {
	if !svc.opt.Enabled {
		return
	}

	var (
		res    = args.ResourceType()
		script *Script

		ok    bool
		runAs string
	)

	if len(scriptName) == 0 {
		return errors.Errorf("script name not provided (%q)", scriptName)
	}

	if _, ok = svc.explicit[scriptName]; !ok {
		return errors.Errorf("unregistered explicit script %q", scriptName)
	}

	if _, ok = svc.explicit[scriptName][res]; !ok {
		return errors.Errorf("unregistered explicit script %q for resource %q", scriptName, res)
	}

	if script = svc.sScripts.FindByName(scriptName); script == nil {
		return errors.Errorf("nonexistent script (%q)", scriptName)
	}

	if !svc.canExec(ctx, scriptName) {
		return errors.Errorf("permission to execute %s denied", scriptName)
	}

	if script.Security != nil {
		runAs = script.Security.RunAs
	}

	return svc.exec(ctx, scriptName, runAs, args)
}

// Can current user execute this script
//
// This is used only in case of explicit execution (onManual) and never when
// scripts are executed implicitly (deferred, before/after...)
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

	ctx, cancel := context.WithTimeout(ctx, svc.opt.ListTimeout)
	defer cancel()

	if !svc.sScriptsTS.IsZero() {
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
			"if-modified-since": []string{svc.sScriptsTS.Format(time.RFC3339)},
		})
	}

	rsp, err = svc.ssClient.List(ctx, &ServerScriptListRequest{}, grpc.WaitForReady(true))
	if err != nil {
		svc.log.Error("could not load corredor server scripts", zap.Error(err))
		return
	}

	svc.sScriptsTS = time.Now()

	if len(rsp.Scripts) > 0 {
		svc.log.Debug("reloading server scripts")
		svc.registerServerScripts(rsp.Scripts...)
	}
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
	svc.explicit = make(map[string]map[string]bool)

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
			if manual := mapExplicitTriggers(script); len(manual) > 0 {
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

				svc.explicit[script.Name] = manual
			}

			svc.registered[script.Name] = svc.registerTriggers(script)

			svc.log.Debug(
				"script registered",
				zap.String("script", s.Name),
				zap.Int("explicit", len(svc.explicit[script.Name])),
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
		if ops, err = triggerToHandlerOps(script.Triggers[i]); err != nil {
			log.Warn(
				"could not make trigger options",
				zap.Error(err),
			)

			continue
		}

		if len(ops) == 0 {
			continue
		}

		ptr := svc.eventRegistry.Register(func(ctx context.Context, ev eventbus.Event) (err error) {
			// Is this compatible event?
			if ce, ok := ev.(ScriptArgs); ok {
				// Can only work with corteza compatible events
				return svc.exec(ctx, script.Name, runAs, ce)
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
func (svc service) exec(ctx context.Context, script string, runAs string, args ScriptArgs) (err error) {
	var (
		requestId = middleware.GetReqID(ctx)

		rsp *ExecResponse

		invoker auth.Identifiable

		encodedEvent   map[string][]byte
		encodedResults = make(map[string][]byte)

		log = svc.log.With(
			zap.String("script", script),
			zap.String("runAs", runAs),
			zap.String("args", args.EventType()),
			zap.String("resource", args.ResourceType()),
		)
	)

	log.Debug("triggered")

	if encodedEvent, err = args.Encode(); err != nil {
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

	// Resolve/expand invoker user details from the context (if present
	if i := auth.GetIdentityFromContext(ctx); i.Valid() {
		invoker, err = svc.users.FindByAny(i)
		if err != nil {
			return err
		}

		log = log.With(zap.Stringer("invoker", invoker))

		if err = encodeArguments(req.Args, "invoker", invoker); err != nil {
			return
		}
	}

	if len(runAs) > 0 {
		if !svc.opt.RunAsEnabled {
			return errors.New("could not make runner context, run-as disabled")
		}

		var definer auth.Identifiable

		// Run this script as defined user
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

	} else if invoker != nil {
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

	// basic args/event info
	if err = encodeArguments(req.Args, "args", args.EventType()); err != nil {
		return
	}
	if err = encodeArguments(req.Args, "resource", args.ResourceType()); err != nil {
		return
	}

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////

	ctx, cancel := context.WithTimeout(
		// We need a new, independent context here
		// to be sure this is executed safely & fully
		// without any outside interfeance (cancellation, timeouts)
		context.Background(),
		svc.opt.DefaultExecTimeout,
	)
	defer cancel()

	// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// ////

	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"x-request-id": []string{requestId},
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

	// Send results back to the args for decoding
	err = args.Decode(encodedResults)
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

	ctx, cancel := context.WithTimeout(ctx, svc.opt.ListTimeout)
	defer cancel()

	if !svc.sScriptsTS.IsZero() {
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
			"if-modified-since": []string{svc.sScriptsTS.Format(time.RFC3339)},
		})
	}

	rsp, err = svc.csClient.List(ctx, &ClientScriptListRequest{}, grpc.WaitForReady(true))
	if err != nil {
		svc.log.Error("could not load corredor client scripts", zap.Error(err))
		return
	}

	svc.sScriptsTS = time.Now()

	if len(rsp.Scripts) > 0 {
		svc.log.Debug("reloading client scripts")
		svc.registerClientScripts(rsp.Scripts...)
	}
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

func (svc *service) GetBundle(ctx context.Context, name, bType string) *Bundle {
	if !svc.opt.Enabled {
		return nil
	}

	var (
		err error
		rsp *BundleResponse
	)

	ctx, cancel := context.WithTimeout(ctx, svc.opt.ListTimeout)
	defer cancel()

	rsp, err = svc.csClient.Bundle(ctx, &BundleRequest{Name: name}, grpc.WaitForReady(true))
	if err != nil {
		svc.log.Error("could not load client scripts bundle from corredor", zap.Error(err))
		return nil
	}

	for _, b := range rsp.Bundles {
		if b.Type == bType {
			return b
		}
	}

	return nil
}

func (svc *service) Debug() interface{} {
	return map[string]interface{}{
		"registered":     svc.registered,
		"explicit":       svc.explicit,
		"server-scripts": svc.sScripts,
		"client-scripts": svc.cScripts,
	}
}
