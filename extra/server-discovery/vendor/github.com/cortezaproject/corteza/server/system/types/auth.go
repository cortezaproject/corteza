package types

import "github.com/markbates/goth"

type (
	AuthProvider struct {
		Provider string
	}

	ExternalAuthUser struct {
		goth.User
	}
)
