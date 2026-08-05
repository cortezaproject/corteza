package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
)

// Project-state resource lock.
//
// A project's resources are editable only while the project is a draft (this
// includes draft revisions, which are separate draft projects). Once a project
// is active/published/archived/suspended/deprecated its resources are frozen —
// changes are made on a new draft revision and then published.
//
// The generated services call svc.guard(ctx, res) on Update/Delete/Undelete
// (reads are never guarded). The per-resource guard() methods (in the
// hand-written service files, enabled via `service.guard: true` in the .cue)
// delegate here.

func projectLockedErr() *errors.Error {
	return errors.New(
		errors.KindInternal,
		"the project is read-only — only draft projects can be modified",
		errors.Meta("type", "projectLocked"),
		errors.Meta("resource", "compose:project"),
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.locked"),
	)
}

// guardProjectWritable blocks writes to a resource owned by a non-draft project.
// A zero projectID means the resource is not project-scoped and is writable.
func guardProjectWritable(ctx context.Context, s store.Storer, projectID uint64) error {
	if projectID == 0 {
		return nil
	}

	p, err := store.LookupProjectByID(ctx, s, projectID)
	if err != nil {
		if errors.IsNotFound(err) {
			// A resource pointing at a missing project is not something the
			// lock should block on.
			return nil
		}
		return err
	}

	if p.Status != systemTypes.ProjectStatusDraft {
		return projectLockedErr()
	}

	return nil
}

// guardNamespaceWritable resolves the owning project through the namespace —
// compose child resources (module, chart, page, page layout) scope by namespace
// rather than carrying a reliable ProjectID of their own — and guards it.
func guardNamespaceWritable(ctx context.Context, s store.Storer, namespaceID uint64) error {
	if namespaceID == 0 {
		return nil
	}

	ns, err := store.LookupComposeNamespaceByID(ctx, s, namespaceID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return guardProjectWritable(ctx, s, ns.ProjectID)
}
