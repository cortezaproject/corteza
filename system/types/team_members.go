package types

type (
	TeamMember struct {
		TeamID uint64 `db:"rel_team"`
		UserID uint64 `db:"rel_user"`
	}

	TeamMemberFilter struct {
		Query string
	}
)
