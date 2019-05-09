package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/crusttech/crust/internal/permissions"
	"github.com/pkg/errors"
)

type (
	User struct {
		ID       uint64   `json:"userID,string" db:"id"`
		Username string   `json:"username" db:"username"`
		Email    string   `json:"email" db:"email"`
		Name     string   `json:"name" db:"name"`
		Handle   string   `json:"handle" db:"handle"`
		Kind     UserKind `json:"kind" db:"kind"`

		Meta *UserMeta `json:"meta" db:"meta"`

		OrganisationID uint64 `json:"organisationID,string" db:"rel_organisation"`
		RelatedUserID  uint64 `json:"relatedUserID,string" db:"rel_user_id"`
		User           *User  `json:"user" db:"-"`

		EmailConfirmed bool `json:"-" db:"email_confirmed"`

		CreatedAt   time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		SuspendedAt *time.Time `json:"suspendedAt,omitempty" db:"suspended_at"`
		DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

		Roles []*Role `json:"roles,omitempty" db:"-"`
	}

	UserMeta struct {
		Avatar string `json:"avatar,omitempty"`
	}

	UserFilter struct {
		Query    string
		Email    string
		Username string
		OrderBy  string
	}

	UserKind string
)

const (
	NormalUser UserKind = ""
	BotUser             = "bot"
)

func (u *User) Valid() bool {
	return u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (u *User) Identity() uint64 {
	return u.ID
}

// Resource returns a resource ID for this type
func (u *User) PermissionResource() permissions.Resource {
	return UserPermissionResource.AppendID(u.ID)
}

func (meta *UserMeta) Scan(value interface{}) error {
	switch value.(type) {
	case nil:
		*meta = UserMeta{}
		return nil
	case []uint8:
		if err := json.Unmarshal(value.([]byte), meta); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into User.Meta", value)
		}
		return nil
	}
	return errors.Errorf("User.Meta: unknown type %T, expected []uint8", value)
}

func (meta *UserMeta) Value() (driver.Value, error) {
	return json.Marshal(meta)
}
