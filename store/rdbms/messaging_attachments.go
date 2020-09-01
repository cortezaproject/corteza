package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

func (s Store) convertMessagingAttachmentFilter(f types.AttachmentFilter) (query squirrel.SelectBuilder, err error) {
	query = s.composePagesSelectBuilder()

	// @todo join & filter by message

	return
}
