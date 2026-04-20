package dal

import "sort"

type (
	SensitivityLevel struct {
		Handle string
		ID     uint64
		Level  int
	}
	SensitivityLevelSet []SensitivityLevel

	sensitivityLevelIndex struct {
		set SensitivityLevelSet

		byHandle map[string]int
		byID     map[uint64]int
	}

	SensitivityLevelUsage struct {
		connections []map[string]any
		modules     []map[string]any
		fields      []map[string]any
	}
)

func SensitivityLevelIndex(levels ...SensitivityLevel) *sensitivityLevelIndex {
	out := &sensitivityLevelIndex{
		set:      make(SensitivityLevelSet, len(levels)),
		byHandle: make(map[string]int),
		byID:     make(map[uint64]int),
	}

	for i, l := range levels {
		out.set[i] = l
		out.byHandle[l.Handle] = i
		out.byID[l.ID] = i
	}

	return out
}

func (sli sensitivityLevelIndex) with(levels ...SensitivityLevel) *sensitivityLevelIndex {
	slvls := append(sli.set, levels...)
	sort.Sort(slvls)

	return SensitivityLevelIndex(slvls...)
}

func (sli sensitivityLevelIndex) without(levels ...SensitivityLevel) *sensitivityLevelIndex {
	nn := make(SensitivityLevelSet, 0, len(sli.set))

	remIndex := SensitivityLevelIndex(levels...)

	for _, existing := range sli.set {
		if !remIndex.includes(existing.ID) {
			nn = append(nn, existing)
		}
	}

	sort.Sort(nn)

	return SensitivityLevelIndex(nn...)
}

func (sli sensitivityLevelIndex) includes(l uint64) (ok bool) {
	if l == 0 {
		return true
	}

	if sli.byID == nil {
		return false
	}

	_, ok = sli.byID[l]
	return
}

func (sli sensitivityLevelIndex) isSubset(a, b uint64) (ok bool) {
	// Edgecases
	// If A is zero theneverything is possible
	if a == 0 {
		return true
	}
	// If B is zero, then A must also be zero
	if b == 0 {
		return a == 0
	}

	var lvlA, lvlB int

	if lvlA, ok = sli.byID[a]; !ok {
		return false
	}

	if lvlB, ok = sli.byID[b]; !ok {
		return false
	}

	return lvlA <= lvlB
}

func (a SensitivityLevelSet) Len() int           { return len(a) }
func (a SensitivityLevelSet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SensitivityLevelSet) Less(i, j int) bool { return a[i].Level < a[j].Level }

func (u SensitivityLevelUsage) Empty() bool {
	return len(u.connections)+len(u.fields)+len(u.modules) == 0
}
