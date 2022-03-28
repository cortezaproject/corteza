package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestModules(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = sqlite.ConnectInMemory(ctx)

		// ctx    = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		// s, err = sqlite.ConnectInMemoryWithDebug(ctx)

		namespaceID = nextID()
		ns          *types.Namespace
	)

	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	if err = store.Upgrade(ctx, zap.NewNop(), s); err != nil {
		t.Fatalf("failed to upgrade store: %v", err)
	}

	if err = s.TruncateComposeNamespaces(ctx); err != nil {
		t.Fatalf("failed to truncate compose namespaces: %v", err)
	}

	if err = s.TruncateComposeModules(ctx); err != nil {
		t.Fatalf("failed to truncate compose modules: %v", err)
	}

	ns = &types.Namespace{Name: "testing", ID: namespaceID, CreatedAt: *now()}
	if err = store.CreateComposeNamespace(ctx, s, ns); err != nil {
		t.Fatalf("failed to seed namespaces: %v", err)
	}

	t.Run("crud", func(t *testing.T) {
		req := require.New(t)
		svc := module{
			store:    s,
			ac:       &accessControl{rbac: &rbac.ServiceAllowAll{}},
			eventbus: eventbus.New(),
		}
		res, err := svc.Create(ctx, &types.Module{Name: "My first module", NamespaceID: namespaceID})
		req.NoError(err)
		req.NotNil(res)

		res, err = svc.FindByID(ctx, namespaceID, res.ID)
		req.NoError(err)
		req.NotNil(res)

		res, err = svc.FindByHandle(ctx, namespaceID, res.Handle)
		req.NoError(err)
		req.NotNil(res)

		res.Name = "Changed"
		res, err = svc.Update(ctx, res)
		req.NoError(err)
		req.NotNil(res)
		req.NotNil(res.UpdatedAt)
		req.Equal(res.Name, "Changed")

		res, err = svc.FindByID(ctx, namespaceID, res.ID)
		req.NoError(err)
		req.NotNil(res)
		req.Equal(res.Name, "Changed")

		err = svc.DeleteByID(ctx, namespaceID, res.ID)
		req.NoError(err)
		req.NotNil(res)

		// this works because we're allowed to do everything
		res, err = svc.FindByID(ctx, namespaceID, res.ID)
		req.NoError(err)
		req.NotNil(res)
		req.NotNil(res.DeletedAt)
	})

	t.Run("labels", func(t *testing.T) {
		t.Run("search", func(t *testing.T) {
			req := require.New(t)
			svc := module{
				store:    s,
				ac:       &accessControl{rbac: &rbac.ServiceAllowAll{}},
				eventbus: eventbus.New(),
				locale:   ResourceTranslationsManager(locale.Static()),
			}

			makeModule := func(n ...string) *types.Module {
				mod := &types.Module{
					NamespaceID: namespaceID,
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

			findModules := func(labels map[string]string, IDs []uint64) types.ModuleSet {
				f := types.ModuleFilter{NamespaceID: namespaceID, Labels: labels, ModuleID: IDs}
				set, _, err := svc.Find(ctx, f)
				req.NoError(err)

				return set
			}

			makeModule("labeled module 1", "label1", "value1", "label2", "value2")
			m2 := makeModule("labeled module 2", "label1", "value1")
			m3 := makeModule("labeled module 3")

			// return all -- no label/ID filter, return all
			req.Len(findModules(nil, nil), 3)

			// return 2 - both that have label1=valu1
			req.Len(findModules(map[string]string{"label1": "value1"}, nil), 2)

			// return 0 - none have foo=foo
			req.Len(findModules(map[string]string{"missing": "missing"}, nil), 0)

			// one has label2=value2
			req.Len(findModules(map[string]string{"label2": "value2"}, nil), 1)

			// explicit by ID and label
			req.Len(findModules(map[string]string{"label1": "value1"}, []uint64{m2.ID}), 1)

			// none with this combo
			req.Len(findModules(map[string]string{"foo": "foo"}, []uint64{m3.ID}), 0)

			// one with explicit ID (regression) and nil for label filter
			req.Len(findModules(nil, []uint64{m3.ID}), 1)

			// one with explicit ID (regression) and empty map for label filter
			req.Len(findModules(map[string]string{}, []uint64{m3.ID}), 1)
		})

		t.Run("CRUD", func(t *testing.T) {
			req := require.New(t)
			svc := module{
				store:    s,
				ac:       &accessControl{rbac: &rbac.ServiceAllowAll{}},
				eventbus: eventbus.New(),
			}

			findAndReturnLabel := func(id uint64) map[string]string {
				res, err := svc.FindByID(ctx, namespaceID, id)
				req.NoError(err)
				req.NotNil(res)
				return res.Labels
			}

			// create unlabeled module
			res, err := svc.Create(ctx, &types.Module{Name: "unLabeledIDs", NamespaceID: namespaceID})
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

			res, err = svc.Create(ctx, &types.Module{Name: "LabeledIDs", NamespaceID: namespaceID, Labels: map[string]string{"label2": "2nd"}})
			req.NoError(err)
			req.NotNil(res)
			req.Contains(res.Labels, "label2")

			// must contain the added label
			req.Contains(findAndReturnLabel(res.ID), "label2")

			// update with Labels:nil (should keep labels intact)
			res.Labels = nil
			res, err = svc.Update(ctx, res)
			req.NoError(err)

			req.Contains(findAndReturnLabel(res.ID), "label2")

			// update with Labels:empty-map (should remove all labels)
			res.Labels = map[string]string{}
			res, err = svc.Update(ctx, res)
			req.NoError(err)

			req.Empty(findAndReturnLabel(res.ID))
		})
	})
}
