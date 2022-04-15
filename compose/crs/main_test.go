package crs

import (
	"time"

	"context"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/id"
)

type (
	testDriver struct{}

	testStore map[string]*types.Record
)

var (
	gStore = make(testStore)

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}

	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func (d testDriver) Capabilities() capabilities.Set {
	return capabilities.FullCapabilities()
}

func (d testDriver) Can(uri string, cc ...capabilities.Capability) bool {
	return strings.HasPrefix(uri, "noop://") && capabilities.Set(cc).IsSubset(d.Capabilities()...)
}

func (d testDriver) Store(ctx context.Context, uri string) (str Store, err error) {
	return &testStore{}, nil
}

func (d testDriver) Close(ctx context.Context, uri string) (err error) {
	return nil
}

// ---

func (ds *testStore) CreateRecords(ctx context.Context, sch *data.Model, cc ...Getter) error {
	return nil
}

func (ds *testStore) SearchRecords(ctx context.Context, sch *data.Model, filter any) (Loader, error) {
	return nil, nil
}

func (ds *testStore) Models(context.Context) (data.ModelSet, error) {
	return nil, nil
}

func (ds *testStore) AddModel(context.Context, *data.Model, ...*data.Model) error {
	return nil
}

func (ds *testStore) AlterModel(ctx context.Context, old *data.Model, new *data.Model) error {
	return nil
}

func (ds *testStore) AlterModelAttribute(ctx context.Context, sch *data.Model, old data.Attribute, new data.Attribute, trans ...func(*data.Model, data.Attribute, expr.TypedValue) (expr.TypedValue, bool, error)) error {
	return nil
}

func initCRS(ctx context.Context) *composeRecordStore {
	crs, _ := ComposeRecordStore(ctx, CRSConnectionWrap(0, "noop://primary", capabilities.FullCapabilities()...), testDriver{})
	crs.AddStore(ctx, CRSConnectionWrap(defaultExternalCRS, "noop://external", capabilities.FullCapabilities()...))
	return crs
}

func defaultModule() *types.Module {
	return &types.Module{
		ID:    nextID(),
		Store: types.CRSDef{ComposeRecordStoreID: defaultExternalCRS},
		Fields: types.ModuleFieldSet{{
			Name: "first_name",
			Kind: "String",
		}},
	}
}
