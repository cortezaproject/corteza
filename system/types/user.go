package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/internal/permissions"
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

		// Hold list of roles this user is member of.
		// we're using this for auth/identifier purposes, to support Roles() func
		// that satisfies Identifiable interface
		roles []uint64
	}

	UserMeta struct {
		Avatar string `json:"avatar,omitempty"`
	}

	UserFilter struct {
		Query        string   `json:"query"`
		Email        string   `json:"email"`
		Username     string   `json:"username"`
		Kind         UserKind `json:"kind"`
		IncDeleted   bool     `json:"incDeleted"`
		IncSuspended bool     `json:"incSuspended"`

		Page    uint   `json:"page"`
		PerPage uint   `json:"perPage"`
		Sort    string `json:"sort"`
		Count   uint   `json:"count"`
	}

	UserKind string
)

const (
	NormalUser UserKind = ""
	BotUser    UserKind = "bot"
)

func (u *User) Valid() bool {
	return u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (u User) Identity() uint64 {
	return u.ID
}

func (u User) Roles() []uint64 {
	return u.roles
}

func (u *User) SetRoles(rr []uint64) {
	u.roles = rr
}

// Resource returns a resource ID for this type
func (u *User) PermissionResource() permissions.Resource {
	return UserPermissionResource.AppendID(u.ID)
}

func (meta *UserMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
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
