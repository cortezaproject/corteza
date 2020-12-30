package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
	"sync"
)

type (
	session struct {
		store      store.Storer
		actionlog  actionlog.Recorder
		ac         sessionAccessController
		log        *zap.Logger
		mux        *sync.RWMutex
		pool       map[uint64]*wfexec.Session
		spawnQueue chan *spawn
	}

	spawn struct {
		session chan *wfexec.Session
		graph   *wfexec.Graph
	}

	sessionAccessController interface {
		CanSearchSessions(context.Context) bool
		CanManageWorkflowSessions(context.Context, *types.Workflow) bool
	}
)

func Session(log *zap.Logger) *session {
	return &session{
		log:        log,
		actionlog:  DefaultActionlog,
		store:      DefaultStore,
		ac:         DefaultAccessControl,
		mux:        &sync.RWMutex{},
		pool:       make(map[uint64]*wfexec.Session),
		spawnQueue: make(chan *spawn),
	}
}

func (svc *session) Find(ctx context.Context, filter types.SessionFilter) (rr types.SessionSet, f types.SessionFilter, err error) {
	var (
		wap = &sessionActionProps{filter: &filter}
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

	return rr, filter, svc.recordAction(ctx, wap, SessionActionSearch, err)
}

func (svc *session) FindByID(ctx context.Context, sessionID uint64) (res *types.Session, err error) {
	var (
		wap = &sessionActionProps{session: &types.Session{ID: sessionID}}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if res, err = loadSession(ctx, s, sessionID); err != nil {
			return err
		}

		return nil
	})

	return res, svc.recordAction(ctx, wap, SessionActionLookup, err)
}

func (svc *session) resumeAll(ctx context.Context) error {
	return nil
}

func (svc *session) suspendAll(ctx context.Context) error {
	return nil
}

// Start new workflow session on a specific step with a given identity and scope
func (svc *session) Start(g *wfexec.Graph, stepID uint64, i auth.Identifiable, input types.Variables) error {
	svc.log.Debug("spawning")
	var (
		ctx   = auth.SetIdentityToContext(context.Background(), i)
		ses   = svc.spawn(g)
		start wfexec.Step
	)

	svc.log.Debug("spawned")
	if stepID == 0 {
		// starting step is not explicitly workflows on trigger, find orphan step
		switch oo := g.Orphans(); len(oo) {
		case 1:
			start = oo[0]
		case 0:
			return fmt.Errorf("could not find step without parents")
		default:
			return fmt.Errorf("multiple steps without parents")
		}
	} else if start = g.GetStepByIdentifier(stepID); start == nil {
		return fmt.Errorf("trigger staring step references nonexisting step")
	}

	svc.log.Debug("starting new workflow session", zap.Any("input", input))

	return ses.Exec(ctx, start, wfexec.Variables(input))
}

// Resume resumes suspended session/state
//
// Session can only be resumed by knowing session and state ID. Resume is an asynchronous operation
func (svc *session) Resume(sessionID, stateID uint64, i auth.Identifiable, input types.Variables) error {
	defer svc.mux.RUnlock()
	svc.mux.RLock()
	ses := svc.pool[sessionID]
	if ses == nil {
		return errors.NotFound("workflow session not found")
	}

	return ses.Resume(auth.SetIdentityToContext(context.Background(), i), stateID, wfexec.Variables(input))
}

// spawns a new session
//
// We need initial context for the session because we want to catch all cancellations or timeouts from there
// and not from any potential HTTP requests or similar temporary context that can prematurely destroy a workflow session
func (svc *session) spawn(g *wfexec.Graph) *wfexec.Session {
	s := &spawn{make(chan *wfexec.Session, 1), g}
	svc.spawnQueue <- s
	return <-s.session
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
				fresh := wfexec.NewSession(ctx, s.graph)
				svc.mux.Lock()
				svc.pool[fresh.ID()] = fresh
				svc.mux.Unlock()
				s.session <- fresh

				// case time for a pool cleanup
				// @todo cleanup pool when sessions are complete
			}
		}

		// @todo serialize sessions & suspended states

	}()

	svc.log.Debug("watcher initialized")
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
