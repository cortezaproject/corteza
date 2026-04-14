package types

type (
	RoleMember struct {
		RoleID   uint64
		Resource string
	}

	RoleMemberFilter struct {
		RoleID   uint64
		Resource string
		Limit    uint
	}
)
