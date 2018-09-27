package types

import (
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID             uint64          `json:"id" db:"id"`
		Username       string          `json:"username" db:"username"`
		Email          string          `json:"email" db:"email"`
		Name           string          `json:"name" db:"name"`
		Handle         string          `json:"handle" db:"handle"`
		SatosaID       string          `json:"satosaId" db:"satosa_id"`
		Meta           json.RawMessage `json:"-" db:"meta"`
		OrganisationID uint64          `json:"organisationId" db:"rel_organisation"`
		Password       []byte          `json:"-" db:"password"`
		CreatedAt      time.Time       `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt      *time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
		SuspendedAt    *time.Time      `json:"suspendedAt,omitempty" db:"suspended_at"`
		DeletedAt      *time.Time      `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	UserFilter struct {
		Query string
	}

	UserSet []*User
)

func (u *User) Valid() bool {
	return u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (u *User) Identity() uint64 {
	return u.ID
}

func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password)) == nil
}

func (u *User) GeneratePassword(password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = pwd
	return nil
}

func (uu UserSet) Walk(w func(*User) error) (err error) {
	for i := range uu {
		if err = w(uu[i]); err != nil {
			return
		}
	}

	return
}
