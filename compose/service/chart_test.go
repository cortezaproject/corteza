package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCharts(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = sqlite3.ConnectInMemory(ctx)

		namespaceID = nextID()
		ns          *types.Namespace
	)

	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	func(s interface{}) {
		err = s.(interface {
			Upgrade(context.Context, *zap.Logger) error
		}).Upgrade(ctx, zap.NewNop())
		if err != nil {
			t.Fatalf("failed to upgrade store: %v", err)
		}
	}(s)

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
		svc := chart{
			store: s,
			ctx:   context.Background(),
			ac:    AccessControl(&permissions.ServiceAllowAll{}),
		}
		res, err := svc.Create(&types.Chart{Name: "My first chart", NamespaceID: namespaceID})
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		res, err = svc.FindByID(namespaceID, res.ID)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		res, err = svc.FindByHandle(namespaceID, res.Handle)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		res.Name = "Changed"
		res, err = svc.Update(res)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)
		req.NotNil(res.UpdatedAt)
		req.Equal(res.Name, "Changed")

		res, err = svc.FindByID(namespaceID, res.ID)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)
		req.Equal(res.Name, "Changed")

		err = svc.DeleteByID(namespaceID, res.ID)
		req.NoError(unwrapChartInternal(err))
		req.NotNil(res)

		// this works because we're allowed to do everything
		res, err = svc.FindByID(namespaceID, res.ID)
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
