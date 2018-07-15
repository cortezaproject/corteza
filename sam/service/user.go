package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/crusttech/crust/sam/types"
)

type User struct {
	*types.User
}

func (User) New() *User {
	return &User{types.User{}.New()}
}

func (u *User) Set(user *types.User) {
	u.User = user
}

func (u *User) CanLogin() bool {
	return u.User != nil && u.User.ID > 0
}

func (u *User) GeneratePassword(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
}

func (u *User) ValidatePassword(value string) bool {
	return u.User != nil && bcrypt.CompareHashAndPassword(u.User.Password, []byte(value)) == nil
}

func (u *User) ValidateUser() bool {
	return u.User != nil && u.User.ID > 0
}
