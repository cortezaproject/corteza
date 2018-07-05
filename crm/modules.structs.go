package crm

// Modules
type Modules struct {
	id   uint64
	name string

	changed []string
}

func (Modules) new() *Modules {
	return &Modules{}
}

func (m *Modules) Getid() uint64 {
	return m.id
}

func (m *Modules) Setid(value uint64) *Modules {
	if m.id != value {
		m.changed = append(m.changed, "id")
		m.id = value
	}
	return m
}
func (m *Modules) Getname() string {
	return m.name
}

func (m *Modules) Setname(value string) *Modules {
	if m.name != value {
		m.changed = append(m.changed, "name")
		m.name = value
	}
	return m
}
