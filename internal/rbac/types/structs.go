package types

type (
	User struct {
		UserID          string   `json:"userid"`
		Username        string   `json:"username"`
		AssignedRoles   []string `json:"assignedRoles"`
		AuthorizedRoles []string `json:"authorizedRoles"`
	}

	Session struct {
		ID       string   `json:"session"`
		Username string   `json:"username"`
		Roles    []string `json:"roles"`
	}

	// @todo: need to list nested roles,
	// @todo: don't return users=null - return users: []?
	Role struct {
		Name  string   `json:"rolename"`
		Users []string `json:"users"`

		// key = resource name
		Permissions map[string]Operations `json:"permissions"`
	}

	Operations struct {
		Operations []string `json:"operations"`
	}

	// @todo: read resource information
	Resource struct {
	}
)
