package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms"

	"github.com/cortezaproject/corteza/server/compose/service/values"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	sysService "github.com/cortezaproject/corteza/server/system/service"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func makeTestRecordService(t *testing.T, mods ...any) *record {
	var (
		err error
		req = require.New(t)
		log = zap.NewNop()
	)

	for _, m := range mods {
		switch c := m.(type) {

		case *zap.Logger:
			t.Log("using custom Logger to initialize Record service")
			log = c
		}
	}

	var (
		ctx = logger.ContextWithValue(context.Background(), log)
		svc = &record{
			sanitizer: values.Sanitizer(),
			formatter: values.Formatter(),
			eventbus:  eventbus.New(),
		}
	)

	svc.validator = defaultValidator(svc)

	for _, m := range mods {
		switch c := m.(type) {
		case rbacService:
			t.Log("using custom RBAC to initialize Record service")
			svc.ac = &accessControl{rbac: c}
		case store.Storer:
			t.Log("using custom Store to initialize Record service")
			svc.store = c
		case dal.FullService:
			t.Log("using custom DAL to initialize Record service")
			t.Log("make sure you manually reload models!")
			svc.dal = c
		}
	}

	if svc.store == nil {
		svc.store, err = sqlite.ConnectInMemoryWithDebug(ctx)
		req.NoError(err)

		req.NoError(store.Upgrade(ctx, log, svc.store))
		req.NoError(store.TruncateUsers(ctx, svc.store))
		req.NoError(store.TruncateRoles(ctx, svc.store))
		req.NoError(store.TruncateComposeNamespaces(ctx, svc.store))
		req.NoError(store.TruncateComposeModules(ctx, svc.store))
		req.NoError(store.TruncateComposeModuleFields(ctx, svc.store))
		req.NoError(store.TruncateRbacRules(ctx, svc.store))
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
		svc.rbacSvc = rc
		svc.ac = &accessControl{rbac: rc}
	}

	resourceMaker(ctx, t, svc.store, mods...)

	if svc.dal == nil {
		t.Log("adding default DAL service")
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

		// assuming store will always be RDBMS/SQLite and we can just run delete
		//
		// this could be done more "properly" by invoking DAL to truncate records
		// but at the momment we need to provide the right model/module.
		_, err = svc.store.(*rdbms.Store).DB.ExecContext(ctx, "DELETE FROM "+recordsTable)
		req.NoError(err)

		t.Log("reloading DAL models")
		req.NoError(DalModelReload(ctx, svc.store, sysService.DefaultDalSchemaAlteration, dalAux))
	}

	return svc
}

func ModelRef(svc *record, moduleID uint64) *dal.Model {
	type (
		modelFinder interface {
			FindModelByResourceID(connectionID uint64, resourceID uint64) *dal.Model
		}
	)

	var (
		aux any = svc.dal
	)

	return aux.(modelFinder).FindModelByResourceID(1, moduleID)
}

func verifyRecErrSet(t *testing.T, err any) {
	if rves, is := err.(*types.RecordValueErrorSet); is {
		if rves.IsValid() {
			return
		}

		t.Logf("Record errors:")
		for _, e := range rves.Set {
			t.Logf(" => [%s] %s", e.Kind, e.Message)
			t.Logf("    %v", e.Meta)
		}

		t.FailNow()
	}

	if err, is := err.(*errors.Error); is {
		t.Errorf("unexpected error: %s", err.Error())
		for k, v := range err.Meta() {
			t.Errorf("                  %+v: %+v", k, v)
		}
		t.FailNow()
	}

	if err, is := err.(error); is {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
}

func TestGeneralValueSetValidation(t *testing.T) {
	var (
		req = require.New(t)

		module = &types.Module{
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Name: "single1"},
				&types.ModuleField{Name: "multi1", Multi: true},
				&types.ModuleField{Name: "ref1", Kind: "Record"},
				&types.ModuleField{Name: "multiRef1", Kind: "Record", Multi: true},
			},
		}

		rvs types.RecordValueSet
		err error
	)

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}}
	err = RecordValueSanitization(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "unknown", Value: "single"}}
	err = RecordValueSanitization(module, rvs)
	req.True(err != nil, "expecting RecordValueSanitization() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}, {Name: "single1", Value: "single2"}}
	err = RecordValueSanitization(module, rvs)
	req.Error(err)

	rvs = types.RecordValueSet{{Name: "multi1", Value: "multi1"}, {Name: "multi1", Value: "multi1"}}
	err = RecordValueSanitization(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "non numeric value"}}
	err = RecordValueSanitization(module, rvs)
	req.Error(err)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "12345"}}
	err = RecordValueSanitization(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "multiRef1", Value: "12345"}, {Name: "multiRef1", Value: "67890"}}
	err = RecordValueSanitization(module, rvs)
	req.NoError(err)
	req.Len(rvs, 2, "expecting 2 record values after sanitization, got %d", len(rvs))

	rvs = types.RecordValueSet{{Name: "ref1", Value: ""}}
	err = RecordValueSanitization(module, rvs)
	req.NoError(err)
}

