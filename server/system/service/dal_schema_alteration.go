package service

import (
	"context"
	"fmt"

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
		GetConnectionByID(ID uint64) *dal.ConnectionWrap
		ApplyAlteration(ctx context.Context, alts ...*dal.Alteration) (errs []error, err error)
		ReloadModel(ctx context.Context, currentAlts []*dal.Alteration, model *dal.Model) (newAlts []*dal.Alteration, err error)
		FindModelByResourceIdent(connectionID uint64, resourceType, resourceIdent string) *dal.Model
	}

	dalSchemaAlterationAccessController interface {
		CanManageDalSchemaAlterations(ctx context.Context) bool
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

		if !svc.ac.CanManageDalSchemaAlterations(ctx) {
			return DalSchemaAlterationErrNotAllowedToManage()
		}

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

	err = func() error {
		if !svc.ac.CanManageDalSchemaAlterations(ctx) {
			return DalSchemaAlterationErrNotAllowedToManage()
		}

		aa, f, err = store.SearchDalSchemaAlterations(ctx, svc.store, filter)
		return err
	}()

	return aa, f, svc.recordAction(ctx, uaProps, DalSchemaAlterationActionSearch, err)
}

// ModelAlterations returns all non deleted, non completed, and non dismissed alterations for the given model
func (svc dalSchemaAlteration) ModelAlterations(ctx context.Context, m *dal.Model) (out []*dal.Alteration, err error) {
	return svc.modelAlterations(ctx, svc.store, m)
}

