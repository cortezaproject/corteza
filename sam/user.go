package sam

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

var _ = errors.Wrap

const (
	sqlUserScope  = "suspended_at IS NULL AND archived_at IS NULL"
	sqlUserSelect = "SELECT * FROM users WHERE " + sqlUserScope
)

func (*User) Read(r *teamReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	t := User{}.new()
	return t, db.Get(t, sqlUserSelect+" AND id = ?", r.id)
}

// User lookup & login
func (*User) Login(r *userLoginRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	u := &User{}
	if err != db.Get(u, sqlUserSelect+" AND username = ?", r.username) {
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

	if err != db.Get(uu, sqlUserSelect+" AND username LIKE ?", r.query+"%") {
		return nil, err
	}

	return uu, nil
}

func (*User) Remove(r *teamRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := "UPDATE users SET deleted_at = NOW() WHERE deleted_at IS NULL AND id = ?"

	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
		return err
	}()
}

func (*User) Archive(r *teamArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(
		"UPDATE users SET archived_at = NOW() WHERE %s AND id = ?",
		sqlUserScope)

	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
		return err
	}()
}
