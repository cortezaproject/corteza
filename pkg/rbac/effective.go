package rbac

type (
	effective struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Allow     bool   `json:"allow"`
	}

	EffectiveSet []effective

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

func (ee *EffectiveSet) Push(res, op string, allow bool) {
	*ee = append(*ee, effective{
		Resource:  res,
		Operation: op,
		Allow:     allow,
	})
}
