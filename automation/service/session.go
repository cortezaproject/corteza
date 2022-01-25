package service

import (
	"context"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

type (
	promptSender interface {
		Send(kind string, payload interface{}, userIDs ...uint64) error
	}

	session struct {
		store        store.Storer
		actionlog    actionlog.Recorder
		ac           sessionAccessController
		opt          options.WorkflowOpt
		log          *zap.Logger
		mux          sync.RWMutex
		pool         map[uint64]*types.Session
		spawnQueue   chan *spawn
		promptSender promptSender
	}

	spawn struct {
		workflowID uint64
		session    chan *wfexec.Session
		graph      *wfexec.Graph
		invoker    auth.Identifiable
		runner     auth.Identifiable
		trace      bool
		callStack  []uint64
	}

	sessionAccessController interface {
		CanSearchSessions(context.Context) bool
		CanManageSessionsOnWorkflow(context.Context, *types.Workflow) bool
	}

	WaitFn func(ctx context.Context) (*expr.Vars, wfexec.SessionStatus, types.Stacktrace, error)
)

const (
	// when the state changes, state-change-handler is called and for non-fatal,
	// non-interactive or non-delay steps (that are much more frequent) we need
	// to limit how often the store is updated with the updated session info
	//
	// We use the size of the stacktrace and for every F (see the value of the constant)
	// we flush the session info to the store.
	sessionStateFlushFrequency = 1000
)

func Session(log *zap.Logger, opt options.WorkflowOpt, ps promptSender) *session {
	return &session{
		log:          log,
		opt:          opt,
		actionlog:    DefaultActionlog,
		store:        DefaultStore,
		ac:           DefaultAccessControl,
		pool:         make(map[uint64]*types.Session),
		spawnQueue:   make(chan *spawn),
		promptSender: ps,
	}
}

func (svc *session) Search(ctx context.Context, filter types.SessionFilter) (rr types.SessionSet, f types.SessionFilter, err error) {
	var (
		sap = &sessionActionProps{filter: &filter}
	)

	err = func() (err error) {
		if !svc.ac.CanSearchSessions(ctx) {
			return SessionErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchAutomationSessions(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, sap, SessionActionSearch, err)
}

func (svc *session) LookupByID(ctx context.Context, sessionID uint64) (res *types.Session, err error) {
	var (
		sap = &sessionActionProps{session: &types.Session{ID: sessionID}}
		wf  *types.Workflow
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if !svc.ac.CanSearchSessions(ctx) {
			return SessionErrNotAllowedToRead()
		}

		if res, err = loadSession(ctx, s, sessionID); err != nil {
			return err
		}

		if wf, err = loadWorkflow(ctx, s, res.WorkflowID); err != nil {
			return err
		}

		if !svc.ac.CanManageSessionsOnWorkflow(ctx, wf) {
			return SessionErrNotAllowedToManage()
		}

		return nil
	})

	return res, svc.recordAction(ctx, sap, SessionActionLookup, err)
}

func (svc *session) resumeAll(ctx context.Context) error {
	// In theory we could resume active/pending/prompt sessions from persistent store
	// so that they can survive server termination

	// @todo resume active sessions from storage
	//       load all active sessions from store and load them into the pool
	//
	return nil
}

func (svc *session) suspendAll(ctx context.Context) error {
	// In theory we could suspend active/pending/prompt sessions to persistent store
	// so that they can survive server termination

	// @todo suspend active sessions to storage:
	//       stop watcher queue
	//       run gc
	//       stop worker on each session
	//       flush session to store (like we're doing in the status handler
	return nil
}

// PendingPrompts returns all prompts on all sessions owned by current user
func (svc *session) PendingPrompts(ctx context.Context) (pp []*wfexec.PendingPrompt) {
	var (
		i = auth.GetIdentityFromContext(ctx)
	)

	if i == nil {
		return
	}

	svc.mux.RLock()
	defer svc.mux.RUnlock()

	pp = make([]*wfexec.PendingPrompt, 0, len(svc.pool))
	for _, s := range svc.pool {
		pp = append(pp, s.PendingPrompts(i.Identity())...)
	}

	return
}

// Start new workflow session on a specific step with a given identity and scope
//
// Start is an asynchronous operation
//
// Please note that context passed to the function is NOT the the one that is
// used for the execution of the workflow. See watch function!
//
// It does not check user's permissions to execute workflow(s) so it should be used only when !
func (svc *session) Start(ctx context.Context, g *wfexec.Graph, ssp types.SessionStartParams) (wait WaitFn, err error) {
	var (
		start wfexec.Step
	)

	if g == nil {
		return nil, errors.InvalidData("cannot start workflow, uninitialized graph")
	}

	if len(ssp.CallStack) > svc.opt.CallStackSize {
		return nil, WorkflowErrMaximumCallStackSizeExceeded()
	}

	ssp.CallStack = append(ssp.CallStack, ssp.WorkflowID)

	if ssp.Invoker == nil {
		return nil, errors.InvalidData("cannot start workflow without user")
	}

	if ssp.Runner == nil {
		ssp.Runner = ssp.Invoker
	}

	if ssp.StepID == 0 {
		// starting step is not explicitly set
		// find orphan step
		switch oo := g.Orphans(); len(oo) {
		case 1:
			start = oo[0]
		case 0:
			return nil, errors.InvalidData("could not find starting step")
		default:
			return nil, errors.InvalidData("cannot start workflow session multiple starting steps found")
		}
	} else if start = g.StepByID(ssp.StepID); start == nil {
		return nil, errors.InvalidData("trigger staring step references non-existing step")
	} else if len(g.Parents(g.StepByID(ssp.StepID))) > 0 {
		return nil, errors.InvalidData("cannot start workflow on a step with parents")
	}

	var (
		ses = svc.spawn(g, ssp.WorkflowID, ssp.Trace, ssp.CallStack, ssp.Runner, ssp.Invoker)
	)

	ses.CreatedAt = *now()
	ses.CreatedBy = ssp.Invoker.Identity()
	ses.Status = types.SessionStarted
	ses.Apply(ssp)

	_ = ssp.Input.AssignFieldValue("eventType", expr.Must(expr.NewString(ssp.EventType)))
	_ = ssp.Input.AssignFieldValue("resourceType", expr.Must(expr.NewString(ssp.ResourceType)))
	_ = ssp.Input.AssignFieldValue("invoker", expr.Must(expr.NewAny(ssp.Invoker)))
	_ = ssp.Input.AssignFieldValue("runner", expr.Must(expr.NewAny(ssp.Runner)))

	if err = ses.Exec(ctx, start, ssp.Input); err != nil {
		return
	}

	return func(ctx context.Context) (*expr.Vars, wfexec.SessionStatus, types.Stacktrace, error) {
		return ses.WaitResults(ctx)
	}, nil
}

// Resume resumes suspended session/state
//
// Session can only be resumed by knowing session and state ID. Resume is an asynchronous operation
func (svc *session) Resume(sessionID, stateID uint64, i auth.Identifiable, input *expr.Vars) error {
	var (
		ctx = auth.SetIdentityToContext(context.Background(), i)
	)

	svc.mux.RLock()
	defer svc.mux.RUnlock()
	ses := svc.pool[sessionID]
	if ses == nil {
		return errors.NotFound("session not found")
	}

	resPrompt, err := ses.Resume(ctx, stateID, input)
	if err != nil {
		return err
	}

	if err = svc.promptSender.Send("workflowSessionResumed", resPrompt, resPrompt.OwnerId); err != nil {
		svc.log.Error("failed to send prompt resume status to user", zap.Error(err))
	}

	return nil
}

// spawns a new session
//
// We need initial context for the session because we want to catch all cancellations or timeouts from there
// and not from any potential HTTP requests or similar temporary context that can prematurely destroy a workflow session
func (svc *session) spawn(g *wfexec.Graph, workflowID uint64, trace bool, callStack []uint64, runner, invoker auth.Identifiable) (ses *types.Session) {
	s := &spawn{
		workflowID: workflowID,
		session:    make(chan *wfexec.Session, 1),
		graph:      g,
		trace:      trace,
		callStack:  callStack,
		invoker:    invoker,
		runner:     runner,
	}

	// Send new-session request
	svc.spawnQueue <- s

	// blocks until session is set
	ses = types.NewSession(<-s.session)

	svc.mux.Lock()
	svc.pool[ses.ID] = ses
	svc.mux.Unlock()
	return ses
}

// Watch looks over session's spawn queue
func (svc *session) Watch(ctx context.Context) {
	gcTicker := time.NewTicker(time.Second)
	lpTicker := time.NewTicker(time.Second * 30)

	go func() {
		defer sentry.Recover()
		defer gcTicker.Stop()
		defer svc.log.Info("stopped")

		for {
			select {
			case <-ctx.Done():
				return
			case s := <-svc.spawnQueue:
				var execCtx = context.Background()

				opts := []wfexec.SessionOpt{
					wfexec.SetWorkflowID(s.workflowID),
					wfexec.SetCallStack(s.callStack...),
					wfexec.SetHandler(svc.stateChangeHandler(ctx)),
				}

				if svc.opt.ExecDebug {
					log := svc.log.
						Named("exec").
						With(zap.Uint64("workflowID", s.workflowID)).
						With(zap.Uint64("runnerID", s.runner.Identity())).
						With(zap.Uint64s("runnerRoles", s.runner.Roles()))

					opts = append(
						opts,
						wfexec.SetLogger(log),
						wfexec.SetDumpStacktraceOnPanic(true),
					)
				}

				// Encode runner into execution context
				// runner is used as identity and for access control
				execCtx = auth.SetIdentityToContext(execCtx, s.runner)

				// Encode invoker into execution context
				// invoker is used
				execCtx = context.WithValue(execCtx, workflowInvokerCtxKey{}, s.invoker)

				s.session <- wfexec.NewSession(execCtx, s.graph, opts...)
				// case time for a pool cleanup
				// @todo cleanup pool when sessions are complete

			case <-gcTicker.C:
				svc.gc()

			case <-lpTicker.C:
				svc.logPending()
			}
		}

		// @todo serialize sessions & suspended states
		//svc.suspendAll(ctx)
	}()

	svc.log.Debug("watcher initialized")
}

// garbage collection for stale sessions
func (svc *session) gc() {
	svc.mux.Lock()
	defer svc.mux.Unlock()

	for _, s := range svc.pool {
		if s.GC() {
			delete(svc.pool, s.ID)
		}
	}
}

// garbage collection for stale sessions
func (svc *session) logPending() {
	svc.mux.RLock()
	defer svc.mux.RUnlock()

	var (
		total                                    = len(svc.pool)
		pending, pending1m, pending1h, pending1d int
	)

	for _, s := range svc.pool {
		switch {
		case s.CreatedAt.Sub(*now()) > time.Hour*24:
			pending1d++
		case s.CreatedAt.Sub(*now()) > time.Hour:
			pending1h++
		case s.CreatedAt.Sub(*now()) > time.Minute:
			pending1m++
		default:
			pending++
		}
	}

	if total > 0 {
		svc.log.Info(
			"workflow session garbage collector stats",
			zap.Int("total", total),
			zap.Int("pending", pending),
			zap.Int("pending1m", pending1m),
			zap.Int("pending1h", pending1h),
			zap.Int("pending1d", pending1d),
		)
	}
}

// stateChangeHandler keeps track of session status changes and frequently stores session into db
func (svc *session) stateChangeHandler(ctx context.Context) wfexec.StateChangeHandler {
	return func(i wfexec.SessionStatus, state *wfexec.State, s *wfexec.Session) {
		svc.mux.Lock()
		defer svc.mux.Unlock()

		log := svc.log.With(zap.Uint64("sessionID", s.ID()))

		ses := svc.pool[s.ID()]
		if ses == nil {
			log.Warn("could not find session to update")
			return
		}

		log = log.With(zap.Uint64("workflowID", ses.WorkflowID))

		var (
			// By default, we want to update session when new status is prompted, delayed, completed or failed
			// But if status is active, we'll flush it every X frames (sessionStateFlushFrequency)
			update = true

			frame = state.MakeFrame()
		)

		// Stacktrace will be set to !nil if frame collection is needed
		if len(ses.RuntimeStacktrace) > 0 {
			// calculate how long it took to get to this step
			frame.ElapsedTime = uint(frame.CreatedAt.Sub(ses.RuntimeStacktrace[0].CreatedAt) / time.Millisecond)
		}

		ses.AppendRuntimeStacktrace(frame)

		switch i {
		case wfexec.SessionPrompted:
			ses.SuspendedAt = now()
			ses.Status = types.SessionPrompted

			// Send the pending prompts to user
			if svc.promptSender != nil {
				for _, pp := range s.AllPendingPrompts() {
					if err := svc.promptSender.Send("workflowSessionPrompt", pp, pp.OwnerId); err != nil {
						svc.log.Error("failed to send prompt to user", zap.Error(err))
					}
				}
			}

		case wfexec.SessionDelayed:
			ses.SuspendedAt = now()
			ses.Status = types.SessionSuspended

		case wfexec.SessionCompleted:
			ses.SuspendedAt = nil
			ses.CompletedAt = now()
			ses.Status = types.SessionCompleted

		case wfexec.SessionFailed:
			ses.SuspendedAt = nil
			ses.CompletedAt = now()
			ses.Error = state.Error()
			ses.Status = types.SessionFailed

		default:
			// force update on every F new frames (F=sessionStateFlushFrequency) but only when stacktrace is not nil
			update = ses.RuntimeStacktrace != nil && len(ses.RuntimeStacktrace)%sessionStateFlushFrequency == 0
		}

		if !update {
			return
		}

		ses.CopyRuntimeStacktrace()

		if err := svc.store.UpsertAutomationSession(ctx, ses); err != nil {
			log.Error("failed to update session", zap.Error(err))
		}
	}
}

func loadSession(ctx context.Context, s store.Storer, sessionID uint64) (res *types.Session, err error) {
	if sessionID == 0 {
		return nil, SessionErrInvalidID()
	}

	if res, err = store.LookupAutomationSessionByID(ctx, s, sessionID); errors.IsNotFound(err) {
		return nil, SessionErrNotFound()
	}

	return
}
