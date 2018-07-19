package types

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `team.go`, `team.util.go` or `team_test.go` to
	implement your API calls, helper functions and tests. The file `team.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"time"
)

type (
	// Teams - An organisation may have many teams. Teams may have many channels available. Access to channels may be shared between teams.
	Team struct {
		ID         uint64     `db:"id"`
		Name       string     `db:"name"`
		MemberIDs  []uint64   `json:"-" db:"member_i_ds"`
		Members    []User     `json:",omitempty" db:"members"`
		ArchivedAt *time.Time `json:",omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:",omitempty" db:"deleted_at"`

		changed []string
	}
)

// New constructs a new instance of Team
func (Team) New() *Team {
	return &Team{}
}

// Get the value of ID
func (t *Team) GetID() uint64 {
	return t.ID
}

// Set the value of ID
func (t *Team) SetID(value uint64) *Team {
	if t.ID != value {
		t.changed = append(t.changed, "ID")
		t.ID = value
	}
	return t
}

// Get the value of Name
func (t *Team) GetName() string {
	return t.Name
}

// Set the value of Name
func (t *Team) SetName(value string) *Team {
	if t.Name != value {
		t.changed = append(t.changed, "Name")
		t.Name = value
	}
	return t
}

// Get the value of MemberIDs
func (t *Team) GetMemberIDs() []uint64 {
	return t.MemberIDs
}

// Set the value of MemberIDs
func (t *Team) SetMemberIDs(value []uint64) *Team {
	t.changed = append(t.changed, "MemberIDs")
	t.MemberIDs = value
	return t
}

// Get the value of Members
func (t *Team) GetMembers() []User {
	return t.Members
}

// Set the value of Members
func (t *Team) SetMembers(value []User) *Team {
	t.changed = append(t.changed, "Members")
	t.Members = value
	return t
}

// Get the value of ArchivedAt
func (t *Team) GetArchivedAt() *time.Time {
	return t.ArchivedAt
}

// Set the value of ArchivedAt
func (t *Team) SetArchivedAt(value *time.Time) *Team {
	t.changed = append(t.changed, "ArchivedAt")
	t.ArchivedAt = value
	return t
}

// Get the value of DeletedAt
func (t *Team) GetDeletedAt() *time.Time {
	return t.DeletedAt
}

// Set the value of DeletedAt
func (t *Team) SetDeletedAt(value *time.Time) *Team {
	t.changed = append(t.changed, "DeletedAt")
	t.DeletedAt = value
	return t
}

// Changes returns the names of changed fields
func (t *Team) Changes() []string {
	return t.changed
}
