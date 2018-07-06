package sam

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
