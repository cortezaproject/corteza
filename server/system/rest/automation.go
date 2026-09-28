package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/corredor"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/service/event"
)

type (
	Automation struct {
		ac automationAccessController
	}

	automationAccessController interface {
		CanSearchCorredorScripts(context.Context) bool
	}
)

func (Automation) New() *Automation {
	return &Automation{
		ac: service.DefaultAccessControl,
	}
}

func (ctrl *Automation) List(ctx context.Context, r *request.AutomationList) (interface{}, error) {
	if !ctrl.ac.CanSearchCorredorScripts(ctx) {
		return nil, errors.Unauthorized("not allowed to search or list Corredor scripts")
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
		"system",
	)
}

// The bundle is the client scripts' own source and the trigger runs one, so
// both are for the people this instance knows; listing the catalogue asks for
// more than that, and a script's own security still decides who may run it.
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

	sArgs := corredor.ExtendScriptArgs(event.SystemOnManual(), r.Args)
	if err := corredor.Service().Exec(ctx, r.Script, sArgs); err != nil {
		return nil, err
	}
	return r.Args, nil
}
