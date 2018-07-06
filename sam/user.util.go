package sam

import (
	"golang.org/x/crypto/bcrypt"
)

func (u *User) CanLogin() bool {
	return u.ID > 0
}

func (u *User) GeneratePassword(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
}

func (u User) ValidatePassword(value string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(value)) == nil
}
