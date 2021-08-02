package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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

func TestProcUpdateOwnerPreservation(t *testing.T) {
	var (
		ctx = context.Background()
		a   = assert.New(t)

		svc = record{
			sanitizer: values.Sanitizer(),
			validator: values.Validator(),
			formatter: values.Formatter(),
		}

		mod = &types.Module{
			Fields: types.ModuleFieldSet{},
		}

		oldRec = &types.Record{
			OwnedBy: 1,
			Values:  types.RecordValueSet{},
		}
		newRec = &types.Record{
			OwnedBy: 0,
			Values:  types.RecordValueSet{},
		}
	)

	svc.procUpdate(ctx, 10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(1))
	svc.procUpdate(ctx, 10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(1))
}

func TestProcUpdateOwnerChanged(t *testing.T) {
	var (
		ctx = context.Background()
		a   = assert.New(t)

		svc = record{
			sanitizer: values.Sanitizer(),
			validator: values.Validator(),
			formatter: values.Formatter(),
		}

		mod = &types.Module{
			Fields: types.ModuleFieldSet{},
		}

		oldRec = &types.Record{
			OwnedBy: 1,
			Values:  types.RecordValueSet{},
		}
		newRec = &types.Record{
			OwnedBy: 9,
			Values:  types.RecordValueSet{},
		}
	)

	a.True(svc.procUpdate(ctx, 10, mod, newRec, oldRec).IsValid())
	a.Equal(newRec.OwnedBy, uint64(9))
	a.True(svc.procUpdate(ctx, 10, mod, newRec, oldRec).IsValid())
	a.Equal(newRec.OwnedBy, uint64(9))
}

func TestRecord_boolFieldPermissionIssueKBR(t *testing.T) {
	var (
		req = require.New(t)

		// uncomment to enable sql conn debugging
		//ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		ctx    = context.Background()
		s, err = sqlite3.ConnectInMemoryWithDebug(ctx)
	)

	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))
	req.NoError(store.TruncateComposeRecords(ctx, s, nil))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		rbacService = rbac.NewService(
			//zap.NewNop(),
			logger.MakeDebugLogger(),
			nil,
		)
		ac = &accessControl{rbac: rbacService}

		svc = &record{
			sanitizer: values.Sanitizer(),
			formatter: values.Formatter(),
			ac:        ac,
			store:     s,
		}

		userID      = nextID()
		ns          = &types.Namespace{ID: nextID()}
		mod         = &types.Module{ID: nextID(), NamespaceID: ns.ID}
		stringField = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, Name: "string", Kind: "String"}
		boolField   = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, Name: "bool", Kind: "Boolean"}

		authRoleID uint64 = 1

		readerRole = &sysTypes.Role{Name: "reader", ID: nextID()}
		writerRole = &sysTypes.Role{Name: "writer", ID: nextID()}

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

	req.NoError(store.CreateComposeNamespace(ctx, s, ns))
	req.NoError(store.CreateComposeModule(ctx, s, mod))
	req.NoError(store.CreateComposeModuleField(ctx, s, stringField, boolField))

	rbacService.UpdateRoles(
		rbac.CommonRole.Make(readerRole.ID, readerRole.Name),
		rbac.CommonRole.Make(writerRole.ID, writerRole.Name),
		rbac.AuthenticatedRole.Make(authRoleID, "authenticated"),
	)

	rbacService.Grant(ctx,
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
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, writerRole.ID, authRoleID))

		recChecked, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: valChecked})
		req.NoError(errors.Unwrap(err))

		req.NotNil(recChecked.Values.Get("bool", 0), "should be checked")
		req.Equal("1", recChecked.Values.Get("bool", 0).Value)

		recUnchecked, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: valUnchecked})
		req.NoError(err)

		req.Nil(recUnchecked.Values.Get("bool", 0))

		// *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** ***

		// security context w/ reader role
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, readerRole.ID, authRoleID))

		recChecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
		}

		recChecked, err = svc.Update(ctx, recChecked)
		req.NoError(err)

		req.NotNil(recChecked.Values.Get("bool", 0), "should still be checked")
		req.Equal("1", recChecked.Values.Get("bool", 0).Value)

		recUnchecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
		}

		recUnchecked, err = svc.Update(ctx, recUnchecked)
		req.NoError(err)
		req.Nil(recUnchecked.Values.Get("bool", 0), "should not be checked anymore")

		// *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** ***

		// security context w/ writer role
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, writerRole.ID, authRoleID))

		recChecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
			&types.RecordValue{Name: "bool", Value: "1"},
		}

		recChecked, err = svc.Update(ctx, recChecked)
		req.NoError(err)

		req.NotNil(recChecked.Values.Get("bool", 0), "should checked again")
		req.Equal("1", recChecked.Values.Get("bool", 0).Value, "should checked again")
	}
}

