package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
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
		// User's profile avatar photo attachment ID
		AvatarID   uint64 `json:"avatarID,string"`
		AvatarKind string `json:"avatarKind,omitempty"`

		// User's avatar initial text and background color
		AvatarColor   string `json:"avatarColor,omitempty"`
		AvatarBgColor string `json:"avatarBgColor,omitempty"`

		PreferredLanguage string `json:"preferredLanguage"`

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
		UserID   []string `json:"userID"`
		RoleID   []string `json:"roleID"`
		Query    string   `json:"query"`
		Email    string   `json:"email"`
		Username string   `json:"username"`
		Handle   string   `json:"handle"`
		Kind     UserKind `json:"kind"`

		// Set to true if you want to get all kinds/types of users
		AllKinds bool `json:"anyKind"`

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
	SystemUser UserKind = "sys"
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

func (u *User) SetRoles(rr ...uint64) {
	u.roles = rr
}

func (u *User) Clone() *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		Name:           u.Name,
		Handle:         u.Handle,
		Kind:           u.Kind,
		Meta:           u.Meta,
		EmailConfirmed: u.EmailConfirmed,
		Labels:         u.Labels,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		SuspendedAt:    u.SuspendedAt,
		DeletedAt:      u.DeletedAt,
		roles:          u.roles,
	}
}

func (meta *UserMeta) Scan(src any) error           { return sql.ParseJSON(src, meta) }
func (meta *UserMeta) Value() (driver.Value, error) { return json.Marshal(meta) }
