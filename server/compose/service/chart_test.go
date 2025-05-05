package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCharts(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = sqlite.ConnectInMemory(ctx)

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

	if err = s.TruncateComposeCharts(ctx); err != nil {
		t.Fatalf("failed to truncate compose charts: %v", err)
	}

	ns = &types.Namespace{Name: "testing", ID: namespaceID, CreatedAt: *now()}
	if err = store.CreateComposeNamespace(ctx, s, ns); err != nil {
		t.Fatalf("failed to seed namespaces: %v", err)
	}

	t.Run("crud", func(t *testing.T) {
		req := require.New(t)
		svc := &chart{
			store: s,
			ac: &accessControl{rbac: rbac.NoopSvc(rbac.Allow, rbac.Config{
				RuleStorage: s,
				RoleStorage: s,
			})},
		}
		res, err := svc.Create(ctx, &types.Chart{Name: "My first chart", NamespaceID: namespaceID})
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		res, err = svc.FindByID(ctx, namespaceID, res.ID)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		res, err = svc.FindByHandle(ctx, namespaceID, res.Handle)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		res.Name = "Changed"
		res, err = svc.Update(ctx, res)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)
		req.NotNil(res.UpdatedAt)
		req.Equal(res.Name, "Changed")

		res, err = svc.FindByID(ctx, namespaceID, res.ID)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)
		req.Equal(res.Name, "Changed")

		err = svc.DeleteByID(ctx, namespaceID, res.ID)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		// this works because we're allowed to do everything
		res, err = svc.FindByID(ctx, namespaceID, res.ID)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)
		req.NotNil(res.DeletedAt)

	})
}

func unwrapChartInternal(err error) error {
	g := ChartErrGeneric()
	for {
		if errors.Is(err, g) {
			err = errors.Unwrap(err)
			continue
		}

		return err
	}
}