func TestDefaultValueSetting(t *testing.T) {
	var (
		a = assert.New(t)

		mod = &types.Module{
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Name: "single", DefaultValue: []*types.RecordValue{{Value: "s"}}},
				&types.ModuleField{Name: "multi", Multi: true, DefaultValue: []*types.RecordValue{{Value: "m1", Place: 0}, {Value: "m2", Place: 1}}},
			},
		}

		chk = func(vv types.RecordValueSet, f string, p uint, v string) {
			a.True(vv.Has("multi", p))
			a.Equal(v, vv.Get(f, p).Value)
		}
	)

	out := RecordValueDefaults(mod, nil)
	chk(out, "single", 0, "s")
	chk(out, "multi", 0, "m1")
	chk(out, "multi", 1, "m2")
}

func TestRecord_boolFieldPermissionIssueKBR(t *testing.T) {
	var (
		err error

		ctx = context.Background()
		req = require.New(t)

		u = &sysTypes.User{ID: nextID()}

		ns = &types.Namespace{ID: nextID()}

		modConf = types.ModuleConfig{DAL: types.ModuleConfigDAL{ConnectionID: 1}}

		mod         = &types.Module{ID: nextID(), NamespaceID: ns.ID, Config: modConf}
		stringField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "string", Kind: "String"}
		boolField   = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "bool", Kind: "Boolean"}

		authRoleID uint64 = 1

		readerRole = &sysTypes.Role{Name: "reader", ID: nextID()}
		writerRole = &sysTypes.Role{Name: "writer", ID: nextID()}

		//

		svc = makeTestRecordService(t,
			logger.MakeDebugLogger(),
			u,
			ns,
			mod,
			stringField,
			boolField,
		)

		recChecked, recUnchecked *types.Record

		valChecked = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
			&types.RecordValue{Name: "bool", Value: "1"},
		}

		valUnchecked = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
		}
	)

	svc.validator = defaultValidator(svc)

	svc.rbacSvc.UpdateRoles(
		rbac.CommonRole.Make(readerRole.ID, readerRole.Name),
		rbac.CommonRole.Make(writerRole.ID, writerRole.Name),
		rbac.AuthenticatedRole.Make(authRoleID, "authenticated"),
	)

	svc.rbacSvc.Grant(ctx,
		// base permissions
		rbac.AllowRule(authRoleID, mod.RbacResource(), "record.create"),
		rbac.AllowRule(authRoleID, types.RecordRbacResource(0, 0, 0), "read"),
		rbac.AllowRule(authRoleID, types.RecordRbacResource(0, 0, 0), "update"),
		rbac.DenyRule(authRoleID, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read"),
		rbac.DenyRule(authRoleID, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update"),

		// special perm for writer
		rbac.AllowRule(writerRole.ID, stringField.RbacResource(), "record.value.read"),
		rbac.AllowRule(writerRole.ID, stringField.RbacResource(), "record.value.update"),
		rbac.AllowRule(writerRole.ID, boolField.RbacResource(), "record.value.read"),
		rbac.AllowRule(writerRole.ID, boolField.RbacResource(), "record.value.update"),

		// special perm for reader
		rbac.AllowRule(readerRole.ID, stringField.RbacResource(), "record.value.read"),
		rbac.AllowRule(readerRole.ID, boolField.RbacResource(), "record.value.read"),
		rbac.DenyRule(readerRole.ID, boolField.RbacResource(), "record.value.update"),
	)

	{
		// security context w/ writer role
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(u.ID, writerRole.ID, authRoleID))

		recChecked, _, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: valChecked})
		verifyRecErrSet(t, err)

		req.NotNil(recChecked.Values.Get("bool", 0), "should be checked")
		req.Equal("1", recChecked.Values.Get("bool", 0).Value)

		recUnchecked, _, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: valUnchecked})
		req.NoError(err)

		req.Nil(recUnchecked.Values.Get("bool", 0))

		// *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** ***

		// security context w/ reader role
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(u.ID, readerRole.ID, authRoleID))

		recChecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
		}

		recChecked, _, err = svc.Update(ctx, recChecked)
		req.NoError(err)

		req.NotNil(recChecked.Values.Get("bool", 0), "should still be checked")
		req.Equal("1", recChecked.Values.Get("bool", 0).Value)

		recUnchecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
		}

		recUnchecked, _, err = svc.Update(ctx, recUnchecked)
		req.NoError(err)
		req.Nil(recUnchecked.Values.Get("bool", 0), "should not be checked anymore")

		// *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** ***

		// security context w/ writer role
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(u.ID, writerRole.ID, authRoleID))

		recChecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
			&types.RecordValue{Name: "bool", Value: "1"},
		}

		recChecked, _, err = svc.Update(ctx, recChecked)
		req.NoError(err)

		req.NotNil(recChecked.Values.Get("bool", 0), "should checked again")
		req.Equal("1", recChecked.Values.Get("bool", 0).Value, "should checked again")
	}
}

