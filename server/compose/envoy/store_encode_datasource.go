package envoy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/service/values"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	systemEnvoy "github.com/cortezaproject/corteza/server/system/envoy"
	"github.com/spf13/cast"
)

const (
	recordBatchMaxChunk = 100
)

var (
	rvSanitizer = values.Sanitizer()
	rvValidator = values.Validator()
	rvFormatter = values.Formatter()
)

func (e StoreEncoder) prepareRecordDatasource(ctx context.Context, p envoyx.EncodeParams, s store.Storer, dl dal.FullService, nn envoyx.NodeSet) (err error) {
	// @todo match existing records; for now use just the ID like V1

	for _, n := range nn {
		if n.Datasource == nil {
			panic("unexpected state: cannot call prepareRecordDatasource with nodes without a defined Datasource")
		}

		ds, ok := n.Datasource.(*RecordDatasource)
		if !ok {
			panic("unexpected datasource type: node expecting type of RecordDatasource")
		}

		err = e.prepareRecords(ctx, p, s, dl, ds, nn)
		if err != nil {
			return
		}
		err = ds.Reset(ctx)
		if err != nil {
			return
		}
	}

	return
}

func (e StoreEncoder) prepareRecords(ctx context.Context, p envoyx.EncodeParams, s store.Storer, dl dal.FullService, ds *RecordDatasource, nn envoyx.NodeSet) (err error) {
	var (
		aux   = make(map[string]string)
		more  bool
		ident []string
	)

	ds.refToID = make(map[string]uint64)
	ds.existingIDs = make(map[uint64]bool)

	// Just so we don't need to do a branch later down the line
	cex := ds.CheckExisting
	if cex == nil || len(ds.Mapping.KeyField) == 0 {
		cex = func(ctx context.Context, ref ...[]string) ([]uint64, error) {
			return make([]uint64, len(ref)), nil
		}
	}

	for {
		ident, more, err = ds.Next(ctx, aux)
		if err != nil || !more {
			return
		}

		// @todo we'll need to batch these up
		existing, err := cex(ctx, ident)
		if err != nil {
			return err
		}

		if existing[0] != 0 {
			ds.existingIDs[existing[0]] = true
		}

		if existing[0] == 0 {
			existing[0] = id.Next()
		}

		ds.AddRef(existing[0], ident...)

		// Construct a simple record for basic validation/preprocessing
		rec := types.Record{}
		for k, v := range aux {
			// Ignore errors at this point since some values can't yet be properly casted
			rec.SetValue(k, 0, v)
		}

		// @note defaults and validation will have to happen again when encoding
		//       since we won't persist it.
		//       Consider supporting updating the datasource's data.
		err = e.setRecordDefaults(&rec)
		if err != nil {
			return err
		}

		err = e.validateRecord(&rec)
		if err != nil {
			return err
		}
	}
}

func (e StoreEncoder) encodeRecordDatasources(ctx context.Context, p envoyx.EncodeParams, s store.Storer, dl dal.FullService, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeRecordDatasource(ctx, p, s, dl, n, tree)
		if err != nil {
			return
		}
	}

	return
}

func (e StoreEncoder) encodeRecordDatasource(ctx context.Context, p envoyx.EncodeParams, s store.Storer, dl dal.FullService, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	var (
		auxRec = make(map[string]string)
		more   bool
		ident  []string

		nsNode  *envoyx.Node
		ns      *types.Namespace
		modNode *envoyx.Node
		mod     *types.Module

		// This was already validated so we can blindly cast
		ds = n.Datasource.(*RecordDatasource)
	)

	// Get the parent namespace
	nsNode = tree.ParentForRef(n, n.References["NamespaceID"])
	if nsNode == nil {
		err = fmt.Errorf("cannot encode record datasource: missing reference for NamespaceID")
		return
	}
	ns = nsNode.Resource.(*types.Namespace)

	// Get the parent module
	modNode = tree.ParentForRef(n, n.References["ModuleID"])
	if modNode == nil {
		err = fmt.Errorf("cannot encode record datasource: missing reference for ModuleID")
		return
	}
	mod = modNode.Resource.(*types.Module)

	// Prepare getters for related resources
	recordGetters := e.makeRecordGetters(dl, tree, n)
	userGetters := e.makeUserGetters(s, dl, tree, mod, n)

	maykr := e.recordMaker(ns, mod, recordGetters, userGetters)

	var (
		rec     types.Record
		creates types.RecordSet
		updates types.RecordSet
		rve     *types.RecordValueErrorSet
	)
	for {
		ident, more, err = ds.Next(ctx, auxRec)
		if err != nil {
			return
		}
		if !more {
			break
		}

		err = func() (err error) {
			rec, err = maykr(ctx, auxRec)
			if err != nil {
				return
			}

			// Do these at the end so they can't be overwritten
			rec.NamespaceID = ns.ID
			rec.ModuleID = mod.ID
			rec.CreatedAt = time.Now()
			rec.OwnedBy = service.CalcRecordOwner(0, rec.OwnedBy, p.Encoder.DefaultUserID)
			rec.ID, err = ds.ResolveRefS(ident...)
			if err != nil {
				return err
			}

			// Standard record processing
			rec.Values.SetUpdatedFlag(true)
			//
			rve = service.RecordValueUpdateOpCheck(ctx, nil, mod, rec.Values)
			if !rve.IsValid() {
				return rve
			}
			//
			rve = service.RecordPreparer(ctx, s, rvSanitizer, rvValidator, rvFormatter, mod, &rec)
			if !rve.IsValid() {
				return rve
			}

			ax := rec

			if ds.existingIDs[rec.ID] {
				updates = append(updates, &ax)
			} else {
				creates = append(creates, &ax)
			}

			if len(creates) > recordBatchMaxChunk {
				err = dalutils.ComposeRecordCreate(ctx, dl, mod, creates...)
				if err != nil {
					return
				}

				creates = make(types.RecordSet, 0, recordBatchMaxChunk/2)
			}

			if len(updates) > recordBatchMaxChunk {
				err = dalutils.ComposeRecordUpdate(ctx, dl, mod, updates...)
				if err != nil {
					return
				}

				updates = make(types.RecordSet, 0, recordBatchMaxChunk/2)
			}
			return
		}()

		if p.Defer != nil {
			p.Defer()
		}
		if err != nil {
			if p.DeferNok != nil {
				err = p.DeferNok(err)
			}
			if err != nil {
				return
			}
		} else if p.DeferOk != nil {
			p.DeferOk()
		}
	}

	if len(creates) > 0 {
		err = dalutils.ComposeRecordCreate(ctx, dl, mod, creates...)
	}

	if len(updates) > 0 {
		err = dalutils.ComposeRecordUpdate(ctx, dl, mod, updates...)
	}
	return
}

