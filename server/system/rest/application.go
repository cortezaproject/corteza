package rest

import (
	"context"
	"strconv"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/corredor"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/flag"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
)

type (
	Application struct {
		application applicationService
		attachment  service.AttachmentService
		ac          applicationAccessController
	}

	applicationService interface {
		FindByID(ctx context.Context, ID uint64) (app *types.Application, err error)
		Search(ctx context.Context, filter types.ApplicationFilter) (aa types.ApplicationSet, f types.ApplicationFilter, err error)
		Create(ctx context.Context, new *types.Application) (app *types.Application, err error)
		Update(ctx context.Context, upd *types.Application) (app *types.Application, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		UndeleteByID(ctx context.Context, ID uint64) (err error)
		Reorder(ctx context.Context, order []uint64) (err error)
		Flag(ctx context.Context, app *types.Application, ownedBy uint64, f string) error
		Unflag(ctx context.Context, app *types.Application, ownedBy uint64, f string) error
	}

	applicationAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateApplication(context.Context, *types.Application) bool
		CanDeleteApplication(context.Context, *types.Application) bool
	}

	applicationPayload struct {
		*types.Application

		CanGrant             bool `json:"canGrant"`
		CanUpdateApplication bool `json:"canUpdateApplication"`
		CanDeleteApplication bool `json:"canDeleteApplication"`
	}

	applicationSetPayload struct {
		Filter types.ApplicationFilter `json:"filter"`
		Set    []*applicationPayload   `json:"set"`
	}
)

func (Application) New() *Application {
	return &Application{
		application: service.DefaultApplication,
		attachment:  service.DefaultAttachment,
		ac:          service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl Application) makeFilter(ctx context.Context, r *request.ApplicationList) (types.ApplicationFilter, error) {
	var (
		err error
		f   = types.ApplicationFilter{
			Name:     r.Name,
			Query:    r.Query,
			Labels:   r.Labels,
			Flags:    r.Flags,
			IncFlags: r.IncFlags,

			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl Application) beforeCreate(ctx context.Context, res *types.Application, r *request.ApplicationCreate) error {
	res.Labels = r.Labels

	if r.Unify != nil {
		res.Unify = &types.ApplicationUnify{}
		if err := r.Unify.Unmarshal(res.Unify); err != nil {
			return err
		}
	}

	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl Application) beforeUpdate(ctx context.Context, res *types.Application, r *request.ApplicationUpdate) error {
	res.Labels = r.Labels

	if r.Unify != nil {
		res.Unify = &types.ApplicationUnify{}
		if err := r.Unify.Unmarshal(res.Unify); err != nil {
			return err
		}
	}

	return nil
}

// Read is kept as a custom method because the generated Read cannot
// reproduce the flag.Load post-load enrichment: it depends on r.IncFlags
// (only available on the request) and cannot be folded into makePayload,
// which has no access to the request. The generator's Read does not call
// an afterRead hook, so this enrichment must stay here to preserve behavior.
func (ctrl *Application) Read(ctx context.Context, r *request.ApplicationRead) (interface{}, error) {
	app, err := ctrl.application.FindByID(ctx, r.ApplicationID)
	if err != nil {
		return ctrl.makePayload(ctx, app, err)
	}

	err = flag.Load(ctx, service.DefaultStore, r.IncFlags, auth.GetIdentityFromContext(ctx).Identity(), app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Application) Upload(ctx context.Context, r *request.ApplicationUpload) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	a, err := ctrl.attachment.CreateApplicationAttachment(
		ctx,
		r.Upload.Filename,
		r.Upload.Size,
		file,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl *Application) TriggerScript(ctx context.Context, r *request.ApplicationTriggerScript) (rsp interface{}, err error) {
	var (
		application *types.Application
	)

	if application, err = ctrl.application.FindByID(ctx, r.ApplicationID); err != nil {
		return
	}

	// @todo implement same behaviour as we have on record - Application+oldApplication
	err = corredor.Service().Exec(ctx, r.Script, corredor.ExtendScriptArgs(event.ApplicationOnManual(application, application), r.Args))
	return application, err
}

func (ctrl *Application) Reorder(ctx context.Context, r *request.ApplicationReorder) (interface{}, error) {
	order := make([]uint64, len(r.ApplicationIDs))

	for i, aid := range r.ApplicationIDs {
		parsed, err := strconv.ParseUint(aid, 10, 64)
		if err != nil {
			return nil, err
		}

		order[i] = parsed
	}

	return api.OK(), ctrl.application.Reorder(ctx, order)
}

func (ctrl *Application) FlagCreate(ctx context.Context, r *request.ApplicationFlagCreate) (interface{}, error) {
	app, err := ctrl.application.FindByID(ctx, r.ApplicationID)
	if err != nil {
		return nil, err
	}

	return api.OK(), ctrl.application.Flag(ctx, app, r.OwnedBy, r.Flag)

}

func (ctrl *Application) FlagDelete(ctx context.Context, r *request.ApplicationFlagDelete) (interface{}, error) {
	app, err := ctrl.application.FindByID(ctx, r.ApplicationID)
	if err != nil {
		return nil, err
	}

	return api.OK(), ctrl.application.Unflag(ctx, app, r.OwnedBy, r.Flag)
}

func (ctrl Application) makePayload(ctx context.Context, m *types.Application, err error) (*applicationPayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &applicationPayload{
		Application: m,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateApplication: ctrl.ac.CanUpdateApplication(ctx, m),
		CanDeleteApplication: ctrl.ac.CanDeleteApplication(ctx, m),
	}, nil
}

func (ctrl Application) makeFilterPayload(ctx context.Context, nn types.ApplicationSet, f types.ApplicationFilter, err error) (*applicationSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &applicationSetPayload{Filter: f, Set: make([]*applicationPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
