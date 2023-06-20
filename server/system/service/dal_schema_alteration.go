package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	dalSchemaAlteration struct {
		actionlog actionlog.Recorder
		ac        dalSchemaAlterationAccessController
		dal       dalAltManager

		store store.Storer
	}

	dalAltManager interface {
		ApplyAlteration(ctx context.Context, alts ...*dal.Alteration) (errs []error, err error)
	}

	dalSchemaAlterationAccessController interface {
	}
)

func DalSchemaAlteration(dal dalAltManager) *dalSchemaAlteration {
	return &dalSchemaAlteration{
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
		dal:       dal,
	}
}

func (svc dalSchemaAlteration) FindByID(ctx context.Context, dalSchemaAlterationID uint64) (a *types.DalSchemaAlteration, err error) {
	var (
		uaProps = &dalSchemaAlterationActionProps{dalSchemaAlteration: &types.DalSchemaAlteration{ID: dalSchemaAlterationID}}
	)

	err = func() error {
		a, err = loadDalSchemaAlteration(ctx, svc.store, dalSchemaAlterationID)
		if err != nil {
			return err
		}

		uaProps.setDalSchemaAlteration(a)

		// if !svc.ac.CanReadDalSchemaAlteration(ctx, u) {
		// 	return DalSchemaAlterationErrNotAllowedToRead()
		// }

		return nil
	}()

	return a, svc.recordAction(ctx, uaProps, DalSchemaAlterationActionLookup, err)
}

