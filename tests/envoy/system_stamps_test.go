package envoy

import (
	"context"
	"testing"
	"time"

	ctypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	stypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestStamps(t *testing.T) {
	var (
		ctx = context.Background()
		s   = initServices(ctx, t)
		req = require.New(t)

		createdAtTs   = "2021-12-01T10:00:00Z"
		updatedAtTs   = "2021-12-02T10:00:00Z"
		deletedAtTs   = "2021-12-03T10:00:00Z"
		suspendedAtTs = "2021-12-04T10:00:00Z"
		archivedAtTs  = "2021-12-05T10:00:00Z"
		ft            = func(t time.Time) string {
			return t.Format(time.RFC3339)
		}

		lcn = func(ctx context.Context, s store.Storer, slg string) (*ctypes.Namespace, error) {
			rr, _, err := store.SearchComposeNamespaces(ctx, s, ctypes.NamespaceFilter{
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

		lcm = func(ctx context.Context, s store.Storer, hnd string) (*ctypes.Module, error) {
			rr, _, err := store.SearchComposeModules(ctx, s, ctypes.ModuleFilter{
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

		lcp = func(ctx context.Context, s store.Storer, hnd string) (*ctypes.Page, error) {
			rr, _, err := store.SearchComposePages(ctx, s, ctypes.PageFilter{
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

		lcc = func(ctx context.Context, s store.Storer, hnd string) (*ctypes.Chart, error) {
			rr, _, err := store.SearchComposeCharts(ctx, s, ctypes.ChartFilter{
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

		lcr = func(ctx context.Context, s store.Storer, m *ctypes.Module) (*ctypes.Record, error) {
			rr, _, err := store.SearchComposeRecords(ctx, s, m, ctypes.RecordFilter{
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

		lu = func(ctx context.Context, s store.Storer, hnd string) (*stypes.User, error) {
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

		la = func(ctx context.Context, s store.Storer, name string) (*stypes.Application, error) {
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

		lr = func(ctx context.Context, s store.Storer, hnd string) (*stypes.Role, error) {
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
	)

	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	nn, err := decodeDirectory(ctx, "sys_stamps")
	req.NoError(err)
	req.NoError(encode(ctx, s, nn))

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
		u, err := lu(ctx, s, "test")
		req.NoError(err)
		req.NotNil(u)

		req.Equal(createdAtTs, ft(a.CreatedAt))
		req.Equal(updatedAtTs, ft(*a.UpdatedAt))
		req.Equal(deletedAtTs, ft(*a.DeletedAt))

		req.Equal(u.ID, a.OwnerID)
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
}
