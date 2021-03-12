package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
	"sync"
	"time"
)

type (
	session struct {
		store      store.Storer
		actionlog  actionlog.Recorder
		ac         sessionAccessController
		log        *zap.Logger
		mux        *sync.RWMutex
		pool       map[uint64]*types.Session
		spawnQueue chan *spawn
	}

	spawn struct {
		session chan *wfexec.Session
		graph   *wfexec.Graph
		trace   bool
	}

	sessionAccessController interface {
		CanSearchSessions(context.Context) bool
		CanManageWorkflowSessions(context.Context, *types.Workflow) bool
	}

	WaitFn func(ctx context.Context) (*expr.Vars, wfexec.SessionStatus, error)
)

func Session(log *zap.Logger) *session {
	return &session{
		log:        log,
		actionlog:  DefaultActionlog,
		store:      DefaultStore,
		ac:         DefaultAccessControl,
		mux:        &sync.RWMutex{},
		pool:       make(map[uint64]*types.Session),
		spawnQueue: make(chan *spawn),
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

	return rr, filter, svc.recordAction(ctx, sap, SessionActionSearch, err)
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

		if !svc.ac.CanManageWorkflowSessions(ctx, wf) {
			return SessionErrNotAllowedToManage()
		}

		return nil
	})

	return res, svc.recordAction(ctx, sap, SessionActionLookup, err)
}

func (svc *session) resumeAll(ctx context.Context) error {
	// @todo resume active sessions from storage
	return nil
}

func (svc *session) suspendAll(ctx context.Context) error {
	// @todo suspend active sessions to storage
	return nil
}

type (
	foo struct {
		SessionID uint64
		StateID   uint64
		CreatedAt time.Time
		Payload   *expr.Vars
	}
)

// PendingPrompts returns all prompts on all sessions owned by current user
func (svc *session) PendingPrompts(ctx context.Context) (pp []*wfexec.PendingPrompt) {
	var (
		i = auth.GetIdentityFromContext(ctx)
	)

	if i == nil {
		return
	}

	defer svc.mux.RUnlock()
	svc.mux.RLock()

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
// It does not check user's permissions to execute workflow(s) so it should be used only when !
func (svc *session) Start(g *wfexec.Graph, i auth.Identifiable, ssp types.SessionStartParams) (wait WaitFn, err error) {
	var (
		start wfexec.Step
	)

	if ssp.StepID == 0 {
		// starting step is not explicitly workflows on trigger, find orphan step
		switch oo := g.Orphans(); len(oo) {
		case 1:
			start = oo[0]
		case 0:
			return nil, errors.InvalidData("could not find starting step")
		default:
			return nil, errors.InvalidData("can not start workflow session multiple starting steps found")
		}
	} else if start = g.StepByID(ssp.StepID); start == nil {
		return nil, errors.InvalidData("trigger staring step references nonexisting step")
	}

	var (
		ctx = auth.SetIdentityToContext(context.Background(), i)
		ses = svc.spawn(g, ssp.Trace)
	)

	ses.CreatedAt = *now()
	ses.CreatedBy = i.Identity()
	ses.Apply(ssp)

	if err = store.CreateAutomationSession(context.TODO(), svc.store, ses); err != nil {
		return
	}

	if err = ses.Exec(ctx, start, ssp.Input); err != nil {
		return
	}

	return func(ctx context.Context) (*expr.Vars, wfexec.SessionStatus, error) { return ses.WaitResults(ctx) }, nil
}

// Resume resumes suspended session/state
//
// Session can only be resumed by knowing session and state ID. Resume is an asynchronous operation
func (svc *session) Resume(sessionID, stateID uint64, i auth.Identifiable, input *expr.Vars) error {
	var (
		ctx = auth.SetIdentityToContext(context.Background(), i)
	)

	defer svc.mux.RUnlock()
	svc.mux.RLock()
	ses := svc.pool[sessionID]
	if ses == nil {
		return errors.NotFound("session not found")
	}

	return ses.Resume(ctx, stateID, input)
}

// spawns a new session
//
// We need initial context for the session because we want to catch all cancellations or timeouts from there
// and not from any potential HTTP requests or similar temporary context that can prematurely destroy a workflow session
func (svc *session) spawn(g *wfexec.Graph, trace bool) (ses *types.Session) {
	s := &spawn{make(chan *wfexec.Session, 1), g, trace}

	// Send new-session request
	svc.spawnQueue <- s

	// blocks until session is set
	ses = types.NewSession(<-s.session)

	svc.mux.Lock()
	svc.pool[ses.ID] = ses
	svc.mux.Unlock()
	return ses
}

func (svc *session) Watch(ctx context.Context) {
	go func() {
		defer sentry.Recover()
		defer svc.log.Info("stopped")

		for {
			select {
			case <-ctx.Done():
				return
			case s := <-svc.spawnQueue:
				//
				s.session <- wfexec.NewSession(ctx,
					s.graph,
					wfexec.SetHandler(svc.stateChangeHandler(ctx)),
					wfexec.SetLogger(svc.log),
				)
				// case time for a pool cleanup
				// @todo cleanup pool when sessions are complete
			}
		}

		// @todo serialize sessions & suspended states
		//svc.suspendAll(ctx)
	}()

	svc.log.Debug("watcher initialized")
}

func (svc *session) stateChangeHandler(ctx context.Context) wfexec.StateChangeHandler {
	return func(i wfexec.SessionStatus, state *wfexec.State, s *wfexec.Session) {
		log := svc.log.With(zap.Uint64("sessionID", s.ID()))

		defer svc.mux.RUnlock()
		svc.mux.RLock()
		ses := svc.pool[s.ID()]
		if ses == nil {
			log.Warn("could not find session to update")
			return
		}

		var (
			update = true
			frame  = state.MakeFrame()
		)

		if ses.Stacktrace != nil {
			// Stacktrace will be set to !nil if frame collection is needed
			ses.Stacktrace = append(ses.Stacktrace, frame)
		}

		switch i {
		case wfexec.SessionPrompted:
			ses.SuspendedAt = now()
			ses.Status = types.SessionPrompted

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
			// force update on every 10 new frames but only when stacktrace is not nil
			update = ses.Stacktrace != nil && len(ses.Stacktrace)%10 == 0
		}

		if !update {
			return
		}

		if err := svc.store.UpdateAutomationSession(ctx, ses); err != nil {
			log.Error("failed to update session", zap.Error(err))
		} else {
			log.Debug("session updated", zap.Stringer("status", ses.Status))
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
