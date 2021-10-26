package automation

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	actionlogHandler struct {
		reg actionlogHandlerRegistry
		svc actionlog.Recorder
	}

	actionSetIterator struct {
		// Item buffer, current item pointer, and total items traversed
		ptr    uint
		buffer actionlog.ActionSet
		total  uint

		// When filter limit is set, this constraints it
		iterLimit    uint
		useIterLimit bool

		// Item loader for additional chunks
		filter actionlog.Filter
		loader func() error
	}
)

func ActionlogHandler(reg actionlogHandlerRegistry, svc actionlog.Recorder) *actionlogHandler {
	h := &actionlogHandler{
		reg: reg,
		svc: svc,
	}

	h.register()
	return h
}

func (h actionlogHandler) search(ctx context.Context, args *actionlogSearchArgs) (results *actionlogSearchResults, err error) {
	results = &actionlogSearchResults{}

	var (
		f = actionlog.Filter{
			FromTimestamp:  args.FromTimestamp,
			ToTimestamp:    args.ToTimestamp,
			BeforeActionID: args.BeforeActionID,
			Origin:         args.Origin,
			Resource:       args.Resource,
			Action:         args.Action,
			Limit:          uint(args.Limit),
		}
	)

	if args.ActorID > 0 {
		f.ActorID = []uint64{args.ActorID}
	}

	results.Actions, _, err = h.svc.Find(ctx, f)
	return
}

func (h actionlogHandler) each(ctx context.Context, args *actionlogEachArgs) (out wfexec.IteratorHandler, err error) {
	var (
		i = &actionSetIterator{}
		f = actionlog.Filter{
			FromTimestamp:  args.FromTimestamp,
			ToTimestamp:    args.ToTimestamp,
			BeforeActionID: args.BeforeActionID,
			Origin:         args.Origin,
			Resource:       args.Resource,
			Action:         args.Action,
		}
	)

	if args.ActorID > 0 {
		f.ActorID = []uint64{args.ActorID}
	}

	if args.hasLimit {
		i.useIterLimit = true
		i.iterLimit = uint(args.Limit)

		if args.Limit > uint64(wfexec.MaxIteratorBufferSize) {
			f.Limit = wfexec.MaxIteratorBufferSize
		}
		i.iterLimit = uint(args.Limit)
	} else {
		f.Limit = wfexec.MaxIteratorBufferSize
	}

	i.filter = f
	i.loader = func() (err error) {
		i.total += i.ptr
		i.ptr = 0

		if len(i.buffer) > 0 {
			i.filter.BeforeActionID = i.buffer[len(i.buffer)-1].ID
		}
		i.buffer, i.filter, err = h.svc.Find(ctx, i.filter)

		return
	}

	// Initial load
	return i, i.loader()
}

func (h actionlogHandler) record(ctx context.Context, args *actionlogRecordArgs) (err error) {
	a := &actionlog.Action{
		Resource:    args.Resource,
		Action:      args.Action,
		Error:       args.Error,
		Severity:    actionlog.NewSeverity(args.Severity),
		Description: args.Description,
	}

	if args.Meta != nil {
		a.Meta = args.Meta.Dict()
	}

	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_Automation)

	h.svc.Record(ctx, a)

	return nil
}

func (i *actionSetIterator) More(context.Context, *Vars) (bool, error) {
	return wfexec.GenericResourceNextCheck(i.useIterLimit, i.ptr, uint(len(i.buffer)), i.total, i.iterLimit, i.iterLimit > uint(len(i.buffer))), nil
}

func (i *actionSetIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *actionSetIterator) Next(context.Context, *Vars) (out *Vars, err error) {
	if len(i.buffer)-int(i.ptr) <= 0 {
		if err = i.loader(); err != nil {
			panic(err)
		}
	}

	out = &Vars{}
	out.Set("user", Must(NewUser(i.buffer[i.ptr])))
	out.Set("index", Must(NewInteger(i.total+i.ptr)))

	i.ptr++
	return out, nil
}
