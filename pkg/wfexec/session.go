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

		// steps graph
		g *Graph

		started time.Time

		// state channel (ie work queue)
		qState chan *State

		// error channel
		qErr chan error

		// locks concurrent executions
		execLock chan struct{}

		// collection of all suspended states
		// map key represents state identifier
		suspended map[uint64]*suspended

		messages []*message

		// how often we check for suspended states and how often idle stat is checked in Wait()
		workerInterval time.Duration

		// holds final result
		result Variables
		err    error

		mux *sync.RWMutex

		// debug logger
		log *zap.SugaredLogger

		eventHandler StateChangeHandler
	}

	StateChangeHandler func(int, *State, *Session)

	suspended struct {
		// when not nil, assuming delayed
		resumeAt *time.Time

		// when true, assuming waiting for input (resumable through Resume())
		input bool

		// state to be resumed
		state *State
	}

	sessionOpt func(*Session)

	// state holds information about Session ID
	State struct {
		created   time.Time
		completed *time.Time

		// state identifier
		stateId uint64

		// Session identifier
		sessionId uint64

		// caller, parent step
		caller Step

		// current step
		step Step

		// step error (if any)
		err error

		// input variables that were sent to resume the session
		// (prompt step)
		input Variables

		// scope
		scope Variables
	}

	// ExecRequest is passed to Exec() functions and contains all information to
	// resume suspended states in a Graph Session
	ExecRequest struct {
		SessionID uint64
		StateID   uint64

		// Current input received when prompting for input
		Input Variables

		// Current scope
		Scope Variables

		// Helps with gateway join/merge steps
		// that needs info about the step it's currently merging
		Caller Step
	}
)

const (
	sessionStateChanBuf   = 512
	sessionConcurrentExec = 32
)

const (
	SessionActive int = iota
	SessionSuspended
	SessionStepSuspended
	SessionCompleted
	SessionFailed
	SessionNewMessage
)

var (
	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}
)

func NewSession(ctx context.Context, wf *Graph, oo ...sessionOpt) *Session {
	s := &Session{
		g:              wf,
		id:             nextID(),
		started:        *now(),
		qState:         make(chan *State, sessionStateChanBuf),
		qErr:           make(chan error, 1),
		execLock:       make(chan struct{}, sessionConcurrentExec),
		suspended:      make(map[uint64]*suspended),
		messages:       make([]*message, 0),
		workerInterval: time.Second,

		mux: &sync.RWMutex{},

		log: zap.NewNop().Sugar(),
		eventHandler: func(int, *State, *Session) {
			// noop
		},
	}

	for _, o := range oo {
		o(s)
	}

	go s.worker(ctx)

	return s
}

func (s Session) Status() int {
	switch {
	case s.err != nil:
		return SessionFailed

	case len(s.execLock) > 0 || len(s.qState) > 0 || len(s.qErr) > 0:
		// active
		return SessionActive

	case len(s.suspended) > 0:
		return SessionSuspended

	default:
		return SessionCompleted
	}
}

func (s Session) ID() uint64 { return s.id }

func (s Session) Idle() bool {
	return s.Status() != SessionActive
}

func (s *Session) Error() error {
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
	if s.g.Len() == 0 {
		return fmt.Errorf("refusing to execute without steps")
	}

	if len(s.g.Parents(step)) > 0 {
		return fmt.Errorf("can not execute step with parents")
	}

	return s.enqueue(ctx, NewState(s, nil, step, scope))
}

func (s *Session) Resume(ctx context.Context, stateId uint64, input Variables) error {
	defer s.mux.Unlock()
	s.mux.Lock()

	resumed, has := s.suspended[stateId]
	if !has {
		return fmt.Errorf("unexisting state")
	}

	if !resumed.input {
		return fmt.Errorf("not input state")
	}

	delete(s.suspended, stateId)

	// setting received input to state
	resumed.state.input = input

	return s.enqueue(ctx, resumed.state)
}

