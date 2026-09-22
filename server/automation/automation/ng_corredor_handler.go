package automation

import (
	"context"

	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
)

type (
	ngCorredorHandler struct {
		reg     constructSvc
		h       *corredorHandler
		enabled bool
	}
)

// NgCorredorHandler registers the corredorExec construct; with Corredor
// turned off the construct is listed but disabled
func NgCorredorHandler(reg constructSvc, svc corredorServiceExecutor, enabled bool) *ngCorredorHandler {
	h := &ngCorredorHandler{
		reg:     reg,
		h:       &corredorHandler{svc: svc},
		enabled: enabled,
	}

	h.register()
	return h
}

func (h ngCorredorHandler) register() {
	h.reg.AddFunctions(
		h.Exec(),
	)
}

// Exec runs a manual Corredor server script with the given arguments and
// returns what the script returned
func (h ngCorredorHandler) Exec() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:      "corredorExec",
		Kind:     "function",
		Groups:   []string{"Corredor"},
		Labels:   map[string]string{"corredor": "step,workflow"},
		Disabled: !h.enabled,
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Run Corredor script",
			Description: "Run a manual Corredor server script and return what it returns",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "code"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "script",
				Types:        []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Script",
					Description: "Name of the server script to run, e.g. /server-scripts/hello.js:default.",
				},
			},
			{
				ArgumentName: "args",
				Types:        []string{"Vars"},
				Meta: &atypes.ParamMeta{
					Label:       "Arguments",
					Description: "Values the script receives as its arguments.",
				},
			},
		},
		Results: []*atypes.Param{
			{
				ArgumentName: "results",
				Types:        []string{"Vars"},
			},
		},
		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "CorredorScriptSelector", Label: "Script", Argument: "script", Required: true}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Arguments", Argument: "args"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &corredorExecArgs{
					hasScript: in.Has("script"),
					hasArgs:   in.Has("args"),
				}

				results *corredorExecResults
			)

			if err = in.Decode(args); err != nil {
				return
			}

			if args.Args == nil {
				args.Args = &expr.Vars{}
			}

			if results, err = h.h.exec(ctx, args); err != nil {
				return
			}

			return expr.NewVars(map[string]interface{}{"results": results.Results})
		},
	}
}
