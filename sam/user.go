package sam

import (
	"github.com/pkg/errors"
)

func (u *User) Login(r *userLoginRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: User.login")
}
func (u *User) Search(r *userSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: User.search")
}
