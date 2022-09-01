package dal

type (
	relIndexBuffer struct {
		// attribute statistics to track
		//
		// These stats can be used when performing operations like sorting and
		// binary searching for specific values.
		track []string

		min map[string]any
		max map[string]any

		// @todo make b-tree for sorting? will probably be small so simple slice sort
		// should be ok ig...
		//
		// Probably don't need to sort yet
		rows []*Row
	}
)

// newRelIndexBuffer initializes a new relIndexBuffer tracking the given attributes
func newRelIndexBuffer(tt ...string) *relIndexBuffer {
	return &relIndexBuffer{
		min:   make(map[string]any),
		max:   make(map[string]any),
		track: tt,
	}
}

// add adds a new *row to the buffer
func (lc *relIndexBuffer) add(r *Row) {
	if len(lc.rows) == 0 {
		for _, ix := range lc.track {
			lc.min[ix] = r.values[ix][0]
			lc.max[ix] = r.values[ix][0]
		}

		lc.rows = append(lc.rows, r)
		return
	}

	lc.updMin(r)
	lc.updMax(r)

	lc.rows = append(lc.rows, r)
}

// updMin updates the min stat
func (lc *relIndexBuffer) updMin(r *Row) {
	for _, ix := range lc.track {
		for i := uint(0); i < r.CountValues()[ix]; i++ {
			v := r.values[ix][i]
			if compareValues(v, lc.min[ix]) < 0 {
				lc.min[ix] = v
			}
		}
	}
}

// updMax updates the max stat
func (lc *relIndexBuffer) updMax(r *Row) {
	for _, ix := range lc.track {
		for i := uint(0); i < r.CountValues()[ix]; i++ {
			v := r.values[ix][i]
			if compareValues(v, lc.max[ix]) > 0 {
				lc.max[ix] = v
			}
		}
	}
}
