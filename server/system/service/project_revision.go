package service

import (
	"context"
	"fmt"
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
	DefaultProject.nsSvc = nsSvc
	DefaultProject.dalSvc = importer
	DefaultProject.dalConns = dalConns
	DefaultProject.recordSvc = recordSvc
}

// CreateRevision clones the project and its namespace into a new draft revision.
func (svc *project) CreateRevision(ctx context.Context, projectID uint64) (rev *types.Project, err error) {
	var parent *types.Project
	if parent, err = loadProject(ctx, svc.store, projectID); err != nil {
		return
	}

	if parent.Status != types.ProjectStatusActive {
		return nil, fmt.Errorf("can only revise active projects")
	}

	// one draft per chain at a time
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
	revSlug := oldNs.Slug + "-rev" + strconv.Itoa(newRevision)
	// Envoy requires a non-empty slug for clone.
	if revSlug == "-rev"+strconv.Itoa(newRevision) {
		revSlug = "rev-" + strconv.FormatUint(rootID, 10) + "-" + strconv.Itoa(newRevision)
	}

	if svc.nsSvc == nil {
		return nil, fmt.Errorf("project revision deps not initialised (nsSvc)")
	}

	dup := &composeTypes.Namespace{
		TenantID:  oldNs.TenantID,
		ProjectID: oldNs.ProjectID,
		Name:      oldNs.Name + " (revision " + strconv.Itoa(newRevision) + ")",
		Slug:      revSlug,
		Enabled:   false, // draft namespaces are inactive
	}
	clonedNs, err := svc.nsSvc.CloneFromStore(ctx, parent.Config.NamespaceID, dup)
	if err != nil {
		return nil, fmt.Errorf("clone namespace for revision: %w", err)
	}

	rev = &types.Project{
		ID:               nextID(),
		TenantID:         parent.TenantID,
		ProjectID:        rootID,
		ParentRevisionID: parent.ID,
		Revision:         newRevision,
		Handle:           parent.Handle + "-rev" + strconv.Itoa(newRevision),
		Status:           types.ProjectStatusDraft,
		Meta:             parent.Meta,
		Config:           parent.Config,
		Governance:       parent.Governance,
		CreatedAt:        *now(),
		CreatedBy:        a.GetIdentityFromContext(ctx).Identity(),
	}
	rev.Config.NamespaceID = clonedNs.ID

	if err = store.CreateProject(ctx, svc.store, rev); err != nil {
		return nil, err
	}
	return rev, nil
}

// ListRevisions returns all projects in the same revision chain, sorted by revision number.
func (svc *project) ListRevisions(ctx context.Context, projectID uint64) (types.ProjectSet, error) {
	p, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}

	rootID := p.RootProjectID()
	set, _, err := store.SearchProjects(ctx, svc.store, types.ProjectFilter{
		RootProjectID: rootID,
		Sorting:       filter.Sorting{Sort: filter.SortExprSet{{Column: "revision"}}},
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

	return set, nil
}

// DeploymentPlan diffs the draft project against its parent and returns the
// deployment plan (changes + risk classification + suggested mappings).
func (svc *project) DeploymentPlan(ctx context.Context, projectID uint64) (*types.ProjectDeploymentPlan, error) {
	draft, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}
	if draft.Status != types.ProjectStatusDraft {
		return nil, fmt.Errorf("deployment plan only available for draft projects")
	}
	if draft.ParentRevisionID == 0 {
		return nil, fmt.Errorf("project has no parent revision")
	}

	parent, err := loadProject(ctx, svc.store, draft.ParentRevisionID)
	if err != nil {
		return nil, err
	}

	return computeDeploymentPlan(ctx, svc.store, parent.Config.NamespaceID, draft.Config.NamespaceID)
}

