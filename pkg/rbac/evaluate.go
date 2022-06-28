package rbac

type (
	Evaluated struct {
		Resource  string      `json:"resource"`
		Operation string      `json:"operation"`
		Access    Access      `json:"-"`
		Can       bool        `json:"can"`
		Step      explanation `json:"step"`

		RoleID uint64 `json:"roleID,string,omitempty"`
		Rule   *Rule  `json:"rule,omitempty"`

		Default *Evaluated `json:"default,omitempty"`
	}

	EvaluatedSet []Evaluated

	explanation string
)

const (
	stepIntegrity explanation = "integrity"
	stepBypass    explanation = "bypass"
	stepRuleless  explanation = "ruleless"
	stepEvaluated explanation = "evaluated"
)
