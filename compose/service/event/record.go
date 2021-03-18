package event

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"strings"
	"time"
)

const (
	recordMatchValues = "record.values."
)

// Match returns false if given conditions do not match event & resource internals
func (res recordBase) Match(c eventbus.ConstraintMatcher) bool {
	return eventbus.MatchFirst(
		func() bool { return recordMatch(res.record, c) },
		func() bool { return moduleMatch(res.module, c) },
		func() bool { return namespaceMatch(res.namespace, c) },
	)
}

func recordMatch(r *types.Record, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "record.updatedAt", "record.updated-at":
		return c.Match(r.UpdatedAt.Format(time.RFC3339))
	case "record.createdAt", "record.created-at":
		return c.Match(r.CreatedAt.Format(time.RFC3339))
	case "record.deletedAt", "record.deleted-at":
		return c.Match(r.DeletedAt.Format(time.RFC3339))
	}

	if strings.HasPrefix(c.Name(), recordMatchValues) {
		fieldName := c.Name()[len(recordMatchValues):]
		for _, v := range r.Values.FilterByName(fieldName) {
			if c.Match(v.Value) {
				return true
			}
		}
	}

	return false
}
