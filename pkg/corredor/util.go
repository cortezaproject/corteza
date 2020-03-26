package corredor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/pkg/errors"
	"net/http"
)

type (
	automationListSetPayload struct {
		Filter Filter    `json:"filter"`
		Set    []*Script `json:"set"`
	}
)

// mapManualTriggers removes all manual triggers from the list of script's triggers
//
// and returns a hash map with resources from these manual triggers
func mapManualTriggers(script *ServerScript) map[string]bool {
	var (
		hash = make(map[string]bool)
	)

	for i := range script.Triggers {
		// We're modifying trigger in the loop,
		// so let's make a copy we can play with
		trigger := script.Triggers[i]

		if slice.HasString(trigger.EventTypes, onManualEventType) {
			for _, res := range trigger.ResourceTypes {
				hash[res] = true
			}
		}
	}

	return hash
}

// converts trigger's constraint to eventbus' constraint options
func triggerToHandlerOps(t *Trigger) (oo []eventbus.HandlerRegOp, err error) {
	if len(t.ResourceTypes) == 0 {
		return nil, fmt.Errorf("can not generate event handler without at least one resource")
	}

	if len(t.EventTypes) == 0 {
		return nil, fmt.Errorf("can not generate event handler without at least one events")
	}

	// Make a copy of event types slice so that we do not modify it
	types := slice.PluckString(t.EventTypes, onManualEventType)
	if len(types) == 0 && len(t.EventTypes) > 0 {
		return
	}

	oo = append(oo, eventbus.For(t.ResourceTypes...))
	oo = append(oo, eventbus.On(types...))

	if cc, err := constraintsToHandlerOps(t.Constraints); err != nil {
		return nil, err
	} else {
		oo = append(oo, cc...)
	}

	return
}

func constraintsToHandlerOps(cc []*TConstraint) (oo []eventbus.HandlerRegOp, err error) {
	for _, raw := range cc {
		if c, err := eventbus.ConstraintMaker(raw.Name, raw.Op, raw.Value...); err != nil {
			return nil, errors.Wrap(err, "can not generate constraints")
		} else {
			oo = append(oo, eventbus.Constraint(c))
		}
	}

	return
}

// encode adds entry (with json encoded value) to hash map
// used to prepare data for transmission
func encodeArguments(args map[string]string, key string, val interface{}) (err error) {
	var tmp []byte

	if tmp, err = json.Marshal(val); err != nil {
		return
	}

	args[key] = string(tmp)
	return
}

// GenericListHandler returns filtered list of scripts
func GenericListHandler(ctx context.Context, svc *service, f Filter, resourcePrefix string) (p *automationListSetPayload, err error) {
	f.procRTPrefixes(resourcePrefix)

	p = &automationListSetPayload{}
	p.Set, p.Filter, err = svc.Find(ctx, f)
	return p, err
}

func GenericBundleHandler(ctx context.Context, svc *service, bundleName, bundleType, ext string) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Serve bundle directly for now
		bundle := svc.GetBundle(ctx, bundleName, bundleType)
		if bundle == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(bundle.Code))
	}, nil
}
