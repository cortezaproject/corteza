package service

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/logger"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	sysService "github.com/cortezaproject/corteza/server/system/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func makeTestModuleService(t *testing.T, mods ...any) *module {

	var (
		err error
		req = require.New(t)
		log = zap.NewNop()
	)

	for _, m := range mods {
		switch c := m.(type) {

		case *zap.Logger:
			t.Log("using custom Logger to initialize Module service")
			log = c
		}
	}

	var (
		ctx = logger.ContextWithValue(context.Background(), log)
		svc = &module{
			eventbus: eventbus.New(),
		}
	)

	for _, m := range mods {
		switch c := m.(type) {
		case rbacService:
			t.Log("using custom RBAC to initialize Module service")
			svc.ac = &accessControl{rbac: c}
		case store.Storer:
			t.Log("using custom Store to initialize Module service")
			svc.store = c
		case dal.FullService:
			t.Log("using custom DAL to initialize Module service")
			t.Log("make sure you manually reload models!")
			svc.dal = c
		}
	}

	if svc.store == nil {
		t.Log("using SQLite in-memory Store")
		svc.store, err = sqlite.ConnectInMemoryWithDebug(ctx)
		req.NoError(err)

		t.Log("upgrading store")
		req.NoError(store.Upgrade(ctx, log, svc.store))

		t.Log("data cleanup")
		req.NoError(store.TruncateUsers(ctx, svc.store))
		req.NoError(store.TruncateRoles(ctx, svc.store))
		req.NoError(store.TruncateComposeNamespaces(ctx, svc.store))
		req.NoError(store.TruncateComposeModules(ctx, svc.store))
		req.NoError(store.TruncateComposeModuleFields(ctx, svc.store))
		req.NoError(store.TruncateRbacRules(ctx, svc.store))
		req.NoError(store.TruncateLabels(ctx, svc.store))
	}

	if svc.ac == nil {
		rc, err := rbac.NewService(ctx, log, svc.store, rbac.Config{
			Synchronous:        true,
			DecayInterval:      time.Hour * 2,
			CleanupInterval:    time.Hour * 2,
			ReindexInterval:    time.Hour * 2,
			IndexFlushInterval: time.Hour * 2,
			RuleStorage:        svc.store,
			RoleStorage:        svc.store,
		})
		require.NoError(t, err)
		svc.ac = &accessControl{rbac: rc}
	}

	resourceMaker(ctx, t, svc.store, mods...)

	if svc.dal == nil {
		dalAux, err := dal.New(zap.NewNop(), true)
		req.NoError(err)

		const (
			recordsTable = "compose_record"
		)

		req.NoError(
			dalAux.ReplaceConnection(
				ctx,
				dal.MakeConnection(1, svc.store.ToDalConn(),
					dal.ConnectionParams{},
					dal.ConnectionConfig{ModelIdent: recordsTable},
				),
				true,
			),
		)

		svc.dal = dalAux

		t.Log("reloading DAL models")
		req.NoError(DalModelReload(ctx, svc.store, sysService.DefaultDalSchemaAlteration, dalAux))
	}

	return svc
}

func TestModules(t *testing.T) {
	var (
		ctx = context.Background()
		ns  = &types.Namespace{Name: "testing", ID: nextID(), CreatedAt: *now()}
	)

	t.Run("crud", func(t *testing.T) {
		req := require.New(t)

		svc := makeTestModuleService(t,
			ns,
			rbac.NoopSvc(rbac.Allow, rbac.Config{}),
		)

		res, err := svc.Create(ctx, &types.Module{Name: "My first module", NamespaceID: ns.ID})
		req.NoError(err)
		req.NotNil(res)

		res, err = svc.FindByID(ctx, ns.ID, res.ID)
		req.NoError(err)
		req.NotNil(res)

		res, err = svc.FindByHandle(ctx, ns.ID, res.Handle)
		req.NoError(err)
		req.NotNil(res)

		res.Name = "Changed"
		res, err = svc.Update(ctx, res)
		req.NoError(err)
		req.NotNil(res)
		req.NotNil(res.UpdatedAt)
		req.Equal(res.Name, "Changed")

		res, err = svc.FindByID(ctx, ns.ID, res.ID)
		req.NoError(err)
		req.NotNil(res)
		req.Equal(res.Name, "Changed")

		err = svc.DeleteByID(ctx, ns.ID, res.ID)
		req.NoError(err)
		req.NotNil(res)

		// this works because we're allowed to do everything
		res, err = svc.FindByID(ctx, ns.ID, res.ID)
		req.NoError(err)
		req.NotNil(res)
		req.NotNil(res.DeletedAt)
	})
}

