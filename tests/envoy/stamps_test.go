package envoy

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	mtypes "github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	stypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestStamps(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = initStore(ctx)

		createdAtTs   = "2021-12-01T10:00:00Z"
		updatedAtTs   = "2021-12-02T10:00:00Z"
		deletedAtTs   = "2021-12-03T10:00:00Z"
		suspendedAtTs = "2021-12-04T10:00:00Z"
		archivedAtTs  = "2021-12-05T10:00:00Z"
		ft            = func(t time.Time) string {
			return t.Format(time.RFC3339)
		}
	)
	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	prepare := func(ctx context.Context, s store.Storer, t *testing.T, suite string) (*require.Assertions, error) {
		req := require.New(t)

		nn, err := dd(ctx, suite)
		req.NoError(err)

		return req, encode(ctx, s, nn)
	}

	// Prepare
	s, err = initStore(ctx)
	err = ce(
		err,

		s.TruncateRbacRules(ctx),
		s.TruncateRoles(ctx),
		s.TruncateUsers(ctx),
		s.TruncateActionlogs(ctx),
		s.TruncateApplications(ctx),
		s.TruncateAttachments(ctx),
		s.TruncateComposeAttachments(ctx),
		s.TruncateComposeCharts(ctx),
		s.TruncateComposeNamespaces(ctx),
		s.TruncateComposeModules(ctx),
		s.TruncateComposeModuleFields(ctx),
		s.TruncateComposePages(ctx),
		s.TruncateComposeRecords(ctx, nil),
		s.TruncateMessagingChannels(ctx),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := prepare(ctx, s, t, "sys_stamps")
	req.NoError(err)

	lcn := func(ctx context.Context, s store.Storer, slg string) (*types.Namespace, error) {
		rr, _, err := store.SearchComposeNamespaces(ctx, s, types.NamespaceFilter{
			Slug:    slg,
			Deleted: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	lcm := func(ctx context.Context, s store.Storer, hnd string) (*types.Module, error) {
		rr, _, err := store.SearchComposeModules(ctx, s, types.ModuleFilter{
			Handle:  hnd,
			Deleted: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	lcp := func(ctx context.Context, s store.Storer, hnd string) (*types.Page, error) {
		rr, _, err := store.SearchComposePages(ctx, s, types.PageFilter{
			Handle:  hnd,
			Deleted: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	lcc := func(ctx context.Context, s store.Storer, hnd string) (*types.Chart, error) {
		rr, _, err := store.SearchComposeCharts(ctx, s, types.ChartFilter{
			Handle:  hnd,
			Deleted: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	lcr := func(ctx context.Context, s store.Storer, m *types.Module) (*types.Record, error) {
		rr, _, err := store.SearchComposeRecords(ctx, s, m, types.RecordFilter{
			Deleted: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	lu := func(ctx context.Context, s store.Storer, hnd string) (*stypes.User, error) {
		rr, _, err := store.SearchUsers(ctx, s, stypes.UserFilter{
			Handle:    hnd,
			Deleted:   filter.StateInclusive,
			Suspended: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	la := func(ctx context.Context, s store.Storer, name string) (*stypes.Application, error) {
		rr, _, err := store.SearchApplications(ctx, s, stypes.ApplicationFilter{
			Name:    name,
			Deleted: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	lr := func(ctx context.Context, s store.Storer, hnd string) (*stypes.Role, error) {
		rr, _, err := store.SearchRoles(ctx, s, stypes.RoleFilter{
			Handle:   hnd,
			Deleted:  filter.StateInclusive,
			Archived: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	lmc := func(ctx context.Context, s store.Storer, name string) (*mtypes.Channel, error) {
		rr, _, err := store.SearchMessagingChannels(ctx, s, mtypes.ChannelFilter{
			IncludeDeleted: true,
		})
		if err != nil {
			return nil, err
		}
		if len(rr) > 0 {
			return rr[0], nil
		}
		return nil, nil
	}

	t.Run("compose namespace", func(t *testing.T) {
		ns, err := lcn(ctx, s, "ns1")
		req.NoError(err)
		req.NotNil(ns)

		req.Equal(createdAtTs, ft(ns.CreatedAt))
		req.Equal(updatedAtTs, ft(*ns.UpdatedAt))
		req.Equal(deletedAtTs, ft(*ns.DeletedAt))
	})

	t.Run("compose modules", func(t *testing.T) {
		mod, err := lcm(ctx, s, "mod1")
		req.NoError(err)
		req.NotNil(mod)

		req.Equal(createdAtTs, ft(mod.CreatedAt))
		req.Equal(updatedAtTs, ft(*mod.UpdatedAt))
		req.Equal(deletedAtTs, ft(*mod.DeletedAt))
	})

	t.Run("compose pages", func(t *testing.T) {
		pg, err := lcp(ctx, s, "pg1")
		req.NoError(err)
		req.NotNil(pg)

		req.Equal(createdAtTs, ft(pg.CreatedAt))
		req.Equal(updatedAtTs, ft(*pg.UpdatedAt))
		req.Equal(deletedAtTs, ft(*pg.DeletedAt))
	})

	t.Run("compose charts", func(t *testing.T) {
		chr, err := lcc(ctx, s, "chr1")
		req.NoError(err)
		req.NotNil(chr)

		req.Equal(createdAtTs, ft(chr.CreatedAt))
		req.Equal(updatedAtTs, ft(*chr.UpdatedAt))
		req.Equal(deletedAtTs, ft(*chr.DeletedAt))
	})

	t.Run("compose records", func(t *testing.T) {
		mod, err := lcm(ctx, s, "mod1")
		req.NoError(err)
		req.NotNil(mod)
		r, err := lcr(ctx, s, mod)
		req.NoError(err)
		req.NotNil(r)
		u, err := lu(ctx, s, "test")
		req.NoError(err)
		req.NotNil(u)

		req.Equal(createdAtTs, ft(r.CreatedAt))
		req.Equal(updatedAtTs, ft(*r.UpdatedAt))
		req.Equal(deletedAtTs, ft(*r.DeletedAt))

		req.Equal(u.ID, r.CreatedBy)
		req.Equal(u.ID, r.UpdatedBy)
		req.Equal(u.ID, r.DeletedBy)
	})

	t.Run("users", func(t *testing.T) {
		u, err := lu(ctx, s, "test")
		req.NoError(err)
		req.NotNil(u)

		req.Equal(createdAtTs, ft(u.CreatedAt))
		req.Equal(updatedAtTs, ft(*u.UpdatedAt))
		req.Equal(deletedAtTs, ft(*u.DeletedAt))
		req.Equal(suspendedAtTs, ft(*u.SuspendedAt))
	})

	t.Run("applications", func(t *testing.T) {
		a, err := la(ctx, s, "app1")
		req.NoError(err)
		req.NotNil(a)

		req.Equal(createdAtTs, ft(a.CreatedAt))
		req.Equal(updatedAtTs, ft(*a.UpdatedAt))
		req.Equal(deletedAtTs, ft(*a.DeletedAt))
	})

	t.Run("roles", func(t *testing.T) {
		r, err := lr(ctx, s, "r1")
		req.NoError(err)
		req.NotNil(r)

		req.Equal(createdAtTs, ft(r.CreatedAt))
		req.Equal(updatedAtTs, ft(*r.UpdatedAt))
		req.Equal(deletedAtTs, ft(*r.DeletedAt))
		req.Equal(archivedAtTs, ft(*r.ArchivedAt))
	})

	t.Run("messaging channels", func(t *testing.T) {
		r, err := lmc(ctx, s, "ch1")
		req.NoError(err)
		req.NotNil(r)
		u, err := lu(ctx, s, "test")
		req.NoError(err)
		req.NotNil(u)

		req.Equal(createdAtTs, ft(r.CreatedAt))
		req.Equal(updatedAtTs, ft(*r.UpdatedAt))
		req.Equal(deletedAtTs, ft(*r.DeletedAt))
		// Can't get archived channels
		// req.Equal(archivedAtTs, ft(*r.ArchivedAt))

		req.Equal(u.ID, r.CreatorID)
	})
}
