package sam

// Teams
type Team struct {
	ID        uint64
	Name      string
	MemberIDs []uint64 `json:"-"`
	Members   []User   `json:",omitempty"`
}

func (Team) new() *Team {
	return &Team{}
}

func (t *Team) GetID() uint64 {
	return t.ID
}

func (t *Team) SetID(value uint64) *Team {
	t.ID = value
	return t
}
func (t *Team) GetName() string {
	return t.Name
}

func (t *Team) SetName(value string) *Team {
	t.Name = value
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
