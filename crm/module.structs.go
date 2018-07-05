package crm

type (
	// Modules
	Module struct {
		ID   uint64
		Name string

		changed []string
	}
)

/* Constructors */
func (Module) new() *Module {
	return &Module{}
}

/* Getters/setters */
func (m *Module) GetID() uint64 {
	return m.ID
}

func (m *Module) SetID(value uint64) *Module {
	if m.ID != value {
		m.changed = append(m.changed, "id")
		m.ID = value
	}
	return m
}
func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) SetName(value string) *Module {
	if m.Name != value {
		m.changed = append(m.changed, "name")
		m.Name = value
	}
	return m
}
