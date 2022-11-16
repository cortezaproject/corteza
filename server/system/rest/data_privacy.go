package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	DataPrivacy struct {
		privacy service.DataPrivacyService
		ac      privacyAccessController
	}

	privacyAccessController interface {
		CanGrant(context.Context) bool
	}

	privacyConnectionSetPayload struct {
		Filter types.DalConnectionFilter     `json:"filter"`
		Set    types.PrivacyDalConnectionSet `json:"set"`
	}
)

func (DataPrivacy) New() *DataPrivacy {
	return &DataPrivacy{
		privacy: service.DefaultDataPrivacy,
		ac:      service.DefaultAccessControl,
	}
}

func (ctrl DataPrivacy) ConnectionList(ctx context.Context, r *request.DataPrivacyConnectionList) (interface{}, error) {
	var (
		err error
		set types.PrivacyDalConnectionSet

		f = types.DalConnectionFilter{
			ConnectionID: payload.ParseUint64s(r.ConnectionID),
			Handle:       r.Handle,
			Type:         r.Type,

			Deleted: r.Deleted,
		}
	)

	set, f, err = ctrl.privacy.FindConnections(ctx, f)
	return ctrl.makeFilterConnectionPayload(ctx, set, f, err)
}

func (ctrl DataPrivacy) makeFilterConnectionPayload(_ context.Context, rr types.PrivacyDalConnectionSet, f types.DalConnectionFilter, err error) (*privacyConnectionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		rr = make([]*types.PrivacyDalConnection, 0)
	}

	return &privacyConnectionSetPayload{Filter: f, Set: rr}, nil
}
