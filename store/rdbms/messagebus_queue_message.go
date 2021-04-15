package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

func (s Store) convertMessagebusQueueMessageFilter(f messagebus.QueueMessageFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagebusQueueMessagesSelectBuilder()

	if f.Queue != "" {
		query = query.Where("mqm.queue = ?", f.Queue)
	}

	query = filter.StateCondition(query, "mqm.processed", f.Processed)
	return
}