func TestRecord_defValueFieldPermissionIssue(t *testing.T) {
	var (
		req    = require.New(t)
		ctx    = context.Background()
		s, err = sqlite3.ConnectInMemoryWithDebug(ctx)
	)

	t.Log("setting up the test environment")

	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))
	req.NoError(store.TruncateComposeRecords(ctx, s, nil))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		rbacService = rbac.NewService(
			//zap.NewNop(),
			logger.MakeDebugLogger(),
			nil,
		)
		ac = &accessControl{rbac: rbacService}

		svc = &record{
			sanitizer: values.Sanitizer(),
			formatter: values.Formatter(),
			ac:        ac,
			store:     s,
		}

		ns            = &types.Namespace{ID: nextID()}
		mod           = &types.Module{ID: nextID(), NamespaceID: ns.ID}
		writableField = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, NamespaceID: ns.ID, Name: "writable", Kind: "String", DefaultValue: types.RecordValueSet{{Value: "def-w"}}}
		readableField = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, NamespaceID: ns.ID, Name: "readable", Kind: "String", DefaultValue: types.RecordValueSet{{Value: "def-r"}}}

		userID     = nextID()
		authRoleID = nextID()
		editorRole = &sysTypes.Role{Name: "editor", ID: nextID()}

		recPartial *types.Record

		valueExtractor = func(rec *types.Record, ff ...string) (out string) {
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

	req.NoError(store.CreateComposeNamespace(ctx, s, ns))
	req.NoError(store.CreateComposeModule(ctx, s, mod))
	req.NoError(store.CreateComposeModuleField(ctx, s, writableField, readableField))

	t.Log("setting up security")

	rbacService.UpdateRoles(
		rbac.CommonRole.Make(editorRole.ID, editorRole.Name),
		rbac.AuthenticatedRole.Make(authRoleID, "authenticated"),
	)

	rbacService.Grant(ctx,
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

		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, authRoleID))

		recPartial, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{}})

		req.NoError(errors.Unwrap(err))
		req.Equal("<def-w><def-r>", valueExtractor(recPartial, "writable", "readable"))

		t.Log("creating record with w/o editor role (must be able to crate & update record and modify both fields)")

		recPartial, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{
			&types.RecordValue{Name: "writable", Value: "w"},
			&types.RecordValue{Name: "readable", Value: "r"},
		}})

		req.NoError(errors.Unwrap(err))
		req.Equal("<w><r>", valueExtractor(recPartial, "writable", "readable"))

		t.Log("updating record removing one of the values")

		recPartial.Values = types.RecordValueSet{&types.RecordValue{Name: "writable", Value: "w2"}}

		recPartial, err = svc.Update(ctx, recPartial)
		req.NoError(errors.Unwrap(err))
		req.Equal("<w2><NULL>", valueExtractor(recPartial, "writable", "readable"))
	}

	{
		t.Log("creating record with editor role (expecting defaults")

		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, authRoleID, editorRole.ID))

		recPartial, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{}})

		req.NoError(errors.Unwrap(err))
		req.Equal("<def-w><def-r>", valueExtractor(recPartial, "writable", "readable"))

		t.Log("creating record with editor role (must be able to crate & update record and modify both fields)")

		recPartial, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: types.RecordValueSet{
			// this is the def. value set
			&types.RecordValue{Name: "writable", Value: "def-w"},
			&types.RecordValue{Name: "readable", Value: "r"},
		}})

		req.NoError(errors.Unwrap(err))
		req.Equal("<def-w><r>", valueExtractor(recPartial, "writable", "readable"))
	}
}

