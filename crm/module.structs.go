package crm

// Modules
type Module struct {
	id   uint64
	name string

	changed []string
}

func (Module) new() *Module {
	return &Module{}
}

func (m *Module) Getid() uint64 {
	return m.id
}

func (m *Module) Setid(value uint64) *Module {
	if m.id != value {
		m.changed = append(m.changed, "id")
		m.id = value
	}
	return m
}
func (m *Module) Getname() string {
	return m.name
}

func (m *Module) Setname(value string) *Module {
	if m.name != value {
		m.changed = append(m.changed, "name")
		m.name = value
	}
	return m
}
