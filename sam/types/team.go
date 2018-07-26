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
		Handle     string     `db:"handle"`
		CreatedAt  time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

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

// Get the value of Handle
func (t *Team) GetHandle() string {
	return t.Handle
}

// Set the value of Handle
func (t *Team) SetHandle(value string) *Team {
	if t.Handle != value {
		t.changed = append(t.changed, "Handle")
		t.Handle = value
	}
	return t
}

// Get the value of CreatedAt
func (t *Team) GetCreatedAt() time.Time {
	return t.CreatedAt
}

// Set the value of CreatedAt
func (t *Team) SetCreatedAt(value time.Time) *Team {
	t.changed = append(t.changed, "CreatedAt")
	t.CreatedAt = value
	return t
}

// Get the value of UpdatedAt
func (t *Team) GetUpdatedAt() *time.Time {
	return t.UpdatedAt
}

// Set the value of UpdatedAt
func (t *Team) SetUpdatedAt(value *time.Time) *Team {
	t.changed = append(t.changed, "UpdatedAt")
	t.UpdatedAt = value
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
