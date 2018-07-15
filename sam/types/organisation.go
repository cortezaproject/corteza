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
	// Organisations
	Organisation struct {
		ID         uint64     `db:"id"`
		Name       string     `db:"name"`
		ArchivedAt *time.Time `json:",omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:",omitempty" db:"deleted_at"`

		changed []string
	}
)

/* Constructors */
func (Organisation) New() *Organisation {
	return &Organisation{}
}

/* Getters/setters */
func (o *Organisation) GetID() uint64 {
	return o.ID
}

func (o *Organisation) SetID(value uint64) *Organisation {
	if o.ID != value {
		o.changed = append(o.changed, "ID")
		o.ID = value
	}
	return o
}
func (o *Organisation) GetName() string {
	return o.Name
}

func (o *Organisation) SetName(value string) *Organisation {
	if o.Name != value {
		o.changed = append(o.changed, "Name")
		o.Name = value
	}
	return o
}
func (o *Organisation) GetArchivedAt() *time.Time {
	return o.ArchivedAt
}

func (o *Organisation) SetArchivedAt(value *time.Time) *Organisation {
	o.ArchivedAt = value
	return o
}
func (o *Organisation) GetDeletedAt() *time.Time {
	return o.DeletedAt
}

func (o *Organisation) SetDeletedAt(value *time.Time) *Organisation {
	o.DeletedAt = value
	return o
}
