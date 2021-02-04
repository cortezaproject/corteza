package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertTemplateFilter(f types.TemplateFilter) (query squirrel.SelectBuilder, err error) {
	query = s.templatesSelectBuilder()

	query = filter.StateCondition(query, "tpl.deleted_at", f.Deleted)

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"tpl.id": f.LabeledIDs})
	}

	if f.Partial {
		query = query.Where(squirrel.Eq{"tpl.partial": true})
	}

	if len(f.TemplateID) > 0 {
		query = query.Where(squirrel.Eq{"tpl.id": f.TemplateID})
	}
	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"tpl.handle": f.Handle})
	}
	if f.Type != "" {
		query = query.Where(squirrel.Eq{"tpl.type": f.Type})
	}
	if f.OwnerID > 0 {
		query = query.Where(squirrel.Eq{"tpl.rel_owner": f.OwnerID})
	}

	return
}