// Search interacts with backend storage and
//
// @todo rename to Search() for consistency
func (svc dalSchemaAlteration) Search(ctx context.Context, filter types.DalSchemaAlterationFilter) (aa types.DalSchemaAlterationSet, f types.DalSchemaAlterationFilter, err error) {
	var (
		uaProps = &dalSchemaAlterationActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	// if !svc.ac.CanReadDalSchemaAlteration(ctx, res) {
	// 	return false, nil
	// }

	// 	return true, nil
	// }

	err = func() error {
		// if !svc.ac.CanSearchDalSchemaAlterations(ctx) {
		// 	return DalSchemaAlterationErrNotAllowedToSearch()
		// }

		aa, f, err = store.SearchDalSchemaAlterations(ctx, svc.store, filter)
		return err
	}()

	return aa, f, svc.recordAction(ctx, uaProps, DalSchemaAlterationActionSearch, err)
}

func (svc dalSchemaAlteration) DeleteByID(ctx context.Context, dalSchemaAlterationID uint64) (err error) {
	var (
		u       *types.DalSchemaAlteration
		uaProps = &dalSchemaAlterationActionProps{dalSchemaAlteration: &types.DalSchemaAlteration{ID: dalSchemaAlterationID}}
	)

	err = func() (err error) {
		if u, err = loadDalSchemaAlteration(ctx, svc.store, dalSchemaAlterationID); err != nil {
			return
		}

		// if !svc.ac.CanDeleteDalSchemaAlteration(ctx, u) {
		// 	return DalSchemaAlterationErrNotAllowedToDelete()
		// }

		// if err = svc.eventbus.WaitFor(ctx, event.DalSchemaAlterationBeforeDelete(nil, u)); err != nil {
		// 	return
		// }

		u.DeletedAt = now()
		if err = store.UpdateDalSchemaAlteration(ctx, svc.store, u); err != nil {
			return
		}

		// _ = svc.eventbus.WaitFor(ctx, event.DalSchemaAlterationAfterDelete(nil, u))
		return nil
	}()

	return svc.recordAction(ctx, uaProps, DalSchemaAlterationActionDelete, err)
}

func (svc dalSchemaAlteration) UndeleteByID(ctx context.Context, dalSchemaAlterationID uint64) (err error) {
	var (
		u       *types.DalSchemaAlteration
		uaProps = &dalSchemaAlterationActionProps{dalSchemaAlteration: &types.DalSchemaAlteration{ID: dalSchemaAlterationID}}
	)

	err = func() (err error) {
		if u, err = loadDalSchemaAlteration(ctx, svc.store, dalSchemaAlterationID); err != nil {
			return
		}

		uaProps.setDalSchemaAlteration(u)

		// if err = uniqueDalSchemaAlterationCheck(ctx, svc.store, u); err != nil {
		// 	return err
		// }

		// if !svc.ac.CanDeleteDalSchemaAlteration(ctx, u) {
		// 	return DalSchemaAlterationErrNotAllowedToDelete()
		// }

		u.DeletedAt = nil
		if err = store.UpdateDalSchemaAlteration(ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, DalSchemaAlterationActionUndelete, err)
}

// ModelAlterations returns all non deleted, non completed, and non dismissed alterations for the given model
func (svc dalSchemaAlteration) ModelAlterations(ctx context.Context, m *dal.Model) (out []*dal.Alteration, err error) {
	// @todo boilerplate code around this

	aux, _, err := store.SearchDalSchemaAlterations(ctx, svc.store, types.DalSchemaAlterationFilter{
		Resource:  []string{m.Resource},
		Deleted:   filter.StateExcluded,
		Completed: filter.StateExcluded,
		Dismissed: filter.StateExcluded,
	})
	if err != nil {
		return nil, err
	}

	for _, a := range aux {
		t := &dal.Alteration{
			ID:           a.ID,
			BatchID:      a.BatchID,
			DependsOn:    a.DependsOn,
			Resource:     a.Resource,
			ResourceType: a.ResourceType,
			ConnectionID: a.ConnectionID,
		}

		switch a.Kind {
		case "attributeAdd":
			t.AttributeAdd = a.Params.AttributeAdd
		case "attributeDelete":
			t.AttributeDelete = a.Params.AttributeDelete
		case "attributeReType":
			t.AttributeReType = a.Params.AttributeReType
		case "attributeReEncode":
			t.AttributeReEncode = a.Params.AttributeReEncode
		case "modelAdd":
			t.ModelAdd = a.Params.ModelAdd
		case "modelDelete":
			t.ModelDelete = a.Params.ModelDelete
		}

		out = append(out, t)
	}

	return
}

// @todo we can probably do some diffing here to make it more efficient; it'll do for now
func (svc dalSchemaAlteration) SetAlterations(ctx context.Context, m *dal.Model, stale []*dal.Alteration, aa ...*dal.Alteration) (err error) {
	// @todo boilerplate code around this

	c := svc.dal.GetConnectionByID(0)
	if m.ConnectionID == c.ID && m.Ident == "compose_record" {
		err = fmt.Errorf("cannot set alterations for default schema")
		return
	}

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		// Firstly get rid of the old ones
		aux := make([]*types.DalSchemaAlteration, len(stale))
		for i, a := range stale {
			aux[i] = &types.DalSchemaAlteration{
				ID: a.ID,
			}
		}
		err = store.DeleteDalSchemaAlteration(ctx, s, aux...)
		if err != nil {
			return
		}

		// Now add the new ones
		cvt := make(types.DalSchemaAlterationSet, len(aa))
		n := *now()
		u := intAuth.GetIdentityFromContext(ctx).Identity()
		for i, a := range aa {
			t := &types.DalSchemaAlteration{
				ID:           a.ID,
				BatchID:      a.BatchID,
				DependsOn:    a.DependsOn,
				ConnectionID: a.ConnectionID,
				Resource:     a.Resource,
				ResourceType: a.ResourceType,

				Params: &types.DalSchemaAlterationParams{},
			}

			// @todo we'd need to merge this with the old one probably just to preserve the metadata
			t.ID = nextID()
			t.CreatedAt = n
			t.CreatedBy = u

			switch {
			case a.AttributeAdd != nil:
				t.Kind = "attributeAdd"
				t.Params.AttributeAdd = a.AttributeAdd

			case a.AttributeDelete != nil:
				t.Kind = "attributeDelete"
				t.Params.AttributeDelete = a.AttributeDelete

			case a.AttributeReType != nil:
				t.Kind = "attributeReType"
				t.Params.AttributeReType = a.AttributeReType

			case a.AttributeReEncode != nil:
				t.Kind = "attributeReEncode"
				t.Params.AttributeReEncode = a.AttributeReEncode

			case a.ModelAdd != nil:
				t.Kind = "modelAdd"
				t.Params.ModelAdd = a.ModelAdd

			case a.ModelDelete != nil:
				t.Kind = "modelDelete"
				t.Params.ModelDelete = a.ModelDelete
			}

			cvt[i] = t
		}

		return store.UpsertDalSchemaAlteration(ctx, svc.store, cvt...)
	})
}

