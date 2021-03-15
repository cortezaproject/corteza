package automation

import (
	"context"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/spf13/cast"
)

type (
	corredorHandler struct {
		reg corredorHandlerRegistry
		svc corredorServiceExecutor
	}

	corredorServiceExecutor interface {
		Exec(ctx context.Context, scriptName string, args corredor.ScriptArgs) (err error)
	}

	scriptArgs struct {
		payload map[string]interface{}
	}
)

func CorredorHandler(reg corredorHandlerRegistry, svc corredorServiceExecutor) *corredorHandler {
	h := &corredorHandler{
		reg: reg,
		svc: svc,
	}

	h.register()
	return h
}

func (h corredorHandler) exec(ctx context.Context, args *corredorExecArgs) (r *corredorExecResults, err error) {
	sArgs := makeScriptArgs(args.Args)

	if err = h.svc.Exec(ctx, args.Script, sArgs); err != nil {
		return
	}

	return &corredorExecResults{Results: sArgs.payload}, nil

}

func makeScriptArgs(in interface{}) *scriptArgs {
	return &scriptArgs{
		payload: cast.ToStringMap(in),
	}
}

// mimic onManual event on system:

func (scriptArgs) ResourceType() string                  { return event.SystemOnManual().ResourceType() }
func (scriptArgs) EventType() string                     { return event.SystemOnManual().EventType() }
func (scriptArgs) Match(eventbus.ConstraintMatcher) bool { return false }

func (a *scriptArgs) Encode() (enc map[string][]byte, err error) {
	enc = make(map[string][]byte)
	for k, v := range a.payload {
		enc[k], err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (a *scriptArgs) Decode(enc map[string][]byte) (err error) {
	a.payload = make(map[string]interface{})

	for k, v := range enc {
		var aux interface{}
		if err = json.Unmarshal(v, &aux); err != nil {
			return err
		}

		a.payload[k] = aux
	}

	return
}
