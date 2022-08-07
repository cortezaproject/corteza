package dal

import (
	"context"
	"fmt"
	"os"
	"testing"

	composeService "github.com/cortezaproject/corteza-server/compose/service"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/dal/setup/mysql"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	driver struct {
		name  string
		dsn   string
		setup func(ctx context.Context, t aaaa, dsn string)
	}

	dalService interface {
		Drivers() (drivers []dal.Driver)

		ReplaceSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
		RemoveSensitivityLevel(levelIDs ...uint64) (err error)

		ReplaceConnection(ctx context.Context, cw *dal.ConnectionWrap, isDefault bool) (err error)
		RemoveConnection(ctx context.Context, ID uint64) (err error)

		SearchModels(ctx context.Context) (out dal.ModelSet, err error)
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
		RemoveModel(ctx context.Context, connectionID, ID uint64) (err error)
		ReplaceModelAttribute(ctx context.Context, model *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) (err error)
		FindModelByResourceID(connectionID uint64, resourceID uint64) *dal.Model
		FindModelByResourceIdent(connectionID uint64, resourceType, resourceIdent string) *dal.Model
		FindModelByIdent(connectionID uint64, ident string) *dal.Model

		Create(ctx context.Context, mf dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) (err error)
		Update(ctx context.Context, mf dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) (err error)
		Search(ctx context.Context, mf dal.ModelRef, operations dal.OperationSet, f filter.Filter) (iter dal.Iterator, err error)
		Lookup(ctx context.Context, mf dal.ModelRef, operations dal.OperationSet, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, mf dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) (err error)
		Truncate(ctx context.Context, mf dal.ModelRef, operations dal.OperationSet) (err error)

		SearchConnectionIssues(connectionID uint64) (out []error)
		SearchModelIssues(connectionID, resourceID uint64) (out []error)
	}

	dalConnectionRestRsp struct {
		Response *types.DalConnection
	}

	composeModuleRestRsp struct {
		Response *composeTypes.Module
	}

	composeRecordRestRsp struct {
		Response *composeTypes.Record
	}
	composeRecordSearchRestRsp struct {
		Response struct {
			Set composeTypes.RecordSet
		}
	}

	aaaa interface {
		require.TestingT
		Name() string
	}

	modelMetaMaker func(ident string) *dal.Model
	attributeMaker func() dal.AttributeSet
)

const (
	sysID          = "ID"
	sysNamespaceID = "namespaceID"
	sysModuleID    = "moduleID"
	sysCreatedAt   = "createdAt"
	sysCreatedBy   = "createdBy"
	sysUpdatedAt   = "updatedAt"
	sysUpdatedBy   = "updatedBy"
	sysDeletedAt   = "deletedAt"
	sysDeletedBy   = "deletedBy"
	sysOwnedBy     = "ownedBy"

	colSysID          = "id"
	colSysNamespaceID = "rel_namespace"
	colSysModuleID    = "module_id"
	colSysCreatedAt   = "created_at"
	colSysCreatedBy   = "created_by"
	colSysUpdatedAt   = "updated_at"
	colSysUpdatedBy   = "updated_by"
	colSysDeletedAt   = "deleted_at"
	colSysDeletedBy   = "deleted_by"
	colSysOwnedBy     = "owned_by"
)

func (h helper) cleanupDal() {
	ds := service.DefaultStore
	ctx := context.Background()

	dd := dal.Service()
	models, err := dd.SearchModels(ctx)
	h.a.NoError(err)
	for _, model := range models {
		dd.Truncate(ctx, model.ToFilter(), nil)
	}
	h.a.NoError(store.TruncateComposeModuleFields(ctx, ds))
	h.a.NoError(store.TruncateComposeModules(ctx, ds))
	h.a.NoError(store.TruncateComposeNamespaces(ctx, ds))
	h.a.NoError(store.TruncateDalSensitivityLevels(ctx, ds))

	cc, _, err := store.SearchDalConnections(ctx, ds, types.DalConnectionFilter{})
	h.a.NoError(err)
	for _, c := range cc {
		if c.Type == types.DalPrimaryConnectionResourceType {
			continue
		}
		h.a.NoError(store.DeleteDalConnectionByID(ctx, ds, c.ID))
	}
	dd.Purge(ctx)

	h.a.NoError(service.DefaultDalConnection.ReloadConnections(ctx))
	h.a.NoError(composeService.DefaultModule.ReloadDALModels(ctx))
}

