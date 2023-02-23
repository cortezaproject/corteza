package envoy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
)

func (e StoreEncoder) prepareRecordDatasource(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// @todo match existing records; for now use just the ID like V1

	for _, n := range nn {
		if n.Datasource == nil {
			panic("unexpected state: cannot call prepareRecordDatasource with nodes without a defined Datasource")
		}

		ds, ok := n.Datasource.(*RecordDatasource)
		if !ok {
			panic("unexpected datasource type: node expecting type of RecordDatasource")
		}

		err = e.prepareRecords(ctx, p, s, ds)
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

func (e StoreEncoder) prepareRecords(ctx context.Context, p envoyx.EncodeParams, s store.Storer, ds *RecordDatasource) (err error) {
	var (
		aux   = make(map[string]string)
		more  bool
		ident string
		rec   types.Record
	)

	ds.refToID = make(map[string]uint64)

	for {
		ident, more, err = ds.Next(ctx, aux)
		if err != nil || !more {
			return
		}

		ds.refToID[ident] = id.Next()

		rec, err = e.auxToRecord(aux)
		if err != nil {
			return err
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
		ident  string
		rec    types.Record

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

	// Prepare getters for reference fields
	// @todo user refs
	var (
		modIndex = make(map[string]envoyx.Ref)
		dsIndex  = make(map[string]envoyx.Ref)
	)
	// - first pass collects all the references
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
	getters := make(map[string]*recordGetter)
	for k := range modIndex {
		getters[k] = makeRecordGetter(dl, tree, n, modIndex[k], dsIndex[k])
	}

	// Iterate and encode
	//
	// @todo utilize batching
	for {
		ident, more, err = ds.Next(ctx, auxRec)
		if err != nil || !more {
			return
		}

		rec, err = e.auxToRecord(auxRec)
		if err != nil {
			return err
		}

		rec.NamespaceID = ns.ID
		rec.ModuleID = mod.ID

		// @todo temp
		rec.CreatedAt = time.Now()

		// Values and refs
		rec.ID = ds.refToID[ident]
		for i, v := range rec.Values {
			if getters[v.Name] == nil {
				continue
			}

			v.Ref, err = getters[v.Name].resolve(ctx, v.Value)
			if err != nil {
				return err
			}

			rec.Values[i] = v
		}

		// Save it
		err = dalutils.ComposeRecordCreate(ctx, dl, mod, &rec)
		if err != nil {
			return
		}
	}
}

func (e StoreEncoder) auxToRecord(aux map[string]string) (out types.Record, err error) {
	for k, v := range aux {
		err = out.SetValue(k, 0, v)
		if err != nil {
			return
		}
	}

	return
}
