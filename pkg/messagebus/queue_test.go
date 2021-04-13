package messagebus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_QueueSet(t *testing.T) {
	req := require.New(t)

	fooQueue := &Queue{}
	barQueue := &Queue{}

	qs := QueueSet{
		"foo": fooQueue,
		"bar": barQueue,
	}

	req.Equal([]*Queue{fooQueue, barQueue}, qs.toSlice())
}
