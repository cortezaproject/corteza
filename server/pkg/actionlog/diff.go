package actionlog

import "github.com/crusttech/human/server/pkg/revisions"

var diffSkip = []string{
	"createdAt", "updatedAt", "deletedAt",
	"createdBy", "updatedBy", "deletedBy",
	"createdByAgent", "ownedBy", "issues",
}

// DiffResourceState computes a field-level diff between old and updated resource states.
func DiffResourceState(old, updated any) []*revisions.Change {
	changes, _ := revisions.DiffAny(updated, old, diffSkip...)
	return changes
}
