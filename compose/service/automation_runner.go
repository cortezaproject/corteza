package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cortezaproject/corteza-server/compose/proto"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/automation/corredor"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	automationRunner struct {
		opt          AutomationRunnerOpt
		ac           automationRunnerAccessControler
		logger       *zap.Logger
		runner       corredor.ScriptRunnerClient
		scriptFinder automationScriptsFinder
		jwtEncoder   auth.TokenEncoder
	}

	automationScriptsFinder interface {
		Watch(ctx context.Context)
		WatchScheduled(ctx context.Context, runner automation.DeferredAutomationRunner)
		FindRunnableScripts(resource, event string, cc ...automation.TriggerConditionChecker) automation.ScriptSet
	}

	automationRunnerAccessControler interface {
		CanRunAutomationTrigger(ctx context.Context, r *automation.Trigger) bool
	}

	AutomationRunnerOpt struct {
		ApiBaseURLSystem    string
		ApiBaseURLMessaging string
		ApiBaseURLCompose   string
	}
)

const (
	AutomationResourceRecord = "compose:record"
)

func AutomationRunner(opt AutomationRunnerOpt, f automationScriptsFinder, r corredor.ScriptRunnerClient) automationRunner {
	var svc = automationRunner{
		opt: opt,

		ac: DefaultAccessControl,

		scriptFinder: f,
		runner:       r,

		logger:     DefaultLogger.Named("automationRunner"),
		jwtEncoder: auth.DefaultJwtHandler,
	}

	return svc
}

func (svc automationRunner) Watch(ctx context.Context) {
	svc.scriptFinder.Watch(ctx)
	svc.scriptFinder.WatchScheduled(ctx, svc)
}

// BeforeRecordCreate - run scripts before record is created
//
// This is implicitly called, no extra security check is needed
func (svc automationRunner) BeforeRecordCreate(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	return svc.findRecordScripts("beforeCreate", m.ID).Walk(
		svc.makeRecordScriptRunner(ctx, ns, m, r, false),
	)
}

// AfterRecordCreate - run scripts before record is created
//
// This is implicitly called, no extra security check is needed
func (svc automationRunner) AfterRecordCreate(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	return svc.findRecordScripts("afterCreate", m.ID).Walk(
		svc.makeRecordScriptRunner(ctx, ns, m, r, true),
	)
}

// BeforeRecordUpdate - run scripts before record is updated
//
// This is implicitly called, no extra security check is needed
func (svc automationRunner) BeforeRecordUpdate(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	return svc.findRecordScripts("beforeUpdate", m.ID).Walk(
		svc.makeRecordScriptRunner(ctx, ns, m, r, false),
	)
}

// AfterRecordUpdate - run scripts before record is updated
//
// This is implicitly called, no extra security check is needed
func (svc automationRunner) AfterRecordUpdate(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	return svc.findRecordScripts("afterUpdate", m.ID).Walk(
		svc.makeRecordScriptRunner(ctx, ns, m, r, true),
	)
}

// BeforeRecordDelete - run scripts before record is deleted
//
// This is implicitly called, no extra security check is needed
func (svc automationRunner) BeforeRecordDelete(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	return svc.findRecordScripts("beforeDelete", m.ID).Walk(
		svc.makeRecordScriptRunner(ctx, ns, m, r, false),
	)
}

// AfterRecordDelete - run scripts after record is deleted
//
// This is implicitly called, no extra security check is needed
func (svc automationRunner) AfterRecordDelete(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	return svc.findRecordScripts("afterDelete", m.ID).Walk(
		svc.makeRecordScriptRunner(ctx, ns, m, r, true),
	)
}

// Finds all scripts that are implicitly triggered by backend actions before/after
func (svc automationRunner) findRecordScripts(event string, moduleID uint64) automation.ScriptSet {
	ss, _ := svc.scriptFinder.FindRunnableScripts(AutomationResourceRecord, event, automation.MakeMatcherIDCondition(moduleID)).
		Filter(func(script *automation.Script) (bool, error) {
			// Filter out user-agent scripts
			return !script.RunInUA, nil
		})

	return ss
}

// UserScripts - collect all scripts runnable by users, appends compatible triggers
//
// So, either in their browser (RunInUA) or by running backend scripts explicitly (event:manual)
// All triggers are permission-checked for "run" operation.
//
func (svc automationRunner) UserScripts(ctx context.Context) automation.ScriptSet {
	var ss = automation.ScriptSet{}

	_ = svc.scriptFinder.FindRunnableScripts("", "").Walk(func(script *automation.Script) error {
		var tt = []*automation.Trigger{}

		for _, t := range script.Triggers() {
			if (script.RunInUA || t.Event == "manual") && svc.ac.CanRunAutomationTrigger(ctx, t) {
				// Making a copy so that we do not corrupt the
				tt = append(tt, &(*t))
			}
		}

		// Have any triggers left?
		if len(tt) > 0 {
			var sc = &automation.Script{}

			*sc = *script

			// Replace triggers with a new set
			sc.AddTrigger(automation.STMS_REPLACE, tt...)

			// andd append t
			ss = append(ss, sc)
		}

		return nil
	})

	return ss
}

