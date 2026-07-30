package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	composeTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/dml"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

// SetProjectRevisionDeps wires deps needed for the revision/publish flow.
// Called from app/boot_levels.go after compose services are initialised.
func SetProjectRevisionDeps(
	nsSvc projectNamespaceSvc,
	importer projectDALImporter,
	dalConns projectDALConnSvc,
	recordSvc projectRecordSvc,
) {
	if DefaultProject == nil {
		return
	}
	DefaultProject.services.nsSvc = nsSvc
	DefaultProject.services.dalSvc = importer
	DefaultProject.services.dalConns = dalConns
	DefaultProject.services.recordSvc = recordSvc
}

// CreateRevision clones the project and its namespace into a new draft revision.
//
// ACCESS CONTROL IS EXPLICIT HERE. The whole revision/publish surface is
// hand-written rather than generated CRUD, so it gets none of the generated
// wrapper's checkScope + CanXProject calls (compare project.gen.go) -- until
// 2026-07-30 that meant any authenticated caller could branch or publish any
// project, verified against a user holding no roles at all. Every entry point
// in this file now names the permission it needs.
func (svc *project) CreateRevision(ctx context.Context, projectID uint64) (rev *types.Project, err error) {
	var parent *types.Project
	if parent, err = loadProject(ctx, svc.store, projectID); err != nil {
		return
	}

	if !svc.ac.CanReviseProject(ctx, parent) {
		return nil, ProjectErrNotAllowedToRevise()
	}

	if parent.Status != types.ProjectStatusActive {
		return nil, fmt.Errorf("can only revise active projects")
	}

	// One draft per chain at a time.
	//
	// This gate used to match EVERY row in the chain: it passed Status here,
	// but project.cue's filter did not list status in byValue, so the store
	// emitted no predicate for it and silently answered a different question.
	// After the first publish the chain always holds at least the active
	// revision, so every later branch was refused with "a draft revision
	// already exists" -- against a draft that did not exist. Combined with the
	// read-only guard on non-draft projects, that made a published project
	// permanently unchangeable. Deleted rows are excluded by the filter's
	// default, so discarding a draft frees the chain again.
	rootID := parent.RootProjectID()
	existing, _, err := store.SearchProjects(ctx, svc.store, types.ProjectFilter{
		RootProjectID: rootID,
		Status:        types.ProjectStatusDraft,
	})
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, errors.DuplicateData("a draft revision already exists for this project")
	}

	// Load the parent's namespace so we can clone it.
	oldNs, err := store.LookupComposeNamespaceByID(ctx, svc.store, parent.Config.NamespaceID)
	if err != nil {
		return nil, fmt.Errorf("load namespace for project: %w", err)
	}

	newRevision := parent.Revision + 1
	revSuffix := "-rev" + strconv.Itoa(newRevision)
	// Root-derived fallback for an empty prefix. Projects are routinely created
	// without a handle or slug (the create flow only collects name and
	// description), and both of these are unique — so without this, EVERY
	// project's first revision would want the same "-rev1" and only one project
	// in the whole system could ever be revised.
	rootFallback := "rev-" + strconv.FormatUint(rootID, 10) + revSuffix

	revSlug := oldNs.Slug + revSuffix
	// Envoy requires a non-empty slug for clone.
	if revSlug == revSuffix {
		revSlug = rootFallback
	}

	revHandle := parent.Handle + revSuffix
	if revHandle == revSuffix {
		revHandle = rootFallback
	}

	if svc.services.nsSvc == nil {
		return nil, fmt.Errorf("project revision deps not initialised (nsSvc)")
	}

	// Allocated up front (rather than inside the *types.Project literal below)
	// so the cloned namespace can be repointed at the new draft revision's own
	// ID right after cloning -- see the comment below.
	revID := nextID()

	dup := &composeTypes.Namespace{
		TenantID: oldNs.TenantID,
		Name:     oldNs.Name + " (revision " + strconv.Itoa(newRevision) + ")",
		Slug:     revSlug,
	}
	clonedNs, err := svc.services.nsSvc.CloneFromStore(ctx, parent.Config.NamespaceID, dup)
	if err != nil {
		return nil, fmt.Errorf("clone namespace for revision: %w", err)
	}

	// CloneFromStore's namespace.envoyRun (compose/service/namespace.go) only
	// ever applies dup's Name and Slug to the clone: it decodes the source
	// namespace resource wholesale and overwrites just those two fields, so
	// ProjectID, TenantID and Enabled always come through as the *source's*
	// values, no matter what dup carries -- setting them on dup above (as an
	// earlier version of this fix did) is silently a no-op. Left uncorrected,
	// the draft's namespace would keep the still-Active parent's ProjectID:
	// every compose write is gated by guardNamespaceWritable/
	// guardProjectWritable (compose/service/guard.go), which resolves the
	// *namespace's* ProjectID and rejects writes unless that project is a
	// draft, so every module/page/chart edit against the new revision -- the
	// entire reason a draft revision exists -- would fail with
	// project.errors.locked from the moment it's created. It would also stay
	// Enabled, contrary to the "draft namespaces are inactive" intent. Fix
	// both up explicitly post-clone.
	clonedNs.ProjectID = revID
	clonedNs.Enabled = false
	if err = store.UpdateComposeNamespace(ctx, svc.store, clonedNs); err != nil {
		return nil, fmt.Errorf("repoint cloned namespace to draft revision: %w", err)
	}

	rev = &types.Project{
		ID:               revID,
		TenantID:         parent.TenantID,
		ProjectID:        rootID,
		ParentRevisionID: parent.ID,
		Revision:         newRevision,
		Handle:           revHandle,
		Status:           types.ProjectStatusDraft,
		Meta:             parent.Meta,
		Config:           parent.Config,
		CreatedAt:        *now(),
		CreatedBy:        a.GetIdentityFromContext(ctx).Identity(),
	}
	rev.Config.NamespaceID = clonedNs.ID

	if err = store.CreateProject(ctx, svc.store, rev); err != nil {
		return nil, err
	}

	// The namespace clone above carried the compose side (modules, pages,
	// layouts, charts). Everything scoped by project rather than by namespace
	// has to be copied separately -- see project_revision_clone.go for why
	// envoy cannot do it. Runs after CreateProject so the draft row exists to
	// scope the copies to.
	//
	// On failure the draft is torn down again. Without this the half-built
	// revision survives -- the project row, its cloned namespace and whatever
	// resources were copied before the error -- and because the chain allows
	// only ONE draft at a time, every later attempt fails with "a draft
	// revision already exists" against a revision the user never saw created.
	// The compose clone and the copies below do not share a transaction (
	// CloneFromStore opens and commits its own), so this is a compensating
	// delete rather than a rollback.
	if err = svc.cloneProjectResources(ctx, svc.store, parent, rev); err != nil {
		svc.discardHalfBuiltRevision(ctx, rev)
		return nil, err
	}

	return rev, nil
}

