package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

func (s Store) convertMessagingFlagFilter(f types.MessageFlagFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagingFlagsSelectBuilder()

	if f.Flag != "" {
		query = query.Where(squirrel.Eq{"flag": f.Flag})
	}

	if len(f.MessageID) > 0 {
		query = query.Where(squirrel.Eq{"rel_message": f.MessageID})
	}

	return query, nil
}
