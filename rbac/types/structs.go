package types

type (
	User struct {
		Username string `json:"username"`
		AssignedRoles []string `json:"assignedRoles"`
		AuthorizedRoles []string `json:"authorizedRoles"`
	}

	Session struct {
		ID string `json:"session"`
		Username string `json:"username"`
		Roles []string `json:"roles"`
	}
)
