package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	ResourceActivitySet []*ResourceActivity
)

func (set ResourceActivitySet) Walk(w func(*ResourceActivity) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ResourceActivitySet) Filter(f func(*ResourceActivity) (bool, error)) (out ResourceActivitySet, err error) {
	var ok bool
	out = ResourceActivitySet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ResourceActivitySet) FindByID(ID uint64) *ResourceActivity {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ResourceActivitySet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
