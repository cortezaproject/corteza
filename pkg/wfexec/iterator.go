package wfexec

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	// Iterator can be returned from Exec fn as ExecResponse
	//
	// It helps session's exec fn() to properly navigate through graph
	// by calling is/break/iterator/next function
	Iterator interface {
		// Is the given step this iterator step
		Is(Step) bool

		// Initialize iterator
		Start(context.Context, *expr.Vars) error

		// Break fn is called when loop is forcefully broken
		Break() Step

		Iterator() Step

		// Next is called before each iteration and returns
		// 1st step of the iteration branch and variables that are added to the scope
		Next(context.Context, *expr.Vars) (Step, *expr.Vars, error)
	}

	ResultEvaluator interface {
		EvalResults(ctx context.Context, results *expr.Vars) (out *expr.Vars, err error)
	}

	IteratorHandler interface {
		Start(context.Context, *expr.Vars) error
		More(context.Context, *expr.Vars) (bool, error)
		Next(context.Context, *expr.Vars) (*expr.Vars, error)
	}

	// Handles communication between Session's exec() fn and iterator handler
	genericIterator struct {
		iter, next, exit Step

		h IteratorHandler
	}
)

// GenericIterator creates a wrapper around IteratorHandler and
// returns genericIterator that implements Iterator interface
func GenericIterator(iter, next, exit Step, h IteratorHandler) Iterator {
	return &genericIterator{
		iter: iter,
		next: next,
		exit: exit,
		h:    h,
	}
}

func (i *genericIterator) Is(s Step) bool                                { return i.iter == s }
func (i *genericIterator) Start(ctx context.Context, s *expr.Vars) error { return i.h.Start(ctx, s) }
func (i *genericIterator) Break() Step                                   { return i.exit }
func (i *genericIterator) Iterator() Step                                { return i.iter }

// Next calls More and Next functions on iterator handler.
//
// If iterator step (iter field) implements ResultEvaluator it calls
// EvalResults on it before returning it. If iterator step does not implement it,
// results are omitted.
func (i *genericIterator) Next(ctx context.Context, scope *expr.Vars) (next Step, out *expr.Vars, err error) {
	var (
		more    bool
		results *expr.Vars
	)
	if more, err = i.h.More(ctx, scope); err != nil || !more {
		return
	}

	if results, err = i.h.Next(ctx, scope); err != nil {
		return
	}

	if re, is := i.iter.(ResultEvaluator); is {
		if out, err = re.EvalResults(ctx, results); err != nil {
			return
		}
	}

	next = i.next

	return
}
