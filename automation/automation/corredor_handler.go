package automation

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/system/service/event"
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

	r = &corredorExecResults{}
	if r.Results, err = expr.NewVars(sArgs.payload); err != nil {
		return nil, err
	}

	return

}

func makeScriptArgs(in *expr.Vars) *scriptArgs {
	payload := in.Dict()
	sanitizeMapStringInterface(payload)

	return &scriptArgs{payload: payload}
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

// Sanitizing all uint64 values that are encoded into JSON
func sanitizeMapStringInterface(m map[string]interface{}) {
	var sw func(interface{}) string
	sw = func(i interface{}) (s string) {
		switch v := i.(type) {
		case uint64:
			// make sure uint64 values on fields ending with ID
			// are properly encoded as strings
			s = strconv.FormatUint(v, 10)

		case map[string]interface{}:
			sanitizeMapStringInterface(v)
		case []interface{}:
			for _, vv := range i.([]interface{}) {
				sw(vv)
			}
		}
		return
	}
	for k := range m {
		if s := sw(m[k]); len(s) > 0 {
			m[k] = sw(m[k])
		}
	}
}