func TestRecord_defValueFieldPermissionIssue(t *testing.T) {
	var (
		err error
		req = require.New(t)
		ctx = context.Background()

		user = &sysTypes.User{ID: nextID()}

		modConf = types.ModuleConfig{DAL: types.ModuleConfigDAL{ConnectionID: 1}}

		ns            = &types.Namespace{ID: nextID()}
		mod           = &types.Module{ID: nextID(), NamespaceID: ns.ID, Config: modConf}
		writableField = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, NamespaceID: ns.ID, Name: "writable", Kind: "String", DefaultValue: types.RecordValueSet{{Value: "def-w"}}}
		readableField = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, NamespaceID: ns.ID, Name: "readable", Kind: "String", DefaultValue: types.RecordValueSet{{Value: "def-r"}}}

		svc = makeTestRecordService(t,
			user,
			ns,
			mod,
			writableField,
			readableField,
		)

		authRoleID = nextID()
		editorRole = &sysTypes.Role{Name: "editor", ID: nextID()}

		recPartial *types.Record

		valueExtractor = func(rec *types.Record, ff ...string) (out string) {
			if rec == nil {
				return "<NIL RECORD>"
			}

			if rec.Values == nil {
				return "<NIL RECORD VALUES>"
			}

			for _, f := range ff {
				out += "<"
				if v := rec.Values.Get(f, 0); v != nil {
					out += v.Value
				} else {
					out += "NULL"
				}
				out += ">"
			}

			return
		}
	)

	t.Log("creating namespace, module and fields")

	svc.validator = defaultValidator(svc)

	t.Log("setting up security")

	svc.rbacSvc.UpdateRoles(
		rbac.CommonRole.Make(editorRole.ID, editorRole.Name),
		rbac.AuthenticatedRole.Make(authRoleID, "authenticated"),
	)

	svc.rbacSvc.Grant(ctx,
		// base permissions
		rbac.AllowRule(authRoleID, mod.RbacResource(), "record.create"),
		rbac.AllowRule(authRoleID, types.RecordRbacResource(0, 0, 0), "read"),
		rbac.AllowRule(authRoleID, types.RecordRbacResource(0, 0, 0), "update"),

		rbac.AllowRule(authRoleID, types.ModuleFieldRbacResource(0, 0, 0), "record.value.read"),
		rbac.AllowRule(authRoleID, types.ModuleFieldRbacResource(0, 0, 0), "record.value.update"),

		// expl. deny value updates for editor on readable field (still allowed to write on writable field)
		rbac.DenyRule(editorRole.ID, writableField.RbacResource(), "record.value.update"),
	)

	{
		t.Log("creating record with w/o editor role (expecting defaults")

		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(user.ID, authRoleID))

		recPartial, _, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{}})
		verifyRecErrSet(t, err)
		req.Equal("<def-w><def-r>", valueExtractor(recPartial, "writable", "readable"))

		t.Log("creating record with w/o editor role (must be able to crate & update record and modify both fields)")

		recPartial, _, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{
			&types.RecordValue{Name: "writable", Value: "w"},
			&types.RecordValue{Name: "readable", Value: "r"},
		}})

		verifyRecErrSet(t, err)
		req.Equal("<w><r>", valueExtractor(recPartial, "writable", "readable"))

		t.Log("updating record removing one of the values")

		recPartial.Values = types.RecordValueSet{&types.RecordValue{Name: "writable", Value: "w2"}}

		recPartial, _, err = svc.Update(ctx, recPartial)
		verifyRecErrSet(t, err)
		req.Equal("<w2><NULL>", valueExtractor(recPartial, "writable", "readable"))
	}

	{
		t.Log("creating record with editor role (expecting defaults")

		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(user.ID, authRoleID, editorRole.ID))

		recPartial, _, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{}})
		verifyRecErrSet(t, err)
		req.Equal("<def-w><def-r>", valueExtractor(recPartial, "writable", "readable"))

		t.Log("creating record with editor role (must be able to crate & update record and modify both fields)")

		recPartial, _, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{
			// this is the def. value set
			&types.RecordValue{Name: "writable", Value: "def-w"},
			&types.RecordValue{Name: "readable", Value: "r"},
		}})

		verifyRecErrSet(t, err)
		req.Equal("<def-w><r>", valueExtractor(recPartial, "writable", "readable"))
	}
}

