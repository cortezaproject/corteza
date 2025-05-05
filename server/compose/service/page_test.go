package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
)

func TestPageDeleting(t *testing.T) {
	var (
		//ctx    = context.Background()
		//s, err = sqlite.ConnectInMemory(ctx)

		ctx    = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())
		s, err = sqlite.ConnectInMemoryWithDebug(ctx)

		namespaceID = nextID()
		ns          *types.Namespace

		pages = types.PageSet{
			// should be deleted w/o a problem
			&types.Page{ID: 1},
			&types.Page{ID: 2},
			&types.Page{ID: 3, SelfID: 2},
			//&types.Page{ID: 4},
			//&types.Page{ID: 5, SelfID: 4},

			// will be used for rebase delete
			&types.Page{ID: 10},
			&types.Page{ID: 11, SelfID: 10},
			&types.Page{ID: 12, SelfID: 10},
			&types.Page{ID: 121, SelfID: 12},
			&types.Page{ID: 122, SelfID: 12},

			// will be used for cascade delete
			&types.Page{ID: 20},
			&types.Page{ID: 21, SelfID: 20},
			&types.Page{ID: 22, SelfID: 20},
			&types.Page{ID: 221, SelfID: 22},
			&types.Page{ID: 222, SelfID: 22},
		}

		svc = &page{
			store: s,
			ac: &accessControl{rbac: rbac.NoopSvc(rbac.Allow, rbac.Config{
				RuleStorage: s,
				RoleStorage: s,
			})},
			eventbus: eventbus.New(),
			locale:   ResourceTranslationsManager(locale.Static()),
		}

		pageLookup = func(t *testing.T, pageID uint64) *types.Page {
			p, err := store.LookupComposePageByID(ctx, s, pageID)
			require.NoError(t, err)
			require.NotNil(t, p)
			return p
		}
	)

	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	//if err = store.Upgrade(ctx, zap.NewNop(), s); err != nil {
	if err = store.Upgrade(ctx, logger.MakeDebugLogger(), s); err != nil {
		t.Fatalf("failed to upgrade store: %v", err)
	}

	if err = s.TruncateComposeNamespaces(ctx); err != nil {
		t.Fatalf("failed to truncate compose namespaces: %v", err)
	}

	if err = s.TruncateComposeModules(ctx); err != nil {
		t.Fatalf("failed to truncate compose modules: %v", err)
	}

	//
	ns = &types.Namespace{Name: "testing", ID: namespaceID, CreatedAt: *now()}
	if err = store.CreateComposeNamespace(ctx, s, ns); err != nil {
		t.Fatalf("failed to seed namespaces: %v", err)
	}

	_ = pages.Walk(func(p *types.Page) error {
		p.NamespaceID = ns.ID
		return nil
	})

	if err = store.CreateComposePage(ctx, s, pages...); err != nil {
		t.Fatalf("failed to seed pages: %v", err)
	}

	t.Run("remove page without children", func(t *testing.T) {
		require.NoError(t, svc.DeleteByID(ctx, ns.ID, 1, types.PageChildrenOnDeleteAbort))
		require.NotNil(t, pageLookup(t, 1).DeletedAt, "parent should be deleted")
	})

	t.Run("abort when children are present", func(t *testing.T) {
		require.ErrorIs(t, svc.DeleteByID(ctx, ns.ID, 2, types.PageChildrenOnDeleteAbort), PageErrDeleteAbortedForPageWithSubpages())
		require.Nil(t, pageLookup(t, 2).DeletedAt, "parent should be deleted")
		require.Nil(t, pageLookup(t, 3).DeletedAt, "child should be deleted")
	})

	t.Run("remove only parent when forced", func(t *testing.T) {
		require.NoError(t, svc.DeleteByID(ctx, ns.ID, 2, types.PageChildrenOnDeleteForce))
		require.NotNil(t, pageLookup(t, 2).DeletedAt, "parent should be deleted")
		require.Nil(t, pageLookup(t, 3).DeletedAt, "child should not be deleted")
	})

	t.Run("delete parent and rebase children", func(t *testing.T) {
		require.NoError(t, svc.DeleteByID(ctx, ns.ID, 10, types.PageChildrenOnDeleteRebase))
		require.NotNil(t, pageLookup(t, 10).DeletedAt)

		for _, pageID := range []uint64{11, 12} {
			require.Nil(t, pageLookup(t, pageID).DeletedAt, "child page should not be deleted")
			require.Equal(t, uint64(0), pageLookup(t, pageID).SelfID, "child page should be moved one level lower")
		}

		for _, pageID := range []uint64{121, 122} {
			require.Nil(t, pageLookup(t, pageID).DeletedAt, "grandchild page should not be deleted")
			require.Equal(t, uint64(12), pageLookup(t, pageID).SelfID, "grand child page should stay where it is")
		}
	})

	t.Run("delete parent and all children", func(t *testing.T) {
		require.NoError(t, svc.DeleteByID(ctx, ns.ID, 20, types.PageChildrenOnDeleteCascade))
		require.NotNil(t, pageLookup(t, 20).DeletedAt, "parent page should be deleted")

		for _, pageID := range []uint64{21, 22} {
			require.NotNil(t, pageLookup(t, pageID).DeletedAt, "child page should not be deleted")
			require.Equal(t, uint64(20), pageLookup(t, pageID).SelfID, "child page should be stay where it is")
		}

		for _, pageID := range []uint64{221, 222} {
			require.NotNil(t, pageLookup(t, pageID).DeletedAt, "grandchild page should not be deleted")
			require.Equal(t, uint64(22), pageLookup(t, pageID).SelfID, "grandchild page should be stay where it is")
		}

	})
}
