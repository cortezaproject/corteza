package sam

// Organisations
type Organisation struct {
	ID   uint64
	Name string
}

func (Organisation) new() *Organisation {
	return &Organisation{}
}

func (o *Organisation) GetID() uint64 {
	return o.ID
}

func (o *Organisation) SetID(value uint64) *Organisation {
	o.ID = value
	return o
}
func (o *Organisation) GetName() string {
	return o.Name
}

func (o *Organisation) SetName(value string) *Organisation {
	o.Name = value
	return o
}