func TestRecord_refAccessControl(t *testing.T) {
	var (
		req = require.New(t)

		// uncomment to enable sql conn debugging
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		//ctx    = context.Background()
		s, err = sqlite.ConnectInMemoryWithDebug(ctx)
	)

	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		nextIDi uint64 = 1
		nextID         = func() uint64 {
			nextIDi++
			return nextIDi
		}

		modConf = types.ModuleConfig{DAL: types.ModuleConfigDAL{ConnectionID: 1}}

		user         = &sysTypes.User{ID: nextID()}
		ns           = &types.Namespace{ID: nextID()}
		mod1         = &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "mod one", Config: modConf}
		mod2         = &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "mod two", Config: modConf}
		mod1strField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod1.ID, Name: "str", Kind: "String"}
		mod2refField = &types.ModuleField{
			ID:          nextID(),
			NamespaceID: ns.ID,
			ModuleID:    mod2.ID,
			Name:        "ref",
			Kind:        "Record",
			Options: types.ModuleFieldOptions{
				"moduleID": mod1.ID,
			},
		}

		testerRole = &sysTypes.Role{Name: "tester", ID: nextID()}

		svc = makeTestRecordService(t,
			user,
			ns,
			mod1,
			mod2,
			mod1strField,
			mod2refField,
		)

		mod1rec1 = &types.Record{
			NamespaceID: ns.ID,
			ModuleID:    mod1.ID,
			Values: types.RecordValueSet{
				&types.RecordValue{Name: "str", Value: "abc"},
			},
		}

		mod2rec1 = &types.Record{NamespaceID: ns.ID, ModuleID: mod2.ID}
	)

	svc.validator = defaultValidator(svc)

	t.Log("inform rbac service about new roles")
	svc.rbacSvc.UpdateRoles(
		rbac.CommonRole.Make(testerRole.ID, testerRole.Name),
	)

	t.Log("log-in with test user ")
	ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(user.ID, testerRole.ID))

	_ = mod2rec1

	{
		t.Log("creating record on 1st module; should failed because we do not have permissions to create records")
		_, _, err = svc.Create(ctx, mod1rec1)
		req.EqualError(err, "not allowed to create records")

		t.Logf("granting permissions to create records on this module")
		req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, mod1.RbacResource(), "record.create")))

		t.Log("retry creating record on 1st module; should fail because we do not have permissions to update field")
		_, _, err = svc.Create(ctx, mod1rec1)
		req.Error(err)
		req.True(types.IsRecordValueErrorSet(err).HasKind("updateDenied"))

		t.Logf("granting permissions to update records values on module field")
		req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, mod1strField.RbacResource(), "record.value.update")))

		t.Log("retry creating record on 1st module; should succeed")
		mod1rec1, _, err = svc.Create(ctx, mod1rec1)
		req.NoError(err)
	}
	{
		t.Log("can record be read")
		_, _, err = svc.FindByID(ctx, mod1rec1.NamespaceID, mod1rec1.ModuleID, mod1rec1.ID)
		req.EqualError(err, "not allowed to read this record")
	}
	{
		t.Log("link 2nd record to 1st one")
		mod2rec1.Values = mod2rec1.Values.Set(&types.RecordValue{Name: "ref", Value: fmt.Sprintf("%d", mod1rec1.ID)})

		t.Log("create record on 2nd module with ref to record on the 1st module; must fail, no create perm")
		_, _, err = svc.Create(ctx, mod2rec1)
		req.EqualError(err, "not allowed to create records")

		t.Log("grant record.create on namespace level")
		req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, types.ModuleRbacResource(ns.ID, 0), "record.create")))

		t.Log("grant record.value.update on namespace level")
		req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, types.ModuleFieldRbacResource(ns.ID, 0, 0), "record.value.update")))

		t.Log("create record on 2nd module with ref to record on the 1st module; most fail, not allowed to read (referenced) mod1rec1")
		_, _, err = svc.Create(ctx, mod2rec1)
		req.EqualError(err, "invalid record value input")

		t.Log("grant read on record")
		req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, mod1rec1.RbacResource(), "read")))

		t.Log("create record on 2nd module with ref to record on the 1st module")
		mod2rec1, _, err = svc.Create(ctx, mod2rec1)
		verifyRecErrSet(t, err)
	}
	{
		t.Log("update record on 2nd module with unchanged values; must fail, no update permissions")
		_, _, err = svc.Update(ctx, mod2rec1)
		req.EqualError(err, "not allowed to update this record")

		t.Log("grant update on namespace level")
		req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, types.RecordRbacResource(ns.ID, 0, 0), "update")))

		t.Log("update record on 2nd module with unchanged values")
		mod2rec1, _, err = svc.Update(ctx, mod2rec1)
		verifyRecErrSet(t, err)

		t.Log("update record on 2nd module with unchanged values; unset record value")
		mod2rec1.Values = nil
		mod2rec1, _, err = svc.Update(ctx, mod2rec1)
		verifyRecErrSet(t, err)

		t.Log("link 2nd record to 1st one again")
		mod2rec1.Values = mod2rec1.Values.Set(&types.RecordValue{Name: "ref", Value: fmt.Sprintf("%d", mod1rec1.ID)})
		mod2rec1, _, err = svc.Update(ctx, mod2rec1)
		verifyRecErrSet(t, err)
	}
	{
		t.Log("revoke read on record")
		req.NoError(svc.rbacSvc.Grant(ctx, rbac.DenyRule(testerRole.ID, mod1rec1.RbacResource(), "read")))

		t.Log("link 2nd record to 1st one again but w/o permissions; must work, value did not change")
		mod2rec1.Values = mod2rec1.Values.Set(&types.RecordValue{Name: "ref", Value: fmt.Sprintf("%d", mod1rec1.ID)})
		mod2rec1, _, err = svc.Update(ctx, mod2rec1)
		verifyRecErrSet(t, err)
	}
}

