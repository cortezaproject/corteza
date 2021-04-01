package automation

import (
	"bufio"
	"context"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	// iterates from start to stop by step
	sequenceIterator struct {
		counter, cFirst, cLast, cStep int64
	}
)

func (i *sequenceIterator) More(context.Context, *Vars) (bool, error) {
	return i.more(), nil
}

func (i *sequenceIterator) more() bool {
	return i.counter*(i.cStep/i.cStep) < i.cLast*(i.cStep/i.cStep)
}

func (i *sequenceIterator) Start(context.Context, *Vars) error { return nil }

func (i *sequenceIterator) Next(context.Context, *Vars) (*Vars, error) {
	out := &Vars{}
	out.Set("counter", i.counter)
	out.Set("isFirst", i.counter == i.cFirst)
	out.Set("isLast", !i.more())
	i.counter = i.counter + i.cStep
	return out, nil
}

type (
	// iterates from start to stop by step
	conditionIterator struct {
		expr Evaluable
	}
)

func (i *conditionIterator) More(ctx context.Context, scope *Vars) (bool, error) {
	return i.expr.Test(ctx, scope)
}

func (i *conditionIterator) Start(context.Context, *Vars) error { return nil }

func (i *conditionIterator) Next(context.Context, *Vars) (*Vars, error) {
	return &Vars{}, nil
}

type (
	// iterates from start to stop by step
	collectionIterator struct {
		ptr int
		set []TypedValue
	}
)

func (i *collectionIterator) More(context.Context, *Vars) (bool, error) {
	return i.ptr < len(i.set), nil
}

func (i *collectionIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *collectionIterator) Next(context.Context, *Vars) (out *Vars, err error) {
	out = &Vars{}
	out.Set("line", i.set[i.ptr])
	i.ptr++
	return out, nil
}

type (
	// iterates from start to stop by step
	lineIterator struct {
		s *bufio.Scanner
	}
)

func (i *lineIterator) More(context.Context, *Vars) (bool, error) {
	return i.s.Scan(), nil
}

func (i *lineIterator) Start(context.Context, *Vars) error {
	return nil
}

func (i *lineIterator) Next(context.Context, *Vars) (*Vars, error) {
	if err := i.s.Err(); err != nil {
		return nil, err
	}

	out := &Vars{}
	out.Set("line", i.s.Text())
	return out, nil
}
