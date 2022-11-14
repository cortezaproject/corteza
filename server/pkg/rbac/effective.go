package rbac

type (
	effective struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Allow     bool   `json:"allow"`
	}

	EffectiveSet []effective
)

func (ee *EffectiveSet) Push(res, op string, allow bool) {
	*ee = append(*ee, effective{
		Resource:  res,
		Operation: op,
		Allow:     allow,
	})
}
