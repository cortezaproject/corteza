package auth

type (
	Identity struct {
		id       uint64
		memberOf []uint64
	}
)

func NewIdentity(id uint64, rr ...uint64) *Identity {
	return &Identity{
		id:       id,
		memberOf: rr,
	}
}

func (i Identity) Identity() uint64 {
	return i.id
}

func (i Identity) Roles() []uint64 {
	return i.memberOf
}

func (i Identity) Valid() bool {
	return i.id > 0
}
