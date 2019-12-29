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

type Automation struct {
	// xxx service.XXXService
}

func (Automation) New() *Automation {
	return &Automation{}
}

func (ctrl *Automation) List(ctx context.Context, r *request.AutomationList) (interface{}, error) {
	f := corredor.ManualScriptFilter{
		ResourceTypes:        r.ResourceTypes,
		ExcludeServerScripts: r.ExcludeServerScripts,
		ExcludeClientScripts: r.ExcludeClientScripts,
	}

	f.PrefixResource("system")

	scripts, _, err := corredor.Service().FindOnManual(f)

	return scripts, err
}

func (ctrl *Automation) Trigger(ctx context.Context, r *request.AutomationTrigger) (interface{}, error) {
	return resputil.OK(), corredor.Service().ExecOnManual(ctx, r.Script, event.SystemOnManual())
}
