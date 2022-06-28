package envoy

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/json"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/provision"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
	stypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	createdAt, _   = time.Parse(time.RFC3339, "2021-01-01T11:10:09Z")
	updatedAt, _   = time.Parse(time.RFC3339, "2021-01-02T11:10:09Z")
	suspendedAt, _ = time.Parse(time.RFC3339, "2021-01-03T11:10:09Z")

	initOnce = &sync.Once{}
	s        store.Storer
)

// // // // // // Resource helpers

func decodeYaml(ctx context.Context, suite, fname string) ([]resource.Interface, error) {
	fp := path.Join("testdata", suite, fname)

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := yaml.Decoder()
	return d.Decode(ctx, f, nil)
}

func decodeDirectory(ctx context.Context, suite string) ([]resource.Interface, error) {
	fp := path.Join("testdata", suite)
	yd := yaml.Decoder()
	cd := csv.Decoder()
	jd := json.Decoder()

	return directory.Decode(ctx, fp, yd, cd, jd)
}

func encode(ctx context.Context, s store.Storer, nn []resource.Interface) error {
	se := es.NewStoreEncoder(s, dal.Service(), &es.EncoderConfig{})
	g, err := envoy.NewSafeBuilder(se).Build(ctx, nn...)
	if err != nil && err != envoy.BuilderErrUnresolvedReferences {
		return err
	}

	if err == envoy.BuilderErrUnresolvedReferences {
		md := g.MissingDeps().Unique()
		df := es.NewDecodeFilter().FromRef(md...)

		sd := es.Decoder()
		mm, err := sd.Decode(ctx, s, dal.Service(), df)
		if err != nil {
			return err
		}

		for _, m := range mm {
			m.MarkPlaceholder()
		}

		g, err = envoy.NewBuilder(se).Build(ctx, append(nn, mm...)...)
		if err != nil {
			return err
		}
	}

	return envoy.Encode(ctx, g, se)
}

// // // // // // Store helpers

func truncateStore(ctx context.Context, s store.Storer, t *testing.T) {
	err := collect(
		store.TruncateComposeNamespaces(ctx, s),
		store.TruncateComposeModules(ctx, s),
		store.TruncateComposeModuleFields(ctx, s),
		store.TruncateComposePages(ctx, s),
		store.TruncateComposeCharts(ctx, s),

		store.TruncateRoles(ctx, s),
		store.TruncateUsers(ctx, s),
		store.TruncateTemplates(ctx, s),
		store.TruncateApplications(ctx, s),
		store.TruncateApigwRoutes(ctx, s),
		store.TruncateApigwFilters(ctx, s),
		store.TruncateReports(ctx, s),
		store.TruncateSettingValues(ctx, s),
		store.TruncateRbacRules(ctx, s),

		store.TruncateAutomationWorkflows(ctx, s),
		store.TruncateAutomationTriggers(ctx, s),
	)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// returns store & authorized context
//
// @todo this should be refactored to use TestMain (see https://golang.org/pkg/testing/#hdr-Main)
func initServices(ctx context.Context, t *testing.T) store.Storer {
	initOnce.Do(func() {
		var (
			log = zap.NewNop()
			err error
		)

		s, err = sqlite.ConnectInMemoryWithDebug(ctx)
		if err != nil {
			t.Fatalf("failed to init sqlite in-memory db: %v", err)
		}

		err = store.Upgrade(ctx, log, s)
		if err != nil {
			t.Fatalf("failed to upgrade sqlite in-memory db: %v", err)
		}

		rr, err := provision.SystemRoles(ctx, log, s)
		if err != nil {
			t.Fatalf("failed to provision system roles: %v", err)
		}

		uu, err := provision.SystemUsers(ctx, log, s)
		if err != nil {
			t.Fatalf("failed to provision system users: %v", err)
		}

		{
			// uncomment for verbose logging (db & rbac)
			//logger.SetDefault(logger.MakeDebugLogger())
		}

		if rbac.Global() == nil {
			// make sure this is done only once
			auth.SetSystemRoles(rr)
			auth.SetSystemUsers(uu, rr)

			rbac.SetGlobal(rbac.NewService(logger.Default(), s))
			rbac.Global().UpdateRoles(func() (out []*rbac.Role) {
				for _, r := range auth.BypassRoles() {
					out = append(out, rbac.BypassRole.Make(r.ID, r.Handle))
				}
				return
			}()...)
		}
	})

	return s
}

func storeComposeNamespace(ctx context.Context, s store.Storer, nsID uint64, ss ...string) error {
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

func loadComposeModuleFull(ctx context.Context, s store.Storer, req *require.Assertions, nsID uint64, handle string) (*types.Module, error) {
	mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, nsID, handle)
	req.NoError(err)
	req.NotNil(mod)

	mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
	req.NoError(err)
	req.NotNil(mod.Fields)
	return mod, err
}

func storeComposeModule(ctx context.Context, s store.Storer, nsID, modID uint64, ss ...string) error {
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

func storeComposeModuleField(ctx context.Context, s store.Storer, modID, fieldID uint64, ss ...string) error {
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

func storeComposeRecord(ctx context.Context, s store.Storer, nsID, moduleID, recordID uint64, vv ...string) error {
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

func truncateStoreRecords(ctx context.Context, s store.Storer, t *testing.T) {
	err := collect(
		s.TruncateComposeRecords(ctx, nil),
	)
	if err != nil {
		t.Fatal(err.Error())
	}
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

func storeUser(ctx context.Context, s store.Storer, rID uint64, ss ...string) error {
	u := &stypes.User{
		ID: rID,
	}
	if len(ss) > 0 {
		u.Handle = ss[0]
	}
	if len(ss) > 1 {
		u.Name = ss[1]
	}
	return store.CreateUser(ctx, s, u)
}

// // // // // // Misc. helpers

// collect collects all errors from different call responses
func collect(ee ...error) error {
	for _, e := range ee {
		if e != nil {
			return e
		}
	}
	return nil
}

func parseTime(t *testing.T, ts string) *time.Time {
	tt, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Fatal(err.Error())
	}
	return &tt
}
