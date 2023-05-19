package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	DalSchemaAlteration struct {
		svc           alterationService
		federationSvc federationNodeService

		alterationAc alterationAccessController
	}

	alterationPayload struct {
		*types.DalSchemaAlteration

		CanGrant                  bool `json:"canGrant"`
		CanUpdateSchemaAlteration bool `json:"canUpdateSchemaAlteration"`
		CanDeleteSchemaAlteration bool `json:"canDeleteSchemaAlteration"`
		CanManageDalConfig        bool `json:"canManageDalConfig"`
	}

	alterationSetPayload struct {
		Filter types.DalSchemaAlterationFilter `json:"filter"`
		Set    []*alterationPayload            `json:"set"`
	}

	alterationAccessController interface {
		CanGrant(context.Context) bool
		// @todo?
	}

	alterationService interface {
		Search(ctx context.Context, filter types.DalSchemaAlterationFilter) (types.DalSchemaAlterationSet, types.DalSchemaAlterationFilter, error)
		FindByID(ctx context.Context, ID uint64) (*types.DalSchemaAlteration, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
	}
)

func (DalSchemaAlteration) New() *DalSchemaAlteration {
	return &DalSchemaAlteration{
		svc: service.DefaultDalSchemaAlteration,

		alterationAc: service.DefaultAccessControl,
	}
}

func (ctrl DalSchemaAlteration) List(ctx context.Context, r *request.DalSchemaAlterationList) (interface{}, error) {
	var (
		err error
		set types.DalSchemaAlterationSet

		f = types.DalSchemaAlterationFilter{
			AlterationID: r.AlterationID,
			BatchID:      r.BatchID,
			Kind:         r.Kind,

			Deleted:   filter.State(r.Deleted),
			Completed: filter.State(r.Completed),
		}
	)

	if f.Deleted == 0 {
		f.Deleted = filter.StateExcluded
	}
	if f.Completed == 0 {
		f.Completed = filter.StateExcluded
	}

	f.IncTotal = r.IncTotal

	set, f, err = ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl DalSchemaAlteration) Read(ctx context.Context, r *request.DalSchemaAlterationRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.AlterationID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl DalSchemaAlteration) Delete(ctx context.Context, r *request.DalSchemaAlterationDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.AlterationID)
}

func (ctrl DalSchemaAlteration) Undelete(ctx context.Context, r *request.DalSchemaAlterationUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.AlterationID)
}

func (ctrl DalSchemaAlteration) makePayload(ctx context.Context, res *types.DalSchemaAlteration, err error) (*alterationPayload, error) {
	if err != nil || res == nil {
		return nil, err
	}

	pl := &alterationPayload{
		DalSchemaAlteration: res,

		CanGrant: ctrl.alterationAc.CanGrant(ctx),
	}

	return pl, nil
}

func (ctrl DalSchemaAlteration) makeFilterPayload(ctx context.Context, rr types.DalSchemaAlterationSet, f types.DalSchemaAlterationFilter, err error) (*alterationSetPayload, error) {
	if err != nil {
		return nil, err
	}

	out := &alterationSetPayload{Filter: f, Set: make([]*alterationPayload, len(rr))}
	for i := range rr {
		out.Set[i], _ = ctrl.makePayload(ctx, rr[i], nil)
	}

	return out, nil
}