func TestRecord_searchAccessControl(t *testing.T) {
	var (
		err error
		req = require.New(t)
		ctx = context.Background()

		nextIDi uint64 = 1
		nextID         = func() uint64 {
			nextIDi++
			return nextIDi
		}

		modConf = types.ModuleConfig{DAL: types.ModuleConfigDAL{ConnectionID: 1}}

		user     = &sysTypes.User{ID: nextID()}
		ns       = &types.Namespace{ID: nextID()}
		mod      = &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "mod one", Config: modConf}
		strField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "str", Kind: "String"}

		svc = makeTestRecordService(t,
			user,
			ns,
			mod,
			strField,
		)

		testerRole = &sysTypes.Role{Name: "tester", ID: nextID()}

		rr   = make([]*types.Record, 10)
		hits []*types.Record

		f = types.RecordFilter{
			ModuleID:    mod.ID,
			NamespaceID: ns.ID,
		}

		// search for module's model to be used as a filter
		model = ModelRef(svc, mod.ID)
	)

	req.NotNil(model)

	t.Log("create test records")
	for i := 0; i < cap(rr); i++ {
		rr[i] = &types.Record{
			ID:          nextID(),
			NamespaceID: ns.ID,
			ModuleID:    mod.ID,
		}

		req.NoError(svc.dal.Create(ctx, model.ToFilter(), nil, rr[i]))
	}

	t.Log("inform rbac service about new roles")
	svc.rbacSvc.UpdateRoles(
		rbac.CommonRole.Make(testerRole.ID, testerRole.Name),
	)

	t.Log("log-in with test user ")
	ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(user.ID, testerRole.ID))
	req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, mod.RbacResource(), "records.search")))

	t.Log("search for the newly created records; should not find any (all denied)")
	f.IncTotal = true
	hits, f, err = svc.Find(ctx, f)
	req.NoError(err)
	req.Len(hits, 0)
	req.Equal(uint(0), f.Total)

	t.Log("allow read access for two records")
	req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, rr[3].RbacResource(), "read")))
	req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(testerRole.ID, rr[6].RbacResource(), "read")))

	t.Log("search for the newly created records; should find 2 we're allowed to read")
	f.IncTotal = true
	hits, f, err = svc.Find(ctx, f)
	req.NoError(err)
	req.Len(hits, 2)
	req.Equal(uint(2), f.Total)
}

