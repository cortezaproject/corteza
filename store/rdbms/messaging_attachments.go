package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

func (s Store) convertMessagingAttachmentFilter(f types.AttachmentFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagingAttachmentsSelectBuilder()

	if len(f.AttachmentID) > 0 {
		query = query.Where(squirrel.Eq{"att.id": f.AttachmentID})
	}

	return
}
