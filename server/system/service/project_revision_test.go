package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

// fakeProjectNamespaceCloner is a cut-down projectNamespaceSvc stub. The real
// implementation (compose/service namespace.CloneFromStore) drives envoy
// decode/encode and DAL services to clone every module/page/chart in the
// namespace — far more than CreateRevision's own logic needs to be exercised.
//
// It reproduces exactly the two properties the code under test depends on:
// the clone is PERSISTED (cloneProjectResources loads both namespaces to build
// its ID map), and modules come across with fresh IDs but their handles
// UNCHANGED — which is the only thing the two sides share, and so the only
// thing a reference remap can match on. Everything else about the source
// namespace is left out.
type fakeProjectNamespaceCloner struct{ s store.Storer }

func (c fakeProjectNamespaceCloner) CloneFromStore(ctx context.Context, sourceNsID uint64, dup *composeTypes.Namespace) (*composeTypes.Namespace, error) {
	dup.ID = nextID()
	if err := store.CreateComposeNamespace(ctx, c.s, dup); err != nil {
		return nil, err
	}

	mm, _, err := store.SearchComposeModules(ctx, c.s, composeTypes.ModuleFilter{NamespaceID: sourceNsID})
	if err != nil {
		return nil, err
	}

	for _, m := range mm {
		cp := *m
		cp.ID = nextID()
		cp.NamespaceID = dup.ID
		if err = store.CreateComposeModule(ctx, c.s, &cp); err != nil {
			return nil, err
		}
	}

	return dup, nil
}

// newTestProjectRevisionService wires a project service against a real
// in-memory sqlite store (schema applied via store.Upgrade) plus the fake
// namespace cloner above. This is the store.Storer seam: CreateRevision calls
// straight through to package-level store.* functions and svc.services.nsSvc,
// so there's no interface boundary to fake the store itself out at — but a
// sqlite-backed store runs in-process, in CI, with no live database, and
// still exercises the real uniqueness-constraint check (checkProjectConstraints
// in store/adapters/rdbms/rdbms.gen.go) that the bug this test locks in was
// tripping over.
func newTestProjectRevisionService(t *testing.T) (*project, store.Storer) {
	t.Helper()
	ctx := context.Background()

	s, err := sqlite.ConnectInMemory(ctx)
	require.NoError(t, err)
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), s))

	svc := &project{
		store:    s,
		services: &projectServices{nsSvc: fakeProjectNamespaceCloner{s: s}},
	}
	return svc, s
}

// seedRevisionProject creates an active project plus the compose namespace
// CreateRevision loads via parent.Config.NamespaceID — the two preconditions
// CreateRevision needs before it can run at all.
func seedRevisionProject(t *testing.T, s store.Storer, handle string) *types.Project {
	t.Helper()
	ctx := context.Background()

	nsID := nextID()
	require.NoError(t, store.CreateComposeNamespace(ctx, s, &composeTypes.Namespace{
		ID:   nsID,
		Slug: "ns-" + strconv.FormatUint(nsID, 10),
		Name: "namespace",
	}))

	p := &types.Project{
		ID:     nextID(),
		Handle: handle,
		Status: types.ProjectStatusActive,
		Config: types.ProjectConfig{NamespaceID: nsID},
	}
	require.NoError(t, store.CreateProject(ctx, s, p))
	return p
}

// TestCreateRevision_EmptyHandleProjectsDoNotCollide is the exact scenario
// from the bug report: projects are routinely created with only a name and
// description, so their Handle is empty. CreateRevision used to build the
// first revision's handle as parent.Handle + "-rev1", which collapsed to the
// literal string "-rev1" for every empty-handle project — so only the first
// project ever revised in a database could succeed; every subsequent one hit
// a "not unique" error from checkProjectConstraints. Two distinct
// empty-handle projects must now each get a distinct, root-ID-derived handle
// and both revisions must succeed.
func TestCreateRevision_EmptyHandleProjectsDoNotCollide(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	p1 := seedRevisionProject(t, s, "")
	p2 := seedRevisionProject(t, s, "")

	rev1, err := svc.CreateRevision(ctx, p1.ID)
	require.NoError(t, err, "first project's first revision")

	rev2, err := svc.CreateRevision(ctx, p2.ID)
	require.NoError(t, err, "second project's first revision must not collide with the first")

	require.NotEmpty(t, rev1.Handle)
	require.NotEmpty(t, rev2.Handle)
	require.NotEqual(t, rev1.Handle, rev2.Handle, "handles must be derived per-project, not both collapse to \"-rev1\"")

	// Pin the exact fallback shape (rev-<rootID>-rev<n>) so a regression back
	// to the bare "-rev1" concatenation fails loudly here, not just via the
	// generic non-collision check above.
	require.Equal(t, "rev-"+strconv.FormatUint(p1.ID, 10)+"-rev1", rev1.Handle)
	require.Equal(t, "rev-"+strconv.FormatUint(p2.ID, 10)+"-rev1", rev2.Handle)
}

