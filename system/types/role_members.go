package types

type (
	RoleMember struct {
		RoleID uint64 `db:"rel_role"`
		UserID uint64 `db:"rel_user"`
	}

	RoleMemberFilter struct {
		Query string
	}
)
