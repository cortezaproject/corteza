package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

// permissiveProjectIncidentAccess answers yes to everything: these tests are
// about which project id a read lands on, not about who may run it.
type permissiveProjectIncidentAccess struct{}

func (permissiveProjectIncidentAccess) CanCreateProjectIncident(context.Context) bool  { return true }
func (permissiveProjectIncidentAccess) CanSearchProjectIncidents(context.Context) bool { return true }
func (permissiveProjectIncidentAccess) CanReadProjectIncident(context.Context, *types.ProjectIncident) bool {
	return true
}
func (permissiveProjectIncidentAccess) CanUpdateProjectIncident(context.Context, *types.ProjectIncident) bool {
	return true
}
func (permissiveProjectIncidentAccess) CanDeleteProjectIncident(context.Context, *types.ProjectIncident) bool {
	return true
}

// newTestWorkItemStore is the same sqlite-backed store seam
// newTestProjectRevisionService uses: the work-item services call straight
// through to package-level store.* functions, so there is no interface to fake
// the store out at — but an in-memory sqlite store runs in-process, in CI, and
// exercises the real rel_project predicate that this bug lives in.
func newTestWorkItemStore(t *testing.T) store.Storer {
	t.Helper()
	ctx := context.Background()

	s, err := sqlite.ConnectInMemory(ctx)
	require.NoError(t, err)
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), s))
	return s
}

// seedRevisionChain creates a root project and one revision branched off it,
// the shape every dashboard read runs against once a project has been revised:
// the chain HEAD (the revision) is what the UI routes on, the chain ROOT is
// what work items are filed under.
func seedRevisionChain(t *testing.T, s store.Storer) (root, rev *types.Project) {
	t.Helper()
	ctx := context.Background()

	// sqlite.ConnectInMemory hands every caller the same in-memory database, so
	// handles have to be unique across the whole package's tests, not just
	// within one of them.
	root = &types.Project{ID: nextID(), Status: types.ProjectStatusDeprecated}
	root.Handle = "work-item-scope-root-" + strconv.FormatUint(root.ID, 10)
	require.NoError(t, store.CreateProject(ctx, s, root))

	rev = &types.Project{ID: nextID(), Status: types.ProjectStatusActive, ProjectID: root.ID, Revision: 1}
	rev.Handle = "work-item-scope-rev-" + strconv.FormatUint(rev.ID, 10)
	require.NoError(t, store.CreateProject(ctx, s, rev))

	return root, rev
}

// TestProjectIncidentSearchNormalisesToChainRoot is the "every dashboard panel
// goes empty the moment a project is revised" bug. Work items are WRITTEN
// against the chain root (beforeCreate), but the read path passed the caller's
// project id straight into the store filter, where rel_project is a real
// predicate. The dashboard routes on the chain HEAD, which stops being the root
// after the first revision — so reports, boards and panels all returned nothing
// while the nav badges (which already normalised) showed real counts beside
// them.
func TestProjectIncidentSearchNormalisesToChainRoot(t *testing.T) {
	ctx := context.Background()
	s := newTestWorkItemStore(t)
	root, rev := seedRevisionChain(t, s)

	svc := &projectIncident{store: s, ac: permissiveProjectIncidentAccess{}}

	// Filed against the root, as the frontend does before any revision exists.
	_, err := svc.Create(ctx, &types.ProjectIncident{ProjectID: root.ID, Title: "filed on the root"})
	require.NoError(t, err)

	// Filed against the revision: the write path normalises this to the root
	// too, so both items belong to the one chain-wide pool.
	filedOnRev, err := svc.Create(ctx, &types.ProjectIncident{ProjectID: rev.ID, Title: "filed on the revision"})
	require.NoError(t, err)
	require.Equal(t, root.ID, filedOnRev.ProjectID, "the write path must file against the chain root")

	// An unrelated project's item must never leak into either read below.
	other := &types.Project{ID: nextID(), Status: types.ProjectStatusActive}
	other.Handle = "work-item-scope-other-" + strconv.FormatUint(other.ID, 10)
	require.NoError(t, store.CreateProject(ctx, s, other))
	_, err = svc.Create(ctx, &types.ProjectIncident{ProjectID: other.ID, Title: "someone else's"})
	require.NoError(t, err)

	t.Run("reading the chain head sees the chain's items", func(t *testing.T) {
		set, _, err := svc.Search(ctx, types.ProjectIncidentFilter{ProjectID: rev.ID})
		require.NoError(t, err)
		require.Len(t, set, 2, "a revision reads the whole chain-wide item pool, not an empty board")
	})

	t.Run("reading the chain root is unchanged", func(t *testing.T) {
		set, _, err := svc.Search(ctx, types.ProjectIncidentFilter{ProjectID: root.ID})
		require.NoError(t, err)
		require.Len(t, set, 2)
	})

	t.Run("no project id stays unfiltered", func(t *testing.T) {
		set, _, err := svc.Search(ctx, types.ProjectIncidentFilter{})
		require.NoError(t, err)
		require.Len(t, set, 3, "a filter with no project scope must not acquire one")
	})

	t.Run("an unknown project id returns nothing rather than erroring", func(t *testing.T) {
		set, _, err := svc.Search(ctx, types.ProjectIncidentFilter{ProjectID: nextID()})
		require.NoError(t, err)
		require.Empty(t, set)
	})
}

