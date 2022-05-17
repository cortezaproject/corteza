package wfexec

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"go.uber.org/zap"
)

type (
	Session struct {
		// Session identifier
		id uint64

		workflowID uint64

		// steps graph
		g *Graph

		started time.Time

		// state channel (ie work queue)
		qState chan *State

		// crash channel
		qErr chan error

		// locks concurrent executions
		execLock chan struct{}

		// delayed states (waiting for the right time)
		delayed map[uint64]*delayed

		// prompted
		prompted map[uint64]*prompted

		// how often we check for delayed states and how often idle stat is checked in Wait()
		workerInterval time.Duration

		// only one worker routine per session
		workerLock chan struct{}

		statusChange chan int

		// holds final result
		result *expr.Vars
		err    error

		mux sync.RWMutex

		// debug logger
		log *zap.Logger

		dumpStacktraceOnPanic bool

		eventHandler StateChangeHandler

		callStack []uint64
	}

	StateChangeHandler func(SessionStatus, *State, *Session)

	SessionOpt func(*Session)

	Frame struct {
		CreatedAt time.Time  `json:"createdAt"`
		SessionID uint64     `json:"sessionID"`
		StateID   uint64     `json:"stateID"`
		Input     *expr.Vars `json:"input"`
		Scope     *expr.Vars `json:"scope"`
		Results   *expr.Vars `json:"results"`
		ParentID  uint64     `json:"parentID"`
		StepID    uint64     `json:"stepID"`
		NextSteps []uint64   `json:"nextSteps"`

		// How much time from the 1st step to the start of this step in milliseconds
		ElapsedTime uint `json:"elapsedTime"`

		// How much time it took to execute this step in milliseconds
		StepTime uint `json:"stepTime"`

		Action string `json:"action,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	// ExecRequest is passed to Exec() functions and contains all information
	// for step execution
	ExecRequest struct {
		SessionID uint64
		StateID   uint64

		// Current input received on session resume
		Input *expr.Vars

		// Current scope
		Scope *expr.Vars

		// Helps with gateway join/merge steps
		// that needs info about the step it's currently merging
		Parent Step
	}

	SessionStatus int

	callStackCtxKey struct{}
)

const (
	sessionStateChanBuf   = 512
	sessionConcurrentExec = 32
)

const (
	SessionActive SessionStatus = iota
	SessionPrompted
	SessionDelayed
	SessionFailed
	SessionCompleted
)

var (
	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now()
		return &c
	}
)

func (s SessionStatus) String() string {
	switch s {
	case SessionActive:
		return "active"
	case SessionPrompted:
		return "prompted"
	case SessionDelayed:
		return "delayed"
	case SessionFailed:
		return "failed"
	case SessionCompleted:
		return "completed"
	}

	return "UNKNOWN-SESSION-STATUS"
}

func NewSession(ctx context.Context, g *Graph, oo ...SessionOpt) *Session {
	s := &Session{
		g:        g,
		id:       nextID(),
		started:  *now(),
		qState:   make(chan *State, sessionStateChanBuf),
		qErr:     make(chan error, 1),
		execLock: make(chan struct{}, sessionConcurrentExec),
		delayed:  make(map[uint64]*delayed),
		prompted: make(map[uint64]*prompted),

		//workerInterval: time.Millisecond,
		workerInterval: time.Millisecond * 250, // debug mode rate
		workerLock:     make(chan struct{}, 1),

		log: zap.NewNop(),

		eventHandler: func(SessionStatus, *State, *Session) {
			// noop
		},
	}

	for _, o := range oo {
		o(s)
	}

	s.log = s.log.
		With(zap.Uint64("sessionID", s.id))

	s.callStack = append(s.callStack, s.id)

	go s.worker(ctx)

	return s
}

func (s *Session) Status() SessionStatus {
	s.mux.RLock()
	defer s.mux.RUnlock()

	switch {
	case s.err != nil:
		return SessionFailed

	case len(s.prompted) > 0:
		return SessionPrompted

	case len(s.delayed) > 0:
		return SessionDelayed

	case s.result == nil:
		return SessionActive

	default:
		return SessionCompleted
	}
}

func (s *Session) ID() uint64 {
	return s.id
}

func (s *Session) Idle() bool {
	return s.Status() != SessionActive
}

func (s *Session) Error() error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.err
}

func (s *Session) Result() *expr.Vars {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.result
}

func (s *Session) Exec(ctx context.Context, step Step, scope *expr.Vars) error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	err := func() error {
		if s.g.Len() == 0 {
			return fmt.Errorf("refusing to execute without steps")
		}

		if len(s.g.Parents(step)) > 0 {
			return fmt.Errorf("cannot execute step with parents")
		}
		return nil
	}()

	if err != nil {
		// send nil to error queue to trigger worker shutdown
		// session error must be set to update session status
		s.qErr <- err
		return err
	}

	if scope == nil {
		scope, _ = expr.NewVars(nil)
	}

	return s.enqueue(ctx, NewState(s, auth.GetIdentityFromContext(ctx), nil, step, scope))
}

// UserPendingPrompts prompts fn returns all owner's pending prompts on this session
func (s *Session) UserPendingPrompts(ownerId uint64) (out []*PendingPrompt) {
	if ownerId == 0 {
		return
	}

	defer s.mux.RUnlock()
	s.mux.RLock()

	out = make([]*PendingPrompt, 0, len(s.prompted))

	for _, p := range s.prompted {
		if p.ownerId != ownerId {
			continue
		}

		pending := p.toPending()
		pending.SessionID = s.id
		out = append(out, pending)
	}

	return
}

// AllPendingPrompts returns all pending prompts for all user
func (s *Session) AllPendingPrompts() (out []*PendingPrompt) {
	defer s.mux.RUnlock()
	s.mux.RLock()

	return s.pendingPrompts(s.prompted)
}

// UnsentPendingPrompts returns unsent pending prompts for all user
func (s *Session) UnsentPendingPrompts() (out []*PendingPrompt) {
	defer s.mux.RUnlock()
	s.mux.RLock()

	aux := s.pendingPrompts(s.prompted)
	for _, p := range aux {
		if p.Original.sent {
			continue
		}

		out = append(out, p)
	}

	return
}

func (s *Session) pendingPrompts(prompted map[uint64]*prompted) (out []*PendingPrompt) {
	out = make([]*PendingPrompt, 0, len(prompted))

	for _, p := range prompted {
		pending := p.toPending()
		pending.SessionID = s.id
		out = append(out, pending)
	}

	return
}

func (s *Session) Resume(ctx context.Context, stateId uint64, input *expr.Vars) (*ResumedPrompt, error) {
	defer s.mux.Unlock()
	s.mux.Lock()

	var (
		i      = auth.GetIdentityFromContext(ctx)
		p, has = s.prompted[stateId]
	)
	if !has {
		return nil, fmt.Errorf("unexisting state")
	}

	if i == nil || p.ownerId != i.Identity() {
		return nil, fmt.Errorf("state access denied")
	}

	delete(s.prompted, stateId)

	// setting received input to state
	p.state.input = input

	if err := s.enqueue(ctx, p.state); err != nil {
		return nil, err
	}

	return p.toResumed(), nil
}

func (s *Session) canEnqueue(st *State) error {
	if st == nil {
		return fmt.Errorf("state is nil")
	}

	// when the step is completed right away, it is considered as special
	if st.step == nil && st.completed == nil {
		return fmt.Errorf("state step is nil")
	}

	return nil
}

func (s *Session) enqueue(ctx context.Context, st *State) error {
	if err := s.canEnqueue(st); err != nil {
		return err
	}

	if st.stateId == 0 {
		st.stateId = nextID()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case s.qState <- st:
		s.log.Debug("add step to queue")
		return nil
	}
}

// Wait does not wait for the whole wf to be complete but until:
//  - context timeout
//  - idle state
//  - error in error queue
func (s *Session) Wait(ctx context.Context) error {
	return s.WaitUntil(ctx, SessionFailed, SessionDelayed, SessionCompleted)
}

// WaitUntil blocks until workflow session gets into expected status
//
func (s *Session) WaitUntil(ctx context.Context, expected ...SessionStatus) error {
	indexed := make(map[SessionStatus]bool)
	for _, status := range expected {
		indexed[status] = true
	}

	// already at the expected status
	if indexed[s.Status()] {
		return s.err
	}

	s.log.Debug(
		"waiting for status change",
		zap.Any("expecting", expected),
		zap.Duration("interval", s.workerInterval),
	)

	waitCheck := time.NewTicker(s.workerInterval)
	defer waitCheck.Stop()

	for {
		select {
		case <-waitCheck.C:
			if indexed[s.Status()] {
				s.log.Debug("waiting complete", zap.Stringer("status", s.Status()))
				// nothing in the pipeline
				return s.err
			}

		case <-ctx.Done():
			s.log.Debug("wait context done", zap.Error(ctx.Err()))
			s.Cancel()
			return s.err
		}
	}
}

func (s *Session) worker(ctx context.Context) {
	defer s.Stop()

	// making sure
	defer close(s.workerLock)
	s.workerLock <- struct{}{}

	workerTicker := time.NewTicker(s.workerInterval)

	defer workerTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Debug("worker context done", zap.Error(ctx.Err()))
			return

		case <-workerTicker.C:
			s.queueScheduledSuspended()

		case st := <-s.qState:
			if st == nil {
				// stop worker
				s.log.Debug("completed")
				return
			}

			s.log.Debug("pulled state from queue", zap.Uint64("stateID", st.stateId))
			if st.step == nil {
				// We should not terminate if the session contains any delayed or prompted steps.
				status := s.Status()
				if status == SessionPrompted || status == SessionDelayed {
					break
				}

				s.log.Debug("done, setting results and stopping the worker")

				func() {
					// mini lambda fn to ensure we can properly unlock with defer
					s.mux.Lock()
					defer s.mux.Unlock()

					// with merge we are making sure
					// that result != nil even if state scope is
					s.result = (&expr.Vars{}).MustMerge(st.scope)
				}()

				// Call event handler with completed status
				s.eventHandler(SessionCompleted, st, s)

				return
			}

			// add empty struct to chan to lock and to have control over number of concurrent go processes
			// this will block if number of items in execLock chan reached value of sessionConcurrentExec
			s.execLock <- struct{}{}

			go func() {
				defer func() {
					// remove protection that prevents multiple
					// steps executing at the same time
					<-s.execLock
				}()

				var (
					err error
					log = s.log.With(zap.Uint64("stateID", st.stateId))
				)

				nxt, err := s.exec(ctx, log, st)
				if err != nil && st.err == nil {
					// override the error from the execution
					st.err = err
				}

				st.completed = now()

				status := s.Status()
				if st.err != nil {
					st.err = fmt.Errorf(
						"workflow %d step %d execution failed: %w",
						s.workflowID,
						st.step.ID(),
						st.err,
					)

					s.mux.Lock()

					// when the err handler is defined, the error was handled and should not kill the workflow
					if !st.errHandled {
						// We need to force failed session status
						// because it's not set early enough to pick it up with s.Status()
						status = SessionFailed

						// pushing step execution error into error queue
						// to break worker loop
						s.qErr <- st.err
					}

					s.mux.Unlock()
				}

				s.log.Debug(
					"executed",
					zap.Uint64("stateID", st.stateId),
					zap.Stringer("status", status),
					zap.Error(st.err),
				)

				s.eventHandler(status, st, s)

				for _, n := range nxt {
					if n.step != nil {
						log.Debug("next step queued", zap.Uint64("nextStepId", n.step.ID()))
					} else {
						log.Debug("next step queued", zap.Uint64("nextStepId", 0))
					}
					if err = s.enqueue(ctx, n); err != nil {
						log.Error("unable to enqueue", zap.Error(err))
						return
					}
				}
			}()

		case err := <-s.qErr:
			s.mux.Lock()
			defer s.mux.Unlock()

			if err == nil {
				// stop worker
				return
			}

			// set final error on session
			s.err = err
			return
		}
	}
}

func (s *Session) Cancel() {
	s.log.Debug("canceling")
	s.qErr <- fmt.Errorf("canceled")
}

func (s *Session) Stop() {
	s.log.Debug("stopping worker")
	s.qErr <- nil
}

func (s *Session) Suspended() bool {
	defer s.mux.RUnlock()
	s.mux.RLock()
	return len(s.delayed) > 0
}

func (s *Session) queueScheduledSuspended() {
	defer s.mux.Unlock()
	s.mux.Lock()

	for id, sus := range s.delayed {
		if !sus.resumeAt.IsZero() && sus.resumeAt.After(*now()) {
			continue
		}

		delete(s.delayed, id)

		// Set state input when step is resumed
		sus.state.input = &expr.Vars{}
		sus.state.input.Set("resumed", true)
		sus.state.input.Set("resumeAt", sus.resumeAt)
		s.qState <- sus.state
	}
}

// executes single step, resolves response and schedule following steps for execution
func (s *Session) exec(ctx context.Context, log *zap.Logger, st *State) (nxt []*State, err error) {
	st.created = *now()

	defer func() {
		var reason interface{}
		reason = recover()
		if reason == nil {
			return
		}

		var perr error

		// normalize error and set it to state
		switch reason := reason.(type) {
		case error:
			perr = fmt.Errorf("step %d crashed: %w", st.step.ID(), reason)
		default:
			perr = fmt.Errorf("step %d crashed: %v", st.step.ID(), reason)
		}

		if s.dumpStacktraceOnPanic {
			fmt.Printf("Error: %v\n", perr)
			println(string(debug.Stack()))
		}

		s.qErr <- perr
	}()

	var (
		result ExecResponse
		scope  = (&expr.Vars{}).MustMerge(st.scope)

		currLoop = st.loopCurr()
	)

	if st.step != nil {
		log = log.With(zap.Uint64("stepID", st.step.ID()))
	}

	{
		if currLoop != nil && currLoop.Is(st.step) {
			result = currLoop
		} else {
			// push logger to context but raise the stacktrace level to panic
			// to prevent overly verbose traces
			ctx = logger.ContextWithValue(ctx, log)
			stepCtx := SetContextCallStack(ctx, s.callStack)

			result, st.err = st.step.Exec(stepCtx, st.MakeRequest())

			if iterator, isIterator := result.(Iterator); isIterator && st.err == nil {
				// Exec fn returned an iterator, adding loop to stack
				st.newLoop(iterator)
				if err = iterator.Start(ctx, scope); err != nil {
					return
				}
			}
		}

		if st.err != nil {
			if st.errHandler == nil {
				// no error handler set
				return nil, st.err
			}

			// handling error with error handling
			// step set in one of the previous steps
			log.Warn("step execution error handled",
				zap.Uint64("errorHandlerStepId", st.errHandler.ID()),
				zap.Error(st.err),
			)

			_ = expr.Assign(scope, "error", expr.Must(expr.NewString(st.err.Error())))

			// copy error handler & disable it on state to prevent inf. loop
			// in case of another error in the error-handling branch
			eh := st.errHandler
			st.errHandler = nil
			st.errHandled = true
			return []*State{st.Next(eh, scope)}, nil
		}

		switch l := result.(type) {
		case Iterator:
			st.action = "iterator initialized"
			// add looper to state
			var (
				n Step
			)
			n, result, st.err = l.Next(ctx, scope)
			if st.err != nil {
				return nil, st.err
			}

			if n == nil {
				st.next = st.loopEnd()
			} else {
				st.next = Steps{n}
			}
		}

		log.Debug("step executed", zap.String("resultType", fmt.Sprintf("%T", result)))
		switch result := result.(type) {
		case *expr.Vars:
			// most common (successful) result
			// session will continue with configured child steps
			st.results = result
			scope = scope.MustMerge(st.results)

		case *errHandler:
			st.action = "error handler initialized"
			// this step sets error handling step on current state
			// and continues on the current path
			st.errHandler = result.handler

			// find step that's not error handler and
			// use it for the next step
			for _, c := range s.g.Children(st.step) {
				if c != st.errHandler {
					st.next = Steps{c}
					break
				}
			}

		case *loopBreak:
			st.action = "loop break"
			if currLoop == nil {
				return nil, fmt.Errorf("break step not inside a loop")
			}

			// jump out of the loop
			st.next = st.loopEnd()
			log.Debug("breaking from iterator")

		case *loopContinue:
			st.action = "loop continue"
			if currLoop == nil {
				return nil, fmt.Errorf("continue step not inside a loop")
			}

			// jump back to iterator
			st.next = Steps{currLoop.Iterator()}
			log.Debug("continuing with next iteration")

		case *partial:
			st.action = "partial"
			// *partial is returned when step needs to be executed again
			// it's used mainly for join gateway step that should be called multiple times (one for each parent path)
			return

		case *termination:
			st.action = "termination"
			// terminate all activities, all delayed tasks and exit right away
			log.Debug("termination", zap.Int("delayed", len(s.delayed)))
			s.mux.Lock()
			s.delayed = nil
			s.mux.Unlock()
			return []*State{FinalState(s, scope)}, nil

		case *delayed:
			st.action = "delayed"
			log.Debug("session delayed", zap.Time("at", result.resumeAt))

			result.state = st
			s.mux.Lock()
			s.delayed[st.stateId] = result
			s.mux.Unlock()
			return

		case *resumed:
			st.action = "resumed"
			log.Debug("session resumed")

		case *prompted:
			st.action = "prompted"
			if result.ownerId == 0 {
				return nil, fmt.Errorf("without an owner")
			}

			result.state = st
			s.mux.Lock()
			s.prompted[st.stateId] = result
			s.mux.Unlock()
			return

		case Steps:
			st.action = "next-steps"
			// session continues with set of specified steps
			// steps MUST be configured in a graph as step's children
			st.next = result

		case Step:
			st.action = "next-step"
			// session continues with a specified step
			// step MUST be configured in a graph as step's child
			st.next = Steps{result}

		default:
			return nil, fmt.Errorf("unknown exec response type %T", result)
		}
	}

	if len(st.next) == 0 {
		// step's exec did not return next steps (only gateway steps, iterators and loops controls usually do that)
		//
		// rely on graph and get next (children) steps from there
		st.next = s.g.Children(st.step)
	} else {
		// children returned from step's exec
		// do a quick sanity check
		cc := s.g.Children(st.step)
		if len(cc) > 0 && !cc.Contains(st.next...) {
			return nil, fmt.Errorf("inconsistent relationship")
		}
	}

	if currLoop != nil && len(st.next) == 0 {
		// gracefully handling last step of iteration branch
		// that does not point back to the iterator step
		st.next = Steps{currLoop.Iterator()}
		log.Debug("last step in iteration branch, going back", zap.Uint64("backStepId", st.next[0].ID()))
	}

	if len(st.next) == 0 {
		log.Debug("zero paths, finalizing")
		// using state to transport results and complete the worker loop
		return []*State{FinalState(s, scope)}, nil
	}

	nxt = make([]*State, len(st.next))
	for i, step := range st.next {
		nn := st.Next(step, scope)
		if err = s.canEnqueue(nn); err != nil {
			log.Error("unable to queue", zap.Error(err))
			return
		}

		nxt[i] = nn
	}

	return nxt, nil
}

func SetWorkerInterval(i time.Duration) SessionOpt {
	return func(s *Session) {
		s.workerInterval = i
	}
}

func SetHandler(fn StateChangeHandler) SessionOpt {
	return func(s *Session) {
		s.eventHandler = fn
	}
}

func SetWorkflowID(workflowID uint64) SessionOpt {
	return func(s *Session) {
		s.workflowID = workflowID
	}
}

func SetLogger(log *zap.Logger) SessionOpt {
	return func(s *Session) {
		s.log = log
	}
}

func SetDumpStacktraceOnPanic(dump bool) SessionOpt {
	return func(s *Session) {
		s.dumpStacktraceOnPanic = dump
	}
}

func SetCallStack(id ...uint64) SessionOpt {
	return func(s *Session) {
		s.callStack = id
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

func (ss Steps) IDs() []uint64 {
	if len(ss) == 0 {
		return nil
	}

	var ids = make([]uint64, len(ss))
	for i := range ss {
		ids[i] = ss[i].ID()
	}

	return ids
}

func SetContextCallStack(ctx context.Context, ss []uint64) context.Context {
	return context.WithValue(ctx, callStackCtxKey{}, ss)
}

func GetContextCallStack(ctx context.Context) []uint64 {
	v := ctx.Value(callStackCtxKey{})
	if v == nil {
		return nil
	}

	return v.([]uint64)
}
