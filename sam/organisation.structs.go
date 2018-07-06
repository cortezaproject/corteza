package sam

import (
	"time"
)

type (
	// Organisations
	Organisation struct {
		ID         uint64
		Name       string
		ArchivedAt *time.Time `json:",omitempty"`
		DeletedAt  *time.Time `json:",omitempty"`

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
