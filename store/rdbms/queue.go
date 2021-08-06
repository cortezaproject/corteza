package rdbms

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertQueueFilter(f types.QueueFilter) (query squirrel.SelectBuilder, err error) {
	query = s.queuesSelectBuilder()
	query = filter.StateCondition(query, "mqs.deleted_at", f.Deleted)

	if f.Queue != "" {
		query = query.Where(squirrel.Like{"LOWER(mqs.queue)": strings.ToLower(f.Queue)})
	}

	return
}
