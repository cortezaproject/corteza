package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

func (s Store) convertMessagebusQueuesettingFilter(f messagebus.QueueSettingsFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagebusQueuesettingsSelectBuilder()
	query = filter.StateCondition(query, "mqs.deleted_at", f.Deleted)
	return
}
