package sam

import (
	"time"
)

// Organisations
type Organisation struct {
	ID   uint64
	Name string

	ArchivedAt *time.Time
	DeletedAt  *time.Time

	changed []string
}

func (Organisation) new() *Organisation {
	return &Organisation{}
}

func (o *Organisation) GetID() uint64 {
	return o.ID
}

func (o *Organisation) SetID(value uint64) *Organisation {
	if o.ID != value {
		o.changed = append(o.changed, "id")
		o.ID = value
	}
	return o
}
func (o *Organisation) GetName() string {
	return o.Name
}

func (o *Organisation) SetName(value string) *Organisation {
	if o.Name != value {
		o.changed = append(o.changed, "name")
		o.Name = value
	}
	return o
}
