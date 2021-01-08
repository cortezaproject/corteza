package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

func (s Store) convertMessagingMessageAttachmentFilter(f types.MessageAttachmentFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagingMessageAttachmentsSelectBuilder()

	if len(f.MessageID) > 0 {
		query = query.Where(squirrel.Eq{"mma.rel_message": f.MessageID})
	}

	return
}
