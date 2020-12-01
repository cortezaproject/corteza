package envoy

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/json"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	stypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	tc struct {
		name  string
		suite string
		file  string

		// Before the data gets processed
		pre func() error
		// After the data gets processed
		post func(req *require.Assertions, err error)
		// Data assertions
		check func(req *require.Assertions)
	}
)

func initStore(ctx context.Context) (store.Storer, error) {
	s, err := sqlite3.ConnectInMemoryWithDebug(ctx)
	if err != nil {
		return nil, err
	}

	err = store.Upgrade(ctx, zap.NewNop(), s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func yd(ctx context.Context, suite, fname string) ([]resource.Interface, error) {
	fp := path.Join("testdata", suite, fname)

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := yaml.Decoder()
	return d.Decode(ctx, f, nil)
}

func dd(ctx context.Context, suite string) ([]resource.Interface, error) {
	fp := path.Join("testdata", suite)
	yd := yaml.Decoder()
	cd := csv.Decoder()
	jd := json.Decoder()

	return directory.Decode(ctx, fp, yd, cd, jd)
}

func encode(ctx context.Context, s store.Storer, nn []resource.Interface) error {
	return encodeC(ctx, s, nn, nil)
}

func encodeC(ctx context.Context, s store.Storer, nn []resource.Interface, cfg *es.EncoderConfig) error {
	se := es.NewStoreEncoder(s, cfg)
	bld := envoy.NewBuilder(se)
	g, err := bld.Build(ctx, nn...)
	if err != nil {
		return err
	}

	return envoy.Encode(ctx, g, se)
}

func storeNamespace(ctx context.Context, s store.Storer, nsID uint64, ss ...string) error {
	ns := &types.Namespace{
		ID: nsID,
	}
	if len(ss) > 0 {
		ns.Slug = ss[0]
	}
	if len(ss) > 1 {
		ns.Name = ss[1]
	}
	return store.CreateComposeNamespace(ctx, s, ns)
}

func storeModule(ctx context.Context, s store.Storer, nsID, modID uint64, ss ...string) error {
	mod := &types.Module{
		ID:          modID,
		NamespaceID: nsID,
	}
	if len(ss) > 0 {
		mod.Handle = ss[0]
	}
	if len(ss) > 1 {
		mod.Name = ss[1]
	}
	return store.CreateComposeModule(ctx, s, mod)
}

func storeModuleField(ctx context.Context, s store.Storer, modID, fieldID uint64, ss ...string) error {
	f := &types.ModuleField{
		ID:       fieldID,
		ModuleID: modID,
	}
	if len(ss) > 0 {
		f.Name = ss[0]
	}
	if len(ss) > 1 {
		f.Label = ss[1]
	}
	return store.CreateComposeModuleField(ctx, s, f)
}

func storeRecord(ctx context.Context, s store.Storer, nsID, moduleID, recordID uint64, vv ...string) error {
	r := &types.Record{
		ID:          recordID,
		ModuleID:    moduleID,
		NamespaceID: nsID,
		Values:      make(types.RecordValueSet, 0, len(vv)),
	}

	mod := &types.Module{
		ID:          moduleID,
		NamespaceID: nsID,
		Fields:      make(types.ModuleFieldSet, 0, len(vv)),
	}

	for i, v := range vv {
		r.Values = append(r.Values, &types.RecordValue{
			RecordID: recordID,
			Name:     fmt.Sprintf("f%d", i+1),
			Value:    v,
		})

		mod.Fields = append(mod.Fields, &types.ModuleField{
			ModuleID: moduleID,
			Kind:     "String",
			Name:     fmt.Sprintf("f%d", i+1),
		})
	}

	return store.CreateComposeRecord(ctx, s, mod, r)
}

func storeRole(ctx context.Context, s store.Storer, rID uint64, ss ...string) error {
	r := &stypes.Role{
		ID: rID,
	}
	if len(ss) > 0 {
		r.Handle = ss[0]
	}
	if len(ss) > 1 {
		r.Name = ss[1]
	}
	return store.CreateRole(ctx, s, r)
}

// Helper to collect resulting errors, returning the first one
func ce(ee ...error) error {
	for _, e := range ee {
		if e != nil {
			return e
		}
	}
	return nil
}

func fullModLoad(ctx context.Context, s store.Storer, req *require.Assertions, nsID uint64, handle string) (*types.Module, error) {
	mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, nsID, handle)
	req.NoError(err)
	req.NotNil(mod)

	mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
	req.NoError(err)
	req.NotNil(mod.Fields)
	return mod, err
}
