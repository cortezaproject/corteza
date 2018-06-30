package crm

// Types
type Types struct {
	ID string
}

func (Types) new() *Types {
	return &Types{}
}

func (t *Types) GetID() string {
	return t.ID
}

func (t *Types) SetID(value string) *Types {
	t.ID = value
	return t
}