func TestRecord_contextualRolesAccessControl(t *testing.T) {
	var (
		err error
		req = require.New(t)
		ctx = context.Background()

		//log = zap.NewNop()
		log = logger.MakeDebugLogger()

		nextIDi uint64 = 1
		nextID         = func() uint64 {
			nextIDi++
			return nextIDi
		}

		user = &sysTypes.User{ID: nextID()}
		ns   = &types.Namespace{ID: nextID()}

		mod = &types.Module{
			ID:          nextID(),
			NamespaceID: ns.ID,
			Name:        "mod one",
			Config: types.ModuleConfig{
				DAL: types.ModuleConfigDAL{ConnectionID: 1},
			},
		}

		numField  = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "num", Kind: "String"}
		boolField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "yes", Kind: "String"}

		svc = makeTestRecordService(t,
			log,
			user,
			ns,
			mod,
			numField,
			boolField,
		)

		baseRole   = &sysTypes.Role{Name: "base", ID: nextID()}
		ownerRole  = &sysTypes.Role{Name: "owner", ID: nextID()}
		truthyRole = &sysTypes.Role{Name: "whenBoolTrue", ID: nextID()}
		tttRole    = &sysTypes.Role{Name: "whenNum333", ID: nextID()}

		rr   = make([]*types.Record, 10)
		hits []*types.Record

		f = types.RecordFilter{
			ModuleID:    mod.ID,
			NamespaceID: ns.ID,
		}

		// search for module's model to be used as a filter
		model = ModelRef(svc, mod.ID)

		// setting up rbac context role expression parser
		roleCheckFnMaker = func(expression string) func(scope map[string]interface{}) bool {
			p := expr.NewParser()

			return func(scope map[string]interface{}) bool {
				v, _ := expr.NewVars(scope)
				if e, err := p.Parse(expression); err != nil {
					t.Logf("could not parse expression: %v", err)
				} else if c, err := e.Test(ctx, v); err != nil {
					t.Logf("could not exec expression: %v", err)
				} else {
					return c
				}

				return false
			}

		}
	)

	t.Log("create test records")
	for i := 0; i < cap(rr); i++ {
		rr[i] = &types.Record{
			ID:          nextID(),
			NamespaceID: ns.ID,
			ModuleID:    mod.ID,
		}

		rr[i].SetModule(mod)

		if i%2 == 0 {
			// let's own half of the records
			rr[i].OwnedBy = user.ID
		}

		if i%3 == 0 {
			// set 333 to num on every 3rd record
			rr[i].Values = rr[i].Values.Set(&types.RecordValue{Name: "num", Value: "333"})
		}

		if i >= 5 {
			// and set true to bool field on the last 5 records
			rr[i].Values = rr[i].Values.Set(&types.RecordValue{Name: "yes", Value: "1"})
		}

		req.NoError(svc.dal.Create(ctx, model.ToFilter(), nil, rr[i]))

	}

	// result
	// i      0 1 2 3 4 5 6 7 8 9
	// --------------------------
	// owner  x   x   x   x   x
	//  bool            x x x x x
	//   333  x     x     x     x
	// --------------------------
	// read:  x   x x x x x x x x (all but one)

	t.Log("inform rbac service about new roles")
	svc.rbacSvc.UpdateRoles(
		rbac.CommonRole.Make(baseRole.ID, baseRole.Name),
		rbac.MakeContextRole(ownerRole.ID, ownerRole.Name, roleCheckFnMaker("resource.ownedBy == userID"), types.RecordResourceType),
		rbac.MakeContextRole(truthyRole.ID, truthyRole.Name, roleCheckFnMaker(`has(resource.values, "yes") ? resource.values.yes : false`), types.RecordResourceType),
		rbac.MakeContextRole(tttRole.ID, tttRole.Name, roleCheckFnMaker(`has(resource.values, "num") ? resource.values.num == 333 : false`), types.RecordResourceType),
	)

	req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(baseRole.ID, types.ModuleRbacResource(0, 0), "records.search")))

	t.Log("log-in with test user")
	ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(user.ID, baseRole.ID))

	t.Log("expecting not find any (all denied)")
	hits, _, err = svc.Find(ctx, f)
	req.NoError(err)
	req.Len(hits, 0)

	t.Log("expecting to find 5 records (owned by us)")
	req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(ownerRole.ID, types.RecordRbacResource(0, 0, 0), "read")))
	hits, _, err = svc.Find(ctx, f)
	req.NoError(err)
	req.Len(hits, 5)

	t.Log("expecting to find 2 records (owned by us and with true value for 'yes' field)")
	req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(truthyRole.ID, types.RecordRbacResource(0, 0, 0), "read")))
	hits, _, err = svc.Find(ctx, f)
	req.NoError(err)
	req.Len(hits, 8)

	t.Log("expecting to find 2 records (owned by us and with true value for 'yes' field + 333 for num)")
	req.NoError(svc.rbacSvc.Grant(ctx, rbac.AllowRule(tttRole.ID, types.RecordRbacResource(0, 0, 0), "read")))
	hits, _, err = svc.Find(ctx, f)
	req.NoError(err)
	req.Len(hits, 9)
}

