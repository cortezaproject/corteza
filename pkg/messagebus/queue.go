package messagebus

type (
	message struct {
		q string
		p []byte
	}

	Queue struct {
		consumer Consumer

		// settings are used in messagebus for handling
		// throttling, polling settings
		settings QueueSettings
	}

	QueueSet map[string]*Queue
)