// ManualRecordRun - Manual trigger run
//
// This is explicitly called, extra security  check is needed
func (svc automationRunner) RecordManual(ctx context.Context, scriptID uint64, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	// This scripts are all prechecked & filtered
	script := svc.UserScripts(ctx).FindByID(scriptID)

	if script == nil {
		return errors.New("can not find compatible script")
	}

	// Do not execute UA scripts
	if script.RunInUA {
		return errors.New("can not execute user-agent scripts")
	}

	// Make record script runner and
	runner := svc.makeRecordScriptRunner(ctx, ns, m, r, false)

	// Run it with a script
	//
	// Successfully executed record scripts can have an effect on given record value (r)
	return runner(script)
}

// RecordDeferred - Deferred trigger run
func (svc automationRunner) RecordDeferred(ctx context.Context, script *automation.Script, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	if script == nil {
		return errors.New("can not find compatible script")
	}

	// Do not execute scripts without a user
	if script.RunAs <= 0 {
		return errors.New("can not execute deferred scripts without a bound user")
	}

	// Do not execute UA scripts
	if script.RunInUA {
		return errors.New("can not execute user-agent scripts")
	}

	// Make record script runner and
	runner := svc.makeRecordScriptRunner(ctx, ns, m, r, true)

	// Run it with a script
	//
	// Successfully executed record scripts can have an effect on given record value (r)
	return runner(script)
}

func (svc automationRunner) RecordScriptTester(ctx context.Context, source string, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	// Make record script runner and
	runner := svc.makeRecordScriptRunner(ctx, ns, m, r, false)

	return runner(&automation.Script{
		ID:        0,
		Name:      "test",
		SourceRef: "test",
		Source:    source,
		Async:     false,
		RunAs:     0,
		RunInUA:   false,
		Timeout:   0,
		Critical:  true,
		Enabled:   false,
	})
}

// Runs record script
//
// We set-up script-running environment: security (definer / invoker), async, critical
// and copying values from the run to the given Record
//
func (svc automationRunner) makeRecordScriptRunner(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record, discard bool) func(script *automation.Script) error {
	// Static request params (record gets updated
	var req = &corredor.RunRecordRequest{
		Namespace: proto.FromNamespace(ns),
		Module:    proto.FromModule(m),
		Record:    proto.FromRecord(r),
	}

	return func(script *automation.Script) error {
		if svc.runner == nil {
			return errors.New("can not run corredor script: not connected")
		}

		// This could be executed in a goroutine (by *after triggers,
		// so we need to rewire the sentry panic recovery
		defer sentry.Recover()

		ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		// Add invoker's or defined credentials/jwt
		req.Config = map[string]string{
			"api.jwt": svc.getJWT(ctx, script),

			// Let the script know where the API is
			"api.baseURL.system":    svc.opt.ApiBaseURLSystem,
			"api.baseURL.compose":   svc.opt.ApiBaseURLCompose,
			"api.baseURL.messaging": svc.opt.ApiBaseURLMessaging,
		}

		// Add script info
		req.Script = corredor.FromScript(script)

		rsp, err := svc.runner.Record(ctx, req, grpc.WaitForReady(script.Critical))

		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				svc.logger.Error("unexpected error type", zap.Error(err))
				return err
			}

			switch s.Code() {
			case codes.FailedPrecondition:
				// Sent on syntax errors:
				err = errors.New(s.Message())
			case codes.Aborted:
				err = errors.New(s.Message())
			case codes.InvalidArgument:
				err = errors.New("invalid argument")
			case codes.Internal:
				err = errors.New("internal corredor error")
			default:
			}

			svc.logger.Info("script executed with errors", zap.Error(err))

			if !script.Critical {
				// This was not a critical call and we do not care about
				// errors from script running service.
				return nil
			}

			return err
		}

		if script.Async || discard {
			// Discard returned values (in case of async call or when forced)
			//
			// Backend is still returning, so we do not
			// need to handle multiple gRPC endpoints
			return nil
		}

		if rsp.Record == nil {
			// Script did not return any results
			// This means we should stop with the execution
			return errors.New("aborted")
		}

		// Convert from proto and copy record owner & values from the result
		result := proto.ToRecord(rsp.Record)
		r.OwnedBy, r.Values = result.OwnedBy, result.Values

		// Let's copy module, namespace and other values if they are missing
		if m != nil {
			if r.ModuleID == 0 {
				r.ModuleID = m.ID
				r.NamespaceID = m.NamespaceID
			}
		}

		// Who modified/created/owns the Record
		var currentUserID = auth.GetIdentityFromContext(ctx).Identity()
		currentUserID = script.RunAs

		if r.OwnedBy == 0 {
			r.OwnedBy = currentUserID
		}

		if r.CreatedAt.IsZero() {
			r.CreatedAt = time.Now()
		}

		if r.CreatedBy == 0 {
			r.CreatedBy = currentUserID
		}

		return nil
	}
}

// Creates a new JWT for
func (svc automationRunner) getJWT(ctx context.Context, script *automation.Script) string {
	if script.RunAsDefined() {
		// @todo implement this
		//       at the moment we do not he the ability fetch user info from non-system service
		//       extend/implement this feature when our services will know how to communicate with each-other
		return script.Credentials()
	}

	return svc.jwtEncoder.Encode(auth.GetIdentityFromContext(ctx))
}
