package resource

type (
	IndexNode struct {
		children map[string]*IndexNode

		values []interface{}
	}
)

func NewIndex() *IndexNode {
	return &IndexNode{
		children: make(map[string]*IndexNode),
		values:   make([]interface{}, 0),
	}
}

func (in *IndexNode) Add(value interface{}, pp ...[]string) {
	if in == nil {
		*in = *NewIndex()
	}

	if len(pp) == 1 {
		for _, identifier := range pp[0] {
			if identifier == "*" {
				in.values = append(in.values, value)
			} else {
				nn := in.children[identifier]
				if nn == nil {
					nn = NewIndex()
					in.children[identifier] = nn
				}

				nn.values = append(nn.values, value)
			}
		}
		return
	}

	p := pp[0]
	for _, identifier := range p {
		if identifier == "*" {
			in.values = append(in.values, value)
			continue
		}

		if _, ok := in.children[identifier]; !ok {
			in.children[identifier] = NewIndex()
		}

		in.children[identifier].Add(value, pp[1:]...)
	}
}

func (di *IndexNode) Collect(pp ...[]string) (out []interface{}) {
	if len(pp) == 0 {
		if len(di.values) > 0 {
			out = append(out, di.values...)
		}
		return
	}

	p := pp[0]
	for _, identifier := range p {
		skip := false
		if _, ok := di.children[identifier]; ok {
			aux := di.children[identifier].Collect(pp[1:]...)
			out = append(out, aux...)

			skip = len(aux) > 0
		}

		out = append(out, di.values...)
		if skip {
			break
		}
	}

	return
}
