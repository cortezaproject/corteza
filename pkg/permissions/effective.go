package permissions

type (
	effective struct {
		Resource  Resource  `json:"resource"`
		Operation Operation `json:"operation"`
		Allow     bool      `json:"allow"`
	}

	EffectiveSet []effective
)

func (ee *EffectiveSet) Push(res Resource, op Operation, allow bool) {
	*ee = append(*ee, effective{
		Resource:  res,
		Operation: op,
		Allow:     allow,
	})
}
