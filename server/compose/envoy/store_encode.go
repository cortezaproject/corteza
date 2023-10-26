package envoy

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/pkg/errors"
)

func (e StoreEncoder) encode(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	dl, err := e.grabDal(p)
	if err != nil {
		return
	}

	switch rt {
	case ComposeRecordDatasourceAuxType:
		err = e.encodeRecordDatasources(ctx, p, s, dl, nn, tree)
		if err != nil {
			return
		}
	}
	return
}

func (e StoreEncoder) setChartDefaults(res *types.Chart) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateChart(*types.Chart) (err error) {
	return
}

func (e StoreEncoder) setModuleDefaults(res *types.Module) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateModule(*types.Module) (err error) {
	return
}

func (e StoreEncoder) setModuleFieldDefaults(res *types.ModuleField) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	if res.Kind == "" {
		res.Kind = "String"
	}

	// Update validator ID
	maxValidatorID := uint64(0)
	for _, v := range res.Expressions.Validators {
		if v.ValidatorID > maxValidatorID {
			maxValidatorID = v.ValidatorID
		}
	}

	for _, v := range res.Expressions.Validators {
		if v.ValidatorID == 0 {
			v.ValidatorID = maxValidatorID + 1
			maxValidatorID++
		}
	}

	return
}

func (e StoreEncoder) sanitizeModuleFieldBeforeSave(f *types.ModuleField) (err error) {
	delete(f.Options, "module")
	return
}

func (e StoreEncoder) validateModuleField(*types.ModuleField) (err error) {
	return
}

func (e StoreEncoder) setNamespaceDefaults(res *types.Namespace) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateNamespace(*types.Namespace) (err error) {
	return
}

func (e StoreEncoder) setPageDefaults(res *types.Page) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	if res.Title == "" {
		res.Title = res.Handle
	}

	// Update pageblock ID
	maxPageBlockID := uint64(0)
	for _, b := range res.Blocks {
		if b.BlockID > maxPageBlockID {
			maxPageBlockID = b.BlockID
		}
	}

	for i, b := range res.Blocks {
		if b.BlockID == 0 {
			b.BlockID = maxPageBlockID + 1
			maxPageBlockID++
			res.Blocks[i] = b
		}
	}

	return
}

func (e StoreEncoder) validatePage(*types.Page) (err error) {
	return
}

func (e StoreEncoder) setPageLayoutDefaults(res *types.PageLayout) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	// @note these pageblocks reference the ones define on the page so nothing
	//       to do here.

	return
}

func (e StoreEncoder) validatePageLayout(*types.PageLayout) (err error) {
	return
}

func (e StoreEncoder) setRecordDefaults(*types.Record) (err error) {
	return
}

func (e StoreEncoder) validateRecord(*types.Record) (err error) {
	return
}

func (e StoreEncoder) prepare(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet) (err error) {
	switch rt {
	case ComposeRecordDatasourceAuxType:
		return e.prepareRecordDatasource(ctx, p, s, nn)
	}

	return
}

func (e StoreEncoder) encodeModuleExtend(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, nested envoyx.NodeSet, tree envoyx.Traverser) (err error) {

	// Push fields under mod
	mod := n.Resource.(*types.Module)
	for _, n := range nested {
		if n.ResourceType != types.ModuleFieldResourceType {
			continue
		}

		mod.Fields = append(mod.Fields, n.Resource.(*types.ModuleField))
	}

	// Register to DAL
	dl, err := e.grabDal(p)
	if err != nil {
		return
	}

	nsNode := tree.ParentForRef(n, n.References["NamespaceID"])
	ns := nsNode.Resource.(*types.Namespace)

	// @todo get connection and things from there
	models, err := service.ModulesToModelSet(dl, ns, mod)
	if err != nil {
		return err
	}

	// @note there is only one model so this is ok
	_, err = dl.ReplaceModel(ctx, nil, models[0])
	return
}

func (e StoreEncoder) postModulesEncode(ctx context.Context, p envoyx.EncodeParams, s store.Storer, tree envoyx.Traverser, nn envoyx.NodeSet) (err error) {
	dl, err := e.grabDal(p)
	if err != nil {
		return
	}

	for _, n := range nn {
		cc := tree.ChildrenForResourceType(n, ComposeRecordDatasourceAuxType)
		err = e.encodeRecordDatasources(ctx, p, s, dl, cc, tree)
		if err != nil {
			return
		}
	}

	return
}

func (e *StoreEncoder) grabDal(p envoyx.EncodeParams) (dl dal.FullService, err error) {
	auxs, ok := p.Params[paramsKeyDAL]
	if !ok {
		err = errors.Errorf("store encoder expects a dal conforming to dal.FullService interface")
		return
	}

	dl, ok = auxs.(dal.FullService)
	if !ok {
		err = errors.Errorf("store encoder expects a dal conforming to dal.FullService interface")

		return
	}

	return
}