// TestCreateRevision_HandledProjectKeepsSuffixedHandle covers the non-empty
// path the fix must leave unchanged: a project with a real handle still
// derives "<handle>-rev<n>", not the root-ID fallback.
func TestCreateRevision_HandledProjectKeepsSuffixedHandle(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	// ConnectInMemory uses a fixed "cache=shared" DSN (see
	// store/adapters/rdbms/drivers/sqlite/connect.go), so every in-memory
	// store opened within this test binary shares one physical sqlite
	// database — including the one project_resolver_test.go's fixtures use.
	// Handle must be globally unique across the package's test fixtures, not
	// just within this test.
	p := seedRevisionProject(t, s, "project-revision-test-handled")

	rev, err := svc.CreateRevision(ctx, p.ID)
	require.NoError(t, err)
	require.Equal(t, "project-revision-test-handled-rev1", rev.Handle)
}

// TestDiffSourcesReportsWhatHappened guards the ambiguity this diff was
// reshaped to remove: with only {path, risk}, an added and a removed resource
// were indistinguishable, so the publish screen could not say which it had.
func TestDiffSourcesReportsWhatHappened(t *testing.T) {
	before := []*GraphSource{
		{Kind: "page", ID: 1, Name: "Claims", Handle: "claims"},
		{Kind: "automation", ID: 2, Name: "Notify", Handle: "notify"},
	}
	after := []*GraphSource{
		{Kind: "page", ID: 9, Name: "Claims", Handle: "claims"},
		{Kind: "page", ID: 10, Name: "Invoices", Handle: "invoices"},
	}

	require.Equal(t, []types.ProjectChange{
		{Op: types.ProjectChangeOpRemoved, Kind: "automation", Name: "Notify", Path: "automation.notify", Risk: types.ProjectChangeRiskSafe},
		{Op: types.ProjectChangeOpAdded, Kind: "page", Name: "Invoices", Path: "page.invoices", Risk: types.ProjectChangeRiskSafe},
	}, diffSources(before, after))
}

// A revision clone mints fresh IDs, so matching on ID would report every
// resource in the project as removed-and-re-added on every single publish.
func TestDiffSourcesMatchesOnHandleNotID(t *testing.T) {
	before := []*GraphSource{{Kind: "page", ID: 1, Name: "Claims", Handle: "claims"}}
	after := []*GraphSource{{Kind: "page", ID: 9999, Name: "Claims", Handle: "claims"}}

	require.Empty(t, diffSources(before, after))
}

// Modules are diffed field-by-field elsewhere; users are project membership
// rather than a deployment artefact; knowledge bases are loaded unscoped by the
// graph, so they are identical on both sides by construction.
func TestDiffSourcesSkipsKindsItDoesNotOwn(t *testing.T) {
	ss := []*GraphSource{
		{Kind: "module", ID: 1, Name: "Customer", Handle: "customer"},
		{Kind: "user", ID: 2, Name: "Ana"},
		{Kind: "knowledge-base", ID: 3, Name: "Handbook", Handle: "handbook"},
	}

	require.Empty(t, diffSources(ss, nil))
	require.Empty(t, diffSources(nil, ss))
}

func TestDiffSourcesFallsBackToNameWithoutHandle(t *testing.T) {
	before := []*GraphSource{{Kind: "connection", ID: 1, Name: "Billing DB"}}

	require.Equal(t, []types.ProjectChange{{
		Op:   types.ProjectChangeOpRemoved,
		Kind: "connection",
		Name: "Billing DB",
		Path: "connection.Billing DB",
		Risk: types.ProjectChangeRiskSafe,
	}}, diffSources(before, nil))
}

// Map iteration order is random; this list is read by a person and must not
// reshuffle between two loads of the same screen.
func TestDiffSourcesIsDeterministic(t *testing.T) {
	after := []*GraphSource{
		{Kind: "page", ID: 1, Handle: "zeta", Name: "Zeta"},
		{Kind: "automation", ID: 2, Handle: "alpha", Name: "Alpha"},
		{Kind: "page", ID: 3, Handle: "mid", Name: "Mid"},
	}

	first := diffSources(nil, after)
	for i := 0; i < 25; i++ {
		require.Equal(t, first, diffSources(nil, after))
	}
	require.Equal(t,
		[]string{"automation.alpha", "page.mid", "page.zeta"},
		[]string{first[0].Path, first[1].Path, first[2].Path},
	)
}

func TestAddChangeRaisesPlanRisk(t *testing.T) {
	plan := &types.ProjectDeploymentPlan{Risk: types.ProjectChangeRiskSafe}

	addChange(plan, types.ProjectChange{Op: types.ProjectChangeOpAdded, Risk: types.ProjectChangeRiskSafe})
	require.Equal(t, types.ProjectChangeRiskSafe, plan.Risk)

	addChange(plan, types.ProjectChange{Op: types.ProjectChangeOpRemoved, Risk: types.ProjectChangeRiskDangerous})
	require.Equal(t, types.ProjectChangeRiskDangerous, plan.Risk)

	// One dangerous change is enough — a later safe one must not clear it.
	addChange(plan, types.ProjectChange{Op: types.ProjectChangeOpAdded, Risk: types.ProjectChangeRiskSafe})
	require.Equal(t, types.ProjectChangeRiskDangerous, plan.Risk)
	require.Len(t, plan.Changes, 3)
}
