package sam

import (
	"github.com/pkg/errors"
)

func (*User) Login(r *userLoginRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: User.login")
}
func (*User) Search(r *userSearchRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: User.search")
}