func TestSetRecordOwner(t *testing.T) {
	var (
		req = require.New(t)

		// uncomment to enable sql conn debugging
		//ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		ctx    = context.Background()
		s, err = sqlite.ConnectInMemoryWithDebug(ctx)
	)

	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		rbacService = rbac.NewServiceMust(ctx, zap.NewNop(), s, rbac.Config{
			Synchronous:        true,
			DecayInterval:      time.Hour * 2,
			CleanupInterval:    time.Hour * 2,
			ReindexInterval:    time.Hour * 2,
			IndexFlushInterval: time.Hour * 2,
			RuleStorage:        s,
			RoleStorage:        s,
		})
		ac = &accessControl{rbac: rbacService}

		invoker          = &sysTypes.User{ID: 1001}
		originalOwner    = &sysTypes.User{ID: 2001}
		alternativeOwner = &sysTypes.User{ID: 2002}

		role = &sysTypes.Role{Name: "role-with-ownership-change-permission", ID: 3000}

		mod      = &types.Module{ID: 3, NamespaceID: 71624}
		old, upd *types.Record
	)

	rbacService.UpdateRoles(
		rbac.CommonRole.Make(role.ID, role.Handle),
	)

	req.NoError(rbacService.Grant(ctx, rbac.AllowRule(role.ID, types.RecordRbacResource(0, 0, 0), "owner.manage")))
	req.NoError(rbacService.Grant(ctx, rbac.AllowRule(role.ID, types.ModuleRbacResource(0, 0), "owned-record.create")))
	req.NoError(store.CreateUser(ctx, s, invoker, originalOwner, alternativeOwner))

	t.Run("invalid input", func(t *testing.T) {
		old, upd = nil, nil
		verifyRecErrSet(t, SetRecordOwner(ctx, ac, s, old, upd, invoker.ID))
	})

	t.Run("new record owner is invoker", func(t *testing.T) {
		old, upd = nil, &types.Record{ID: 1, NamespaceID: 2, ModuleID: mod.ID}
		upd.SetModule(mod)

		// this must work, invoker can always set owner to self on a new record
		verifyRecErrSet(t, SetRecordOwner(ctx, ac, s, old, upd, invoker.ID))

		req.Equal(upd.OwnedBy, invoker.ID)
	})

	t.Run("deny setting to alternative owner on create", func(t *testing.T) {
		old, upd = nil, &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: alternativeOwner.ID}
		upd.SetModule(mod)

		// this must work, invoker can always set owner to self on a new record
		if rvse := SetRecordOwner(ctx, ac, s, old, upd, invoker.ID); rvse.IsValid() {
			t.Fatalf("expecting error")
		}
	})

	t.Run("deny setting to alternative owner on update", func(t *testing.T) {
		old = &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: originalOwner.ID}
		old.SetModule(mod)
		upd = &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: alternativeOwner.ID}
		upd.SetModule(mod)

		// this must work, invoker can always set owner to self on a new record
		if rvse := SetRecordOwner(ctx, ac, s, old, upd, invoker.ID); rvse.IsValid() {
			t.Fatalf("expecting error")
		}
	})

	t.Run("allow setting to alternative owner on create", func(t *testing.T) {
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(invoker.ID, role.ID))
		old, upd = nil, &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: alternativeOwner.ID}
		upd.SetModule(mod)

		verifyRecErrSet(t, SetRecordOwner(ctx, ac, s, old, upd, invoker.ID))
		req.Equal(upd.OwnedBy, alternativeOwner.ID)
	})

	t.Run("allow setting to alternative owner on update", func(t *testing.T) {
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(invoker.ID, role.ID))
		old = &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: originalOwner.ID}
		old.SetModule(mod)
		upd = &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: alternativeOwner.ID}
		upd.SetModule(mod)

		verifyRecErrSet(t, SetRecordOwner(ctx, ac, s, old, upd, invoker.ID))
		req.Equal(upd.OwnedBy, alternativeOwner.ID)
	})

	t.Run("allow setting to new owner to zerp", func(t *testing.T) {
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(invoker.ID, role.ID))
		old = &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: originalOwner.ID}
		old.SetModule(mod)
		upd = &types.Record{ID: 1, NamespaceID: 2, ModuleID: 3, OwnedBy: 0}
		upd.SetModule(mod)

		verifyRecErrSet(t, SetRecordOwner(ctx, ac, s, old, upd, invoker.ID))
		req.Equal(upd.OwnedBy, invoker.ID)
	})
}

