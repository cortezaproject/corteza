package rdbms

import (
	"github.com/Masterminds/squirrel"
	"sort"
	"strconv"
)

func (s Store) getMessagingChannelMembersQuery(memberIDs ...uint64) squirrel.SelectBuilder {
	if len(memberIDs) == 0 {
		return squirrel.
			Select("null")
	}

	// Make sure members are sorted
	sort.Slice(memberIDs, func(i, j int) bool {
		return memberIDs[i] < memberIDs[j]
	})

	// Concatentating members fore
	list := ""
	for i := range memberIDs {
		if i > 0 {
			list += ","
		}
		list += strconv.FormatUint(memberIDs[i], 10)
	}

	return s.messagingChannelMembersSelectBuilder().
		GroupBy("cm.rel_channel").
		Having(squirrel.Eq{
			"COUNT(*)": len(memberIDs),
			"GROUP_CONCAT(cm.rel_user ORDER BY 1 ASC SEPARATOR ',')": list,
		})
}
