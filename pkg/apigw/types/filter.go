package types

const (
	PreFilter  FilterKind = "prefilter"
	PostFilter FilterKind = "postfilter"
	Processer  FilterKind = "processer"
)

type (
	FilterKind string

	FilterMeta struct {
		Wgt   int              `json:"-"`
		Name  string           `json:"name"`
		Label string           `json:"label"`
		Kind  FilterKind       `json:"kind"`
		Args  []*FilterMetaArg `json:"params,omitempty"`
	}

	FilterMetaList []*FilterMeta

	FilterMetaArg struct {
		Label   string                 `json:"label"`
		Type    string                 `json:"type"`
		Example string                 `json:"example"`
		Options map[string]interface{} `json:"options"`
	}
)

func (fm FilterMetaList) Filter(f func(*FilterMeta) (bool, error)) (out FilterMetaList, err error) {
	var ok bool
	out = FilterMetaList{}

	for i := range fm {
		if ok, err = f(fm[i]); err != nil {
			return
		} else if ok {
			out = append(out, fm[i])
		}
	}

	return
}
