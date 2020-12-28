package wfexec

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"go.uber.org/zap"
	"sync"
	"time"
)

type (
	Session struct {
		// Session identifier
		id uint64

		// workflow ref, will help us with traversing through Graph
		workflow *Graph

		started time.Time

		// state channel (ie work queue)
		qState chan *state

		// error channel
		qErr chan error

		// locks concurrent executions
		execLock chan struct{}

		// collection of all suspended states
		// map key represents state identifier
		suspended map[uint64]*suspended

		// how often we check for suspended states and how often idle stat is checked in Wait()
		workerInterval time.Duration

		// holds final result
		result Variables
		err    error

		mux *sync.RWMutex

		// debug logger
		log *zap.SugaredLogger
	}

	suspended struct {
		// when not nil, assuming delayed
		resumeAt *time.Time

		// when true, assuming waiting for input (resumable through Resume())
		input bool

		// state to be resumed
		state *state
	}

	sessionOpt func(*Session)

	// state holds information about Session ID
	state struct {
		created time.Time

		// state identifier
		stateId uint64

		// Session identifier
		sessionId uint64

		// caller, parent step
		caller Step

		// current step
		step Step

		// scope
		scope Variables
	}

	// ExecRequest is passed to Exec() functions and contains all information to
	// resume suspended states in a Graph Session
	ExecRequest struct {
		SessionID uint64
		StateID   uint64

		// Current scope (group of variables)
		Scope Variables

		// Helps with gateway join/merge steps
		// that needs info about the step it's currently merging
		Caller Step
	}
)

const (
	sessionStateChanBuf = 512
	sessionMaxExec      = 32
)

var (
	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func NewSession(ctx context.Context, wf *Graph, oo ...sessionOpt) *Session {
	s := &Session{
		workflow:  wf,
		id:        nextID(),
		started:   time.Now(),
		qState:    make(chan *state, sessionStateChanBuf),
		qErr:      make(chan error, 1),
		execLock:  make(chan struct{}, sessionMaxExec),
		suspended: make(map[uint64]*suspended),

		workerInterval: time.Second,

		mux: &sync.RWMutex{},

		log: zap.NewNop().Sugar(),
	}

	for _, o := range oo {
		o(s)
	}

	go s.worker(ctx)

	return s
}

func (s Session) ID() uint64 { return s.id }

func (s Session) Idle() bool {
	if len(s.execLock) > 0 || len(s.qState) > 0 || len(s.qErr) > 0 {
		s.log.Debugf("Session(%d).Idle() => pending work: execLock: %d / qState: %d / qErr: %d\n", s.id, len(s.execLock), len(s.qState), len(s.qErr))
		return false
	}

	return true
}

func (s *Session) FinalError() error {
	defer s.mux.RUnlock()
	s.mux.RLock()

	return s.err
}

func (s *Session) Result() Variables {
	defer s.mux.RUnlock()
	s.mux.RLock()

	return s.result
}

func (s *Session) Exec(ctx context.Context, step Step, scope Variables) error {
	if len(s.workflow.Parents(step)) > 0 {
		return fmt.Errorf("can not execute step with parents")
	}

	return s.enqueue(ctx, State(s, nil, step, scope))
}

func (s *Session) Resume(ctx context.Context, stateId uint64, scope Variables) error {
	defer s.mux.Unlock()
	s.mux.Lock()

	suspended, has := s.suspended[stateId]
	if !has {
		return fmt.Errorf("unexisting state")
	}

	if !suspended.input {
		return fmt.Errorf("not input state")
	}

	delete(s.suspended, stateId)

	suspended.state.scope = suspended.state.scope.Merge(scope)

	return s.enqueue(ctx, suspended.state)
}

func (s *Session) enqueue(ctx context.Context, st *state) error {
	if st == nil {
		return fmt.Errorf("state is nil")
	}

	if st.step == nil {
		return fmt.Errorf("state step is nil")
	}

	if st.stateId == 0 {
		st.stateId = nextID()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case s.qState <- st:
		s.log.Debugf("Session(%d).Run() => added step to qState\n", s.id)
		return nil
	}
}

// does not wait for the whole wf to be complete but until:
//  - context timeout
//  - idle state
//  - error in error queue
func (s *Session) Wait(ctx context.Context) {
	waitCheck := time.NewTicker(s.workerInterval)
	defer waitCheck.Stop()

	for {
		select {
		case <-waitCheck.C:
			if s.Idle() {
				s.log.Debugf("Session(%d).Wait() => complete\n", s.id)
				// nothing in the pipeline
				return
			}

		case <-ctx.Done():
			s.log.Debugf("Session(%d).Wait() => ctx.Done() (err: %s)\n", s.id, ctx.Err())
			return

		case err := <-s.qErr:
			if err == nil {
				// execution complete
				s.log.Debugf("Session(%d).Wait() => done\n", s.id)
				return
			}

			if err != nil {
				defer s.mux.Unlock()
				s.mux.Lock()
				s.err = err
			}

			s.log.Debugf("Session(%d).Wait() => got error (by execution) (err: %s)\n", s.id, s.err)
		}
	}
}

func (s *Session) worker(ctx context.Context) {
	defer s.Close()

	suspCheck := time.NewTicker(s.workerInterval)
	defer suspCheck.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Debugf("Session(%d).worker() => ctx.Done(): %v\n", s.id, ctx.Err())
			return

		case <-suspCheck.C:
			s.log.Debugf("Session(%d).worker() => checking for scheduled suspended\n", s.id)
			s.queueScheduledSuspended()

		case st := <-s.qState:
			if st.step == nil {
				s.log.Debugf("Session(%d).worker() => got step==nil state; stopping & setting result!\n", s.id)
				defer s.mux.Unlock()
				s.mux.Lock()

				s.result = st.scope
				return
			}

			s.execLock <- struct{}{}
			s.log.Debugf("Session(%d).worker() => got state [stateId:%d]; execute!\n", s.id, st.stateId)
			go func() {
				s.exec(ctx, st)
				<-s.execLock
			}()

		}
	}
}

