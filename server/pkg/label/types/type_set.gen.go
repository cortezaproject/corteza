package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	LabelSet []*Label
)

func (set LabelSet) Walk(w func(*Label) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set LabelSet) Filter(f func(*Label) (bool, error)) (out LabelSet, err error) {
	var ok bool
	out = LabelSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
