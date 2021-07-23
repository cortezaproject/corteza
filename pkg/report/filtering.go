package report

type (
	RowDefinition struct {
		And   []*RowDefinition           `json:"and"`
		Or    []*RowDefinition           `json:"or"`
		Cells map[string]*CellDefinition `json:"cells"`
	}

	CellDefinition struct {
		Op    string `json:"op"`
		Value string `json:"value"`
	}
)

func (base *RowDefinition) MergeAnd(merge *RowDefinition) *RowDefinition {
	// 1. merge the two
	rr := &RowDefinition{
		And: make([]*RowDefinition, 0, 2),
	}
	if base != nil {
		rr.And = append(rr.And, base)
	}
	if merge != nil {
		rr.And = append(rr.And, merge)
	}

	// 2. flatten the tree
	// @todo do some more in-depth processing
	if len(rr.And) == 1 {
		return rr.And[0]
	}
	if len(rr.And)+len(rr.Cells)+len(rr.Or) == 0 {
		return nil
	}

	return rr
}

func (base *RowDefinition) MergeOr(merge *RowDefinition) *RowDefinition {
	// 1. merge the two
	rr := &RowDefinition{
		Or: make([]*RowDefinition, 0, 2),
	}
	if base != nil {
		rr.Or = append(rr.Or, base)
	}
	if merge != nil {
		rr.Or = append(rr.Or, merge)
	}

	// 2. flatten the tree
	// @todo do some more in-depth processing
	if len(rr.Or) == 1 {
		return rr.Or[0]
	}
	if len(rr.And)+len(rr.Cells)+len(rr.Or) == 0 {
		return nil
	}

	return rr
}

func (base *RowDefinition) Clone() (out *RowDefinition) {
	if base == nil {
		return
	}
	out = &RowDefinition{}

	if base.Cells != nil {
		out.Cells = make(map[string]*CellDefinition)
		for k, v := range base.Cells {
			out.Cells[k] = &CellDefinition{
				Op:    v.Op,
				Value: v.Value,
			}
		}
	}

	if base.And != nil {
		out.And = make([]*RowDefinition, len(base.And))
		for i, def := range base.And {
			out.And[i] = def.Clone()
		}
	}
	if base.Or != nil {
		out.Or = make([]*RowDefinition, len(base.Or))
		for i, def := range base.Or {
			out.Or[i] = def.Clone()
		}
	}

	return
}
