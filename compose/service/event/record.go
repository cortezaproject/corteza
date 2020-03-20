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
	return recordMatch(res.record, c, namespaceMatch(res.namespace, c, moduleMatch(res.module, c, false)))
}

func recordMatch(r *types.Record, c eventbus.ConstraintMatcher, def bool) bool {
	switch c.Name() {
	case "record.updatedAt":
		return c.Match(r.UpdatedAt.Format(time.RFC3339))
	case "record.createdAt":
		return c.Match(r.CreatedAt.Format(time.RFC3339))
	case "record.deletedAt":
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

	return def
}
