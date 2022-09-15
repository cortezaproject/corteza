package types

type (
	Label struct {
		// Kind of the labeled resource
		Kind string

		// ID of the labeled resource
		ResourceID uint64

		Name  string
		Value string
	}

	LabelFilter struct {
		Kind       string
		ResourceID []uint64
		Filter     map[string]string
		Limit      uint
	}
)

const (
	LabelResourceType = "corteza::generic:label"
)

func (set LabelSet) ResourceIDs() (rr []uint64) {
	rr = make([]uint64, len(set))
	for r := range set {
		rr[r] = set[r].ResourceID
	}

	return
}

func (set LabelSet) FilterByResource(kind string, ID uint64) map[string]string {
	var kv = make(map[string]string)
	for _, label := range set {
		if kind == label.Kind && ID == label.ResourceID {
			kv[label.Name] = label.Value
		}
	}

	return kv
}
