package slice

type (
	Circular[V any] struct {
		size   int
		slc    []V
		head   int
		cycled bool
	}
)

func NewCircular[V any](size int) *Circular[V] {
	return &Circular[V]{
		slc:  make([]V, size),
		size: size,
	}
}

func (cs *Circular[V]) Add(v V) {
	if cs.head >= cs.size {
		cs.head = 0
		cs.cycled = true
	}

	cs.slc[cs.head] = v
	cs.head++
}

func (cs *Circular[V]) Slice() (out []V) {
	if cs == nil {
		return
	}

	if !cs.cycled {
		return cs.sliceUncycled()
	}
	return cs.sliceCycled()
}

func (cs *Circular[V]) sliceUncycled() (out []V) {
	for i := 0; i < cs.head; i++ {
		out = append(out, cs.slc[i])
	}

	return
}

func (cs *Circular[V]) sliceCycled() (out []V) {
	for i := 0; i < cs.size; i++ {
		xo := (cs.head + i) % cs.size
		out = append(out, cs.slc[xo])
	}

	return
}
