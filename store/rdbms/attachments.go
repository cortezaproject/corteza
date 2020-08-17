package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertAttachmentFilter(f types.AttachmentFilter) (query squirrel.SelectBuilder, err error) {
	if f.Sort == "" {
		f.Sort = "id"
	}

	query = s.QueryAttachments()

	if f.Kind != "" {
		query = query.Where(squirrel.Eq{"att.kind": f.Kind})
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, s.AttachmentColumns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	return
}
