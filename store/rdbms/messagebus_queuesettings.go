package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/federation/types"
)

func (s Store) convertMessagebusQueuesettingsFilter(f types.ExposedModuleFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagebusQueuesettingsSelectBuilder()
	return
}