func TestModule_LabelSearch(t *testing.T) {
	var (
		ns  = &types.Namespace{Name: "testing", ID: nextID(), CreatedAt: *now()}
		req = require.New(t)
		svc = makeTestModuleService(t,
			ns,
			rbac.NoopSvc(rbac.Allow, rbac.Config{}),
		)

		ctx = context.Background()
		//ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

		makeModule = func(n ...string) *types.Module {
			mod := &types.Module{
				NamespaceID: ns.ID,
				Name:        n[0],
				Labels:      map[string]string{},
			}

			for i := 1; i < len(n); i += 2 {
				mod.Labels[n[i]] = n[i+1]
			}

			out, err := svc.Create(ctx, mod)

			req.NoError(err)
			return out
		}

		findModules = func(labels map[string]string, IDs []string) types.ModuleSet {
			f := types.ModuleFilter{NamespaceID: ns.ID, Labels: labels, ModuleID: IDs}
			set, _, err := svc.Find(ctx, f)
			req.NoError(err)

			return set
		}
	)

	makeModule("labeled module 1", "label1", "value1", "label2", "value2")
	m2 := makeModule("labeled module 2", "label1", "value1")
	m3 := makeModule("labeled module 3")

	//return all -- no label/ID filter, return all
	req.Len(findModules(nil, nil), 3)

	// return 2 - both that have label1=valu1
	req.Len(findModules(map[string]string{"label1": "value1"}, nil), 2)

	// return 0 - none have foo=foo
	req.Len(findModules(map[string]string{"missing": "missing"}, nil), 0)

	// one has label2=value2
	req.Len(findModules(map[string]string{"label2": "value2"}, nil), 1)

	// explicit by ID and label
	req.Len(findModules(map[string]string{"label1": "value1"}, id.Strings(m2.ID)), 1)

	// none with this combo
	req.Len(findModules(map[string]string{"foo": "foo"}, id.Strings(m3.ID)), 0)

	// one with explicit ID (regression) and nil for label filter
	req.Len(findModules(nil, id.Strings(m3.ID)), 1)

	// one with explicit ID (regression) and empty map for label filter
	req.Len(findModules(map[string]string{}, id.Strings(m3.ID)), 1)

}

func TestModule_LabelCRUD(t *testing.T) {
	var (
		ctx = context.Background()
		//ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

		ns = &types.Namespace{Name: "testing", ID: nextID(), CreatedAt: *now()}

		req = require.New(t)
		svc = makeTestModuleService(t,
			ns,
			rbac.NoopSvc(rbac.Allow, rbac.Config{}),
		)

		findAndReturnLabel = func(id uint64) map[string]string {
			res, err := svc.FindByID(ctx, ns.ID, id)
			req.NoError(err)
			req.NotNil(res)
			return res.Labels
		}
	)

	// create unlabeled module
	res, err := svc.Create(ctx, &types.Module{Name: "unLabeledIDs", NamespaceID: ns.ID})
	req.NoError(err)
	req.NotNil(res)
	req.Nil(res.Labels)

	// no labels should be present
	req.Nil(findAndReturnLabel(res.ID))

	// update the module with labels
	res.Labels = map[string]string{"label1": "1st"}
	res, err = svc.Update(ctx, res)
	req.NoError(err)
	req.NotNil(res)
	req.Contains(res.Labels, "label1")

	// must contain the added label
	req.Contains(findAndReturnLabel(res.ID), "label1")

	res, err = svc.Create(ctx, &types.Module{Name: "LabeledIDs", NamespaceID: ns.ID, Labels: map[string]string{"label2": "2nd"}})
	req.NoError(err)
	req.NotNil(res)
	req.Contains(res.Labels, "label2")

	// must contain the added label
	req.Contains(findAndReturnLabel(res.ID), "label2")

	// update with Meta:nil (should keep labels intact)
	res.Labels = nil
	res, err = svc.Update(ctx, res)
	req.NoError(err)

	req.Contains(findAndReturnLabel(res.ID), "label2")

	// update with Meta:empty-map (should remove all labels)
	res.Labels = map[string]string{}
	res, err = svc.Update(ctx, res)
	req.NoError(err)

	req.Empty(findAndReturnLabel(res.ID))
}

func TestModuleToModel(t *testing.T) {
	var (
		req   = require.New(t)
		model *dal.Model
		err   error

		m = &types.Module{
			ID:        1,
			Handle:    "model-handle",
			Config:    types.ModuleConfig{},
			CreatedAt: time.Time{},
			Fields:    []*types.ModuleField{},
		}
	)

	t.Log("ident on DAL config not set, use ident from connection config")
	model, err = ModuleToModel(nil, m, "ident-from-conn-config")
	req.NoError(err)
	req.Equal("ident-from-conn-config", model.Ident)

	t.Log("explicit ident in module's DAL config should override the handle")
	m.Config.DAL.Ident = "explicit-ident"
	model, err = ModuleToModel(nil, m, "ident-from-conn-config")
	req.NoError(err)
	req.Equal("explicit-ident", model.Ident)
}
