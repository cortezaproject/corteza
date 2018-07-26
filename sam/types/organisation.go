package types

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `organisation.go`, `organisation.util.go` or `organisation_test.go` to
	implement your API calls, helper functions and tests. The file `organisation.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"time"
)

type (
	// Organisations - Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.
	Organisation struct {
		ID         uint64     `db:"id"`
		FQN        string     `db:"fqn"`
		Name       string     `db:"name"`
		CreatedAt  time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

		changed []string
	}
)

// New constructs a new instance of Organisation
func (Organisation) New() *Organisation {
	return &Organisation{}
}

// Get the value of ID
func (o *Organisation) GetID() uint64 {
	return o.ID
}

// Set the value of ID
func (o *Organisation) SetID(value uint64) *Organisation {
	if o.ID != value {
		o.changed = append(o.changed, "ID")
		o.ID = value
	}
	return o
}

// Get the value of FQN
func (o *Organisation) GetFQN() string {
	return o.FQN
}

// Set the value of FQN
func (o *Organisation) SetFQN(value string) *Organisation {
	if o.FQN != value {
		o.changed = append(o.changed, "FQN")
		o.FQN = value
	}
	return o
}

// Get the value of Name
func (o *Organisation) GetName() string {
	return o.Name
}

// Set the value of Name
func (o *Organisation) SetName(value string) *Organisation {
	if o.Name != value {
		o.changed = append(o.changed, "Name")
		o.Name = value
	}
	return o
}

// Get the value of CreatedAt
func (o *Organisation) GetCreatedAt() time.Time {
	return o.CreatedAt
}

// Set the value of CreatedAt
func (o *Organisation) SetCreatedAt(value time.Time) *Organisation {
	o.changed = append(o.changed, "CreatedAt")
	o.CreatedAt = value
	return o
}

// Get the value of UpdatedAt
func (o *Organisation) GetUpdatedAt() *time.Time {
	return o.UpdatedAt
}

// Set the value of UpdatedAt
func (o *Organisation) SetUpdatedAt(value *time.Time) *Organisation {
	o.changed = append(o.changed, "UpdatedAt")
	o.UpdatedAt = value
	return o
}

// Get the value of ArchivedAt
func (o *Organisation) GetArchivedAt() *time.Time {
	return o.ArchivedAt
}

// Set the value of ArchivedAt
func (o *Organisation) SetArchivedAt(value *time.Time) *Organisation {
	o.changed = append(o.changed, "ArchivedAt")
	o.ArchivedAt = value
	return o
}

// Get the value of DeletedAt
func (o *Organisation) GetDeletedAt() *time.Time {
	return o.DeletedAt
}

// Set the value of DeletedAt
func (o *Organisation) SetDeletedAt(value *time.Time) *Organisation {
	o.changed = append(o.changed, "DeletedAt")
	o.DeletedAt = value
	return o
}

// Changes returns the names of changed fields
func (o *Organisation) Changes() []string {
	return o.changed
}
