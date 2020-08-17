package types

type (
	RoleMember struct {
		RoleID uint64
		UserID uint64
	}

	RoleMemberFilter struct {
		RoleID uint64
		UserID uint64
	}
)
