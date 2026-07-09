package corredor

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	ScriptSet []*Script
)

func (set ScriptSet) Walk(w func(*Script) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ScriptSet) Filter(f func(*Script) (bool, error)) (out ScriptSet, err error) {
	var ok bool
	out = ScriptSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
