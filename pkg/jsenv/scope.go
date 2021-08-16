package jsenv

type (
	globalScope map[string]interface{}
)

func (gs globalScope) Set(k string, i interface{}) {
	gs[k] = i
}

func (gs globalScope) Get(k string) interface{} {
	if v, ok := gs[k]; ok {
		return v
	}

	return nil
}
