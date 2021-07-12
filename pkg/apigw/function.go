package apigw

const (
	FunctionKindVerifier  FunctionKind = "verifier"
	FunctionKindValidator FunctionKind = "validator"
	FunctionKindProcesser FunctionKind = "processer"
	FunctionKindExpediter FunctionKind = "expediter"
)

type (
	FunctionKind string

	Handler interface {
		Execer
		Stringer

		Merge([]byte) (Handler, error)
		Meta() functionMeta
	}

	functionMeta struct {
		Step   int                `json:"step"`
		Weight int                `json:"-"`
		Name   string             `json:"name"`
		Label  string             `json:"label"`
		Kind   FunctionKind       `json:"kind"`
		Args   []*functionMetaArg `json:"params,omitempty"`
	}

	functionMetaList []*functionMeta

	functionMetaArg struct {
		Label   string                 `json:"label"`
		Type    string                 `json:"type"`
		Example string                 `json:"example"`
		Options map[string]interface{} `json:"options"`
	}
)

// func (ff functionHandler) Weight() int {
// 	// if there's gonna be more than 1000 funcs
// 	// per step, we're doing something wrong
// 	return ff.step*1000 + ff.weight
// }

func (fm functionMetaList) Filter(f func(*functionMeta) (bool, error)) (out functionMetaList, err error) {
	var ok bool
	out = functionMetaList{}

	for i := range fm {
		if ok, err = f(fm[i]); err != nil {
			return
		} else if ok {
			out = append(out, fm[i])
		}
	}

	return
}
