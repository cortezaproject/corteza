package dal

type (
	// relIndex is a generic struct for indexing data which join/link can use
	//
	// The current index implementation utilizes a series of hashmaps based on what
	// type of values we're indexing on.
	// As we currently only support single value predicate, the hashmaps proved
	// to be a bit more efficient then b-trees (other consideration) when testing
	// on larger datasets.
	//
	// @todo when we support multiple join predicates, this should probably change
	// into a b-tree as it might be a bit faster thennested hashmaps.
	//
	// @todo do some benchmarks in regards to using generics for key
	relIndex struct {
		track []string

		ints    map[int64]*relIndexBuffer
		strings map[string]*relIndexBuffer
		ids     map[uint64]*relIndexBuffer
	}
)

// newRelIndex initializes a new relIndex with the specified tracked attributes
// @todo benchmark with generics as key
func newRelIndex(tt ...string) *relIndex {
	return &relIndex{
		track:   tt,
		ints:    make(map[int64]*relIndexBuffer),
		strings: make(map[string]*relIndexBuffer),
		ids:     make(map[uint64]*relIndexBuffer),
	}
}

// AddInt adds a new row under the int key
func (ri *relIndex) AddInt(k int64, r *Row) {
	c, ok := ri.GetInt(k)
	if !ok {
		c = newRelIndexBuffer(ri.track...)
		ri.ints[k] = c
	}
	c.add(r)
}

func (ri *relIndex) GetInt(k int64) (out *relIndexBuffer, ok bool) {
	out, ok = ri.ints[k]
	return
}

func (ri *relIndex) AddString(k string, r *Row) {
	c, ok := ri.GetString(k)
	if !ok {
		c = newRelIndexBuffer(ri.track...)
		ri.strings[k] = c
	}
	c.add(r)
}

func (ri *relIndex) GetString(k string) (out *relIndexBuffer, ok bool) {
	out, ok = ri.strings[k]
	return
}

func (ri *relIndex) AddID(k uint64, r *Row) {
	c, ok := ri.GetID(k)
	if !ok {
		c = newRelIndexBuffer(ri.track...)
		ri.ids[k] = c
	}
	c.add(r)
}

func (ri *relIndex) GetID(k uint64) (out *relIndexBuffer, ok bool) {
	out, ok = ri.ids[k]
	return
}
