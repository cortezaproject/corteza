package automation

import (
	"context"

	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
)

type (
	ngLogHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *logHandler
	}
)

func NgLogHandler(reg constructSvc, tReg typeRegistry) *ngLogHandler {
	ngh := &ngLogHandler{
		reg:  reg,
		tReg: tReg,
		h:    &logHandler{},
	}

	ngh.register()
	return ngh
}

func (h ngLogHandler) register() {}

func (h ngLogHandler) Debug() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "logDebug",
		Kind:   "function",
		Groups: []string{"Logging"},
		Labels: map[string]string{"debug": "step", "logger": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Log Debug",
			Description: "Write a debug-level log message",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "list"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "message", Types: []string{"String"}, Required: true},
			{ArgumentName: "fields", Types: []string{"KV"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Message", Argument: "message"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Fields", Argument: "fields"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &logDebugArgs{
					hasMessage: in.Has("message"),
					hasFields:  in.Has("fields"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.h.debug(ctx, args)
		},
	}
}

func (h ngLogHandler) Info() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "logInfo",
		Kind:   "function",
		Groups: []string{"Logging"},
		Labels: map[string]string{"debug": "step", "logger": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Log Info",
			Description: "Write an info-level log message",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "list"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "message", Types: []string{"String"}, Required: true},
			{ArgumentName: "fields", Types: []string{"KV"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Message", Argument: "message"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Fields", Argument: "fields"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &logInfoArgs{
					hasMessage: in.Has("message"),
					hasFields:  in.Has("fields"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.h.info(ctx, args)
		},
	}
}

func (h ngLogHandler) Warn() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "logWarn",
		Kind:   "function",
		Groups: []string{"Logging"},
		Labels: map[string]string{"debug": "step", "logger": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Log Warning",
			Description: "Write a warning-level log message",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "list"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "message", Types: []string{"String"}, Required: true},
			{ArgumentName: "fields", Types: []string{"KV"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Message", Argument: "message"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Fields", Argument: "fields"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &logWarnArgs{
					hasMessage: in.Has("message"),
					hasFields:  in.Has("fields"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.h.warn(ctx, args)
		},
	}
}

func (h ngLogHandler) Error() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "logError",
		Kind:   "function",
		Groups: []string{"Logging"},
		Labels: map[string]string{"debug": "step", "logger": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Log Error",
			Description: "Write an error-level log message",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "list"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "message", Types: []string{"String"}, Required: true},
			{ArgumentName: "fields", Types: []string{"KV"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Message", Argument: "message"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Fields", Argument: "fields"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &logErrorArgs{
					hasMessage: in.Has("message"),
					hasFields:  in.Has("fields"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return out, h.h.error(ctx, args)
		},
	}
}