func initSvc(ctx context.Context, d driver) (dalService, error) {
	c := makeConnectionDefinition(d.dsn)

	cm := dal.ConnectionConfig{
		ModelIdent:         c.Config.DAL.ModelIdent,
		AttributeIdent:     c.Config.DAL.AttributeIdent,
		SensitivityLevelID: c.Config.Privacy.SensitivityLevelID,
		Label:              c.Handle,
	}

	svc, err := dal.New(zap.NewNop(), false)
	if err != nil {
		return nil, err
	}

	err = svc.ReplaceConnection(ctx, dal.MakeConnection(c.ID, nil, c.Config.Connection, cm, dal.FullOperations()...), true)
	if err != nil {
		return nil, err
	}

	return svc, err
}

func setup(t *testing.T) (ctx context.Context, h helper, log *zap.Logger) {
	log = zap.NewNop()

	h = newHelperT(t)

	u := &types.User{
		ID: id.Next(),
	}
	u.SetRoles(auth.BypassRoles().IDs()...)

	ctx = auth.SetIdentityToContext(context.Background(), u)

	return ctx, h, log
}

func setupBench(b *testing.B) (ctx context.Context, h helper, log *zap.Logger) {
	log = zap.NewNop()

	h = newHelperB(b)

	u := &types.User{
		ID: id.Next(),
	}
	u.SetRoles(auth.BypassRoles().IDs()...)

	ctx = auth.SetIdentityToContext(context.Background(), u)

	return ctx, h, log
}

func bootstrap(rootTest *testing.T, run func(context.Context, *testing.T, helper, dalService)) {
	var (
		ctx, h, log = setup(rootTest)
		drivers     = collectDrivers()
	)
	_ = log

	for _, driver := range drivers {
		rootTest.Run(driver.name, func(t *testing.T) {
			if driver.dsn == "" {
				t.Skip("DSN for DAL test not set")
			}

			svc, err := initSvc(ctx, driver)
			require.NoError(t, err)

			driver.setup(ctx, rootTest, driver.dsn)

			run(ctx, t, h, svc)
		})
	}
}

func bootstrapWithModel(rootTest *testing.T, ident string, run func(context.Context, *testing.T, helper, dalService, *dal.Model)) {
	var (
		ctx, h, log = setup(rootTest)
		drivers     = collectDrivers()
	)
	_ = log

	for _, driver := range drivers {
		rootTest.Run(driver.name, func(t *testing.T) {
			if driver.dsn == "" {
				t.Skip("DSN for DAL test not set")
			}

			svc, err := initSvc(ctx, driver)
			require.NoError(t, err)

			model := buildModel(ident, basicModelMeta, fullPartitionedSysAttributes)

			driver.setup(ctx, rootTest, driver.dsn)

			require.NoError(t, svc.ReplaceModel(ctx, model))

			run(ctx, t, h, svc, model)
		})
	}
}

func bootstrapBenchmark(rootTest *testing.B, run func(context.Context, *testing.B, helper, dalService)) {
	var (
		ctx, h, log = setupBench(rootTest)
		drivers     = collectDrivers()
	)
	_ = log

	for _, driver := range drivers {
		rootTest.Run(driver.name, func(b *testing.B) {
			if driver.dsn == "" {
				b.Skip("DSN for DAL test not set")
			}

			svc, err := initSvc(ctx, driver)
			require.NoError(b, err)

			driver.setup(ctx, rootTest, driver.dsn)

			run(ctx, b, h, svc)
		})
	}
}

