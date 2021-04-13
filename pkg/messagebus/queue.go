package messagebus

type (
	Queue struct {
		handler Handler

		// main data channels
		in  chan []byte
		out chan []byte
		err chan error

		// notification channel for eventbus
		dispatch chan []byte

		// processed messages in-memory temporary
		processed chan *QueueMessage

		// settings are used in messagebus for handling
		// throttling, polling settings
		settings QueueSettings
	}

	QueueSet map[string]*Queue
)

func (qs *QueueSet) toSlice() []*Queue {
	l := []*Queue{}

	for _, q := range *qs {
		l = append(l, q)
	}

	return l
}
