package rules

type (
	scope struct {
		providers []ScopeItem
	}

	ScopeItem struct {
		Scope       string           `json:"scope"`
		Permissions []OperationGroup `json:"permissions"`
	}

	ScopeProvider interface {
		Resource() Resource
		Permissions() []OperationGroup
	}

	ScopeInterface interface {
		Add(ScopeProvider)
		List() []ScopeItem
	}
)

func NewScope() ScopeInterface {
	return &scope{}
}

func (s *scope) Add(p ScopeProvider) {
	s.providers = append(s.providers, ScopeItem{
		Scope:       p.Resource().Scope,
		Permissions: p.Permissions(),
	})
}

func (s *scope) List() []ScopeItem {
	return s.providers
}
