package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/service/dml"
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

// permissiveProjectAccess answers yes to everything. These tests are about the
// revision mechanics, not about who may run them; the access checks themselves
// are exercised where they belong (and the point of them is that they are in
// the service at all — before 2026-07-30 CreateRevision and Publish asked
// nothing, so a nil controller here would have panicked on nothing).
type permissiveProjectAccess struct{}

func (permissiveProjectAccess) CanCreateProject(context.Context) bool                  { return true }
func (permissiveProjectAccess) CanSearchProjects(context.Context) bool                 { return true }
func (permissiveProjectAccess) CanReadProject(context.Context, *types.Project) bool    { return true }
func (permissiveProjectAccess) CanUpdateProject(context.Context, *types.Project) bool  { return true }
func (permissiveProjectAccess) CanDeleteProject(context.Context, *types.Project) bool  { return true }
func (permissiveProjectAccess) CanReviseProject(context.Context, *types.Project) bool  { return true }
func (permissiveProjectAccess) CanPublishProject(context.Context, *types.Project) bool { return true }
func (permissiveProjectAccess) CanManageMembersOnProject(context.Context, *types.Project) bool {
	return true
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
		ac:       permissiveProjectAccess{},
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

// fakeProjectPublishDeps stands in for the DAL and record services Publish
// wires up for the record migration. The publish tests below carry no records,
// so nothing here is ever called — migrateRecords returns on the empty mapping
// set — but it refuses to run at all while the deps are unset, so they have to
// be something.
type fakeProjectPublishDeps struct{}

func (fakeProjectPublishDeps) NewMigration() *dml.Migration { return nil }

func (fakeProjectPublishDeps) ReplaceConnection(context.Context, *dal.ConnectionWrap, bool) error {
	return nil
}
func (fakeProjectPublishDeps) RemoveConnection(context.Context, uint64) error { return nil }

func (fakeProjectPublishDeps) Bulk(context.Context, bool, ...*composeTypes.RecordBulkOperation) ([]composeTypes.RecordBulkOperationResult, error) {
	return nil, nil
}

func (fakeProjectPublishDeps) Search(context.Context, composeTypes.RecordFilter) (composeTypes.RecordSet, composeTypes.RecordFilter, error) {
	return nil, composeTypes.RecordFilter{}, nil
}

// TestPublishKeepsTheChainRootsIdentity locks the ruling that a chain's name and
// URL belong to its ROOT.
//
// Both were derived from the row being branched from or published over, so both
// grew on every publish: a third revision was handled
// "myproj-rev1-rev2-rev3" (unbounded growth against a 64-char column) and the
// live namespace slug walked "myproj" -> "myproj-rev1" -> "myproj-rev1-rev2",
// moving the URL of a deployed app on every release and leaving "myproj"
// resolving to the deprecated row. Three publishes is the smallest chain that
// shows it: the first one looked correct even before the fix.
func TestPublishKeepsTheChainRootsIdentity(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	svc.services.dalSvc = fakeProjectPublishDeps{}
	svc.services.dalConns = fakeProjectPublishDeps{}
	svc.services.recordSvc = fakeProjectPublishDeps{}
	ctx := context.Background()

	root := seedRevisionProject(t, s, "project-revision-test-identity")

	// Start the live namespace where onCreate leaves it: slug == handle. That is
	// the slug the app is reachable at and the one every publish must restore.
	rootNs, err := store.LookupComposeNamespaceByID(ctx, s, root.Config.NamespaceID)
	require.NoError(t, err)
	rootNs.Slug, rootNs.Enabled = root.Handle, true
	require.NoError(t, store.UpdateComposeNamespace(ctx, s, rootNs))

	head := root
	for n := 1; n <= 3; n++ {
		draft, cerr := svc.CreateRevision(ctx, head.ID)
		require.NoError(t, cerr)

		require.Equal(t, root.Handle+"-rev"+strconv.Itoa(n), draft.Handle,
			"a revision handle is the ROOT's handle plus one suffix, however deep the chain is")

		head, err = svc.Publish(ctx, draft.ID, types.PublishRequest{Confirm: true, DiscardRecords: true})
		require.NoError(t, err)

		liveNs, lerr := store.LookupComposeNamespaceByID(ctx, s, head.Config.NamespaceID)
		require.NoError(t, lerr)
		require.Equal(t, root.Handle, liveNs.Slug,
			"the published app answers at the root's slug, so its URL never moves")
		require.True(t, liveNs.Enabled, "the revision going live has to be in service")
		require.Equal(t, "namespace (revision "+strconv.Itoa(n)+")", liveNs.Name,
			"the label carries one revision suffix, not one per branch ever taken")
	}

	require.Equal(t, 3, head.Revision)

	// And the slug is genuinely free for the chain to keep taking: the outgoing
	// namespace is renamed out of the way before the incoming one claims it, so
	// the canonical slug names exactly one live namespace.
	live, err := store.LookupComposeNamespaceBySlug(ctx, s, root.Handle)
	require.NoError(t, err)
	require.Equal(t, head.Config.NamespaceID, live.ID,
		"the canonical slug must resolve to the current head, not to a deprecated revision")
}

// TestCreateRevision_ChainCanBeRevisedRepeatedly is the "publish once and the
// project is frozen forever" bug. The one-draft gate passed Status: draft to
// SearchProjects, but project.cue's filter did not list status in byValue, so
// no predicate was emitted and the gate matched every row in the chain —
// including the active revision itself. Every branch after the first was
// refused with "a draft revision already exists" against a draft that did not
// exist, and since non-draft projects are read-only, nothing about the project
// could ever be changed again.
func TestCreateRevision_ChainCanBeRevisedRepeatedly(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-test-chain")

	rev1, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	// A second branch while rev1 is still a draft is the case the gate is FOR.
	_, err = svc.CreateRevision(ctx, parent.ID)
	require.EqualError(t, err, "a draft revision already exists for this project")

	// Publishing rev1 leaves the chain with a deprecated parent and an active
	// head, and no draft — so it must be branchable again.
	rev1.Status = types.ProjectStatusActive
	require.NoError(t, store.UpdateProject(ctx, s, rev1))
	parent.Status = types.ProjectStatusDeprecated
	require.NoError(t, store.UpdateProject(ctx, s, parent))

	rev2, err := svc.CreateRevision(ctx, rev1.ID)
	require.NoError(t, err, "a chain that has published must still be revisable")
	require.Equal(t, 2, rev2.Revision)
	require.Equal(t, parent.ID, rev2.RootProjectID())

	// Shelving a draft must NOT free the chain: archived is a listing state,
	// not a lifecycle one, and the filter's default of excluding archived rows
	// would otherwise let a second draft through.
	_, err = svc.Archive(ctx, rev2.ID)
	require.NoError(t, err)

	_, err = svc.CreateRevision(ctx, rev1.ID)
	require.EqualError(t, err, "a draft revision already exists for this project",
		"an archived draft is still a draft")

	require.Contains(t,
		func() []uint64 {
			set, lerr := svc.ListRevisions(ctx, rev1.ID)
			require.NoError(t, lerr)
			ids := make([]uint64, len(set))
			for i, p := range set {
				ids[i] = p.ID
			}
			return ids
		}(),
		rev2.ID, "an archived revision must still appear in its chain")

	_, err = svc.Unarchive(ctx, rev2.ID)
	require.NoError(t, err)

	// And discarding a draft frees the chain rather than wedging it: the gate
	// excludes deleted rows. Deleted through the service, which also releases
	// the draft namespace's slug — a raw row delete would leave it squatting
	// and the next branch would fail on the namespace instead of the gate.
	require.NoError(t, svc.DeleteByID(ctx, rev2.ID))

	rev3, err := svc.CreateRevision(ctx, rev1.ID)
	require.NoError(t, err, "a discarded draft must not block the next branch")
	require.Equal(t, 2, rev3.Revision)
}

// TestProjectStatusIsNotClientWritable locks the other half of the split: a
// live project could be flipped back to draft with a plain update, which
// unlocked its schema for editing while its records were live. Status is now
// written only by the lifecycle endpoints, and archiving — the one status
// change a user makes directly — has its own field.
func TestProjectStatusIsNotClientWritable(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	p := seedRevisionProject(t, s, "project-revision-test-status")

	upd := *p
	upd.Status = types.ProjectStatusDraft
	upd.Meta.Short = "renamed"

	got, err := svc.Update(ctx, &upd)
	require.NoError(t, err)
	require.Equal(t, types.ProjectStatusActive, got.Status, "update must not move the lifecycle")
	require.Equal(t, "renamed", got.Meta.Short, "the rest of the update still applies")

	archived, err := svc.Archive(ctx, p.ID)
	require.NoError(t, err)
	require.NotNil(t, archived.ArchivedAt)
	require.Equal(t, types.ProjectStatusActive, archived.Status,
		"archiving is a shelf state; it must not overwrite the lifecycle status")

	restored, err := svc.Unarchive(ctx, p.ID)
	require.NoError(t, err)
	require.Nil(t, restored.ArchivedAt)
	require.Equal(t, types.ProjectStatusActive, restored.Status,
		"unarchive restores the real status because it was never lost")
}

// TestPublishRefusesDeletedRevisions covers the pair of holes that let a
// publish destroy a live project. loadProject returns deleted rows by design,
// and deleting a project soft-deletes its namespace with it — so publishing a
// discarded draft swapped a deleted namespace into place and soft-deleted the
// live one. The app 404'd and the chain showed only a deprecated parent whose
// namespace was already gone.
func TestPublishRefusesDeletedRevisions(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-test-deleted")

	draft, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	require.NoError(t, svc.DeleteByID(ctx, draft.ID))

	_, err = svc.Publish(ctx, draft.ID, types.PublishRequest{Confirm: true})
	require.EqualError(t, err, "cannot publish a deleted revision; restore it first")

	// And the mirror case: the parent is where the records come from, so
	// publishing over a deleted one would migrate nothing and call it success.
	require.NoError(t, svc.UndeleteByID(ctx, draft.ID))
	require.NoError(t, svc.DeleteByID(ctx, parent.ID))

	_, err = svc.Publish(ctx, draft.ID, types.PublishRequest{Confirm: true})
	require.EqualError(t, err, "cannot publish over a deleted revision; restore "+parent.Handle+" first")
}

// denyingProjectAccess refuses everything, for the tests that assert the
// lifecycle endpoints ask at all.
type denyingProjectAccess struct{ permissiveProjectAccess }

func (denyingProjectAccess) CanReadProject(context.Context, *types.Project) bool    { return false }
func (denyingProjectAccess) CanReviseProject(context.Context, *types.Project) bool  { return false }
func (denyingProjectAccess) CanPublishProject(context.Context, *types.Project) bool { return false }

// TestRevisionLifecycleChecksPermissions locks in the fix for the hole an
// unauthorized-user probe found on 2026-07-30: CreateRevision, Publish,
// DeploymentPlan and ListRevisions are hand-written rather than generated CRUD,
// so they inherited none of the generated wrapper's access checks — a user
// holding no roles at all could branch and publish any project in the system.
func TestRevisionLifecycleChecksPermissions(t *testing.T) {
	svc, s := newTestProjectRevisionService(t)
	ctx := context.Background()

	parent := seedRevisionProject(t, s, "project-revision-test-denied")

	// Branch while allowed, so there is a draft to aim the publish checks at.
	draft, err := svc.CreateRevision(ctx, parent.ID)
	require.NoError(t, err)

	svc.ac = denyingProjectAccess{}

	_, err = svc.CreateRevision(ctx, parent.ID)
	require.EqualError(t, err, "not allowed to create a revision of this project")

	_, err = svc.Publish(ctx, draft.ID, types.PublishRequest{Confirm: true})
	require.EqualError(t, err, "not allowed to publish this project")

	_, err = svc.DeploymentPlan(ctx, draft.ID)
	require.EqualError(t, err, "not allowed to read this project")

	_, err = svc.ListRevisions(ctx, parent.ID)
	require.EqualError(t, err, "not allowed to read this project")

	// The publish check has to come before the status check, or the error tells
	// a caller who may not read the project which statuses it is not in.
	_, err = svc.Publish(ctx, parent.ID, types.PublishRequest{Confirm: true})
	require.EqualError(t, err, "not allowed to publish this project")
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