func (svc dalSchemaAlteration) Apply(ctx context.Context, ids ...uint64) (err error) {
	// @todo boilerplate (RBAC and such); we might have special RBAC rules for this;
	// originally, we wanted to hook into ComposeModule resource (or any resource that defined a model)
	aux, _, err := store.SearchDalSchemaAlterations(ctx, svc.store, types.DalSchemaAlterationFilter{
		AlterationID: id.Strings(ids...),
	})
	if err != nil {
		return
	}

	aa, err := svc.toPkgAlterations(ctx, svc.appliableAlterations(aux...)...)
	if err != nil {
		return
	}

	errors, err := svc.dal.ApplyAlteration(ctx, aa...)
	if err != nil {
		return
	}

	for i, e := range errors {
		if e != nil {
			aux[i].Error = e.Error()
		} else {
			aux[i].CompletedAt = now()
			aux[i].CompletedBy = intAuth.GetIdentityFromContext(ctx).Identity()
		}
	}

	return store.UpdateDalSchemaAlteration(ctx, svc.store, aux...)
}

func (svc dalSchemaAlteration) Dismiss(ctx context.Context, ids ...uint64) (err error) {
	// @todo boilerplate (RBAC and such); we might have special RBAC rules for this;
	// originally, we wanted to hook into ComposeModule resource (or any resource that defined a model)

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		alt, _, err := store.SearchDalSchemaAlterations(ctx, svc.store, types.DalSchemaAlterationFilter{
			AlterationID: id.Strings(ids...),
		})
		if err != nil {
			return
		}

		alt = svc.appliableAlterations(alt...)

		for _, a := range alt {
			a.Error = ""
			a.DismissedAt = now()
			a.DismissedBy = intAuth.GetIdentityFromContext(ctx).Identity()
		}

		return store.UpdateDalSchemaAlteration(ctx, s, alt...)
	})
}

func (svc dalSchemaAlteration) appliableAlterations(aa ...*types.DalSchemaAlteration) (out types.DalSchemaAlterationSet) {
	out = make(types.DalSchemaAlterationSet, 0, len(aa))

	altIndex := make(map[uint64]*types.DalSchemaAlteration, len(aa))
	for _, a := range aa {
		altIndex[a.ID] = a
	}

	for _, a := range aa {
		// Already completed
		if a.CompletedAt != nil {
			continue
		}

		// Dismissed for manual thing
		if a.DismissedAt != nil {
			continue
		}

		if a.DependsOn != 0 {
			// Check if the dependency is completed or inside of the index
			if _, ok := altIndex[a.DependsOn]; !ok {
				continue
			}
		}

		out = append(out, a)
	}

	return
}

func loadDalSchemaAlteration(ctx context.Context, s store.DalSchemaAlterations, ID uint64) (res *types.DalSchemaAlteration, err error) {
	if ID == 0 {
		return nil, DalSchemaAlterationErrInvalidID()
	}

	if res, err = store.LookupDalSchemaAlterationByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, DalSchemaAlterationErrNotFound()
	}

	return
}

func (svc dalSchemaAlteration) toPkgAlterations(ctx context.Context, aa ...*types.DalSchemaAlteration) (out []*dal.Alteration, err error) {
	out = make([]*dal.Alteration, len(aa))
	for i, a := range aa {
		t := &dal.Alteration{
			ID:           a.ID,
			BatchID:      a.BatchID,
			DependsOn:    a.DependsOn,
			ConnectionID: a.ConnectionID,
			Resource:     a.Resource,
			ResourceType: a.ResourceType,
		}

		switch a.Kind {
		case "attributeAdd":
			t.AttributeAdd = a.Params.AttributeAdd
		case "attributeDelete":
			t.AttributeDelete = a.Params.AttributeDelete
		case "attributeReType":
			t.AttributeReType = a.Params.AttributeReType
		case "attributeReEncode":
			t.AttributeReEncode = a.Params.AttributeReEncode
		case "modelAdd":
			t.ModelAdd = a.Params.ModelAdd
		case "modelDelete":
			t.ModelDelete = a.Params.ModelDelete
		}

		out[i] = t
	}

	return
}
