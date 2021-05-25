package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	actionlogUserSearcher interface {
		Find(context.Context, types.UserFilter) (types.UserSet, types.UserFilter, error)
	}

	Actionlog struct {
		actionSvc actionlog.Recorder
		userSvc   actionlogUserSearcher
	}

	// Extend actionlog.Action so we can
	// provide user's email
	actionlogActionPayload struct {
		*actionlog.Action
		Actor string `json:"actor,omitempty"`
	}

	actionlogPayload struct {
		Filter actionlog.Filter          `json:"filter"`
		Set    []*actionlogActionPayload `json:"set"`
	}
)

func (Actionlog) New() *Actionlog {
	return &Actionlog{
		actionSvc: service.DefaultActionlog,
		userSvc:   service.DefaultUser,
	}
}

func (ctrl *Actionlog) List(ctx context.Context, r *request.ActionlogList) (interface{}, error) {
	var (
		err error
		f   = actionlog.Filter{
			FromTimestamp:  r.From,
			ToTimestamp:    r.To,
			BeforeActionID: r.BeforeActionID,
			ActorID:        payload.ParseUint64s(r.ActorID),
			Resource:       r.Resource,
			Action:         r.Action,
			Limit:          r.Limit,
		}
	)

	ee, f, err := ctrl.actionSvc.Find(ctx, f)

	return ctrl.makeFilterPayload(ctx, ee, f, err)
}

func (ctrl Actionlog) makeFilterPayload(ctx context.Context, ee []*actionlog.Action, f actionlog.Filter, err error) (*actionlogPayload, error) {
	if err != nil {
		return nil, err
	}

	var (
		pp = make([]*actionlogActionPayload, len(ee))
	)

	// Remap events to payload structs
	for e := range ee {
		pp[e] = &actionlogActionPayload{Action: ee[e]}
	}

	err = userPreloader(
		ctx,
		ctrl.userSvc,
		func(c chan uint64) {
			for e := range ee {
				c <- ee[e].ActorID
			}

			close(c)
		},
		types.UserFilter{
			Deleted:   filter.StateInclusive,
			Suspended: filter.StateInclusive,
		},
		func(u *types.User) error {
			for p := range pp {
				if pp[p].ActorID == u.ID {
					pp[p].Actor = u.Name
					if pp[p].Actor == "" {
						pp[p].Actor = u.Email
					}
				}
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &actionlogPayload{Filter: f, Set: pp}, nil
}

// Preloader collects all ids of users, loads them and sets them back
//
//
// @todo this kind of preloader is useful and can be implemented in bunch
//       of places and replace old code
func userPreloader(ctx context.Context, svc actionlogUserSearcher, g func(chan uint64), f types.UserFilter, s func(*types.User) error) error {
	var (
		// channel that will collect the IDs in the getter
		ch = make(chan uint64, 0)

		// unique index for IDs
		unq = make(map[uint64]bool)
	)

	// Reset the collection of the IDs
	f.UserID = make([]uint64, 0)

	// Call getter and collect the IDs
	go g(ch)

rangeLoop:
	for {
		select {
		case <-ctx.Done():
			close(ch)
			break rangeLoop
		case id, ok := <-ch:
			if !ok {
				// Channel closed
				break rangeLoop
			}

			if !unq[id] {
				unq[id] = true
				f.UserID = append(f.UserID, id)
			}
		}

	}

	// Load all users (even if deleted, suspended) from the given list of IDs
	uu, _, err := svc.Find(ctx, f)

	if err != nil {
		return err
	}

	return uu.Walk(s)
}
