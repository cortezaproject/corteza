package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// guardProjectWritable blocks writes (update/delete/undelete) to a resource
// owned by a non-draft project. Only draft projects (including draft revisions)
// are editable; a zero projectID means the resource is not project-scoped and
// is always writable. See compose/service/guard.go for the compose counterpart.
func guardProjectWritable(ctx context.Context, s store.Storer, projectID uint64) error {
	if projectID == 0 {
		return nil
	}

	p, err := store.LookupProjectByID(ctx, s, projectID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if p.Status != types.ProjectStatusDraft {
		return errors.New(
			errors.KindInternal,
			"project is read-only: only draft projects can be modified",
			errors.Meta("type", "projectLocked"),
			errors.Meta("resource", "system:project"),
			errors.Meta(locale.ErrorMetaNamespace{}, "system"),
			errors.Meta(locale.ErrorMetaKey{}, "project.errors.locked"),
		)
	}

	return nil
}
