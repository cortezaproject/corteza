package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/system/types"
)

func (s Store) convertAttachmentFilter(f types.AttachmentFilter) (query squirrel.SelectBuilder, err error) {
	query = s.attachmentsSelectBuilder()

	if f.Kind != "" {
		query = query.Where(squirrel.Eq{"att.kind": f.Kind})
	}

	return
}
