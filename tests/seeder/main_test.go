package seeder

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/app"
	composeService "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/label"
	lType "github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/seeder"
	"github.com/cortezaproject/corteza-server/store"
	sTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/stretchr/testify/require"
)

type (
	helper struct {
		t *testing.T
		a *require.Assertions

		ctx context.Context
	}
)

var (
	testApp      *app.CortezaApp
	DefaultStore = seeder.DefaultStore
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

func InitTestApp() {
	if testApp == nil {
		ctx := logger.ContextWithValue(cli.Context(), logger.MakeDebugLogger())

		testApp = helpers.NewIntegrationTestApp(ctx, func(app *app.CortezaApp) (err error) {
			DefaultStore = app.Store
			return nil
		})
	}
}

func TestMain(m *testing.M) {
	InitTestApp()
	os.Exit(m.Run())
}

func newHelper(t *testing.T) helper {
	h := helper{
		t: t,
		a: require.New(t),

		ctx: context.Background(),
	}

	return h
}

// Unwraps error before it passes it to the tester
func (h helper) noError(err error) {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	h.a.NoError(err)
}

func (h helper) setLabel(res label.LabeledResource, name, value string) {
	h.a.NoError(store.UpsertLabel(h.ctx, DefaultStore, &lType.Label{
		Kind:       res.LabelResourceKind(),
		ResourceID: res.LabelResourceID(),
		Name:       name,
		Value:      value,
	}))
}

func (h helper) clearUsers() {
	h.noError(store.TruncateUsers(context.Background(), DefaultStore))
}

func (h helper) lookupUsers() (sTypes.UserSet, error) {
	filter := sTypes.UserFilter{Labels: map[string]string{seeder.FakeDataLabel: seeder.FakeDataLabel}}
	users, _, err := DefaultStore.SearchUsers(h.ctx, filter)
	h.noError(err)

	return users, err
}

func (h helper) clearNamespaces() {
	h.noError(store.TruncateComposeNamespaces(context.Background(), DefaultStore))
}

func (h helper) makeNamespace(name string) *types.Namespace {
	ns := &types.Namespace{Name: name, Slug: name}
	ns.ID = id.Next()
	ns.CreatedAt = time.Now()
	h.noError(store.CreateComposeNamespace(context.Background(), DefaultStore, ns))
	return ns
}

func (h helper) clearModules() {
	h.clearNamespaces()
	h.noError(store.TruncateComposeModules(context.Background(), DefaultStore))
	h.noError(store.TruncateComposeModuleFields(context.Background(), DefaultStore))
}

func (h helper) makeModule(ns *types.Namespace, name string, ff ...*types.ModuleField) *types.Module {
	m := h.createModule(&types.Module{
		Name:        name,
		NamespaceID: ns.ID,
		Fields:      ff,
		CreatedAt:   time.Now(),
	})
	composeService.DefaultModule.ReloadDALModels(context.Background())
	return m
}

func (h helper) createModule(res *types.Module) *types.Module {
	res.ID = id.Next()
	res.CreatedAt = time.Now()
	h.noError(store.CreateComposeModule(h.ctx, DefaultStore, res))

	_ = res.Fields.Walk(func(f *types.ModuleField) error {
		f.ID = id.Next()
		f.ModuleID = res.ID
		f.CreatedAt = time.Now()
		return nil
	})

	h.noError(store.CreateComposeModuleField(h.ctx, DefaultStore, res.Fields...))

	return res
}

func (h helper) lookupModuleByID(ID uint64) *types.Module {
	res, err := store.LookupComposeModuleByID(h.ctx, DefaultStore, ID)
	h.noError(err)

	res.Fields, _, err = store.SearchComposeModuleFields(h.ctx, DefaultStore, types.ModuleFieldFilter{ModuleID: []uint64{ID}})
	h.noError(err)

	return res
}
