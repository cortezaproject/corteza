package automation

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/crusttech/human/server/pkg/corredor"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/stretchr/testify/require"
)

// mockCorredor answers every script with the arguments it received plus a
// marker, the way a script returning a plain object does
type mockCorredor struct {
	gotScript string
	gotArgs   map[string]interface{}
}

func (m *mockCorredor) Exec(_ context.Context, scriptName string, args corredor.ScriptArgs) error {
	m.gotScript = scriptName

	enc, err := args.Encode()
	if err != nil {
		return err
	}

	m.gotArgs = map[string]interface{}{}
	for k, v := range enc {
		var aux interface{}
		if err = json.Unmarshal(v, &aux); err != nil {
			return err
		}
		m.gotArgs[k] = aux
	}

	// script returns its arguments and a marker
	enc["pong"], _ = json.Marshal("pong")
	return args.Decode(enc)
}

func TestNgCorredorHandler_Exec(t *testing.T) {
	var (
		req = require.New(t)
		svc = &mockCorredor{}
		h   = ngCorredorHandler{h: &corredorHandler{svc: svc}, enabled: true}
		fn  = h.Exec()
	)

	req.Equal("corredorExec", fn.Ref)
	req.False(fn.Disabled)

	scriptArgs, err := expr.NewVars(map[string]interface{}{"ping": "ping"})
	req.NoError(err)

	in, err := expr.NewVars(map[string]interface{}{
		"script": "/server-scripts/SystemPing.js:default",
		"args":   scriptArgs,
	})
	req.NoError(err)

	out, err := fn.Handler(context.Background(), in)
	req.NoError(err)

	req.Equal("/server-scripts/SystemPing.js:default", svc.gotScript)
	req.Equal("ping", svc.gotArgs["ping"])

	results, err := out.Select("results")
	req.NoError(err)
	rv, ok := results.(*expr.Vars)
	req.True(ok)
	pong, err := rv.Select("pong")
	req.NoError(err)
	req.Equal("pong", pong.Get())
}

func TestNgCorredorHandler_Disabled(t *testing.T) {
	h := ngCorredorHandler{h: &corredorHandler{svc: &mockCorredor{}}, enabled: false}
	require.True(t, h.Exec().Disabled)
}
