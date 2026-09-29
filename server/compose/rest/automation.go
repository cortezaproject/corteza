package rest

import (
	"context"

	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/service/event"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/corredor"
	"github.com/crusttech/human/server/pkg/errors"
)

type (
	Automation struct{}
)

func (Automation) New() *Automation {
	return &Automation{}
}

// The catalogue names every script, what fires it and who may run it, and the
// bundle is the client scripts' own source, so both are for the people this
// instance knows. A script's own security still decides who may run it.
func (ctrl *Automation) List(ctx context.Context, r *request.AutomationList) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.Unauthorized("cannot list Corredor scripts")
	}

	return corredor.GenericListHandler(
		ctx,
		corredor.Service(),
		corredor.Filter{
			ResourceTypePrefixes: r.ResourceTypePrefixes,
			ExcludeInvalid:       r.ExcludeInvalid,
			ResourceTypes:        r.ResourceTypes,
			EventTypes:           r.EventTypes,
			ExcludeServerScripts: r.ExcludeServerScripts,
			ExcludeClientScripts: r.ExcludeClientScripts,
		},
		"compose",
	)
}

func (ctrl *Automation) Bundle(ctx context.Context, r *request.AutomationBundle) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.Unauthorized("cannot read the Corredor client script bundle")
	}

	return corredor.GenericBundleHandler(
		ctx,
		corredor.Service(),
		r.Bundle,
		r.Type,
		r.Ext,
	)
}

func (ctrl *Automation) TriggerScript(ctx context.Context, r *request.AutomationTriggerScript) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.Unauthorized("cannot run a Corredor script")
	}

	if r.Args == nil {
		// Script results are returned through args
		r.Args = map[string]interface{}{}
	}

	sArgs := corredor.ExtendScriptArgs(event.ComposeOnManual(), r.Args)
	if err := corredor.Service().Exec(ctx, r.Script, sArgs); err != nil {
		return nil, err
	}
	return r.Args, nil
}
