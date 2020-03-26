package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service/event"
)

var _ = errors.Wrap

type (
	Automation struct{}
)

func (Automation) New() *Automation {
	return &Automation{}
}

func (ctrl *Automation) List(ctx context.Context, r *request.AutomationList) (interface{}, error) {
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
	return resputil.OK(), corredor.Service().Exec(ctx, r.Script, event.SystemOnManual())
}
