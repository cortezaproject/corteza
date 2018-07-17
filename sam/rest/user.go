package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
)

var _ = errors.Wrap

const (
	sqlUserScope  = "suspended_at IS NULL AND archived_at IS NULL"
	sqlUserSelect = "SELECT * FROM users WHERE " + sqlUserScope
)

type (
	User struct {
		service UserInterface
	}

	UserInterface interface {
		Set(*types.User)
		CanLogin() bool
		GeneratePassword(value string) ([]byte, error)
		ValidatePassword(value string) bool
		ValidateUser() bool
	}
)

func (User) New() *User {
	return &User{service.User{}.New()}
}

/*
func (self *User) Read(ctx context.Context, r *teamReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	t := types.User{}.New()
	return t, db.Get(t, sqlUserSelect+" AND id = ?", r.ID)
}
*/

// User lookup & login
func (self *User) Login(ctx context.Context, r *server.UserLoginRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	u := types.User{}.New()
	if err = db.Get(u, sqlUserSelect+" AND username = ?", r.Username); err != nil {
		return nil, err
	}
	self.service.Set(u)

	if self.service.ValidateUser() || !self.service.ValidatePassword(r.Password) {
		return nil, errors.New("Invalid username and password combination")
	}

	if !self.service.CanLogin() {
		return nil, errors.New("User is not allowed to login")
	}

	return u, nil
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (*User) Search(ctx context.Context, r *server.UserSearchRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	uu := []types.User{}
	if err != db.Get(uu, sqlUserSelect+" AND username LIKE ?", r.Query+"%") {
		return nil, err
	}

	return uu, nil
}

/*
func (self *User) Remove(ctx context.Context, r *teamRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := "UPDATE users SET deleted_at = NOW() WHERE deleted_at IS NULL AND id = ?"

	return nil, func() error {
		_, err := db.Exec(stmt, r.ID)
		return err
	}()
}

func (self *User) Archive(ctx context.Context, r *teamArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(
		"UPDATE users SET archived_at = NOW() WHERE %s AND id = ?",
		sqlUserScope)

	return nil, func() error {
		_, err := db.Exec(stmt, r.ID)
		return err
	}()
}
*/