func TestRecord_refAccessControl(t *testing.T) {
	var (
		req = require.New(t)

		// uncomment to enable sql conn debugging
		//ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		ctx    = context.Background()
		s, err = sqlite3.ConnectInMemoryWithDebug(ctx)
	)

	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))
	req.NoError(store.TruncateComposeRecords(ctx, s, nil))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		rbacService = rbac.NewService(
			//zap.NewNop(),
			logger.MakeDebugLogger(),
			nil,
		)
		ac = &accessControl{rbac: rbacService}

		svc = &record{
			sanitizer: values.Sanitizer(),
			formatter: values.Formatter(),
			ac:        ac,
			store:     s,
		}

		nextIDi uint64 = 1
		nextID         = func() uint64 {
			nextIDi++
			return nextIDi
		}

		userID       = nextID()
		ns           = &types.Namespace{ID: nextID()}
		mod1         = &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "mod one"}
		mod2         = &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "mod two"}
		mod1strField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod1.ID, Name: "str", Kind: "String"}
		mod2refField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod2.ID, Name: "ref", Kind: "Record"}

		testerRole = &sysTypes.Role{Name: "tester", ID: nextID()}

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

	t.Log("create test namespace")
	req.NoError(store.CreateComposeNamespace(ctx, s, ns))

	t.Log("create test modules")
	req.NoError(store.CreateComposeModule(ctx, s, mod1, mod2))

	t.Log("create test module fields")
	req.NoError(store.CreateComposeModuleField(ctx, s, mod1strField, mod2refField))

	t.Log("inform rbac service about new roles")
	rbacService.UpdateRoles(
		rbac.CommonRole.Make(testerRole.ID, testerRole.Name),
	)

	t.Log("log-in with test user ")
	ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, testerRole.ID))

	_ = mod2rec1

	{
		t.Log("creating record on 1st module; should failed because we do not have permissions to create records")
		_, err = svc.Create(ctx, mod1rec1)
		req.EqualError(err, "not allowed to create records")

		t.Logf("granting permissions to create records on this module")
		rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, mod1.RbacResource(), "record.create"))

		t.Log("retry creating record on 1st module; should fail because we do not have permissions to update field")
		_, err = svc.Create(ctx, mod1rec1)
		req.Error(err)
		req.True(types.IsRecordValueErrorSet(err).HasKind("updateDenied"))

		t.Logf("granting permissions to update records values on module field")
		rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, mod1strField.RbacResource(), "record.value.update"))

		t.Log("retry creating record on 1st module; should succeed")
		mod1rec1, err = svc.Create(ctx, mod1rec1)
		req.NoError(err)
	}
	{
		t.Log("can record be read")
		_, err = svc.FindByID(ctx, mod1rec1.NamespaceID, mod1rec1.ModuleID, mod1rec1.ID)
		req.EqualError(err, "not allowed to read this record")
	}
	{
		t.Log("link 2nd record to 1st one")
		mod2rec1.Values = mod2rec1.Values.Set(&types.RecordValue{Name: "ref", Value: fmt.Sprintf("%d", mod1rec1.ID)})

		t.Log("create record on 2nd module with ref to record on the 1st module; must fail, no create perm")
		_, err = svc.Create(ctx, mod2rec1)
		req.EqualError(err, "not allowed to create records")

		t.Log("grant record.create on namespace level")
		rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, types.ModuleRbacResource(ns.ID, 0), "record.create"))

		t.Log("grant record.value.update on namespace level")
		rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, types.ModuleFieldRbacResource(ns.ID, 0, 0), "record.value.update"))

		t.Log("create record on 2nd module with ref to record on the 1st module; most fail, not allowed to read (referenced) mod1rec1")
		_, err = svc.Create(ctx, mod2rec1)
		req.EqualError(err, "invalid record value input")

		t.Log("grant read on record")
		rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, mod1rec1.RbacResource(), "read"))

		t.Log("create record on 2nd module with ref to record on the 1st module")
		mod2rec1, err = svc.Create(ctx, mod2rec1)
		req.NoError(errors.Unwrap(err))
	}
	{
		t.Log("update record on 2nd module with unchanged values; must fail, no update permissions")
		_, err = svc.Update(ctx, mod2rec1)
		req.EqualError(err, "not allowed to update this record")

		t.Log("grant update on namespace level")
		rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, types.RecordRbacResource(ns.ID, 0, 0), "update"))

		t.Log("update record on 2nd module with unchanged values")
		mod2rec1, err = svc.Update(ctx, mod2rec1)
		req.NoError(errors.Unwrap(err))

		t.Log("update record on 2nd module with unchanged values; unset record value")
		mod2rec1.Values = nil
		mod2rec1, err = svc.Update(ctx, mod2rec1)
		req.NoError(errors.Unwrap(err))

		t.Log("link 2nd record to 1st one again")
		mod2rec1.Values = mod2rec1.Values.Set(&types.RecordValue{Name: "ref", Value: fmt.Sprintf("%d", mod1rec1.ID)})
		mod2rec1, err = svc.Update(ctx, mod2rec1)
		req.NoError(errors.Unwrap(err))
	}
	{
		t.Log("revoke read on record")
		rbacService.Grant(ctx, rbac.DenyRule(testerRole.ID, mod1rec1.RbacResource(), "read"))

		t.Log("link 2nd record to 1st one again but w/o permissions; must work, value did not change")
		mod2rec1.Values = mod2rec1.Values.Set(&types.RecordValue{Name: "ref", Value: fmt.Sprintf("%d", mod1rec1.ID)})
		mod2rec1, err = svc.Update(ctx, mod2rec1)
		req.NoError(errors.Unwrap(err))
	}
}

