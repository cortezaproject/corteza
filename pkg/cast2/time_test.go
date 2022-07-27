package cast2

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	var (
		req = require.New(t)

		target1 time.Time
		target2 = &time.Time{}

		in = "2006-01-02T15:04:05" // abuse format as value
	)

	req.NoError(Time(in, &target1))
	req.Equal(in, target1.Format(in))
	req.NotZero(target1)

	req.NoError(TimePtr(in, &target2))
	req.NotNil(target2)
	req.Equal(in, target2.Format(in))

	req.NoError(Time(nil, &target1))
	req.Zero(target1)

	req.NoError(TimePtr(nil, &target2))
	req.Nil(target2)
}