// discardHalfBuiltRevision removes a draft revision whose creation failed part
// way through, so the chain is left free for another attempt.
//
// Best-effort by design: it runs while another error is already on its way up,
// and that error is the one worth reporting. A failure to clean up is logged
// rather than returned, since replacing the real cause with a cleanup error
// would hide why the revision failed in the first place.
func (svc *project) discardHalfBuiltRevision(ctx context.Context, rev *types.Project) {
	if rev == nil || rev.ID == 0 {
		return
	}

	if err := svc.setNamespaceDeleted(ctx, svc.store, rev.Config.NamespaceID, now()); err != nil && DefaultLogger != nil {
		DefaultLogger.Warn(
			"could not discard the namespace of a half-built revision",
			zap.Uint64("namespaceID", rev.Config.NamespaceID),
			zap.Error(err),
		)
	}

	if err := store.DeleteProjectByID(ctx, svc.store, rev.ID); err != nil && DefaultLogger != nil {
		DefaultLogger.Warn(
			"could not discard a half-built revision",
			zap.Uint64("projectID", rev.ID),
			zap.Error(err),
		)
	}
}

// revisionInChain reports whether revisionID (if non-zero) names a project
// row belonging to the same revision chain as rootProjectID — i.e. its own
// chain root matches. Zero revisionID is always valid: work items are
// deliberately allowed to be unassigned. A non-zero revisionID that fails to
// load is surfaced as-is (loadProject's not-found/invalid-ID error); a
// revision that loads fine but belongs to a different chain reports ok=false
// so the caller can reject it with its own resource-specific error.
func revisionInChain(ctx context.Context, s store.Storer, rootProjectID, revisionID uint64) (ok bool, err error) {
	if revisionID == 0 {
		return true, nil
	}

	rev, err := loadProject(ctx, s, revisionID)
	if err != nil {
		return false, err
	}

	return rev.RootProjectID() == rootProjectID, nil
}

