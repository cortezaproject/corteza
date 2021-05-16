package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/pkg/errors"
)

type (
	User struct {
		ID       uint64   `json:"userID,string"`
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Name     string   `json:"name"`
		Handle   string   `json:"handle"`
		Kind     UserKind `json:"kind"`

		Meta *UserMeta `json:"meta"`

		EmailConfirmed bool `json:"emailConfirmed"`

		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt   time.Time  `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		SuspendedAt *time.Time `json:"suspendedAt,omitempty"`
		DeletedAt   *time.Time `json:"deletedAt,omitempty"`

		// Hold list of roles this user is member of.
		// we're using this for auth/identifier purposes, to support Roles() func
		// that satisfies Identifiable interface
		roles []uint64
	}

	UserMeta struct {
		// @todo remove, obsolete.
		Avatar string `json:"avatar,omitempty"`

		// User's security policy settings
		SecurityPolicy struct {
			// settings for multi-factor authentication
			MFA struct {
				// Enforce OTP on login
				EnforcedEmailOTP bool `json:"enforcedEmailOTP"`

				// Require OTP to be entered every time client is authorized
				//StrictEmailOTP bool `json:"strictEmailOTP"`

				// Is TOTP configured & enforced?
				EnforcedTOTP bool `json:"enforcedTOTP"`

				// Require OTP to be entered every time client is authorized
				//StrictTOTP bool `json:"strictTOTP"`
			} `json:"mfa"`
		} `json:"securityPolicy"`
	}

	UserFilter struct {
		UserID   []uint64 `json:"userID,string"`
		RoleID   []uint64 `json:"roleID,string"`
		Query    string   `json:"query"`
		Email    string   `json:"email"`
		Username string   `json:"username"`
		Handle   string   `json:"handle"`
		Kind     UserKind `json:"kind"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted   filter.State `json:"deleted"`
		Suspended filter.State `json:"suspended"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*User) (bool, error) `json:"-"`

		MaskedEmailsEnabled bool `json:"-"`
		MaskedNamesEnabled  bool `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	UserKind string

	UserMetrics struct {
		Total          uint   `json:"total"`
		Valid          uint   `json:"valid"`
		Deleted        uint   `json:"deleted"`
		Suspended      uint   `json:"suspended"`
		DailyCreated   []uint `json:"dailyCreated"`
		DailyDeleted   []uint `json:"dailyDeleted"`
		DailyUpdated   []uint `json:"dailyUpdated"`
		DailySuspended []uint `json:"dailySuspended"`
	}
)

const (
	NormalUser UserKind = ""
	BotUser    UserKind = "bot"
)

func (u User) String() string {
	return fmt.Sprintf("%d", u.ID)
}

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

func (meta *UserMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*meta = UserMeta{}
		return nil
	case []uint8:
		if err := json.Unmarshal(value.([]byte), meta); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into User.Meta", value)
		}
		return nil
	}
	return errors.Errorf("User.Meta: unknown type %T, expected []uint8", value)
}

func (meta *UserMeta) Value() (driver.Value, error) {
	return json.Marshal(meta)
}
