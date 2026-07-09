package actionlog

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	ActionSet []*Action
)

func (set ActionSet) Walk(w func(*Action) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ActionSet) Filter(f func(*Action) (bool, error)) (out ActionSet, err error) {
	var ok bool
	out = ActionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ActionSet) FindByID(ID uint64) *Action {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ActionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
