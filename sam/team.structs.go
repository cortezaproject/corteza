package sam

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
	// Teams
	Team struct {
		ID         uint64
		Name       string
		MemberIDs  []uint64   `json:"-"`
		Members    []User     `json:",omitempty"`
		ArchivedAt *time.Time `json:",omitempty"`
		DeletedAt  *time.Time `json:",omitempty"`

		changed []string
	}
)

/* Constructors */
func (Team) New() *Team {
	return &Team{}
}

/* Getters/setters */
func (t *Team) GetID() uint64 {
	return t.ID
}

func (t *Team) SetID(value uint64) *Team {
	if t.ID != value {
		t.changed = append(t.changed, "ID")
		t.ID = value
	}
	return t
}
func (t *Team) GetName() string {
	return t.Name
}

func (t *Team) SetName(value string) *Team {
	if t.Name != value {
		t.changed = append(t.changed, "Name")
		t.Name = value
	}
	return t
}
func (t *Team) GetMemberIDs() []uint64 {
	return t.MemberIDs
}

func (t *Team) SetMemberIDs(value []uint64) *Team {
	t.MemberIDs = value
	return t
}
func (t *Team) GetMembers() []User {
	return t.Members
}

func (t *Team) SetMembers(value []User) *Team {
	t.Members = value
	return t
}
func (t *Team) GetArchivedAt() *time.Time {
	return t.ArchivedAt
}

func (t *Team) SetArchivedAt(value *time.Time) *Team {
	t.ArchivedAt = value
	return t
}
func (t *Team) GetDeletedAt() *time.Time {
	return t.DeletedAt
}

func (t *Team) SetDeletedAt(value *time.Time) *Team {
	t.DeletedAt = value
	return t
}
