package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

func (s Store) convertMessagingMentionFilter(f types.MentionFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagingMentionsSelectBuilder()

	if len(f.MessageID) > 0 {
		query.Where(squirrel.Eq{"rel_message": f.MessageID})
	}

	return
}
