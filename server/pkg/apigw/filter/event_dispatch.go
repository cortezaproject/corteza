package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/crusttech/human/server/pkg/apigw/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/j7s"
	sysEvent "github.com/crusttech/human/server/system/service/event"
)

type (
	eventDispatch struct {
		types.FilterMeta

		cfg types.Config

		params struct {
			ConnectionID           uint64              `json:"connectionID,string"`
			ConfiguredConnectionID uint64              `json:"configuredConnectionID,string"`
			EventType              string              `json:"eventType"`
			Mapping                map[string][]string `json:"mapping,omitempty"`
		}
	}
)

func NewEventDispatch(cfg types.Config) *eventDispatch {
	f := &eventDispatch{cfg: cfg}

	f.Name = "eventDispatch"
	f.Label = "Event dispatch processer"
	f.Kind = types.Processer

	f.Args = []*types.FilterMetaArg{
		{Type: "connectionID", Label: "connectionID", Options: map[string]interface{}{}},
		{Type: "configuredConnectionID", Label: "configuredConnectionID", Options: map[string]interface{}{}},
		{Type: "eventType", Label: "eventType", Options: map[string]interface{}{}},
		{Type: "mapping", Label: "mapping", Options: map[string]interface{}{}},
	}

	return f
}

func (h eventDispatch) New(cfg types.Config) types.Handler {
	return NewEventDispatch(cfg)
}

func (h eventDispatch) Enabled() bool { return true }

func (h eventDispatch) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h eventDispatch) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h *eventDispatch) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	h.cfg = cfg

	if err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params); err != nil {
		return h, err
	}

	if h.params.ConnectionID == 0 {
		return h, fmt.Errorf("eventDispatch: connectionID is required")
	}
	if h.params.ConfiguredConnectionID == 0 {
		return h, fmt.Errorf("eventDispatch: configuredConnectionID is required")
	}
	if h.params.EventType == "" {
		return h, fmt.Errorf("eventDispatch: eventType is required")
	}

	return h, nil
}

func (h eventDispatch) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		// Ensure a service identity is present so eventbus handlers can act
		ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

		vars, err := mapBody(r, h.params.Mapping)
		if err != nil {
			return err
		}

		vars["configuredConnectionID"] = h.params.ConfiguredConnectionID

		ev := sysEvent.ConnectionWebhookEvent(h.params.ConnectionID, h.params.ConfiguredConnectionID, h.params.EventType, vars)

		return eventbus.Service().WaitFor(context.WithoutCancel(ctx), ev)
	}
}

// mapBody maps the HTTP response into our internal requested structure
func mapBody(r *http.Request, mapping map[string][]string) (map[string]any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("eventDispatch: could not read request body: %w", err)
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))

	var raw map[string]any
	if len(body) > 0 {
		_ = json.Unmarshal(body, &raw)
	}
	if raw == nil {
		raw = make(map[string]any)
	}

	vars := make(map[string]any, len(mapping))
	for target, selector := range mapping {
		vars[target], _, _ = j7s.GetBySlicePath(raw, selector)
	}

	return vars, nil
}