func collectDrivers() []driver {
	return []driver{
		// {
		// 	name: "sqlite",
		// 	dsn:  "sqlite3+debug://file::memory:?cache=shared&mode=memory",
		// 	setup: func(ctx context.Context, t *testing.T, dsn string) {
		// 		conn, err := mysql.Setup(ctx, dsn)
		// 		require.NoError(t, err)

		// 		_, err = conn.Exec(loadSetupSource(t, "sqlite.sql"))
		// 		require.NoError(t, err)
		// 	},
		// },
		{
			name: "mysql",
			dsn:  os.Getenv("DAL_TEST_DSN_MYSQL"),
			setup: func(ctx context.Context, t aaaa, dsn string) {
				conn, err := mysql.Setup(ctx, dsn)
				require.NoError(t, err)

				_, err = conn.Exec(loadScenarioSources(t, "mysql", "sql"))
				require.NoError(t, err)
			},
		},
		// {
		// 	name: "postgres",
		// 	dsn:  os.Getenv("DAL_TEST_DSN_POSTGRES"),
		// 	setup: func(ctx context.Context, t *testing.T, dsn string) {
		// 		conn, err := mysql.Setup(ctx, dsn)
		// 		require.NoError(t, err)

		// 		_, err = conn.Exec(loadSetupSource(t, "postgres.sql"))
		// 		require.NoError(t, err)
		// 	},
		// },
	}
}

func buildModel(ident string, mm modelMetaMaker, am attributeMaker, aa ...*dal.Attribute) *dal.Model {
	out := mm(ident)
	out.Attributes = aa
	out.Attributes = append(out.Attributes, am()...)

	return out
}

func basicModelMeta(ident string) *dal.Model {
	return &dal.Model{
		ConnectionID: 0,
		Ident:        ident,
		Resource:     fmt.Sprintf("testing-resource/%s", ident),
		ResourceID:   0,
	}
}

func fullSysAttributes() dal.AttributeSet {
	return dal.AttributeSet{
		dal.PrimaryAttribute(sysID, &dal.CodecAlias{Ident: colSysID}),

		dal.FullAttribute(sysModuleID, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysModuleID}),
		dal.FullAttribute(sysNamespaceID, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysNamespaceID}),

		dal.FullAttribute(sysOwnedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysOwnedBy}),

		dal.FullAttribute(sysCreatedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysCreatedAt}),
		dal.FullAttribute(sysCreatedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysCreatedBy}),

		dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysUpdatedAt}),
		dal.FullAttribute(sysUpdatedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysUpdatedBy}),

		dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysDeletedAt}),
		dal.FullAttribute(sysDeletedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysDeletedBy}),
	}
}

func fullPartitionedSysAttributes() dal.AttributeSet {
	return dal.AttributeSet{
		dal.PrimaryAttribute(sysID, &dal.CodecAlias{Ident: colSysID}),

		// mod and ns omitted here

		dal.FullAttribute(sysOwnedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysOwnedBy}),

		dal.FullAttribute(sysCreatedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysCreatedAt}),
		dal.FullAttribute(sysCreatedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysCreatedBy}),

		dal.FullAttribute(sysUpdatedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysUpdatedAt}),
		dal.FullAttribute(sysUpdatedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysUpdatedBy}),

		dal.FullAttribute(sysDeletedAt, &dal.TypeTimestamp{}, &dal.CodecAlias{Ident: colSysDeletedAt}),
		dal.FullAttribute(sysDeletedBy, &dal.TypeID{}, &dal.CodecAlias{Ident: colSysDeletedBy}),
	}
}

func drain(ctx context.Context, i dal.Iterator) (rr []*composeTypes.Record, err error) {
	var r *composeTypes.Record
	rr = make([]*composeTypes.Record, 0, 100)
	for i.Next(ctx) {
		if i.Err() != nil {
			return nil, i.Err()
		}

		r = new(composeTypes.Record)
		if err = i.Scan(r); err != nil {
			return
		}

		rr = append(rr, r)
	}

	return rr, i.Err()
}
