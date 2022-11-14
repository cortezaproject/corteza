package eventbus

// Debug returns structured debug data
func (b *eventbus) Debug() interface{} {
	var (
		o = make([]*handler, 0, len(b.handlers))
	)

	for _, h := range b.handlers {
		o = append(o, h)
	}

	return o
}