// maceRecordGetters returns a map of getters where the key is the field name
func (e StoreEncoder) makeRecordGetters(dl dal.FullService, tree envoyx.Traverser, n *envoyx.Node) (getters map[string]*recordGetter) {
	var (
		modIndex = make(map[string]envoyx.Ref)
		dsIndex  = make(map[string]envoyx.Ref)
	)

	// First pass collects all the references
	for k, v := range n.References {
		if k == "NamespaceID" || k == "ModuleID" {
			continue
		}

		// Encoded as <fieldName>.<refKind>
		pp := strings.Split(k, ".")
		f := pp[0]
		kind := pp[1]

		switch kind {
		case "module":
			modIndex[f] = v
		case "datasource":
			dsIndex[f] = v
		}
	}

	// - second pass makes the getters
	getters = make(map[string]*recordGetter)
	for k := range modIndex {
		getters[k] = makeRecordGetter(dl, tree, n, modIndex[k], dsIndex[k])
	}

	return
}

func (e StoreEncoder) makeUserGetters(s store.Storer, dl dal.FullService, tree envoyx.Traverser, mod *types.Module, n *envoyx.Node) (getters map[string]*systemEnvoy.UserGetter) {
	userGetter := systemEnvoy.MakeUserGetter(s, tree)
	getters = map[string]*systemEnvoy.UserGetter{
		// These are all of the supported sys user ref fields
		"ownedby":    userGetter,
		"owned_by":   userGetter,
		"createdby":  userGetter,
		"created_by": userGetter,
		"updatedby":  userGetter,
		"updated_by": userGetter,
		"deletedby":  userGetter,
		"deleted_by": userGetter,
	}
	for _, f := range mod.Fields {
		if f.Kind == "User" {
			getters[f.Name] = userGetter
		}
	}

	return
}

func (e StoreEncoder) recordMaker(ns *types.Namespace, mod *types.Module, recordGetters map[string]*recordGetter, userGetters map[string]*systemEnvoy.UserGetter) func(ctx context.Context, auxRec map[string]string) (r types.Record, err error) {
	return func(ctx context.Context, auxRec map[string]string) (rec types.Record, err error) {
		// Iterate mapified record values and populate the provided values
		var auxv *types.RecordValue
		for k, v := range auxRec {
			if recordGetters[k] != nil {
				auxv, err = e.resolveRecordRef(ctx, recordGetters, k, v)
				if err != nil {
					return
				}
				err = rec.SetValue(k, auxv.Place, auxv.Ref)
				if err != nil {
					return
				}
				continue
			}

			if userGetters[k] != nil {
				auxv, err = e.resolveUserRef(ctx, userGetters, k, v)
				if err != nil {
					return
				}
				err = rec.SetValue(k, auxv.Place, auxv.Ref)
				if err != nil {
					return
				}
				continue
			}

			// Default is regular values
			err = rec.SetValue(k, 0, v)
			if err != nil {
				return
			}
		}

		return
	}
}

func (e StoreEncoder) resolveRecordRef(ctx context.Context, getters map[string]*recordGetter, k, v string) (out *types.RecordValue, err error) {
	if getters[k] == nil {
		return nil, nil
	}

	out = &types.RecordValue{Name: k}
	out.Ref, err = getters[out.Name].resolve(ctx, v)
	if err != nil {
		return
	}
	out.Value = cast.ToString(out.Ref)

	return out, nil
}

func (e StoreEncoder) resolveUserRef(ctx context.Context, getters map[string]*systemEnvoy.UserGetter, k, v string) (out *types.RecordValue, err error) {
	k = strings.ToLower(k)
	if getters[k] == nil {
		return nil, nil
	}

	out = &types.RecordValue{Name: k}
	out.Ref, err = getters[out.Name].Resolve(ctx, v)
	if err != nil {
		return
	}
	out.Value = cast.ToString(out.Ref)

	return out, nil
}