func (s *Session) enqueue(ctx context.Context, st *State) error {
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
				s.log.Debugf("Session(%d).Wait() => idle\n", s.id)
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

			// add empty struct to chan to lock and to have control over numver of concurrent go processes
			// this will block if number of items in execLock chan reached value of sessionConcurrentExec
			s.execLock <- struct{}{}
			s.log.Debugf("Session(%d).worker() => got state [stateId:%d]; execute!\n", s.id, st.stateId)
			go func() {
				s.exec(ctx, st)
				st.completed = now()

				// remove single
				<-s.execLock

				// after exec lock is released call event handler with (new) session status
				s.eventHandler(s.Status(), st, s)
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
	return len(s.suspended) > 0
}

func (s Session) Messages() []*message {
	return s.messages
}

func (s Session) MessagesAfter(messageId uint64) []*message {
	var mm = make([]*message, 0, len(s.messages))
	for _, m := range mm {
		if m.ID > messageId {
			mm = append(mm, m)
		}
	}

	return mm
}

func (s *Session) queueScheduledSuspended() {
	defer s.mux.Unlock()
	s.mux.Lock()

	for id, sus := range s.suspended {
		if sus.resumeAt == nil {
			continue
		}

		if sus.resumeAt.After(*now()) {
			continue
		}

		s.log.Debugf("Session(%d).queueScheduledSuspended() => found suspended state [stateId=%d]\n", s.id, id)
		delete(s.suspended, id)
		s.qState <- sus.state
	}
}

// executes single step, resolves response and schedule following steps for execution
func (s *Session) exec(ctx context.Context, st *State) {
	var (
		result ExecResponse
		scope  = st.scope
		next   Steps
	)

	s.eventHandler(SessionActive, st, s)

	{
		result, st.err = st.step.Exec(ctx, st.MakeRequest())
		if st.err != nil {
			s.qErr <- st.err
			return
		}

		switch result := result.(type) {
		case Variables:
			// most common (successful) result
			// session will continue with configured child steps
			s.log.Debugf("Session(%d).exec(%d) => variables: %v\n", s.id, st.stateId, result)
			scope = scope.Merge(result)

		case *partial:
			// *partial is returned when step needs to be executed again
			// it's used mainly for join gateway step that should be called multiple times (one for each parent path)
			s.log.Debugf("Session(%d).exec(%d) => partial\n", s.id, st.stateId)
			return

		case *suspended:
			// suspend execution because of delay or pending user input
			// either way, it breaks execution loop for the current path

			if result == nil {
				// @todo why do we allow this again?
				s.log.Warnf("Session(%d).exec(%d) => suspended with nil\n", s.id, st.stateId)
				return
			}

			s.log.Debugf("Session(%d).exec(%d) => suspending step: %v\n", s.id, st.stateId, result)
			result.state = st
			s.mux.Lock()
			s.suspended[st.stateId] = result
			s.mux.Unlock()
			s.eventHandler(SessionStepSuspended, st, s)
			return

		case *message:
			// step emitted a message, store it and continue as planned
			s.log.Debugf("Session(%d).exec(%d) => message received: %v\n", s.id, st.stateId, result)
			s.mux.Lock()
			s.messages = append(s.messages, result)
			s.mux.Unlock()
			s.eventHandler(SessionNewMessage, st, s)

		case Steps:
			// session continues with set of specified steps
			// steps MUST be configured in a graph as step's children
			s.log.Debugf("Session(%d).exec(%d) => multiple steps: %v\n", s.id, st.stateId, result)
			next = result

		case Step:
			// session continues with a specified step
			// step MUST be configured in a graph as step's child
			s.log.Debugf("Session(%d).exec(%d) => single step: %v\n", s.id, st.stateId, result)
			next = Steps{result}

		default:
			s.log.Debugf("Session(%d).exec(%d) => unknown exec response type: %T\n", s.id, st.stateId, result)
			return
		}
	}

	if len(next) == 0 {
		next = s.g.Children(st.step)
	} else {
		cc := s.g.Children(st.step)
		if len(cc) > 0 && !cc.Contains(next...) {
			s.qErr <- fmt.Errorf("inconsistent relationship")
			return
		}
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
		_ = s.enqueue(ctx, NewState(s, st.step, p, scope))
	}

}

func SetWorkerInterval(i time.Duration) sessionOpt {
	return func(s *Session) {
		s.workerInterval = i
	}
}

func SetHandler(fn StateChangeHandler) sessionOpt {
	return func(s *Session) {
		s.eventHandler = fn
	}
}

func NewState(ses *Session, caller, current Step, scope Variables) *State {
	return &State{
		created:   *now(),
		stateId:   ses.id,
		sessionId: nextID(),
		caller:    caller,
		step:      current,
		scope:     scope,
	}
}

func FinalState(ses *Session, scope Variables) *State {
	return &State{
		created:   *now(),
		completed: now(),
		stateId:   ses.id,
		sessionId: nextID(),
		scope:     scope,
	}
}

func (s State) MakeRequest() *ExecRequest {
	return &ExecRequest{
		SessionID: s.sessionId,
		StateID:   s.stateId,
		Scope:     s.scope,
		Input:     s.input,
		Caller:    s.caller,
	}
}

func (s *State) Error() string {
	if s.err == nil {
		return ""
	}

	return s.err.Error()
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