// ListRevisions returns all projects in the same revision chain, sorted by revision number.
func (svc *project) ListRevisions(ctx context.Context, projectID uint64) (types.ProjectSet, error) {
	p, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}

	// Read on the project the caller named. The chain is one resource as far
	// as permissions go -- every row in it is a version of the same thing.
	if !svc.ac.CanReadProject(ctx, p) {
		return nil, ProjectErrNotAllowedToRead()
	}

	rootID := p.RootProjectID()
	// Deliberately unsorted at the store: `revision` is not declared sortable
	// in project.cue, so asking the store to order by it fails outright with
	// "invalid column name: revision" — which took the whole endpoint down, and
	// with it the UI's revision switcher. A chain is a handful of rows, so we
	// sort in Go below, after the root has been folded in (it has to be sorted
	// after that anyway: the root is fetched separately and would otherwise
	// just be prepended regardless of its revision number).
	set, _, err := store.SearchProjects(ctx, svc.store, types.ProjectFilter{
		RootProjectID: rootID,
	})
	if err != nil {
		return nil, err
	}

	// Also include the root project itself (its RootProjectID == 0, not rootID).
	root, rootErr := store.LookupProjectByID(ctx, svc.store, rootID)
	if rootErr == nil && root != nil {
		// Prepend root if not already in set.
		found := false
		for _, s := range set {
			if s.ID == root.ID {
				found = true
				break
			}
		}
		if !found {
			set = append(types.ProjectSet{root}, set...)
		}
	}

	// Oldest revision first — the order the switcher and any chain UI expect.
	sort.SliceStable(set, func(i, j int) bool { return set[i].Revision < set[j].Revision })

	return set, nil
}

// DeploymentPlan diffs the draft project against its parent and returns the
// deployment plan (changes + risk classification + suggested mappings).
func (svc *project) DeploymentPlan(ctx context.Context, projectID uint64) (*types.ProjectDeploymentPlan, error) {
	draft, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}

	// A plan is a read: it enumerates what the revision holds and what
	// publishing it would change. Gating it on publish rights would keep a
	// reviewer from seeing what they are being asked to approve.
	if !svc.ac.CanReadProject(ctx, draft) {
		return nil, ProjectErrNotAllowedToRead()
	}

	if draft.Status != types.ProjectStatusDraft {
		return nil, fmt.Errorf("deployment plan only available for draft projects")
	}
	// A first revision has nothing to diff against and nothing to migrate. That
	// is an empty plan, not an error — the publish flow asks for the plan before
	// it knows whether a parent exists, and an error would read as a failure.
	if draft.ParentRevisionID == 0 {
		return &types.ProjectDeploymentPlan{Risk: types.ProjectChangeRiskSafe}, nil
	}

	parent, err := loadProject(ctx, svc.store, draft.ParentRevisionID)
	if err != nil {
		return nil, err
	}

	return svc.computeDeploymentPlan(ctx, parent, draft)
}

