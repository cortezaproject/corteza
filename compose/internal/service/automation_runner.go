package service

import (
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cortezaproject/corteza-server/compose/proto"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	automationRunner struct {
		logger       *zap.Logger
		runner       proto.ScriptRunnerClient
		scriptFinder automationScriptsFinder
		jwtEncoder   auth.TokenEncoder
	}

	automationScriptsFinder interface {
		Watch(ctx context.Context)
		FindRunnableScripts(event, resource string, cc ...automation.TriggerConditionChecker) automation.ScriptSet
	}
)

func AutomationRunner(f automationScriptsFinder, r proto.ScriptRunnerClient) automationRunner {
	var svc = automationRunner{
		scriptFinder: f,
		runner:       r,

		logger:     DefaultLogger.Named("automationRunner"),
		jwtEncoder: auth.DefaultJwtHandler,
	}

	return svc
}

func (svc automationRunner) findRecordScripts(event string, moduleID uint64) (ss automation.ScriptSet) {
	const resource = "compose:record"

	// We'll be comparing strings, not uint64!
	var moduleIDs = strconv.FormatUint(moduleID, 10)

	return svc.scriptFinder.FindRunnableScripts(event, resource,
		// ModuleID MUST match
		func(cModuleID string) bool {
			return moduleIDs == cModuleID
		},
	)
}

func (svc automationRunner) Watch(ctx context.Context) {
	svc.scriptFinder.Watch(ctx)
}

// ManualRecordRun - Manual trigger run
//
// This is explicitly called, extra security  check is needed
func (svc automationRunner) ManualRecordRun(ctx context.Context, scriptID uint64, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	// @todo security check (can user run this script (scriptID) manually)

	runner := svc.makeRecordScriptRunner(ctx, ns, m, r, true)

	return svc.findRecordScripts("manual", m.ID).Walk(func(script *automation.Script) error {
		// Interested in a specific script, so skip everything else
		if script.ID != scriptID {
			return nil
		}

		return runner(script)
	})
}

// BeforeRecordCreate - run scripts before record is created
//
// This is implicitly called, no extra security check is needed
func (svc automationRunner) BeforeRecordCreate(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) (err error) {
	return svc.findRecordScripts("beforeCreate", m.ID).Walk(
		svc.makeRecordScriptRunner(ctx, ns, m, r, true),
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
		svc.makeRecordScriptRunner(ctx, ns, m, r, true),
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
		svc.makeRecordScriptRunner(ctx, ns, m, r, true),
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

// Runs record script
//
// We set-up script-running environment: security (definer / invoker), async, critical
// and copying values from the run to the given Record
//
func (svc automationRunner) makeRecordScriptRunner(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record, discard bool) func(script *automation.Script) error {
	// Static request params (record gets updated
	var req = &proto.RunRecordRequest{
		Namespace: proto.FromNamespace(ns),
		Module:    proto.FromModule(m),
		Record:    proto.FromRecord(r),
	}

	svc.logger.Debug("executing script", zap.Any("record", r))

	return func(script *automation.Script) error {
		// This could be executed in a goroutine (by *after triggers,
		// so we need ot rewire the sentry panic recoverty
		defer sentry.Recover()

		ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		// Add invoker's or defined credentials/jwt
		req.JWT = svc.getJWT(ctx, script.RunAs)

		// Add script info
		req.Script = proto.FromAutomationScript(script)

		rsp, err := svc.runner.Record(ctx, req, grpc.WaitForReady(script.Critical))

		svc.logger.Debug("call sent")

		if err != nil {
			// @todo aborted?
			svc.logger.Debug("script executed, did not return record", zap.Error(err))
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
			svc.logger.Debug("script executed / async")
			return nil
		}

		if rsp.Record == nil {
			// Script did not return any results
			// This means we should stop with the execution
			// @todo aborted
			return nil
		}

		svc.logger.Debug("script executed", zap.Any("record", rsp.Record))

		// Convert from proto and copy record owner & values from the result
		result := proto.ToRecord(rsp.Record)
		r.OwnedBy, r.Values = result.OwnedBy, result.Values

		return nil
	}
}

// Creates a new JWT for
func (svc automationRunner) getJWT(ctx context.Context, userID uint64) string {
	if userID > 0 {
		// @todo implement this
		//       at the moment we do not he the ability fetch user info from non-system service
		//       extend/implement this feature when our services will know how to communicate with each-other
	}

	return svc.jwtEncoder.Encode(auth.GetIdentityFromContext(ctx))
}