func TestRecord_searchAccessControl(t *testing.T) {
	var (
		req = require.New(t)

		// uncomment to enable sql conn debugging
		//ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		ctx    = context.Background()
		s, err = sqlite3.ConnectInMemoryWithDebug(ctx)
	)

	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))
	req.NoError(store.TruncateComposeRecords(ctx, s, nil))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		rbacService = rbac.NewService(
			//zap.NewNop(),
			logger.MakeDebugLogger(),
			nil,
		)
		ac = &accessControl{rbac: rbacService}

		svc = &record{
			ac:        ac,
			store:     s,
			sanitizer: values.Sanitizer(),
		}

		nextIDi uint64 = 1
		nextID         = func() uint64 {
			nextIDi++
			return nextIDi
		}

		userID = nextID()
		ns     = &types.Namespace{ID: nextID()}
		mod    = &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "mod one"}
		//strField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "str", Kind: "String"}

		testerRole = &sysTypes.Role{Name: "tester", ID: nextID()}

		rr   = make([]*types.Record, 10)
		hits []*types.Record

		f = types.RecordFilter{
			ModuleID:    mod.ID,
			NamespaceID: ns.ID,
		}
	)

	t.Log("create test namespace")
	req.NoError(store.CreateComposeNamespace(ctx, s, ns))

	t.Log("create test modules")
	req.NoError(store.CreateComposeModule(ctx, s, mod))

	//t.Log("create test module fields")
	//req.NoError(store.CreateComposeModuleField(ctx, s, strField))

	t.Log("create test records")
	for i := 0; i < cap(rr); i++ {
		rr[i] = &types.Record{
			ID:          nextID(),
			NamespaceID: ns.ID,
			ModuleID:    mod.ID,
		}

		req.NoError(store.CreateComposeRecord(ctx, s, mod, rr[i]))
	}

	t.Log("inform rbac service about new roles")
	rbacService.UpdateRoles(
		rbac.CommonRole.Make(testerRole.ID, testerRole.Name),
	)

	t.Log("log-in with test user ")
	ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, testerRole.ID))
	rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, mod.RbacResource(), "records.search"))

	t.Log("search for the newly created records; should not find any (all denied)")
	f.IncTotal = true
	hits, f, err = svc.Find(ctx, f)
	req.Len(hits, 0)
	req.Equal(uint(0), f.Total)

	t.Log("allow read access for two records")
	rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, rr[3].RbacResource(), "read"))
	rbacService.Grant(ctx, rbac.AllowRule(testerRole.ID, rr[6].RbacResource(), "read"))

	t.Log("search for the newly created records; should find 2 we're allowed to read")
	f.IncTotal = true
	hits, f, err = svc.Find(ctx, f)
	req.Len(hits, 2)
	req.Equal(uint(2), f.Total)
}

