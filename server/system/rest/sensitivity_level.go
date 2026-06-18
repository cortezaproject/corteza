package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	DalSensitivityLevel struct {
		dalSensitivityLevel dalSensitivityLevelService
	}

	sensitivityLevelSetPayload struct {
		Filter types.DalSensitivityLevelFilter `json:"filter"`
		Set    []*sensitivityLevelPayload      `json:"set"`
	}

	sensitivityLevelPayload struct {
		*types.DalSensitivityLevel
	}

	dalSensitivityLevelService interface {
		FindByID(ctx context.Context, ID uint64) (*types.DalSensitivityLevel, error)
		Create(ctx context.Context, new *types.DalSensitivityLevel) (*types.DalSensitivityLevel, error)
		Update(ctx context.Context, upd *types.DalSensitivityLevel) (*types.DalSensitivityLevel, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.DalSensitivityLevelFilter) (types.DalSensitivityLevelSet, types.DalSensitivityLevelFilter, error)
	}
)

func (DalSensitivityLevel) New() *DalSensitivityLevel {
	return &DalSensitivityLevel{
		dalSensitivityLevel: service.DefaultDalSensitivityLevel,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl DalSensitivityLevel) makeFilter(ctx context.Context, r *request.DalSensitivityLevelList) (types.DalSensitivityLevelFilter, error) {
	f := types.DalSensitivityLevelFilter{
		DalSensitivityLevelID: r.SensitivityLevelID,

		Deleted: filter.State(r.Deleted),
	}

	if f.Deleted == 0 {
		f.Deleted = filter.StateExcluded
	}

	f.IncTotal = r.IncTotal

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl DalSensitivityLevel) beforeCreate(ctx context.Context, res *types.DalSensitivityLevel, r *request.DalSensitivityLevelCreate) error {
	res.Meta = r.Meta
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl DalSensitivityLevel) beforeUpdate(ctx context.Context, res *types.DalSensitivityLevel, r *request.DalSensitivityLevelUpdate) error {
	res.Meta = r.Meta
	return nil
}

func (ctrl DalSensitivityLevel) makePayload(ctx context.Context, res *types.DalSensitivityLevel, err error) (*sensitivityLevelPayload, error) {
	if err != nil || res == nil {
		return nil, err
	}

	pl := &sensitivityLevelPayload{res}

	return pl, nil
}

func (ctrl DalSensitivityLevel) makeFilterPayload(ctx context.Context, rr types.DalSensitivityLevelSet, f types.DalSensitivityLevelFilter, err error) (*sensitivityLevelSetPayload, error) {
	if err != nil {
		return nil, err
	}

	out := &sensitivityLevelSetPayload{Filter: f, Set: make([]*sensitivityLevelPayload, len(rr))}

	for i := range rr {
		out.Set[i], _ = ctrl.makePayload(ctx, rr[i], nil)
	}
	return out, nil
}

func (ctrl DalSensitivityLevel) serve(ctx context.Context, fn string, archive io.ReadSeeker, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename="+fn)

		http.ServeContent(w, req, fn, time.Now(), archive)
	}, nil
}
