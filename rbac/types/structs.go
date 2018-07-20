package types

type (
	User struct {
		Username string `json:"username"`
	}

	Session struct {
		Username string `json:"username"`
		Roles []string `json:"roles"`
	}
)