func (s *Session) Close() {
	s.log.Debugf("Session(%d).Close()\n", s.id)
	close(s.qErr)
	close(s.qState)
	close(s.execLock)
}

func (s Session) Suspended() bool {
	defer s.mux.RUnlock()
	s.mux.RLock()
	//return false
	return len(s.suspended) > 0
}

func (s *Session) queueScheduledSuspended() {
	defer s.mux.Unlock()
	s.mux.Lock()

	for id, sus := range s.suspended {
		if sus.resumeAt == nil {
			continue
		}

		if sus.resumeAt.After(time.Now()) {
			continue
		}

		s.log.Debugf("Session(%d).queueScheduledSuspended() => found suspended state [stateId=%d]\n", s.id, id)
		delete(s.suspended, id)
		s.qState <- sus.state
	}
}

func (s *Session) exec(ctx context.Context, st *state) {
	var (
		scope = st.scope
		next  Steps
	)

	{
		result, err := st.step.Exec(ctx, st.MakeRequest())
		if err != nil {
			s.qErr <- err
			return
		}

		switch result := result.(type) {
		case *Joined:
			s.log.Debugf("Session(%d).exec(%d) => joined\n", s.id, st.stateId)
			return

		case *suspended:
			if result == nil {
				s.log.Debugf("Session(%d).exec(%d) => suspended with nil\n", s.id, st.stateId)
				return
			}

			s.log.Debugf("Session(%d).exec(%d) => suspending step: %v\n", s.id, st.stateId, result)
			result.state = st
			defer s.mux.Unlock()
			s.mux.Lock()
			s.suspended[st.stateId] = result
			return

		case Variables:
			scope = scope.Merge(result)

		case Step:
			next = Steps{result}

		case Steps:
			next = result

		default:
			s.log.Debugf("Session(%d).exec(%d) => unknown exec response type: %T\n", s.id, st.stateId, result)

		}
	}

	if len(next) == 0 {
		next = s.workflow.Children(st.step)
	} else if !s.workflow.Children(st.step).Contains(next...) {
		s.qErr <- fmt.Errorf("inconsistent relationship")
		return
	}

	if len(next) == 0 {
		s.log.Debugf("Session(%d).exec(%d) => zero paths, completing\n", s.id, st.stateId)
		// using state to transport results and complete the worker loop
		s.qState <- FinalState(s, scope)
		return
	}

	s.log.Debugf("Session(%d).exec(%d) => %d paths\n", s.id, st.stateId, len(next))
	for _, p := range next {
		s.log.Debugf("Session(%d).exec(%d) => queuing step\n", s.id, st.stateId)
		_ = s.enqueue(ctx, State(s, st.step, p, scope))
	}

}

func SetWorkerInterval(i time.Duration) sessionOpt {
	return func(s *Session) {
		s.workerInterval = i
	}
}

func State(ses *Session, caller, current Step, scope Variables) *state {
	return &state{
		created:   time.Now(),
		stateId:   ses.id,
		sessionId: nextID(),
		caller:    caller,
		step:      current,
		scope:     scope,
	}
}

func FinalState(ses *Session, scope Variables) *state {
	return &state{
		created:   time.Now(),
		stateId:   ses.id,
		sessionId: nextID(),
		scope:     scope,
	}
}

func (s *state) MakeRequest() *ExecRequest {
	return &ExecRequest{
		SessionID: s.sessionId,
		StateID:   s.stateId,
		Scope:     s.scope,
		Caller:    s.caller,
	}
}

func (ss Steps) hash() map[Step]bool {
	out := make(map[Step]bool)
	for _, s := range ss {
		out[s] = true
	}

	return out
}

func (ss Steps) Contains(steps ...Step) bool {
	hash := ss.hash()
	for _, s1 := range steps {
		if !hash[s1] {
			return false
		}
	}

	return true
}
