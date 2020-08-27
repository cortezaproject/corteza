package rdbms

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
)

func (s Store) convertComposeAttachmentFilter(f types.AttachmentFilter) (query squirrel.SelectBuilder, err error) {
	query = s.attachmentsSelectBuilder()

	if f.Kind != "" {
		query = query.Where(squirrel.Eq{"att.kind": f.Kind})
	}

	if f.NamespaceID > 0 {
		query = query.Where("a.rel_namespace = ?", f.NamespaceID)
	}

	switch f.Kind {
	case types.PageAttachment:
		// @todo implement filtering by page
		if f.PageID > 0 {
			err = fmt.Errorf("filtering by pageID not implemented")
			return
		}

	case types.RecordAttachment:
		query = query.
			Join("compose_record_value AS v ON (v.ref = a.id)")

		if f.ModuleID > 0 {
			query = query.
				Join("compose_record AS r ON (r.id = v.record_id)").
				Where(squirrel.Eq{"r.module_id": f.ModuleID})
		}

		if f.RecordID > 0 {
			query = query.Where(squirrel.Eq{"v.record_id": f.RecordID})
		}

		if f.FieldName != "" {
			query = query.Where(squirrel.Eq{"v.name": f.FieldName})
		}

	default:
		err = fmt.Errorf("unsupported kind value")
		return
	}

	if f.Filter != "" {
		err = fmt.Errorf("filtering by filter not implemented")
		return
	}

	return
}
