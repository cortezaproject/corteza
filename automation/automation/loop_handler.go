package automation

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	loopHandler struct {
		reg    loopHandlerRegistry
		parser expr.Parsable
	}
)

func LoopHandler(reg loopHandlerRegistry, p expr.Parsable) *loopHandler {
	h := &loopHandler{
		reg:    reg,
		parser: p,
	}

	h.register()
	return h
}

func (h loopHandler) sequence(_ context.Context, args *loopSequenceArgs) (wfexec.IteratorHandler, error) {
	if !args.hasFirst {
		args.First = 0
	}

	if !args.hasLast {
		args.Last = 1
	}

	if !args.hasStep {
		args.Step = 1
	}

	if args.First*(args.Step/args.Step) >= args.Last*(args.Step/args.Step) {
		return nil, fmt.Errorf("failed to initialize counter iterator with first step greater than last")
	}

	i := &sequenceIterator{
		counter: 0,
		cFirst:  args.First,
		cLast:   args.Last,
		cStep:   args.Step,
	}

	return i, nil
}

func (h loopHandler) do(_ context.Context, args *loopDoArgs) (wfexec.IteratorHandler, error) {
	var (
		i   = &conditionIterator{}
		err error
	)

	if i.expr, err = h.parser.Parse(args.While); err != nil {
		return nil, err
	}

	return i, nil
}

func (h loopHandler) each(_ context.Context, args *loopEachArgs) (wfexec.IteratorHandler, error) {
	return &collectionIterator{set: args.Items}, nil
}

func (h loopHandler) lines(_ context.Context, args *loopLinesArgs) (wfexec.IteratorHandler, error) {
	return &lineIterator{s: bufio.NewScanner(args.Stream)}, nil
}
