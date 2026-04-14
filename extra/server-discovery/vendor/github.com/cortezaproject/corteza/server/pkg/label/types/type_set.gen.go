package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/label/types/types.yaml

type (

	// LabelSet slice of Label
	//
	// This type is auto-generated.
	LabelSet []*Label
)

// Walk iterates through every slice item and calls w(Label) err
//
// This function is auto-generated.
func (set LabelSet) Walk(w func(*Label) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Label) (bool, err) and return filtered slice
//
// This function is auto-generated.
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
