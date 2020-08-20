package rdbms

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

func (s Store) convertComposeModuleFieldFilter(f types.ModuleFieldFilter) (query squirrel.SelectBuilder, err error) {
	query = s.composeModuleFieldsSelectBuilder()

	if len(f.ModuleID) == 0 {
		err = fmt.Errorf("can not search for module fields without module IDs")
		return
	}

	query = rh.FilterNullByState(query, "cmf.deleted_at", f.Deleted)
	query = query.Where(squirrel.Eq{"cmf.rel_module": f.ModuleID})

	return
}
