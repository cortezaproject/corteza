package auth

type (
	Identifiable interface {
		Identity() uint64
		Valid() bool
	}

	Identity struct {
		identity uint64
	}
)

func NewIdentity(id uint64) *Identity {
	return &Identity{
		identity: id,
	}
}

func (i Identity) Identity() uint64 {
	return i.identity
}

func (i Identity) Valid() bool {
	return i.identity > 0
}
