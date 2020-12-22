package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"sort"
	"strconv"
	"strings"
)

func (s Store) convertMessagingChannelMemberFilter(f types.ChannelMemberFilter) (query squirrel.SelectBuilder, err error) {
	query = s.messagingChannelMembersSelectBuilder()

	if len(f.ChannelID) > 0 {
		query = query.Where(squirrel.Eq{"rel_channel": f.ChannelID})
	}

	if len(f.MemberID) > 0 {
		query = query.Where(squirrel.Eq{"rel_user": f.MemberID})
	}

	return query, nil
}

// perfectly, we would be able to solve this per-rdbms implementation with query directly
func (s Store) getMessagingChannelMembersQuery(cnd squirrel.Sqlizer, memberIDs ...uint64) squirrel.Sqlizer {
	if len(memberIDs) == 0 {
		return squirrel.
			Select("null")
	}

	if strings.HasPrefix(s.config.DriverName, "mysql") {
		// Make sure members are sorted
		sort.Slice(memberIDs, func(i, j int) bool {
			return memberIDs[i] < memberIDs[j]
		})

		// Concatenating members fore
		list := ""
		for i := range memberIDs {
			if i > 0 {
				list += ","
			}
			list += strconv.FormatUint(memberIDs[i], 10)
		}

		return s.SelectBuilder(s.messagingChannelMemberTable("mcm"), "mcm.rel_channel").
			GroupBy("mcm.rel_channel").
			Having(squirrel.Eq{
				"COUNT(*)": len(memberIDs),
				"GROUP_CONCAT(mcm.rel_user ORDER BY 1 ASC SEPARATOR ',')": list,
			})
	}

	var (
		base = s.SelectBuilder(s.messagingChannelMemberTable("mcm"), "mcm.rel_channel").
			PlaceholderFormat(squirrel.Question)

		// construct SQLs with fitting number of members
		counter = base.
			Where(cnd).
			GroupBy("mcm.rel_channel").
			Having(squirrel.Eq{"COUNT(*)": len(memberIDs)})
	)

	for _, memberID := range memberIDs {
		sql, args, _ := base.Where(squirrel.Eq{"mcm.rel_user": memberID}).ToSql()
		counter = counter.Suffix(" INTERSECT "+sql, args...)
	}

	return counter

}