// Publish migrates records from the parent namespace to the draft, flips statuses,
// and soft-deletes the old namespace. The request must carry confirm=true.
func (svc *project) Publish(ctx context.Context, projectID uint64, req types.PublishRequest) (*types.Project, error) {
	if !req.Confirm {
		return nil, fmt.Errorf("publish requires confirm=true")
	}

	draft, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}
	if draft.Status != types.ProjectStatusDraft {
		return nil, fmt.Errorf("only draft projects can be published")
	}
	if draft.ParentRevisionID == 0 {
		return nil, fmt.Errorf("project has no parent revision")
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
		// Rename old namespace to mark it as a deprecated revision.
		oldNs.Slug = oldNs.Slug + "-deprecated-" + strconv.FormatInt(time.Now().Unix(), 10)
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
	if svc.dalSvc == nil || svc.dalConns == nil || svc.recordSvc == nil {
		return fmt.Errorf("project revision deps not initialised (dal)")
	}
	if len(mappings) == 0 {
		return nil
	}

	// Build and register the ephemeral internal DAL connection.
	connID := nextID()
	internalConn := dml.NewInternalConn(oldNs.ID, svc.store, svc.recordSvc)
	cw := dal.MakeConnection(connID, internalConn, dal.ConnectionParams{
		Type:   "corteza::dal:connection:internal",
		Params: map[string]any{"namespaceID": oldNs.ID},
	}, dal.ConnectionConfig{})
	if err := svc.dalConns.ReplaceConnection(ctx, cw, false); err != nil {
		return fmt.Errorf("register internal dal connection: %w", err)
	}
	defer func() {
		_ = svc.dalConns.RemoveConnection(ctx, connID)
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

	return svc.dalSvc.RunImportForeground(ctx, mp.ID)
}

// computeDeploymentPlan diffs modules between two namespaces.
func computeDeploymentPlan(
	ctx context.Context,
	s store.Storer,
	oldNsID, newNsID uint64,
) (*types.ProjectDeploymentPlan, error) {
	oldMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{NamespaceID: oldNsID})
	if err != nil {
		return nil, fmt.Errorf("load old namespace modules: %w", err)
	}
	newMods, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{NamespaceID: newNsID})
	if err != nil {
		return nil, fmt.Errorf("load new namespace modules: %w", err)
	}

	oldByHandle := make(map[string]*composeTypes.Module, len(oldMods))
	for _, m := range oldMods {
		oldByHandle[m.Handle] = m
	}
	newByHandle := make(map[string]*composeTypes.Module, len(newMods))
	for _, m := range newMods {
		newByHandle[m.Handle] = m
	}

	plan := &types.ProjectDeploymentPlan{}
	overallRisk := types.ProjectChangeRiskSafe

	// Detect deleted modules.
	for handle, om := range oldByHandle {
		if _, exists := newByHandle[handle]; !exists {
			plan.Changes = append(plan.Changes, types.ProjectChange{
				Path: "module." + handle,
				Risk: types.ProjectChangeRiskDangerous,
			})
			overallRisk = types.ProjectChangeRiskDangerous
			_ = om
		}
	}

	// Detect added modules and field changes within shared modules.
	for handle, nm := range newByHandle {
		if _, exists := oldByHandle[handle]; !exists {
			plan.Changes = append(plan.Changes, types.ProjectChange{
				Path: "module." + handle,
				Risk: types.ProjectChangeRiskSafe,
			})
			// New module: suggest an identity mapping.
			plan.SuggestedMappings = append(plan.SuggestedMappings, buildSuggestedMapping(nil, nm))
			continue
		}
		// Shared module: diff fields.
		oldFields, newFields, fieldChanges, risk := diffModuleFields(ctx, s, oldByHandle[handle], nm)
		plan.Changes = append(plan.Changes, fieldChanges...)
		if risk == types.ProjectChangeRiskDangerous {
			overallRisk = types.ProjectChangeRiskDangerous
		}
		plan.SuggestedMappings = append(plan.SuggestedMappings, buildSuggestedMappingFromFields(handle, oldFields, newFields))
	}

	plan.Path = string(overallRisk)
	return plan, nil
}

func diffModuleFields(
	ctx context.Context,
	s store.Storer,
	oldMod, newMod *composeTypes.Module,
) (
	oldFields, newFields composeTypes.ModuleFieldSet,
	changes []types.ProjectChange,
	maxRisk types.ProjectChangeRisk,
) {
	maxRisk = types.ProjectChangeRiskSafe

	oldFields, _, _ = store.SearchComposeModuleFields(ctx, s, composeTypes.ModuleFieldFilter{ModuleID: []uint64{oldMod.ID}})
	newFields, _, _ = store.SearchComposeModuleFields(ctx, s, composeTypes.ModuleFieldFilter{ModuleID: []uint64{newMod.ID}})

	oldByName := make(map[string]*composeTypes.ModuleField)
	for _, f := range oldFields {
		oldByName[f.Name] = f
	}
	newByName := make(map[string]*composeTypes.ModuleField)
	for _, f := range newFields {
		newByName[f.Name] = f
	}

	for name := range oldByName {
		if _, exists := newByName[name]; !exists {
			changes = append(changes, types.ProjectChange{
				Path: "module." + oldMod.Handle + ".field." + name,
				Risk: types.ProjectChangeRiskDangerous,
			})
			maxRisk = types.ProjectChangeRiskDangerous
		}
	}
	for name, nf := range newByName {
		of, exists := oldByName[name]
		if !exists {
			changes = append(changes, types.ProjectChange{
				Path: "module." + newMod.Handle + ".field." + name,
				Risk: types.ProjectChangeRiskSafe,
			})
			continue
		}
		if of.Kind != nf.Kind {
			changes = append(changes, types.ProjectChange{
				Path: "module." + newMod.Handle + ".field." + name,
				Risk: types.ProjectChangeRiskDangerous,
			})
			maxRisk = types.ProjectChangeRiskDangerous
		}
	}

	return
}

func buildSuggestedMapping(oldMod, newMod *composeTypes.Module) types.ModuleMapping {
	mm := types.ModuleMapping{Module: newMod.Handle}
	// For a new module, no field mappings are needed.
	return mm
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
