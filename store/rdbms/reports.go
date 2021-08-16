package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertReportFilter(f types.ReportFilter) (query squirrel.SelectBuilder, err error) {
	query = s.reportsSelectBuilder()

	query = filter.StateCondition(query, "rp.deleted_at", f.Deleted)

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"rp.id": f.LabeledIDs})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"rp.handle": f.Handle})
	}

	return
}
