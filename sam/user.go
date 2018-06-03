package sam

import (
	"github.com/pkg/errors"
)

func (u *User) Login(r *UserLoginRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: User.login")
}
func (u *User) Search(r *UserSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: User.search")
}
