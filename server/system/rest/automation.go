package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
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

func (ctrl *Automation) Bundle(ctx context.Context, r *request.AutomationBundle) (interface{}, error) {
	return corredor.GenericBundleHandler(
		ctx,
		corredor.Service(),
		r.Bundle,
		r.Type,
		r.Ext,
	)
}

func (ctrl *Automation) TriggerScript(ctx context.Context, r *request.AutomationTriggerScript) (interface{}, error) {
	return api.OK(), corredor.Service().Exec(ctx, r.Script, corredor.ExtendScriptArgs(event.SystemOnManual(), r.Args))
}
