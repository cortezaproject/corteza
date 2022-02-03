package conv

import "encoding/json"

type (
	Gig struct{}

	ParamWrap struct {
		Ref    string
		Params map[string]interface{}
	}
	ParamWrapSet []ParamWrap
)

func ParseParamWrap(ss []string) (out ParamWrapSet, err error) {
	for _, s := range ss {
		aux := make(ParamWrapSet, 0, 2)
		err = json.Unmarshal([]byte(s), &aux)
		if err != nil {
			return
		}

		out = append(out, aux...)
	}
	return
}
