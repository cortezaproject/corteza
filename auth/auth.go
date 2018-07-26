package auth

type (
	Identifiable interface {
		GetID() uint64
		Valid() bool
	}

	Identity struct {
		id uint64
	}
)

func NewIdentity(id uint64) *Identity {
	return &Identity{
		id: id,
	}
}

func (i Identity) GetID() uint64 {
	return i.id
}

func (i Identity) Valid() bool {
	return i.id > 0
}
