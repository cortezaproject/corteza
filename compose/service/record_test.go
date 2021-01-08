package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestGeneralValueSetValidation(t *testing.T) {
	var (
		req = require.New(t)

		svc = record{
			ac: AccessControl(&rbac.ServiceAllowAll{}),
		}
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
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "unknown", Value: "single"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.True(err != nil, "expecting generalValueSetValidation() to return an error, got nil")

	rvs = types.RecordValueSet{{Name: "single1", Value: "single"}, {Name: "single1", Value: "single2"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.Error(err)

	rvs = types.RecordValueSet{{Name: "multi1", Value: "multi1"}, {Name: "multi1", Value: "multi1"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "non numeric value"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.Error(err)

	rvs = types.RecordValueSet{{Name: "ref1", Value: "12345"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)

	rvs = types.RecordValueSet{{Name: "multiRef1", Value: "12345"}, {Name: "multiRef1", Value: "67890"}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)
	req.Len(rvs, 2, "expecting 2 record values after sanitization, got %d", len(rvs))

	rvs = types.RecordValueSet{{Name: "ref1", Value: ""}}
	err = svc.generalValueSetValidation(module, rvs)
	req.NoError(err)
}

func TestDefaultValueSetting(t *testing.T) {
	var (
		a = assert.New(t)

		svc = record{
			ac: AccessControl(&rbac.ServiceAllowAll{}),
		}
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

	out := svc.setDefaultValues(mod, nil)
	chk(out, "single", 0, "s")
	chk(out, "multi", 0, "m1")
	chk(out, "multi", 1, "m2")
}

func TestProcUpdateOwnerPreservation(t *testing.T) {
	var (
		ctx   = context.Background()
		store store.Storer

		a = assert.New(t)

		svc = record{
			sanitizer: values.Sanitizer(),
			validator: values.Validator(),
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

	svc.procUpdate(ctx, store, 10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(1))
	svc.procUpdate(ctx, store, 10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(1))
}

func TestProcUpdateOwnerChanged(t *testing.T) {
	var (
		ctx   = context.Background()
		store store.Storer

		a = assert.New(t)

		svc = record{
			sanitizer: values.Sanitizer(),
			validator: values.Validator(),
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

	svc.procUpdate(ctx, store, 10, mod, newRec, oldRec)
	a.Equal(newRec.OwnedBy, uint64(9))
	svc.procUpdate(ctx, store, 10, mod, newRec, oldRec)
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
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))
	req.NoError(store.TruncateComposeRecords(ctx, s, nil))
	req.NoError(store.TruncateRbacRules(ctx, s))

	var (
		rbacService = rbac.NewService(zap.NewNop(), s)
		ac          = AccessControl(rbacService)

		svc = record{
			sanitizer: values.Sanitizer(),
			validator: values.Validator(),
			ac:        ac,
			store:     s,
		}

		userID      = nextID()
		ns          = &types.Namespace{ID: nextID()}
		mod         = &types.Module{ID: nextID(), NamespaceID: ns.ID}
		stringField = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, Name: "string", Kind: "String"}
		boolField   = &types.ModuleField{ID: nextID(), ModuleID: mod.ID, Name: "bool", Kind: "Boolean"}

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

	req.NoError(store.CreateComposeNamespace(ctx, s, ns))
	req.NoError(store.CreateComposeModule(ctx, s, mod))
	req.NoError(store.CreateComposeModuleField(ctx, s, stringField, boolField))

	rbacService.Grant(ctx, ac.Whitelist(),
		rbac.AllowRule(readerRole.ID, mod.RBACResource(), "record.read"),
		rbac.AllowRule(readerRole.ID, mod.RBACResource(), "record.create"),
		rbac.AllowRule(readerRole.ID, mod.RBACResource(), "record.update"),
		rbac.AllowRule(readerRole.ID, stringField.RBACResource(), "record.value.read"),
		rbac.DenyRule(readerRole.ID, boolField.RBACResource(), "record.value.update"),

		rbac.AllowRule(writerRole.ID, mod.RBACResource(), "record.read"),
		rbac.AllowRule(writerRole.ID, mod.RBACResource(), "record.create"),
		rbac.AllowRule(writerRole.ID, mod.RBACResource(), "record.update"),
		rbac.AllowRule(writerRole.ID, stringField.RBACResource(), "record.value.read"),
		rbac.AllowRule(writerRole.ID, stringField.RBACResource(), "record.value.update"),
	)

	{
		// security context w/ writer role
		ctx = auth.SetIdentityToContext(ctx, auth.NewIdentity(userID, writerRole.ID))

		recChecked, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: valChecked})
		req.NoError(err)

		req.NotNil(recChecked.Values.Get("bool", 0))
		req.Equal("1", recChecked.Values.Get("bool", 0).Value)

		recUnchecked, err = svc.Create(ctx, &types.Record{ModuleID: mod.ID, NamespaceID: ns.ID, Values: valUnchecked})
		req.NoError(err)

		req.Nil(recUnchecked.Values.Get("bool", 0))

		// *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** ***

		// security context w/ writer role
		ctx = auth.SetIdentityToContext(ctx, auth.NewIdentity(userID, readerRole.ID))

		recChecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
		}

		recChecked, err = svc.Update(ctx, recChecked)
		req.NoError(err)

		req.NotNil(recChecked.Values.Get("bool", 0))
		req.Equal("1", recChecked.Values.Get("bool", 0).Value)

		recUnchecked.Values = types.RecordValueSet{
			&types.RecordValue{Name: "string", Value: "abc"},
		}

		recUnchecked, err = svc.Update(ctx, recUnchecked)
		req.NoError(err)
		req.Nil(recUnchecked.Values.Get("bool", 0))
	}

}