func TestRecord_contextualRolesAccessControl(t *testing.T) {
	var (
		req = require.New(t)

		// uncomment to enable sql conn debugging
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		//ctx    = context.Background()
		s, err = sqlite3.ConnectInMemoryWithDebug(ctx)
	)

	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))
	req.NoError(store.TruncateComposeRecords(ctx, s, nil))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		rbacService = rbac.NewService(
			//zap.NewNop(),
			logger.MakeDebugLogger(),
			nil,
		)
		ac = &accessControl{rbac: rbacService}

		svc = &record{
			ac:        ac,
			store:     s,
			sanitizer: values.Sanitizer(),
		}

		nextIDi uint64 = 1
		nextID         = func() uint64 {
			nextIDi++
			return nextIDi
		}

		userID    = nextID()
		ns        = &types.Namespace{ID: nextID()}
		mod       = &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "mod one"}
		numField  = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "num", Kind: "String"}
		boolField = &types.ModuleField{ID: nextID(), NamespaceID: ns.ID, ModuleID: mod.ID, Name: "yes", Kind: "String"}

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

	t.Log("create test namespace")
	req.NoError(store.CreateComposeNamespace(ctx, s, ns))

	mod.Fields = types.ModuleFieldSet{numField, boolField}
	numField.ModuleID = mod.ID
	boolField.ModuleID = mod.ID

	t.Log("create test modules")
	req.NoError(store.CreateComposeModule(ctx, s, mod))
	req.NoError(store.CreateComposeModuleField(ctx, s, numField, boolField))

	t.Log("create test records")
	for i := 0; i < cap(rr); i++ {
		rr[i] = &types.Record{
			ID:          nextID(),
			NamespaceID: ns.ID,
			ModuleID:    mod.ID,
		}

		if i%2 == 0 {
			// let's own half of the records
			rr[i].OwnedBy = userID
		}

		if i%3 == 0 {
			// set 333 to num on every 3rd record
			rr[i].Values = rr[i].Values.Set(&types.RecordValue{Name: "num", Value: "333"})
		}

		if i >= 5 {
			// and set true to bool field on the last 5 records
			rr[i].Values = rr[i].Values.Set(&types.RecordValue{Name: "yes", Value: "1"})
		}

		req.NoError(store.CreateComposeRecord(ctx, s, mod, rr[i]))
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
	rbacService.UpdateRoles(
		rbac.CommonRole.Make(baseRole.ID, baseRole.Name),
		rbac.MakeContextRole(ownerRole.ID, ownerRole.Name, roleCheckFnMaker("resource.ownedBy == userID"), types.RecordResourceType),
		rbac.MakeContextRole(truthyRole.ID, truthyRole.Name, roleCheckFnMaker(`has(resource.values, "yes") ? resource.values.yes : false`), types.RecordResourceType),
		rbac.MakeContextRole(tttRole.ID, tttRole.Name, roleCheckFnMaker(`has(resource.values, "num") ? resource.values.num == 333 : false`), types.RecordResourceType),
	)

	rbacService.Grant(ctx, rbac.AllowRule(baseRole.ID, types.ModuleRbacResource(0, 0), "records.search"))

	t.Log("log-in with test user")
	ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(userID, baseRole.ID))

	t.Log("expecting not find any (all denied)")
	hits, _, err = svc.Find(ctx, f)
	req.Len(hits, 0)

	t.Log("expecting to find 5 records (owned by us)")
	rbacService.Grant(ctx, rbac.AllowRule(ownerRole.ID, types.RecordRbacResource(0, 0, 0), "read"))
	hits, _, err = svc.Find(ctx, f)
	req.Len(hits, 5)

	t.Log("expecting to find 2 records (owned by us and with true value for 'yes' field)")
	rbacService.Grant(ctx, rbac.AllowRule(truthyRole.ID, types.RecordRbacResource(0, 0, 0), "read"))
	hits, _, err = svc.Find(ctx, f)
	req.Len(hits, 8)

	t.Log("expecting to find 2 records (owned by us and with true value for 'yes' field + 333 for num)")
	rbacService.Grant(ctx, rbac.AllowRule(tttRole.ID, types.RecordRbacResource(0, 0, 0), "read"))
	hits, _, err = svc.Find(ctx, f)
	req.Len(hits, 9)
}
