package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertQueueMessageFilter(f types.QueueMessageFilter) (query squirrel.SelectBuilder, err error) {
	query = s.queueMessagesSelectBuilder()

	if f.Queue != "" {
		query = query.Where("mqm.queue = ?", f.Queue)
	}

	query = filter.StateCondition(query, "mqm.processed", f.Processed)
	return
}