// Publish migrates records from the parent namespace to the draft, flips statuses,
// and soft-deletes the old namespace. The request must carry confirm=true.
//
// The approval flow that gates publishing is enforced client-side; the backend
// publishes directly once confirm=true and the project is a draft.
func (svc *project) Publish(ctx context.Context, projectID uint64, req types.PublishRequest) (*types.Project, error) {
	if !req.Confirm {
		return nil, fmt.Errorf("publish requires confirm=true")
	}

	draft, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanPublishProject(ctx, draft) {
		return nil, ProjectErrNotAllowedToPublish()
	}

	if draft.Status != types.ProjectStatusDraft {
		return nil, fmt.Errorf("only draft projects can be published")
	}

	// First publish: a project with no parent revision has no prior namespace to
	// migrate from — its namespace is already the live one. Publishing simply
	// promotes the draft to active; no record migration, no namespace swap.
	if draft.ParentRevisionID == 0 {
		draft.Status = types.ProjectStatusActive
		draft.UpdatedAt = now()
		draft.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
		if err = store.UpdateProject(ctx, svc.store, draft); err != nil {
			return nil, err
		}
		return draft, nil
	}

	parent, err := loadProject(ctx, svc.store, draft.ParentRevisionID)
	if err != nil {
		return nil, err
	}

	oldNs, err := store.LookupComposeNamespaceByID(ctx, svc.store, parent.Config.NamespaceID)
	if err != nil {
		return nil, fmt.Errorf("load old namespace: %w", err)
	}
	newNs, err := store.LookupComposeNamespaceByID(ctx, svc.store, draft.Config.NamespaceID)
	if err != nil {
		return nil, fmt.Errorf("load new namespace: %w", err)
	}

	if err = svc.migrateRecords(ctx, oldNs, newNs, req.Mappings); err != nil {
		return nil, err
	}

	// Flip statuses and namespaces in a single transaction.
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		// Rename old namespace to mark it as a deprecated revision. The
		// timestamp alone is NOT a sufficient disambiguator: a project created
		// without a handle has an empty namespace slug (and newNs.Slug below
		// resets it to that same empty handle on every publish), so two such
		// projects publishing within the same second would both want
		// "-deprecated-<ts>" and collide on the slug's unique index. Prefixing
		// the project id makes it unique per project, the same fallback shape
		// CreateRevision uses for an empty handle/slug.
		oldNs.Slug = fmt.Sprintf("%s-deprecated-%d-%d", oldNs.Slug, draft.RootProjectID(), time.Now().Unix())
		oldNs.Enabled = false
		if e := store.UpdateComposeNamespace(ctx, s, oldNs); e != nil {
			return e
		}
		// Soft-delete old namespace (cascades to resources).
		if e := svc.setNamespaceDeleted(ctx, s, oldNs.ID, now()); e != nil {
			return e
		}

		// Rename new namespace to match the original project handle.
		newNs.Slug = parent.Handle
		newNs.Enabled = true
		if e := store.UpdateComposeNamespace(ctx, s, newNs); e != nil {
			return e
		}

		// Deprecate old project, activate new one.
		parent.Status = types.ProjectStatusDeprecated
		parent.UpdatedAt = now()
		parent.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
		if e := store.UpdateProject(ctx, s, parent); e != nil {
			return e
		}

		draft.Status = types.ProjectStatusActive
		draft.UpdatedAt = now()
		draft.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
		draft.Config.NamespaceID = newNs.ID
		return store.UpdateProject(ctx, s, draft)
	})
	if err != nil {
		return nil, err
	}

	return draft, nil
}

// migrateRecords creates a temporary internal DAL connection backed by the old
// namespace, runs DML imports for each module mapping, then cleans up.
func (svc *project) migrateRecords(
	ctx context.Context,
	oldNs, newNs *composeTypes.Namespace,
	mappings []types.ModuleMapping,
) error {
	if svc.services.dalSvc == nil || svc.services.dalConns == nil || svc.services.recordSvc == nil {
		return fmt.Errorf("project revision deps not initialised (dal)")
	}
	if len(mappings) == 0 {
		return nil
	}

	// Build and register the ephemeral internal DAL connection.
	connID := nextID()
	internalConn := dml.NewInternalConn(oldNs.ID, svc.store, svc.services.recordSvc)
	cw := dal.MakeConnection(connID, internalConn, dal.ConnectionParams{
		Type:   "corteza::dal:connection:internal",
		Params: map[string]any{"namespaceID": oldNs.ID},
	}, dal.ConnectionConfig{})
	if err := svc.services.dalConns.ReplaceConnection(ctx, cw, false); err != nil {
		return fmt.Errorf("register internal dal connection: %w", err)
	}
	defer func() {
		_ = svc.services.dalConns.RemoveConnection(ctx, connID)
	}()

	// Store the ephemeral DmlConnection row so the Importer can find it.
	dbConn := &types.DmlConnection{
		ID:        connID,
		Handle:    "internal-publish-" + strconv.FormatUint(connID, 10),
		Label:     "internal publish",
		CreatedAt: time.Now(),
		Params: types.DmlConnectionParams{
			Type:   "corteza::dal:connection:internal",
			Params: map[string]any{"namespaceID": oldNs.ID},
		},
	}
	if err := store.CreateDmlConnection(ctx, svc.store, dbConn); err != nil {
		return fmt.Errorf("store ephemeral dml connection: %w", err)
	}
	defer func() {
		_ = store.DeleteDmlConnectionByID(ctx, svc.store, connID)
	}()

	for _, m := range mappings {
		// A mapping with nothing to copy from is a module the revision ADDED:
		// diffModules emits one for every new module so the caller sees it in
		// the plan, but there is no source table behind it. Running it anyway
		// failed the whole publish with `source table "x" not found on
		// connection` -- the plan's own suggestions could not be published back.
		if !m.HasSource() {
			continue
		}

		if err := svc.runModuleMapping(ctx, connID, newNs.Slug, m); err != nil {
			return err
		}
	}
	return nil
}