func TestRecordReportToDalPipeline(t *testing.T) {
	mod := &types.Module{
		ID:          10,
		NamespaceID: 42,
	}

	t.Run("no additional metrics", func(t *testing.T) {
		pp, _, err := recordReportToDalPipeline(mod, "", "created_at", "")
		require.NoError(t, err)

		require.Len(t, pp, 2)

		agg := pp[1].(*dal.Aggregate)
		require.Len(t, agg.OutAttributes, 1)
		require.Equal(t, "count", agg.OutAttributes[0].Identifier)
	})

	t.Run("additional metrics with alias", func(t *testing.T) {
		pp, _, err := recordReportToDalPipeline(mod, "MAX(numbers) AS   something", "created_at", "")
		require.NoError(t, err)

		require.Len(t, pp, 2)

		agg := pp[1].(*dal.Aggregate)
		require.Len(t, agg.OutAttributes, 2)
		require.Equal(t, "something", agg.OutAttributes[1].Identifier)
		require.Equal(t, "MAX(numbers)", agg.OutAttributes[1].RawExpr)
	})

	t.Run("additional metrics without alias", func(t *testing.T) {
		pp, _, err := recordReportToDalPipeline(mod, "MAX(numbers)", "created_at", "")
		require.NoError(t, err)

		require.Len(t, pp, 2)

		agg := pp[1].(*dal.Aggregate)
		require.Len(t, agg.OutAttributes, 2)
		require.Equal(t, "MAX(numbers)", agg.OutAttributes[1].Identifier)
		require.Equal(t, "MAX(numbers)", agg.OutAttributes[1].RawExpr)
	})
}
