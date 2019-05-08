package outgoing

type (
	User struct {
		// Channel to part (nil) for ALL channels
		ID       uint64 `json:"userID,string"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Handle   string `json:"handle"`
	}

	UserSet []*User
)