// TestProjectReportNormalisesToChainRoot pins the aggregate read too: the
// report service owns no project filtering of its own, it pages the work-item
// services' Search — so it is fixed by, and stays fixed with, the same hook.
func TestProjectReportNormalisesToChainRoot(t *testing.T) {
	ctx := context.Background()
	s := newTestWorkItemStore(t)
	root, rev := seedRevisionChain(t, s)

	inc := &projectIncident{store: s, ac: permissiveProjectIncidentAccess{}}
	_, err := inc.Create(ctx, &types.ProjectIncident{ProjectID: root.ID, Title: "x", Status: "Open"})
	require.NoError(t, err)

	res, err := (&projectReport{incident: inc}).Report(ctx, &types.ProjectReportRequest{
		ProjectID:  rev.ID,
		Resource:   "incident",
		Dimensions: []string{"status"},
	})
	require.NoError(t, err)
	require.Len(t, res.Set, 1)
	require.Equal(t, float64(1), res.Set[0].Metrics["count"])
}

// TestWorkItemServicesNormaliseSearchScope covers the other five kinds: every
// work-item service has to normalise, not just the one with an end-to-end test
// above, because they all share the chain-root write rule.
func TestWorkItemServicesNormaliseSearchScope(t *testing.T) {
	ctx := context.Background()
	s := newTestWorkItemStore(t)
	root, rev := seedRevisionChain(t, s)
	unknown := nextID()

	// Each entry runs its service's beforeSearch over a filter and reports the
	// project id the store would then be asked for.
	normalise := map[string]func(uint64) uint64{
		"incident": func(id uint64) uint64 {
			f := types.ProjectIncidentFilter{ProjectID: id}
			require.NoError(t, (&projectIncident{store: s}).beforeSearch(ctx, &f))
			return f.ProjectID
		},
		"task": func(id uint64) uint64 {
			f := types.ProjectTaskFilter{ProjectID: id}
			require.NoError(t, (&projectTask{store: s}).beforeSearch(ctx, &f))
			return f.ProjectID
		},
		"feature": func(id uint64) uint64 {
			f := types.ProjectFeatureFilter{ProjectID: id}
			require.NoError(t, (&projectFeature{store: s}).beforeSearch(ctx, &f))
			return f.ProjectID
		},
		"privacy": func(id uint64) uint64 {
			f := types.ProjectPrivacyFilter{ProjectID: id}
			require.NoError(t, (&projectPrivacy{store: s}).beforeSearch(ctx, &f))
			return f.ProjectID
		},
		"review": func(id uint64) uint64 {
			f := types.ProjectReviewFilter{ProjectID: id}
			require.NoError(t, (&projectReview{store: s}).beforeSearch(ctx, &f))
			return f.ProjectID
		},
		"backlogItem": func(id uint64) uint64 {
			f := types.ProjectBacklogItemFilter{ProjectID: id}
			require.NoError(t, (&projectBacklogItem{store: s}).beforeSearch(ctx, &f))
			return f.ProjectID
		},
	}

	for kind, fn := range normalise {
		t.Run(kind, func(t *testing.T) {
			require.Equal(t, root.ID, fn(rev.ID), "a revision must read the chain root")
			require.Equal(t, root.ID, fn(root.ID), "the root normalises to itself")
			require.Equal(t, uint64(0), fn(0), "no scope must stay no scope")
			require.Equal(t, unknown, fn(unknown), "an unresolvable id passes through, it does not error")
		})
	}
}
