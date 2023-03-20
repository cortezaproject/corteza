package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/actionlog_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"time"
)

var _ wfexec.ExecResponse

type (
	actionlogHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h actionlogHandler) register() {
	h.reg.AddFunctions(
		h.Search(),
		h.Each(),
		h.Record(),
	)
}

type (
	actionlogSearchArgs struct {
		hasFromTimestamp bool
		FromTimestamp    *time.Time

		hasToTimestamp bool
		ToTimestamp    *time.Time

		hasBeforeActionID bool
		BeforeActionID    uint64

		hasActorID bool
		ActorID    uint64

		hasOrigin bool
		Origin    string

		hasResource bool
		Resource    string

		hasAction bool
		Action    string

		hasLimit bool
		Limit    uint64
	}

	actionlogSearchResults struct {
		Actions []*actionlog.Action
	}
)

// Search function Action log search
//
// expects implementation of search function:
// func (h actionlogHandler) search(ctx context.Context, args *actionlogSearchArgs) (results *actionlogSearchResults, err error) {
//    return
// }
func (h actionlogHandler) Search() *atypes.Function {
	return &atypes.Function{
		Ref:    "actionlogSearch",
		Kind:   "function",
		Labels: map[string]string{"actionlog": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Action log search",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "fromTimestamp",
				Types: []string{"DateTime"},
			},
			{
				Name:  "toTimestamp",
				Types: []string{"DateTime"},
			},
			{
				Name:  "beforeActionID",
				Types: []string{"ID"},
			},
			{
				Name:  "actorID",
				Types: []string{"ID"},
			},
			{
				Name:  "origin",
				Types: []string{"String"},
			},
			{
				Name:  "resource",
				Types: []string{"String"},
			},
			{
				Name:  "action",
				Types: []string{"String"},
			},
			{
				Name:  "limit",
				Types: []string{"UnsignedInteger"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:    "actions",
				Types:   []string{"Action"},
				IsArray: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &actionlogSearchArgs{
					hasFromTimestamp:  in.Has("fromTimestamp"),
					hasToTimestamp:    in.Has("toTimestamp"),
					hasBeforeActionID: in.Has("beforeActionID"),
					hasActorID:        in.Has("actorID"),
					hasOrigin:         in.Has("origin"),
					hasResource:       in.Has("resource"),
					hasAction:         in.Has("action"),
					hasLimit:          in.Has("limit"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *actionlogSearchResults
			if results, err = h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Actions (*actionlog.Action) to Array (of Action)
				var (
					tval expr.TypedValue
					tarr = make([]expr.TypedValue, len(results.Actions))
				)

				for i := range results.Actions {
					if tarr[i], err = h.reg.Type("Action").Cast(results.Actions[i]); err != nil {
						return
					}
				}

				if tval, err = expr.NewArray(tarr); err != nil {
					return
				} else if err = expr.Assign(out, "actions", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	actionlogEachArgs struct {
		hasFromTimestamp bool
		FromTimestamp    *time.Time

		hasToTimestamp bool
		ToTimestamp    *time.Time

		hasBeforeActionID bool
		BeforeActionID    uint64

		hasActorID bool
		ActorID    uint64

		hasOrigin bool
		Origin    string

		hasResource bool
		Resource    string

		hasAction bool
		Action    string

		hasLimit bool
		Limit    uint64
	}

	actionlogEachResults struct {
		Action *actionlog.Action
	}
)

// Each function Action log
//
// expects implementation of each function:
// func (h actionlogHandler) each(ctx context.Context, args *actionlogEachArgs) (results *actionlogEachResults, err error) {
//    return
// }
func (h actionlogHandler) Each() *atypes.Function {
	return &atypes.Function{
		Ref:    "actionlogEach",
		Kind:   "iterator",
		Labels: map[string]string{"actionlog": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Action log",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "fromTimestamp",
				Types: []string{"DateTime"},
			},
			{
				Name:  "toTimestamp",
				Types: []string{"DateTime"},
			},
			{
				Name:  "beforeActionID",
				Types: []string{"ID"},
			},
			{
				Name:  "actorID",
				Types: []string{"ID"},
			},
			{
				Name:  "origin",
				Types: []string{"String"},
			},
			{
				Name:  "resource",
				Types: []string{"String"},
			},
			{
				Name:  "action",
				Types: []string{"String"},
			},
			{
				Name:  "limit",
				Types: []string{"UnsignedInteger"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "action",
				Types: []string{"Action"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &actionlogEachArgs{
					hasFromTimestamp:  in.Has("fromTimestamp"),
					hasToTimestamp:    in.Has("toTimestamp"),
					hasBeforeActionID: in.Has("beforeActionID"),
					hasActorID:        in.Has("actorID"),
					hasOrigin:         in.Has("origin"),
					hasResource:       in.Has("resource"),
					hasAction:         in.Has("action"),
					hasLimit:          in.Has("limit"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return h.each(ctx, args)
		},
	}
}

type (
	actionlogRecordArgs struct {
		hasAction bool
		Action    string

		hasResource bool
		Resource    string

		hasError bool
		Error    string

		hasSeverity bool
		Severity    string

		hasDescription bool
		Description    string

		hasMeta bool
		Meta    *expr.Vars
	}
)

// Record function Record action into action log
//
// expects implementation of record function:
// func (h actionlogHandler) record(ctx context.Context, args *actionlogRecordArgs) (err error) {
//    return
// }
func (h actionlogHandler) Record() *atypes.Function {
	return &atypes.Function{
		Ref:    "actionlogRecord",
		Kind:   "function",
		Labels: map[string]string{"actionlog": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Record action into action log",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "action",
				Types: []string{"String"},
			},
			{
				Name:  "resource",
				Types: []string{"String"},
			},
			{
				Name:  "error",
				Types: []string{"String"},
			},
			{
				Name:  "severity",
				Types: []string{"String"},
				Meta: &atypes.ParamMeta{
					Visual: map[string]interface{}{"options": []interface{}{"emergency", "alert", "critical", "err", "warning", "notice", "info", "debug"}},
				},
			},
			{
				Name:  "description",
				Types: []string{"String"},
			},
			{
				Name:  "meta",
				Types: []string{"Vars"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &actionlogRecordArgs{
					hasAction:      in.Has("action"),
					hasResource:    in.Has("resource"),
					hasError:       in.Has("error"),
					hasSeverity:    in.Has("severity"),
					hasDescription: in.Has("description"),
					hasMeta:        in.Has("meta"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.record(ctx, args)
		},
	}
}