func (svc dalSchemaAlteration) modelAlterations(ctx context.Context, s store.Storer, m *dal.Model) (out []*dal.Alteration, err error) {
	aux, _, err := store.SearchDalSchemaAlterations(ctx, s, types.DalSchemaAlterationFilter{
		Resource:     []string{m.Resource},
		ResourceType: m.ResourceType,
		Deleted:      filter.StateExcluded,
		Completed:    filter.StateExcluded,
		Dismissed:    filter.StateExcluded,
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

// SetAlterations updates the DB state to reflect the given alterations
//
// This function should only be invoked by internal proceses so it doesn't need
// to check for permissions.
func (svc dalSchemaAlteration) SetAlterations(ctx context.Context, s store.Storer, m *dal.Model, stale []*dal.Alteration, aa ...*dal.Alteration) (err error) {
	if len(stale)+len(aa) == 0 {
		return
	}

	var (
		n = *now()
		u = intAuth.GetIdentityFromContext(ctx).Identity()
	)

	// @todo this won't work entirely; if someone defines a dal connection to the same DSN as the primary one,
	//       they can easily bypass this.
	//       We'll need to do some checking on the DSN; potentially when defining the connection itself.
	c := svc.dal.GetConnectionByID(0)
	if m.ConnectionID == c.ID && m.Ident == "compose_record" {
		err = fmt.Errorf("cannot set alterations for default schema")
		return
	}

	// Delete current ones
	// @todo we might be able to do some diffing to preserve the metadata/ids
	//       but for now this should be just fine.
	auxStale := make([]*types.DalSchemaAlteration, len(stale))
	for i, a := range stale {
		auxStale[i] = &types.DalSchemaAlteration{
			ID: a.ID,
		}
	}
	err = store.DeleteDalSchemaAlteration(ctx, s, auxStale...)
	if err != nil {
		return
	}

	// Set new ones
	cvt := make(types.DalSchemaAlterationSet, len(aa))
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

		default:
			panic(fmt.Sprintf("unknown alteration type %v", a))
		}

		cvt[i] = t
	}

	return store.UpsertDalSchemaAlteration(ctx, svc.store, cvt...)
}

func (svc dalSchemaAlteration) Apply(ctx context.Context, ids ...uint64) (err error) {
	var (
		uaProps = &dalSchemaAlterationActionProps{}
	)

	err = func() (err error) {
		if !svc.ac.CanManageDalSchemaAlterations(ctx) {
			return DalSchemaAlterationErrNotAllowedToManage()
		}

		aux, _, err := store.SearchDalSchemaAlterations(ctx, svc.store, types.DalSchemaAlterationFilter{
			AlterationID: id.Strings(ids...),
		})
		if err != nil {
			return
		}

		alts := svc.appliableAlterations(aux...)
		pkgAlts, err := svc.toPkgAlterations(ctx, alts...)
		if err != nil {
			return
		}

		ii := make([]uint64, len(alts))
		for i, a := range alts {
			ii[i] = a.ID
		}

		uaProps.setApply(ii)

		errors, err := svc.dal.ApplyAlteration(ctx, pkgAlts...)
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

		err = store.UpdateDalSchemaAlteration(ctx, svc.store, aux...)
		if err != nil {
			return
		}

		return svc.reloadAlteredModels(ctx, svc.store, alts)
	}()

	return svc.recordAction(ctx, uaProps, DalSchemaAlterationActionApply, err)

}

func (svc dalSchemaAlteration) Dismiss(ctx context.Context, ids ...uint64) (err error) {
	var (
		uaProps = &dalSchemaAlterationActionProps{}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.ac.CanManageDalSchemaAlterations(ctx) {
			return DalSchemaAlterationErrNotAllowedToManage()
		}

		alt, _, err := store.SearchDalSchemaAlterations(ctx, s, types.DalSchemaAlterationFilter{
			AlterationID: id.Strings(ids...),
		})
		if err != nil {
			return
		}

		ii := make([]uint64, len(alt))
		for i, a := range alt {
			ii[i] = a.ID
		}

		uaProps.setApply(ii)

		alt = svc.appliableAlterations(alt...)
		identity := intAuth.GetIdentityFromContext(ctx).Identity()
		for _, a := range alt {
			a.Error = ""
			a.DismissedAt = now()
			a.DismissedBy = identity
		}

		err = store.UpdateDalSchemaAlteration(ctx, s, alt...)
		if err != nil {
			return
		}

		return svc.reloadAlteredModels(ctx, s, alt)
	})

	return svc.recordAction(ctx, uaProps, DalSchemaAlterationActionDismiss, err)
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

func (svc dalSchemaAlteration) reloadAlteredModels(ctx context.Context, s store.Storer, alts types.DalSchemaAlterationSet) (err error) {
	// Skip any models whish were already reloaded by some alterations.
	// These might be mixed up so we'll need to do it like so.
	processed := make(map[string]bool, 3)
	mkKey := func(a *types.DalSchemaAlteration) string {
		return fmt.Sprintf("%s;%s", a.ResourceType, a.Resource)
	}

	for _, a := range alts {
		k := mkKey(a)
		if processed[k] {
			continue
		}
		processed[k] = true

		err = svc.reloadAlteredModel(ctx, s, a)
		if err != nil {
			return
		}
	}

	return
}

func (svc dalSchemaAlteration) reloadAlteredModel(ctx context.Context, s store.Storer, alt *types.DalSchemaAlteration) (err error) {
	// Fetch current alterations to see if there are any left over
	_, f, err := store.SearchDalSchemaAlterations(ctx, s, types.DalSchemaAlterationFilter{
		Resource:     []string{alt.Resource},
		ResourceType: alt.ResourceType,
		Deleted:      filter.StateExcluded,
		Completed:    filter.StateExcluded,
		Dismissed:    filter.StateExcluded,

		Paging: filter.Paging{
			IncTotal: true,
		},
	})
	if err != nil {
		return err
	}

	// There are some alterations left so we can't reload the model
	if f.Total > 0 {
		return
	}

	model := svc.dal.FindModelByResourceIdent(alt.ConnectionID, alt.ResourceType, alt.Resource)
	if model == nil {
		err = fmt.Errorf("cannot find model for resource %s", alt.Resource)
		return
	}

	currentAlts, err := svc.modelAlterations(ctx, s, model)
	if err != nil {
		return
	}

	newAlts, err := svc.dal.ReloadModel(ctx, currentAlts, model)
	if err != nil {
		return
	}

	err = svc.SetAlterations(ctx, s, model, currentAlts, newAlts...)
	if err != nil {
		return
	}

	return
}
