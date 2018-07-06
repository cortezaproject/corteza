package sam

import (
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

var _ = errors.Wrap

// User lookup & login
func (*User) Login(r *userLoginRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	u := &User{}
	if err != db.Get(u, "SELECT * FROM users WHERE username = ?", r.username) {
		return nil, err
	}

	if u.ID == 0 || !u.ValidatePassword(r.password) {
		return nil, errors.New("Invalid username and password combination")
	}

	if !u.CanLogin() {
		return nil, errors.New("User is not allowed to login")
	}

	return u, nil
}

// Searches the users table in the database to find users by matching (by-prefix) their username
func (*User) Search(r *userSearchRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	uu := []*User{}

	if err != db.Get(uu, "SELECT * FROM users WHERE username LIKE ?", r.query+"%") {
		return nil, err
	}

	return uu, nil
}
