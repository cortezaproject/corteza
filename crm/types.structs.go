package crm

// Types
type Types struct {
	changed []string
}

func (Types) new() *Types {
	return &Types{}
}