func (svc *project) runModuleMapping(
	ctx context.Context,
	connID uint64,
	newNsHandle string,
	m types.ModuleMapping,
) error {
	cols := make(types.DmlColumnMapSet, 0, len(m.Fields))
	for _, f := range m.Fields {
		cols = append(cols, &types.DmlColumnMap{
			SourceIdent: f.SourceField,
			FieldName:   f.TargetField,
		})
	}

	mp := &types.DmlMapping{
		ID:              nextID(),
		ConnectionID:    connID,
		NamespaceHandle: newNsHandle,
		SourceIdent:     m.Module,
		ModuleHandle:    m.Module,
		Columns:         cols,
		CreatedAt:       time.Now(),
	}
	if err := store.CreateDmlMapping(ctx, svc.store, mp); err != nil {
		return fmt.Errorf("store dml mapping for module %q: %w", m.Module, err)
	}
	defer func() {
		_ = store.DeleteDmlMappingByID(ctx, svc.store, mp.ID)
	}()

	return svc.services.dalSvc.RunImportForeground(ctx, mp.ID)
}

// computeDeploymentPlan diffs a draft revision against its parent: every
// resource kind the project graph enumerates, plus field-level detail and the
// record-migration mapping publish needs for compose modules.
func (svc *project) computeDeploymentPlan(
	ctx context.Context,
	parent, draft *types.Project,
) (*types.ProjectDeploymentPlan, error) {
	plan := &types.ProjectDeploymentPlan{Risk: types.ProjectChangeRiskSafe}

	if err := svc.diffResources(ctx, parent, draft, plan); err != nil {
		return nil, err
	}
	if err := svc.diffModules(ctx, parent, draft, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// addChange records a change and carries its risk up to the plan.
func addChange(plan *types.ProjectDeploymentPlan, c types.ProjectChange) {
	plan.Changes = append(plan.Changes, c)
	if c.Risk == types.ProjectChangeRiskDangerous {
		plan.Risk = types.ProjectChangeRiskDangerous
	}
}

// resourceDiffSkip lists the kinds diffResources leaves alone: modules are
// diffed field-by-field by diffModules; users are project membership rather
// than something a publish deploys; knowledge bases are loaded unscoped by the
// graph, so they are identical on both sides by construction.
var resourceDiffSkip = map[string]bool{
	"module":         true,
	"user":           true,
	"knowledge-base": true,
}

// diffResources compares the two revisions across every other kind, reusing the
// enumeration the resource graph is built from so the publish screen and the
// Build canvas can never disagree about what a revision contains.
func (svc *project) diffResources(
	ctx context.Context,
	parent, draft *types.Project,
	plan *types.ProjectDeploymentPlan,
) error {
	g := ProjectGraph(svc.store)

	_, oldSrc, err := g.fetch(ctx, parent.ID)
	if err != nil {
		return fmt.Errorf("load parent revision resources: %w", err)
	}
	_, newSrc, err := g.fetch(ctx, draft.ID)
	if err != nil {
		return fmt.Errorf("load draft revision resources: %w", err)
	}

	for _, c := range diffSources(oldSrc, newSrc) {
		addChange(plan, c)
	}

	return nil
}

// diffSources is the resource diff itself, kept free of the store so it can be
// tested directly.
//
// Identity is kind + handle: IDs are freshly minted by the revision clone, so
// they never match across revisions. A rename therefore reads as a removal plus
// an addition — the honest reading, since the migration cannot tell those apart
// either.
func diffSources(oldSrc, newSrc []*GraphSource) []types.ProjectChange {
	index := func(ss []*GraphSource) map[string]*GraphSource {
		out := make(map[string]*GraphSource, len(ss))
		for _, s := range ss {
			if resourceDiffSkip[s.Kind] {
				continue
			}
			out[s.Kind+"."+firstNonEmpty(s.Handle, s.Name)] = s
		}
		return out
	}

	oldByKey, newByKey := index(oldSrc), index(newSrc)

	keys := make([]string, 0, len(oldByKey)+len(newByKey))
	for k := range oldByKey {
		keys = append(keys, k)
	}
	for k := range newByKey {
		if _, both := oldByKey[k]; !both {
			keys = append(keys, k)
		}
	}
	// Map iteration is unordered and a human reads this; fix an order.
	sort.Strings(keys)

	out := make([]types.ProjectChange, 0, len(keys))
	for _, k := range keys {
		o, inOld := oldByKey[k]
		n, inNew := newByKey[k]

		switch {
		case inOld && inNew:
			// On both sides. Whether its innards changed is a question for the
			// resource's own editor, not for a migration plan.
		case inNew:
			out = append(out, types.ProjectChange{
				Op:   types.ProjectChangeOpAdded,
				Kind: n.Kind,
				Name: firstNonEmpty(n.Name, n.Handle),
				Path: k,
				Risk: types.ProjectChangeRiskSafe,
			})
		default:
			// Safe on purpose: dropping a page or an automation costs no
			// records. Only module data is at stake, and diffModules owns that.
			out = append(out, types.ProjectChange{
				Op:   types.ProjectChangeOpRemoved,
				Kind: o.Kind,
				Name: firstNonEmpty(o.Name, o.Handle),
				Path: k,
				Risk: types.ProjectChangeRiskSafe,
			})
		}
	}

	return out
}

// diffModules compares compose modules between the two namespaces and builds
// the per-module mapping publish migrates records with. Namespace-scoped rather
// than project-scoped because the migration itself runs namespace to namespace.
func (svc *project) diffModules(
	ctx context.Context,
	parent, draft *types.Project,
	plan *types.ProjectDeploymentPlan,
) error {
	oldNsID, newNsID := parent.Config.NamespaceID, draft.Config.NamespaceID

	oldMods, _, err := store.SearchComposeModules(ctx, svc.store, composeTypes.ModuleFilter{NamespaceID: oldNsID})
	if err != nil {
		return fmt.Errorf("load old namespace modules: %w", err)
	}
	newMods, _, err := store.SearchComposeModules(ctx, svc.store, composeTypes.ModuleFilter{NamespaceID: newNsID})
	if err != nil {
		return fmt.Errorf("load new namespace modules: %w", err)
	}

	oldByHandle := make(map[string]*composeTypes.Module, len(oldMods))
	for _, m := range oldMods {
		oldByHandle[m.Handle] = m
	}
	newByHandle := make(map[string]*composeTypes.Module, len(newMods))
	for _, m := range newMods {
		newByHandle[m.Handle] = m
	}

	for _, om := range oldMods {
		if _, kept := newByHandle[om.Handle]; kept {
			continue
		}
		addChange(plan, types.ProjectChange{
			Op:      types.ProjectChangeOpRemoved,
			Kind:    "module",
			Name:    firstNonEmpty(om.Name, om.Handle),
			Path:    "module." + om.Handle,
			Risk:    types.ProjectChangeRiskDangerous,
			Records: svc.moduleRecordCount(ctx, oldNsID, om),
		})
	}

	for _, nm := range newMods {
		om, shared := oldByHandle[nm.Handle]
		if !shared {
			addChange(plan, types.ProjectChange{
				Op:   types.ProjectChangeOpAdded,
				Kind: "module",
				Name: firstNonEmpty(nm.Name, nm.Handle),
				Path: "module." + nm.Handle,
				Risk: types.ProjectChangeRiskSafe,
			})
			// A new module starts empty, so there is nothing to map into it.
			plan.SuggestedMappings = append(plan.SuggestedMappings, types.ModuleMapping{Module: nm.Handle})
			continue
		}

		oldFields, newFields := svc.diffModuleFields(ctx, oldNsID, om, nm, plan)
		plan.SuggestedMappings = append(
			plan.SuggestedMappings,
			buildSuggestedMappingFromFields(nm.Handle, oldFields, newFields),
		)
	}

	return nil
}

// diffModuleFields records what changed between two versions of one module and
// returns both field sets for the mapping builder.
func (svc *project) diffModuleFields(
	ctx context.Context,
	oldNsID uint64,
	om, nm *composeTypes.Module,
	plan *types.ProjectDeploymentPlan,
) (oldFields, newFields composeTypes.ModuleFieldSet) {
	oldFields, _, _ = store.SearchComposeModuleFields(ctx, svc.store, composeTypes.ModuleFieldFilter{ModuleID: []uint64{om.ID}})
	newFields, _, _ = store.SearchComposeModuleFields(ctx, svc.store, composeTypes.ModuleFieldFilter{ModuleID: []uint64{nm.ID}})

	oldByName := make(map[string]*composeTypes.ModuleField, len(oldFields))
	for _, f := range oldFields {
		oldByName[f.Name] = f
	}
	newByName := make(map[string]*composeTypes.ModuleField, len(newFields))
	for _, f := range newFields {
		newByName[f.Name] = f
	}

	// Counted at most once per module, and only if something destructive turns
	// up: the count exists to size a warning, and most revisions never raise one.
	records, counted := uint64(0), false
	atStake := func() uint64 {
		if !counted {
			records, counted = svc.moduleRecordCount(ctx, oldNsID, om), true
		}
		return records
	}

	for _, of := range oldFields {
		if _, kept := newByName[of.Name]; kept {
			continue
		}
		addChange(plan, types.ProjectChange{
			Op:      types.ProjectChangeOpRemoved,
			Kind:    "field",
			Module:  om.Handle,
			Name:    of.Name,
			Path:    "module." + om.Handle + ".field." + of.Name,
			Risk:    types.ProjectChangeRiskDangerous,
			Records: atStake(),
		})
	}

	for _, nf := range newFields {
		of, shared := oldByName[nf.Name]
		if !shared {
			addChange(plan, types.ProjectChange{
				Op:     types.ProjectChangeOpAdded,
				Kind:   "field",
				Module: nm.Handle,
				Name:   nf.Name,
				Path:   "module." + nm.Handle + ".field." + nf.Name,
				Risk:   types.ProjectChangeRiskSafe,
			})
			continue
		}
		if of.Kind != nf.Kind {
			addChange(plan, types.ProjectChange{
				Op:      types.ProjectChangeOpChanged,
				Kind:    "field",
				Module:  nm.Handle,
				Name:    nf.Name,
				Detail:  of.Kind + " → " + nf.Kind,
				Path:    "module." + nm.Handle + ".field." + nf.Name,
				Risk:    types.ProjectChangeRiskDangerous,
				Records: atStake(),
			})
		}
	}

	return
}

// moduleRecordCount is how many records a module holds — what a removed or
// retyped field puts at stake. Best-effort by design: the count only sharpens
// the warning, so a failure returns 0 ("unknown") instead of failing the plan.
func (svc *project) moduleRecordCount(ctx context.Context, nsID uint64, m *composeTypes.Module) uint64 {
	if svc.services.recordSvc == nil || m == nil {
		return 0
	}

	_, f, err := svc.services.recordSvc.Search(ctx, composeTypes.RecordFilter{
		ModuleID:    m.ID,
		NamespaceID: nsID,
		Paging:      filter.Paging{Limit: 1, IncTotal: true},
	})
	if err != nil {
		return 0
	}

	return uint64(f.Total)
}

func buildSuggestedMappingFromFields(
	handle string,
	oldFields, newFields composeTypes.ModuleFieldSet,
) types.ModuleMapping {
	mm := types.ModuleMapping{Module: handle}

	oldByName := make(map[string]*composeTypes.ModuleField)
	for _, f := range oldFields {
		oldByName[f.Name] = f
	}

	for _, nf := range newFields {
		of, exists := oldByName[nf.Name]
		if !exists {
			// New field: default value (empty).
			mm.Fields = append(mm.Fields, types.ModuleFieldMapping{
				TargetField: nf.Name,
				Op:          "default",
				Value:       "",
			})
			continue
		}
		if of.Kind != nf.Kind {
			// Type change: suggest cast.
			mm.Fields = append(mm.Fields, types.ModuleFieldMapping{
				SourceField: of.Name,
				TargetField: nf.Name,
				Op:          "cast",
				OnError:     "null",
			})
		} else {
			// Same type: copy.
			mm.Fields = append(mm.Fields, types.ModuleFieldMapping{
				SourceField: of.Name,
				TargetField: nf.Name,
				Op:          "copy",
			})
		}
	}

	return mm
}
