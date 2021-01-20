package types

type (
	identifiableStep struct{ id uint64 }
)

func (i *identifiableStep) ID() uint64      { return i.id }
func (i *identifiableStep) SetID(id uint64) { i.id = id }
